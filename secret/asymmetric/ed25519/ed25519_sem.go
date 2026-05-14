package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/aid297/aid/v2/secret"
)

type Ed25519SemImpl struct {
	priKey secret.SemenerPriKey
	pubKey secret.SemenerPubKey
}

// NewSem 实例化：种子
func NewSem(attrs ...secret.SemenerAttr) (secret.Semener, error) {
	var (
		err error
		sem secret.Semener = &Ed25519SemImpl{}
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
func (my *Ed25519SemImpl) SetAttrs(attrs ...secret.SemenerAttr) error {
	for _, attr := range attrs {
		if err := attr(my); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePriKey 生成私钥
func (my *Ed25519SemImpl) GeneratePriKey() (err error) {
	_, priKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	my.priKey = priKey
	return nil
}

// SetPubKey 设置公钥
func (my *Ed25519SemImpl) SetPubKey(pubKey secret.SemenerPubKey) (err error) { my.pubKey = pubKey; return }

// SetPubKeyBytes 设置公钥：bytes
func (my *Ed25519SemImpl) SetPubKeyBytes(pubKeyBytes []byte) (err error) {
	pubKey, err := parseEd25519PublicKey(pubKeyBytes)
	if err != nil {
		return errors.New("公钥(bytes)内容不合法")
	}
	my.pubKey = pubKey
	return nil
}

// SetPubKeyBase64 设置公钥：base64
func (my *Ed25519SemImpl) SetPubKeyBase64(pubKeyBase64 string) (err error) {
	var pubKeyBytes []byte

	if pubKeyBytes, err = base64.StdEncoding.DecodeString(pubKeyBase64); err != nil {
		return
	}

	err = my.SetPubKeyBytes(pubKeyBytes)

	return
}

// GetPubKey 获取公钥：如果公钥存在则返回公钥，如果公钥不存在则使用私钥返回公钥
func (my *Ed25519SemImpl) GetPubKey() secret.SemenerPubKey {
	if my.pubKey != nil {
		return my.pubKey
	}

	if my.priKey != nil {
		return my.priKey.(ed25519.PrivateKey).Public().(ed25519.PublicKey)
	}

	return nil
}

// GetPubKeyBytes 获取公钥：bytes
func (my *Ed25519SemImpl) GetPubKeyBytes() ([]byte, error) {
	pubKey := my.GetPubKey()
	if pubKey == nil {
		return nil, errors.New("公钥不存在")
	}

	return pubKey.(ed25519.PublicKey), nil
}

// GetPubKeyBase64 获取公钥：base64
func (my *Ed25519SemImpl) GetPubKeyBase64() (string, error) {
	var (
		err         error
		pubKeyBytes []byte
	)

	if pubKeyBytes, err = my.GetPubKeyBytes(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pubKeyBytes), nil
}

// SetPriKey 设置私钥
func (my *Ed25519SemImpl) SetPriKey(priKey secret.SemenerPriKey) (err error) { my.priKey = priKey; return }

// SetPriKeyBytes 设置私钥：bytes
func (my *Ed25519SemImpl) SetPriKeyBytes(priKeyBytes []byte) (err error) {
	priKey, err := parseEd25519PrivateKey(priKeyBytes)
	if err != nil {
		return errors.New("私钥(bytes)内容不合法")
	}
	my.priKey = priKey
	return nil
}

// SetPriKeyBase64 设置私钥：base64
func (my *Ed25519SemImpl) SetPriKeyBase64(priKeyBase64 string) (err error) {
	var priKeyBytes []byte

	if priKeyBytes, err = base64.StdEncoding.DecodeString(priKeyBase64); err != nil {
		return
	}

	err = my.SetPriKeyBytes(priKeyBytes)

	return
}

// GetPriKey 获取私钥
func (my *Ed25519SemImpl) GetPriKey() secret.SemenerPriKey { return my.priKey }

// GetPriKeyBytes 获取私钥：bytes
func (my *Ed25519SemImpl) GetPriKeyBytes() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	return my.priKey.(ed25519.PrivateKey), nil
}

// GetPriKeyBase64 获取私钥：base64
func (my *Ed25519SemImpl) GetPriKeyBase64() (string, error) {
	var (
		err         error
		priKeyBytes []byte
	)

	if priKeyBytes, err = my.GetPriKeyBytes(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(priKeyBytes), nil
}

// PubKey 设置公钥属性
func PubKey(pubKey secret.SemenerPubKey) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPubKey(pubKey) }
}

// PubKeyBytes 设置公钥属性：bytes
func PubKeyBytes(pubKeyBytes []byte) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPubKeyBytes(pubKeyBytes) }
}

// PubKeyBase64 设置公钥属性：base64
func PubKeyBase64(pubKeyBase64 string) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPubKeyBase64(pubKeyBase64) }
}

// PriKey 设置私钥属性
func PriKey(priKey secret.SemenerPriKey) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPriKey(priKey) }
}

// PriKeyBytes 设置私钥属性：bytes
func PriKeyBytes(priKeyBytes []byte) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPriKeyBytes(priKeyBytes) }
}

// PriKeyBase64 设置私钥属性：base64
func PriKeyBase64(priKeyBase64 string) secret.SemenerAttr {
	return func(sem secret.Semener) error { return sem.SetPriKeyBase64(priKeyBase64) }
}

// MustGeneratePriKey 生成私钥并验证
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

// ==================== 内部工具函数 ====================

func parseEd25519PublicKey(pubKeyBytes []byte) (ed25519.PublicKey, error) {
	// 尝试解析 PEM 格式
	if block, _ := pem.Decode(pubKeyBytes); block != nil {
		pubKeyBytes = block.Bytes
	}

	// 尝试解析 PKIX 格式
	pubAny, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err == nil {
		pub, ok := pubAny.(ed25519.PublicKey)
		if ok {
			return pub, nil
		}
	}

	// 尝试直接作为 Ed25519 公钥解析
	if len(pubKeyBytes) == ed25519.PublicKeySize {
		return ed25519.PublicKey(pubKeyBytes), nil
	}

	return nil, errors.New("无效的 Ed25519 公钥")
}

func parseEd25519PrivateKey(priKeyBytes []byte) (ed25519.PrivateKey, error) {
	// 尝试解析 PEM 格式
	if block, _ := pem.Decode(priKeyBytes); block != nil {
		priKeyBytes = block.Bytes
	}

	// 尝试解析 PKCS8 格式
	priPKCS8, err := x509.ParsePKCS8PrivateKey(priKeyBytes)
	if err == nil {
		pri, ok := priPKCS8.(ed25519.PrivateKey)
		if ok {
			return pri, nil
		}
	}

	// 尝试直接作为 Ed25519 私钥解析
	if len(priKeyBytes) == ed25519.PrivateKeySize {
		return ed25519.PrivateKey(priKeyBytes), nil
	}

	return nil, errors.New("无效的 Ed25519 私钥")
}
