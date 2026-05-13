package rsa

import (
	"crypto/rand"
	cryptorsa "crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"

	"crypto/x509"

	"github.com/aid297/aid/v2/secret"

	"encoding/pem"
)

type RSASemImpl struct{ priKey secret.SemenerPriKey }

// NewSem 实例化：种子
func NewSem(attrs ...secret.SemenerAttr) (secret.Semener, error) {
	var (
		err error
		sem secret.Semener = &RSASemImpl{}
	)

	if err = sem.SetAttrs(attrs...); err != nil {
		return nil, err
	}

	if sem.GetPriKey() == nil {
		if err = MustGeneratePriKey(sem); err != nil {
			return nil, err
		}
	}

	return sem, nil
}

// setAttr 设置属性
func (my *RSASemImpl) SetAttrs(attrs ...secret.SemenerAttr) error {
	for _, attr := range attrs {
		if err := attr(my); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePriKey 生成私钥
func (my *RSASemImpl) GeneratePriKey() (err error) {
	if my.priKey, err = cryptorsa.GenerateKey(rand.Reader, 2048); err != nil {
		return
	}
	return
}

// GetPriKey 获取私钥
func (my *RSASemImpl) GetPriKey() secret.SemenerPriKey { return my.priKey }

// GetPriKeyBytes 获取私钥：bytes
func (my *RSASemImpl) GetPriKeyBytes() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	return x509.MarshalPKCS8PrivateKey(my.priKey.(*cryptorsa.PrivateKey))
}

// GetPriKeyBase64 获取私钥：base64
func (my *RSASemImpl) GetPriKeyBase64() (string, error) {
	var (
		err         error
		priKeyBytes []byte
	)
	if priKeyBytes, err = my.GetPriKeyBytes(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(priKeyBytes), nil
}

// GetPubKey 获取公钥
func (my *RSASemImpl) GetPubKey() secret.SemenerPubKey {
	return &my.priKey.(*cryptorsa.PrivateKey).PublicKey
}

// GetPubKeyBytes 获取公钥：bytes
func (my *RSASemImpl) GetPubKeyBytes() ([]byte, error) {
	var (
		err         error
		pubKeyBytes []byte
	)

	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	if pubKeyBytes, err = x509.MarshalPKIXPublicKey(&my.priKey.(*cryptorsa.PrivateKey).PublicKey); err != nil {
		return nil, err
	}

	return pubKeyBytes, nil
}

// GetPubKeyBase64 获取公钥：base64
func (my *RSASemImpl) GetPubKeyBase64() (string, error) {
	var (
		err         error
		pubKeyBytes []byte
	)

	if pubKeyBytes, err = my.GetPubKeyBytes(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pubKeyBytes), nil
}

// SetPriKey
func (my *RSASemImpl) SetPriKey(priKey secret.SemenerPriKey) (err error) { my.priKey = priKey; return }

// SetPriKeyBytes 设置私钥：bytes（PKCS8 DER 格式，与 GetPriKeyBytes 对应）
func (my *RSASemImpl) SetPriKeyBytes(priKeyBytes []byte) (err error) {
	var (
		pemBlock  *pem.Block
		priPKCS8  any
		ok        bool
		priKeyPEM []byte
	)
	if priKeyPEM, err = decodeBytesKeyToPEM(priKeyBytes, false); err != nil {
		return errors.New("私钥(bytes)内容不合法")
	}

	pemBlock, _ = pem.Decode(priKeyPEM)
	if pemBlock == nil {
		return errors.New("私钥(PEM)内容不合法")
	}
	if my.priKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err == nil {
		return fmt.Errorf("解析PEM失败：%v", err)
	}
	if priPKCS8, err = x509.ParsePKCS8PrivateKey(pemBlock.Bytes); err != nil {
		return fmt.Errorf("解析PKCS8失败：%v", err)
	}
	if my.priKey, ok = priPKCS8.(*cryptorsa.PrivateKey); !ok {
		return errors.New("私钥格式不合法：PKCS8")
	}

	return
}

// SetPriKeyBase64 设置私钥：base64
func (my *RSASemImpl) SetPriKeyBase64(priKeyBase64 string) (err error) {
	var priKeyBytes []byte

	if priKeyBytes, err = base64.StdEncoding.DecodeString(priKeyBase64); err != nil {
		return
	}

	err = my.SetPriKeyBytes(priKeyBytes)

	return
}

func PriKey(priKey secret.SemenerPriKey) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPriKey(priKey) }
}

func PriKeyBytes(priKeyBytes []byte) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPriKeyBytes(priKeyBytes) }
}

func PriKeyBase64(priKeyBase64 string) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPriKeyBase64(priKeyBase64) }
}

func MustGeneratePriKey(sem secret.Semener) (err error) {
	var (
		pubKeyBase64, priKeyBase64 string
		pubKeyBytes, priKeyBytes   []byte
	)

	if err = sem.GeneratePriKey(); err != nil {
		return err
	}

	if pubKeyBase64, err = sem.GetPubKeyBase64(); pubKeyBase64 == "" {
		return errors.New("公钥长度为空：base64")
	}

	if pubKeyBytes, err = sem.GetPubKeyBytes(); len(pubKeyBytes) == 0 {
		return errors.New("公钥长度为空：bytes")
	}

	if priKeyBase64, err = sem.GetPriKeyBase64(); priKeyBase64 == "" {
		return errors.New("私钥长度为空：base64")
	}

	if priKeyBytes, err = sem.GetPriKeyBytes(); len(priKeyBytes) == 0 {
		return errors.New("私钥长度为空：bytes")
	}

	return nil
}

func decodeBytesKeyToPEM(bytesKey []byte, isPublic bool) ([]byte, error) {
	if block, _ := pem.Decode(bytesKey); block != nil {
		return bytesKey, nil
	}
	if isPublic {
		if _, err := x509.ParsePKIXPublicKey(bytesKey); err != nil {
			return nil, fmt.Errorf("rsa: invalid public key der: %w", err)
		}
		return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: bytesKey}), nil
	}
	if _, err := x509.ParsePKCS8PrivateKey(bytesKey); err == nil {
		return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: bytesKey}), nil
	}
	if _, err := x509.ParsePKCS1PrivateKey(bytesKey); err == nil {
		return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: bytesKey}), nil
	}
	return nil, errors.New("rsa: invalid private key der")
}

func decodeBase64KeyToPEM(base64Key string, isPublic bool) ([]byte, error) {
	derOrPem, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	return decodeBytesKeyToPEM(derOrPem, isPublic)
}
