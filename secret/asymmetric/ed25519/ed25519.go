package ed25519

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"

	"github.com/aid297/aid/v2/secret"
)

var _ secret.Asymmetric = (*Ed25519)(nil)

type Ed25519 struct {
	sem secret.Semen
}

// New 实例化
func New(sem secret.Semen) secret.Asymmetric { return &Ed25519{sem: sem} }

// Encrypt 公钥加密（Ed25519 不直接支持加密，返回错误）
func (my *Ed25519) Encrypt(plainText []byte) (string, error) {
	return "", errors.New("Ed25519 不支持加密操作，请使用 X25519 或混合加密")
}

// Decrypt 私钥解密
func (my *Ed25519) Decrypt(cipherBase64 string) ([]byte, error) {
	return nil, errors.New("Ed25519 不支持解密操作")
}

// Sign 私钥签名（Ed25519），返回 base64 签名
func (my *Ed25519) Sign(data []byte) (string, error) {
	priKey := my.sem.GetPriKey().(ed25519.PrivateKey)
	sig := ed25519.Sign(priKey, data)
	return base64.RawURLEncoding.EncodeToString(sig), nil
}

// Verify 公钥验签（Ed25519），输入 base64 签名
func (my *Ed25519) Verify(data []byte, sigBase64 string) (bool, error) {
	pubKey := my.sem.GetPubKey().(ed25519.PublicKey)
	sig, err := base64.RawURLEncoding.DecodeString(sigBase64)
	if err != nil {
		return false, err
	}

	return ed25519.Verify(pubKey, data, sig), nil
}
