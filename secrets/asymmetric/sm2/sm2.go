package sm2

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"

	"github.com/aid297/aid/v3/secrets"
)

var _ secrets.Asymmetric = (*SM2)(nil)

type SM2 struct{ sem secrets.Semen }

// New 实例化
func New(sem secrets.Semen) secrets.Asymmetric { return &SM2{sem: sem} }

// Encrypt 加密
func (my *SM2) Encrypt(plainText []byte) (string, error) {
	var (
		pubKeyBytes  []byte
		err          error
		pubKeyParsed *sm2.PublicKey
		cipherText   []byte
	)

	if pubKeyBytes, err = my.sem.GetPubKeyBytes(); err != nil {
		return "", err
	}

	if pubKeyParsed, err = x509.ParseSm2PublicKey(pubKeyBytes); err != nil {
		return "", err
	}

	if cipherText, err = sm2.Encrypt(pubKeyParsed, plainText, rand.Reader, sm2.C1C3C2); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt 使用 base64 私钥解密，输入 base64 密文
func (my *SM2) Decrypt(cipherBase64 string) ([]byte, error) {
	var (
		err          error
		priKeyBytes  []byte
		priKeyParsed *sm2.PrivateKey
		cipherText   []byte
	)

	if priKeyBytes, err = my.sem.GetPriKeyBytes(); err != nil {
		return nil, err
	}

	if priKeyParsed, err = x509.ParsePKCS8UnecryptedPrivateKey(priKeyBytes); err != nil {
		return nil, err
	}

	if cipherText, err = base64.StdEncoding.DecodeString(cipherBase64); err != nil {
		return nil, err
	}

	return sm2.Decrypt(priKeyParsed, cipherText, sm2.C1C3C2)
}

// Sign 使用 base64 私钥对数据签名，返回 hex 签名字符串
func (my *SM2) Sign(data []byte) (string, error) {
	var (
		err            error
		priKeyBytes    []byte
		priKeyParsed   *sm2.PrivateKey
		r, s           *big.Int
		sig            = make([]byte, 64)
		rBytes, sBytes []byte
	)

	if priKeyBytes, err = my.sem.GetPriKeyBytes(); err != nil {
		return "", err
	}

	if priKeyParsed, err = x509.ParsePKCS8UnecryptedPrivateKey(priKeyBytes); err != nil {
		return "", err
	}

	if r, s, err = sm2.Sm2Sign(priKeyParsed, data, nil, rand.Reader); err != nil {
		return "", err
	}

	// 将 r, s 编码为 64 字节（各 32 字节）hex 字符串
	rBytes = r.Bytes()
	sBytes = s.Bytes()
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)

	return hex.EncodeToString(sig), nil
}

// Verify 使用 base64 公钥验证 hex 签名
func (my *SM2) Verify(data []byte, sigHex string) (bool, error) {
	var (
		err          error
		pubKeyBytes  []byte
		pubKeyParsed *sm2.PublicKey
		sigBytes     []byte
		r, s         *big.Int
	)

	if pubKeyBytes, err = my.sem.GetPubKeyBytes(); err != nil {
		return false, err
	}

	if pubKeyParsed, err = x509.ParseSm2PublicKey(pubKeyBytes); err != nil {
		return false, err
	}

	if sigBytes, err = hex.DecodeString(sigHex); err != nil {
		return false, err
	}

	if len(sigBytes) != 64 {
		return false, errors.New("验证错误：签名长度错误")
	}

	r = new(big.Int).SetBytes(sigBytes[:32])
	s = new(big.Int).SetBytes(sigBytes[32:])

	return sm2.Sm2Verify(pubKeyParsed, data, nil, r, s), nil
}
