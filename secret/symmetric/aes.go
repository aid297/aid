package symmetric

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

type (
	// AES Aes密钥对象
	AES struct {
		Err     error
		Encrypt *AESEncrypt
		Decrypt *AESDecrypt
		sailStr string
	}

	// AESEncrypt Aes加密密钥对象
	AESEncrypt struct {
		Err      error
		sailStr  string
		sailByte []byte
		randKey  []byte
		aesKey   []byte
		openKey  string
	}

	// AESDecrypt Aes解密密钥对象
	AESDecrypt struct {
		Err      error
		sailStr  string
		sailByte []byte
		randKey  []byte
		aesKey   []byte
		openKey  string
	}
)

// NewAES 实例化：Aes密钥
func NewAES(sail string) *AES { return &AES{sailStr: sail} }

// NewEncrypt 实例化：Aes加密密钥对象
func (my *AES) NewEncrypt() *AES {
	my.Encrypt = NewAESEncrypt(my.sailStr)
	return my
}

// NewDecrypt 实例化：Aes解密密钥对象
func (my *AES) NewDecrypt(openKey string) *AES {
	my.Decrypt = NewAESDecrypt(my.sailStr, openKey)

	return my
}

// GetEncrypt 获取加密密钥
func (my *AES) GetEncrypt() *AESEncrypt { return my.Encrypt }

// GetDecrypt 获取解密密钥
func (my *AES) GetDecrypt() *AESDecrypt { return my.Decrypt }

// NewAESEncrypt 实例化：Aes加密密钥对象
func NewAESEncrypt(sail string) *AESEncrypt {
	aesHelper := &AESEncrypt{
		sailStr:  sail,
		sailByte: make([]byte, 16),
		randKey:  make([]byte, 16),
		aesKey:   make([]byte, 16),
		openKey:  "",
	}

	aesHelper.randKey = make([]byte, 16)
	_, aesHelper.Err = io.ReadFull(rand.Reader, aesHelper.randKey)
	aesHelper.sailByte, aesHelper.Err = base64.StdEncoding.DecodeString(sail)

	return aesHelper.sailByByte()
}

// sailByByte 密码加盐：使用byte盐
func (r *AESEncrypt) sailByByte() *AESEncrypt {
	copy(r.aesKey, r.randKey)

	for i := 0; i < 4; i++ {
		index := int(r.randKey[i]) % 16
		r.aesKey[index] = r.sailByte[index]
	}

	r.openKey = base64.StdEncoding.EncodeToString(r.randKey)

	return r
}

// GetAesKey 获取加盐后的密钥
func (r *AESEncrypt) GetAesKey() []byte { return r.aesKey }

// SetAesKey 设置加盐后的密钥
func (r *AESEncrypt) SetAesKey(aesKey []byte) *AESEncrypt {
	r.aesKey = aesKey
	return r
}

// GetOpenKey 获取公开密码
func (r *AESEncrypt) GetOpenKey() string { return r.openKey }

// NewAESDecrypt 实例化：Aes解密密钥对象
func NewAESDecrypt(sailStr, openKey string) *AESDecrypt {
	aesDecrypt := &AESDecrypt{
		sailStr:  sailStr,
		sailByte: make([]byte, 16),
		randKey:  make([]byte, 16),
		aesKey:   make([]byte, 16),
		openKey:  openKey,
	}

	aesDecrypt.randKey, aesDecrypt.Err = base64.StdEncoding.DecodeString(openKey)
	copy(aesDecrypt.aesKey, aesDecrypt.randKey)
	aesDecrypt.sailByte, aesDecrypt.Err = base64.StdEncoding.DecodeString(sailStr)

	return aesDecrypt.deSailByByte()
}

// deSailByByte 密码解盐：使用byte盐
func (r *AESDecrypt) deSailByByte() *AESDecrypt {
	index := r.randKey[:4]

	// 替换key中的字节
	for _, x := range index {
		i := int(x) % 16
		r.aesKey[i] = r.sailByte[i]
	}

	return r
}

// GetAesKey 获取加盐后的密钥
func (r *AESDecrypt) GetAesKey() []byte {
	return r.aesKey
}

// SetAesKey 设置加盐后的密钥
func (r *AESDecrypt) SetAesKey(aesKey []byte) *AESDecrypt {
	r.aesKey = aesKey
	return r
}

// GetOpenKey 获取公开密码
func (r *AESDecrypt) GetOpenKey() string {
	return r.openKey
}
