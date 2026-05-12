package sm2

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/aid297/aid/v2/secret"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

type (
	SM2SemImpl struct{ priKey *sm2.PrivateKey }
)

// NewSem 实例化：种子
func NewSem(attrs ...secret.SemenerAttr) (secret.Semener, error) {
	var (
		err error
		sem secret.Semener = &SM2SemImpl{}
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
func (my *SM2SemImpl) SetAttrs(attrs ...secret.SemenerAttr) error {
	for _, attr := range attrs {
		if err := attr(my); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePriKey 生成私钥
func (my *SM2SemImpl) GeneratePriKey() (err error) {
	my.priKey, err = sm2.GenerateKey(rand.Reader)
	return
}

// GetPriKey 获取私钥
func (my *SM2SemImpl) GetPriKey() *sm2.PrivateKey { return my.priKey }

// GetPriKeyBytes 获取私钥：bytes
func (my *SM2SemImpl) GetPriKeyBytes() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	return x509.MarshalSm2UnecryptedPrivateKey(my.priKey)
}

// GetPriKeyBase64 获取私钥：base64
func (my *SM2SemImpl) GetPriKeyBase64() (string, error) {
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
func (my *SM2SemImpl) GetPubKey() *sm2.PublicKey { return &my.priKey.PublicKey }

// GetPubKeyBytes 获取公钥：bytes
func (my *SM2SemImpl) GetPubKeyBytes() ([]byte, error) {
	var (
		err         error
		pubKeyBytes []byte
	)

	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	if pubKeyBytes, err = x509.MarshalSm2PublicKey(&my.priKey.PublicKey); err != nil {
		return nil, err
	}

	return pubKeyBytes, nil
}

// GetPubKeyBase64 获取公钥：base64
func (my *SM2SemImpl) GetPubKeyBase64() (string, error) {
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
func (my *SM2SemImpl) SetPriKey(priKey *sm2.PrivateKey) (err error) { my.priKey = priKey; return }

// SetPriKeyBytes 设置私钥：bytes（PKCS8 DER 格式，与 GetPriKeyBytes 对应）
func (my *SM2SemImpl) SetPriKeyBytes(priKeyBytes []byte) (err error) {
	if my.priKey, err = x509.ParsePKCS8UnecryptedPrivateKey(priKeyBytes); err != nil {
		return
	}
	return
}

// SetPriKeyBase64 设置私钥：base64
func (my *SM2SemImpl) SetPriKeyBase64(priKeyBase64 string) (err error) {
	var priKeyBytes []byte

	if priKeyBytes, err = base64.StdEncoding.DecodeString(priKeyBase64); err != nil {
		return
	}

	err = my.SetPriKeyBytes(priKeyBytes)

	return
}

func PriKey(priKey *sm2.PrivateKey) secret.SemenerAttr {
	return func(sm2Sem secret.Semener) error { return sm2Sem.SetPriKey(priKey) }
}

func PriKeyBytes(priKeyBytes []byte) secret.SemenerAttr {
	return func(sm2Sem secret.Semener) error { return sm2Sem.SetPriKeyBytes(priKeyBytes) }
}

func PriKeyBase64(priKeyBase64 string) secret.SemenerAttr {
	return func(sm2Sem secret.Semener) error { return sm2Sem.SetPriKeyBase64(priKeyBase64) }
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
