package ecdsa

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"

	"github.com/aid297/aid/v2/secret"
)

var _ secret.Asymmetric = (*ECDSA)(nil)

type ECDSA struct{ sem secret.Semen }

// New 实例化
func New(sem secret.Semen) secret.Asymmetric { return &ECDSA{sem: sem} }

// Encrypt 公钥加密（ECDSA 不支持加密）
func (my *ECDSA) Encrypt(plainText []byte) (string, error) {
	return "", errors.New("ECDSA 不支持加密操作，请使用 ECDH 或混合加密")
}

// Decrypt 私钥解密
func (my *ECDSA) Decrypt(cipherBase64 string) ([]byte, error) {
	return nil, errors.New("ECDSA 不支持解密操作")
}

// Sign 私钥签名（SHA-256 + ECDSA），返回 base64 签名
func (my *ECDSA) Sign(data []byte) (string, error) {
	priKey := my.sem.GetPriKey().(*ecdsa.PrivateKey)

	// 对数据做 SHA-256 哈希
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, priKey, hash[:])
	if err != nil {
		return "", err
	}

	// 编码为 r || s（DER ASN.1 格式的替代方案）
	size := (priKey.Curve.Params().BitSize + 7) / 8
	sig := make([]byte, 2*size)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(sig[size-len(rBytes):size], rBytes)
	copy(sig[2*size-len(sBytes):2*size], sBytes)

	return base64.RawURLEncoding.EncodeToString(sig), nil
}

// Verify 公钥验签（SHA-256 + ECDSA），输入 base64 签名
func (my *ECDSA) Verify(data []byte, sigBase64 string) (bool, error) {
	pubKey := my.sem.GetPubKey().(*ecdsa.PublicKey)

	sig, err := base64.RawURLEncoding.DecodeString(sigBase64)
	if err != nil {
		return false, err
	}

	size := (pubKey.Curve.Params().BitSize + 7) / 8
	if len(sig) != 2*size {
		return false, errors.New("签名长度错误")
	}

	r := new(big.Int).SetBytes(sig[:size])
	s := new(big.Int).SetBytes(sig[size:])

	// 对数据做 SHA-256 哈希
	hash := sha256.Sum256(data)
	return ecdsa.Verify(pubKey, hash[:], r, s), nil
}

// ==================== ECDSA ASN.1 DER 格式支持（可选） ====================

// SignASN1 私钥签名，返回 DER ASN.1 编码的签名
func (my *ECDSA) SignASN1(data []byte) (string, error) {
	priKey := my.sem.GetPriKey().(*ecdsa.PrivateKey)

	r, err := ecdsa.SignASN1(rand.Reader, priKey, data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(r), nil
}

// VerifyASN1 公钥验签，输入 DER ASN.1 编码的签名
func (my *ECDSA) VerifyASN1(data []byte, sigBase64 string) (bool, error) {
	pubKey := my.sem.GetPubKey().(*ecdsa.PublicKey)

	sig, err := base64.StdEncoding.DecodeString(sigBase64)
	if err != nil {
		return false, err
	}

	return ecdsa.VerifyASN1(pubKey, data, sig), nil
}
