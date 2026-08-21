package sm2

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"

	"github.com/aid297/aid/v3/secrets"
)

type SM2Sem struct {
	priKey secrets.SemenPriKey
	pubKey secrets.SemenPubKey
}

// NewSem 实例化：种子
func NewSem(attrs ...secrets.SemenAttr) (secrets.Semen, error) {
	var (
		err error
		sem secrets.Semen = &SM2Sem{}
	)

	if err = sem.SetAttrs(attrs...); err != nil {
		return nil, err
	}

	if sem.GetPriKey() == nil && sem.GetPubKey() == nil {
		if err = MustGeneratePriKey(sem); err != nil {
			return nil, err
		}
	}

	return sem, nil
}

// setAttr 设置属性
func (my *SM2Sem) SetAttrs(attrs ...secrets.SemenAttr) error {
	for _, attr := range attrs {
		if err := attr(my); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePriKey 生成私钥
func (my *SM2Sem) GeneratePriKey() (err error) {
	my.priKey, err = sm2.GenerateKey(rand.Reader)
	return
}

// SetPubKey 设置公钥：crypto.PublicKey(secrets.SemenPubKey)
func (my *SM2Sem) SetPubKey(pubKey secrets.SemenPubKey) (err error) { my.pubKey = pubKey; return }

// SetPubKeyBytes 设置公钥：bytes
func (my *SM2Sem) SetPubKeyBytes(pubKeyBytes []byte) (err error) {
	my.pubKey, err = x509.ParseSm2PublicKey(pubKeyBytes)
	return
}

// SetPubKeyBase64 设置公钥：base64
func (my *SM2Sem) SetPubKeyBase64(pubKeyBase64 string) (err error) {
	var pubKeyBytes []byte

	if pubKeyBytes, err = base64.StdEncoding.DecodeString(pubKeyBase64); err != nil {
		return
	}

	my.SetPubKeyBytes(pubKeyBytes)

	return
}

// GetPubKey 获取公钥：如果公钥存在则返回公钥，如果公钥不存在则使用私钥返回公钥
func (my *SM2Sem) GetPubKey() secrets.SemenPubKey {
	if my.pubKey != nil {
		return my.pubKey
	}

	if my.priKey != nil {
		return &my.priKey.(*sm2.PrivateKey).PublicKey
	}

	return nil
}

// GetPubKeyBytes 获取公钥：bytes
func (my *SM2Sem) GetPubKeyBytes() ([]byte, error) {
	var (
		err         error
		pubKeyBytes []byte
		pubKey      secrets.SemenPubKey
	)

	if pubKey = my.GetPubKey(); pubKey == nil {
		return nil, errors.New("公钥不存在")
	}

	if pubKeyBytes, err = x509.MarshalSm2PublicKey(pubKey.(*sm2.PublicKey)); err != nil {
		return nil, err
	}

	return pubKeyBytes, nil
}

// GetPubKeyBase64 获取公钥：base64
func (my *SM2Sem) GetPubKeyBase64() (string, error) {
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
func (my *SM2Sem) SetPriKey(priKey secrets.SemenPriKey) (err error) { my.priKey = priKey; return }

// SetPriKeyBytes 设置私钥：bytes（PKCS8 DER 格式，与 GetPriKeyBytes 对应）
func (my *SM2Sem) SetPriKeyBytes(priKeyBytes []byte) (err error) {
	if my.priKey, err = x509.ParsePKCS8UnecryptedPrivateKey(priKeyBytes); err != nil {
		return
	}
	return
}

// SetPriKeyBase64 设置私钥：base64
func (my *SM2Sem) SetPriKeyBase64(priKeyBase64 string) (err error) {
	var priKeyBytes []byte

	if priKeyBytes, err = base64.StdEncoding.DecodeString(priKeyBase64); err != nil {
		return
	}

	err = my.SetPriKeyBytes(priKeyBytes)

	return
}

// GetPriKey 获取私钥
func (my *SM2Sem) GetPriKey() secrets.SemenPriKey { return my.priKey }

// GetPriKeyBytes 获取私钥：bytes
func (my *SM2Sem) GetPriKeyBytes() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	return x509.MarshalSm2UnecryptedPrivateKey(my.priKey.(*sm2.PrivateKey))
}

// GetPriKeyBase64 获取私钥：base64
func (my *SM2Sem) GetPriKeyBase64() (string, error) {
	var (
		err         error
		priKeyBytes []byte
	)
	if priKeyBytes, err = my.GetPriKeyBytes(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(priKeyBytes), nil
}

// GetPubKeyPEM 获取公钥：PEM（PUBLIC KEY）
func (my *SM2Sem) GetPubKeyPEM() ([]byte, error) {
	pub := my.GetPubKey()
	if pub == nil {
		return nil, errors.New("公钥不存在")
	}
	return x509.WritePublicKeyToPem(pub.(*sm2.PublicKey))
}

// GetPriKeyPEM 获取私钥：PEM（PRIVATE KEY，国密 SM2 PKCS#8）
func (my *SM2Sem) GetPriKeyPEM() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}
	return x509.WritePrivateKeyToPem(my.priKey.(*sm2.PrivateKey), nil)
}

func PubKey(pubKey secrets.SemenPubKey) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPubKey(pubKey) }
}

func PubKeyBytes(pubKeyBytes []byte) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPubKeyBytes(pubKeyBytes) }
}

func PubKeyBase64(pubKeyBase64 string) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPubKeyBase64(pubKeyBase64) }
}

func PriKey(priKey secrets.SemenPriKey) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPriKey(priKey) }
}

func PriKeyBytes(priKeyBytes []byte) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPriKeyBytes(priKeyBytes) }
}

func PriKeyBase64(priKeyBase64 string) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPriKeyBase64(priKeyBase64) }
}

func MustGeneratePriKey(sem secrets.Semen) (err error) {
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
