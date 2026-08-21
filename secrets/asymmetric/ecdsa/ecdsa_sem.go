package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"

	"github.com/aid297/aid/v3/secrets"
)

type ECDSASem struct {
	priKey secrets.SemenPriKey
	pubKey secrets.SemenPubKey
}

// NewSem 实例化：种子
// 默认使用 P-256 曲线（ES256）
func NewSem(attrs ...secrets.SemenAttr) (secrets.Semen, error) {
	var (
		err error
		sem secrets.Semen = &ECDSASem{}
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
func (my *ECDSASem) SetAttrs(attrs ...secrets.SemenAttr) error {
	for _, attr := range attrs {
		if err := attr(my); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePriKey 生成私钥（P-256 曲线）
func (my *ECDSASem) GeneratePriKey() (err error) {
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	my.priKey = priKey
	return nil
}

// SetPubKey 设置公钥
func (my *ECDSASem) SetPubKey(pubKey secrets.SemenPubKey) (err error) {
	my.pubKey = pubKey
	return
}

// SetPubKeyBytes 设置公钥：bytes
func (my *ECDSASem) SetPubKeyBytes(pubKeyBytes []byte) (err error) {
	var pubKey *ecdsa.PublicKey
	if pubKey, err = parseECDSAPublicKey(pubKeyBytes); err != nil {
		return errors.New("公钥(bytes)内容不合法")
	}
	my.pubKey = pubKey
	return nil
}

// SetPubKeyBase64 设置公钥：base64
func (my *ECDSASem) SetPubKeyBase64(pubKeyBase64 string) (err error) {
	var pubKeyBytes []byte

	if pubKeyBytes, err = base64.StdEncoding.DecodeString(pubKeyBase64); err != nil {
		return
	}

	err = my.SetPubKeyBytes(pubKeyBytes)

	return
}

// GetPubKey 获取公钥：如果公钥存在则返回公钥，如果公钥不存在则使用私钥返回公钥
func (my *ECDSASem) GetPubKey() secrets.SemenPubKey {
	if my.pubKey != nil {
		return my.pubKey
	}

	if my.priKey != nil {
		return &my.priKey.(*ecdsa.PrivateKey).PublicKey
	}

	return nil
}

// GetPubKeyBytes 获取公钥：bytes
func (my *ECDSASem) GetPubKeyBytes() ([]byte, error) {
	pubKey := my.GetPubKey()
	if pubKey == nil {
		return nil, errors.New("公钥不存在")
	}

	return x509.MarshalPKIXPublicKey(pubKey)
}

// GetPubKeyBase64 获取公钥：base64
func (my *ECDSASem) GetPubKeyBase64() (string, error) {
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
func (my *ECDSASem) SetPriKey(priKey secrets.SemenPriKey) (err error) { my.priKey = priKey; return }

// SetPriKeyBytes 设置私钥：bytes
func (my *ECDSASem) SetPriKeyBytes(priKeyBytes []byte) (err error) {
	var priKey *ecdsa.PrivateKey
	if priKey, err = parseECDSAPrivateKey(priKeyBytes); err != nil {
		return errors.New("私钥(bytes)内容不合法")
	}
	my.priKey = priKey
	return nil
}

// SetPriKeyBase64 设置私钥：base64
func (my *ECDSASem) SetPriKeyBase64(priKeyBase64 string) (err error) {
	var priKeyBytes []byte

	if priKeyBytes, err = base64.StdEncoding.DecodeString(priKeyBase64); err != nil {
		return
	}

	err = my.SetPriKeyBytes(priKeyBytes)

	return
}

// GetPriKey 获取私钥
func (my *ECDSASem) GetPriKey() secrets.SemenPriKey { return my.priKey }

// GetPriKeyBytes 获取私钥：bytes
func (my *ECDSASem) GetPriKeyBytes() ([]byte, error) {
	if my.priKey == nil {
		return nil, errors.New("私钥不能为空")
	}

	return x509.MarshalPKCS8PrivateKey(my.priKey)
}

// GetPriKeyBase64 获取私钥：base64
func (my *ECDSASem) GetPriKeyBase64() (string, error) {
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
func (my *ECDSASem) GetPubKeyPEM() ([]byte, error) {
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
func (my *ECDSASem) GetPriKeyPEM() ([]byte, error) {
	der, err := my.GetPriKeyBytes()
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

func parseECDSAPublicKey(pubKeyBytes []byte) (*ecdsa.PublicKey, error) {
	// 尝试解析 PEM 格式
	if block, _ := pem.Decode(pubKeyBytes); block != nil {
		pubKeyBytes = block.Bytes
	}

	pubAny, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	pub, ok := pubAny.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("不是 ECDSA 公钥")
	}

	return pub, nil
}

func parseECDSAPrivateKey(priKeyBytes []byte) (*ecdsa.PrivateKey, error) {
	// 尝试解析 PEM 格式
	if block, _ := pem.Decode(priKeyBytes); block != nil {
		priKeyBytes = block.Bytes
	}

	priPKCS8, err := x509.ParsePKCS8PrivateKey(priKeyBytes)
	if err != nil {
		return nil, err
	}

	pri, ok := priPKCS8.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("不是 ECDSA 私钥")
	}

	return pri, nil
}

// PointAtInfinity 表示无穷远点
type ecdsaPoint struct {
	X, Y *big.Int
}

var pointAtInfinity = &ecdsaPoint{X: big.NewInt(0), Y: big.NewInt(0)}

func (p *ecdsaPoint) isInfinity() bool {
	return p.X.Sign() == 0 && p.Y.Sign() == 0
}
