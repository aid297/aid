package secret

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"

	"github.com/aid297/aid/compression"

	"github.com/aid297/aid/common"
	"github.com/aid297/aid/secret/symmetric"
	"github.com/aid297/aid/str"
)

func EncryptAuthorization(key, secretKey string, iv []byte, randStr ...string) (string, string, error) {
	var (
		err   error
		uuid  string
		token []byte
	)

	if key == "" {
		return "", "", err
	}

	// 生成随机串
	if len(randStr) > 0 {
		uuid = randStr[0]
	} else {
		uuid, err = MustEncrypt(str.NewRand().GetLetters(10))
		if err != nil {
			return "", "", err
		}
	}

	token, err = symmetric.NewCBC().Encrypt([]byte(key+uuid), []byte(secretKey), iv)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(token), uuid, nil
}

func DecryptAuthorization(token, secretKey string, iv []byte) (string, string, error) {
	var (
		err                   error
		token64, decryptToken []byte
	)

	if token == "" {
		return "", "", errors.New("token 不能为空")
	}

	if token64, err = base64.StdEncoding.DecodeString(token); err != nil {
		return "", "", fmt.Errorf("base64解码token失败：%s", err.Error())
	}
	if decryptToken, err = symmetric.NewCBC().Decrypt(token64, []byte(secretKey), iv); err != nil {
		return "", "", fmt.Errorf("解密失败：%s", err.Error())
	}

	return string(decryptToken[:len(decryptToken)-32]), string(decryptToken[len(decryptToken)-32:]), nil
}

func MustEncrypt(data any) (string, error) {
	var (
		err       error
		dataBytes []byte
		h         hash.Hash
	)
	dataBytes = common.ToBytes(data)

	h = md5.New()
	if _, err = h.Write(dataBytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func ECB16Encrypt(data any, needEncrypt, needZip bool, aes *symmetric.AES) (string, error) {
	var (
		err         error
		jsonByte, b []byte
	)

	// json序列化
	if jsonByte, err = json.Marshal(data); err != nil {
		return "", err
	}

	// 压缩
	if needZip {
		if b, err = compression.NewZlib().Compress(jsonByte); err != nil {
			return "", err
		}
	}

	// 加密
	if needEncrypt {
		if b, err = symmetric.NewECB().Encrypt(b, aes.Encrypt.GetAesKey()); err != nil {
			return "", err
		}
	}

	if !needZip && !needEncrypt {
		return string(b), nil
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func ECB16Decrypt(data string, needEncrypt, needZip bool, aes *symmetric.AES) (any, error) {
	var (
		err                                         error
		r                                           any
		cipherText, decryptedByte, decompressedByte []byte
	)

	if needEncrypt {
		// base64 解码
		if cipherText, err = base64.StdEncoding.DecodeString(data); err != nil {
			return nil, err
		}

		// aes解密：ecb
		if decryptedByte, err = symmetric.NewECB().Decrypt(cipherText, aes.Encrypt.GetAesKey()); err != nil {
			return nil, err
		}

		// 解压
		if needZip {
			if decompressedByte, err = compression.NewZlib().Decompress(decryptedByte); err != nil {
				return nil, err
			}

			if err = json.Unmarshal(decompressedByte, &r); err != nil {
				return nil, err
			}

			return r, nil
		}

		// 将data反序列化
		if err = json.Unmarshal(decryptedByte, &r); err != nil {
			return nil, err
		}

		return r, nil
	}

	if err = json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return r, nil
}
