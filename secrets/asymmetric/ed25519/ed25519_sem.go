package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/aid297/aid/v2/secrets"
)

type Ed25519Sem struct {
	priKey secrets.SemenPriKey
	pubKey secrets.SemenPubKey
}

// NewSem 实例化：种子
func NewSem(attrs ...secrets.SemenAttr) (secrets.Semen, error) {
	var (
		err error
		sem secrets.Semen = &Ed25519Sem{}
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
func (my *Ed25519Sem) SetAttrs(attrs ...secrets.SemenAttr) error {
	for _, attr := range attrs {
		if err := attr(my); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePriKey 生成私钥
func (my *Ed25519Sem) GeneratePriKey() (err error) {
	_, priKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	my.priKey = priKey
	return nil
}

// SetPubKey 设置公钥
func (my *Ed25519Sem) SetPubKey(pubKey secrets.SemenPubKey) (err error) {
	my.pubKey = pubKey
	return
}

// SetPubKeyBytes 设置公钥：bytes
func (my *Ed25519Sem) SetPubKeyBytes(pubKeyBytes []byte) (err error) {
	pubKey, err := parseEd25519PublicKey(pubKeyBytes)
	if err != nil {
		return errors.New("公钥(bytes)内容不合法")
	}
	my.pubKey = pubKey
	return nil
}

// SetPubKeyBase64 设置公钥：base64
func (my *Ed25519Sem) SetPubKeyBase64(pubKeyBase64 string) (err error) {
	var pubKeyBytes []byte

	if pubKeyBytes, err = base64.StdEncoding.DecodeString(pubKeyBase64); err != nil {
		return
	}

	err = my.SetPubKeyBytes(pubKeyBytes)

	return
}

// GetPubKey 获取公钥：如果公钥存在则返回公钥，如果公钥不存在则使用私钥返回公钥
func (my *Ed25519Sem) GetPubKey() secrets.SemenPubKey {
	if my.pubKey != nil {
		return my.pubKey
	}

	if my.priKey != nil {
		return my.priKey.(ed25519.PrivateKey).Public().(ed25519.PublicKey)
	}

	return nil
}

// GetPubKeyBytes 获取公钥：bytes
func (my *Ed25519Sem) GetPubKeyBytes() ([]byte, error) {
	pubKey := my.GetPubKey()
	if pubKey == nil {
		return nil, errors.New("公钥不存在")
	}

	return pubKey.(ed25519.PublicKey), nil
}

// GetPubKeyBase64 获取公钥：base64
func (my *Ed25519Sem) GetPubKeyBase64() (string, error) {
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
func (my *Ed25519Sem) SetPriKey(priKey secrets.SemenPriKey) (err error) {
	my.priKey = priKey
	return
}

// SetPriKeyBytes 设置私钥：bytes
func (my *Ed25519Sem) SetPriKeyBytes(priKeyBytes []byte) (err error) {
	priKey, err := parseEd25519PrivateKey(priKeyBytes)
	if err != nil {
		return errors.New("私钥(bytes)内容不合法")
	}
	my.priKey = priKey
	return nil
}

// SetPriKeyBase64 设置私钥：base64
func (my *Ed25519Sem) SetPriKeyBase64(priKeyBase64 string) (err error) {
	var priKeyBytes []byte

	if priKeyBytes, err = base64.StdEncoding.DecodeString(priKeyBase64); err != nil {
		return
	}

	err = my.SetPriKeyBytes(priKeyBytes)

	return
}

// GetPriKey 获取私钥
func (my *Ed25519Sem) GetPriKey() secrets.SemenPriKey { return my.priKey }

// GetPriKeyBytes 获取私钥：bytes
func (my *Ed25519Sem) GetPriKeyBytes() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	return my.priKey.(ed25519.PrivateKey), nil
}

// GetPriKeyBase64 获取私钥：base64
func (my *Ed25519Sem) GetPriKeyBase64() (string, error) {
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
func (my *Ed25519Sem) GetPubKeyPEM() ([]byte, error) {
	pub := my.GetPubKey()
	if pub == nil {
		return nil, errors.New("公钥不存在")
	}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}

// GetPriKeyPEM 获取私钥：PEM（PRIVATE KEY，PKCS#8）
func (my *Ed25519Sem) GetPriKeyPEM() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}
	der, err := x509.MarshalPKCS8PrivateKey(my.priKey.(ed25519.PrivateKey))
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// PubKey 设置公钥属性
func PubKey(pubKey secrets.SemenPubKey) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPubKey(pubKey) }
}

// PubKeyBytes 设置公钥属性：bytes
func PubKeyBytes(pubKeyBytes []byte) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPubKeyBytes(pubKeyBytes) }
}

// PubKeyBase64 设置公钥属性：base64
func PubKeyBase64(pubKeyBase64 string) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPubKeyBase64(pubKeyBase64) }
}

// PriKey 设置私钥属性
func PriKey(priKey secrets.SemenPriKey) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPriKey(priKey) }
}

// PriKeyBytes 设置私钥属性：bytes
func PriKeyBytes(priKeyBytes []byte) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPriKeyBytes(priKeyBytes) }
}

// PriKeyBase64 设置私钥属性：base64
func PriKeyBase64(priKeyBase64 string) secrets.SemenAttr {
	return func(sem secrets.Semen) error { return sem.SetPriKeyBase64(priKeyBase64) }
}

// MustGeneratePriKey 生成私钥并验证
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
