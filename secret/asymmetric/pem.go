package asymmetric

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type PEMBase64 struct {
	base64PublicKey  string
	base64PrivateKey string
	publicKey        []byte
	privateKey       []byte
}

// NewPEMBase64 实例化
func NewPEMBase64() *PEMBase64 { return &PEMBase64{} }

func (my *PEMBase64) SetBase64PublicKey(base64PublicKey string) *PEMBase64 {
	my.base64PublicKey = base64PublicKey

	return my
}

func (my *PEMBase64) SetBase64PrivateKye(base64PrivateKey string) *PEMBase64 {
	my.base64PrivateKey = base64PrivateKey

	return my
}

func (my *PEMBase64) GetBase64PublicKey() string { return my.base64PublicKey }

func (my *PEMBase64) GetBase64PrivateKey() string { return my.base64PrivateKey }

func (my *PEMBase64) GetPemPublicKey() []byte { return my.publicKey }

func (my *PEMBase64) GeneratePemPublicKey() (*PEMBase64, error) {
	// 解码Base64字符串
	publicKeyBytes, err := base64.StdEncoding.DecodeString(my.base64PublicKey)
	if err != nil {
		return my, fmt.Errorf("解码Base64失败: %v", err)
	}

	// 尝试解析为PEM块
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		// 如果不是PEM格式，则尝试解析为x509公钥并创建一个PEM块
		_, err = x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			return my, fmt.Errorf("解析公钥失败se64失败: %v", err)
		}

		// 创建PEM块
		block = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}
	}

	// 将PEM块编码为内存中的字节切片
	my.publicKey = pem.EncodeToMemory(block)

	return my, nil
}

// GetPemPrivateKey 获取pem私钥
func (my *PEMBase64) GetPemPrivateKey() []byte { return my.privateKey }

// GeneratePemPrivateKey 生成pem密钥
func (my *PEMBase64) GeneratePemPrivateKey() (*PEMBase64, error) {
	// 解码Base64字符串
	privateKeyBytes, err := base64.StdEncoding.DecodeString(my.base64PrivateKey)
	if err != nil {
		return my, fmt.Errorf("解码Base64失败: %v", err)
	}

	// 手动添加PEM头部和尾部
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	my.privateKey = pem.EncodeToMemory(pemBlock)

	// 尝试解析为PEM块
	block, _ := pem.Decode(my.privateKey)
	if block == nil {
		return my, errors.New("不是有效的PEM编码私钥")
	}

	return my, nil
}
