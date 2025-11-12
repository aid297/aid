package symmetric

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

type (
	// Aes Aes密钥对象
	Aes struct {
		Err     error
		Encrypt *AesEncrypt
		Decrypt *AesDecrypt
		sailStr string
	}

	// AesEncrypt Aes加密密钥对象
	AesEncrypt struct {
		Err      error
		sailStr  string
		sailByte []byte
		randKey  []byte
		aesKey   []byte
		openKey  string
	}

	// AesDecrypt Aes解密密钥对象
	AesDecrypt struct {
		Err      error
		sailStr  string
		sailByte []byte
		randKey  []byte
		aesKey   []byte
		openKey  string
	}
)

var AesApp Aes

func (*Aes) New(sail string) *Aes { return &Aes{sailStr: sail} }

// NewAes 实例化：Aes密钥
//
//go:fix 推荐使用：New方法
func NewAes(sail string) *Aes { return &Aes{sailStr: sail} }

// NewEncrypt 实例化：Aes加密密钥对象
func (my *Aes) NewEncrypt() *Aes {
	my.Encrypt = NewAesEncrypt(my.sailStr)

	return my
}

// NewDecrypt 实例化：Aes解密密钥对象
func (my *Aes) NewDecrypt(openKey string) *Aes {
	my.Decrypt = NewAesDecrypt(my.sailStr, openKey)

	return my
}

// GetEncrypt 获取加密密钥
func (my *Aes) GetEncrypt() *AesEncrypt { return my.Encrypt }

// GetDecrypt 获取解密密钥
func (my *Aes) GetDecrypt() *AesDecrypt { return my.Decrypt }

// NewAesEncrypt 实例化：Aes加密密钥对象
func NewAesEncrypt(sail string) *AesEncrypt {
	aesHelper := &AesEncrypt{
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
func (r *AesEncrypt) sailByByte() *AesEncrypt {
	copy(r.aesKey, r.randKey)

	for i := 0; i < 4; i++ {
		index := int(r.randKey[i]) % 16
		r.aesKey[index] = r.sailByte[index]
	}

	r.openKey = base64.StdEncoding.EncodeToString(r.randKey)

	return r
}

// GetAesKey 获取加盐后的密钥
func (r *AesEncrypt) GetAesKey() []byte { return r.aesKey }

// SetAesKey 设置加盐后的密钥
func (r *AesEncrypt) SetAesKey(aesKey []byte) *AesEncrypt {
	r.aesKey = aesKey

	return r
}

// GetOpenKey 获取公开密码
func (r *AesEncrypt) GetOpenKey() string { return r.openKey }

// NewAesDecrypt 实例化：Aes解密密钥对象
func NewAesDecrypt(sailStr, openKey string) *AesDecrypt {
	aesDecrypt := &AesDecrypt{
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
func (r *AesDecrypt) deSailByByte() *AesDecrypt {
	index := r.randKey[:4]

	// 替换key中的字节
	for _, x := range index {
		i := int(x) % 16
		r.aesKey[i] = r.sailByte[i]
	}

	return r
}

// GetAesKey 获取加盐后的密钥
func (r *AesDecrypt) GetAesKey() []byte {
	return r.aesKey
}

// SetAesKey 设置加盐后的密钥
func (r *AesDecrypt) SetAesKey(aesKey []byte) *AesDecrypt {
	r.aesKey = aesKey
	return r
}

// GetOpenKey 获取公开密码
func (r *AesDecrypt) GetOpenKey() string {
	return r.openKey
}
