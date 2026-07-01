// Package rsaoaep 提供基于 RSA-OAEP（SHA-256 + MGF1-SHA-256）的非对称加解密，并实现 secrets.Asymmetric。
//
// 密钥种子与 secrets/asymmetric/rsa 相同，可直接使用本包的 NewSem，或传入 rsa.NewSem 得到的 Semen。
//
// Sign / Verify 使用 SHA-256 + PKCS#1 v1.5（与 rsa 包一致）：OAEP 仅用于加密填充，不用于签名。
package rsaOAEP

import (
	"bytes"
	"crypto"
	"crypto/rand"
	cryptorsa "crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hash"

	"github.com/aid297/aid/v2/secrets"
	myrsa "github.com/aid297/aid/v2/secrets/asymmetric/rsa"
)

var _ secrets.Asymmetric = (*RSAOAEP)(nil)

// RSAOAEP 使用 RSA-OAEP-SHA256 做加解密；签名为 RSA-PKCS1v15 + SHA-256。
type RSAOAEP struct{ sem secrets.Semen }

// NewSem 创建 RSA 密钥种子（与 asymmetric/rsa 行为一致）。
func NewSem(attrs ...secrets.SemenAttr) (secrets.Semen, error) {
	return myrsa.NewSem(attrs...)
}

// New 基于已有 Semen 构造 OAEP 加密的 Asymmetric 实例。
func New(sem secrets.Semen) secrets.Asymmetric { return &RSAOAEP{sem: sem} }

func oaepHash() hash.Hash { return sha256.New() }

// maxOAEPPlainChunk 单段明文最大长度（字节），见 PKCS#1 RSAES-OAEP。
func maxOAEPPlainChunk(pub *cryptorsa.PublicKey) int {
	return pub.Size() - 2*sha256.Size - 2
}

// Encrypt 公钥加密（RSA-OAEP-SHA256，label 为空），返回 base64 密文，自动分段。
func (my *RSAOAEP) Encrypt(plainText []byte) (string, error) {
	pubKey, ok := my.sem.GetPubKey().(*cryptorsa.PublicKey)
	if !ok || pubKey == nil {
		return "", errors.New("rsaoaep: 公钥缺失或不是 RSA")
	}
	maxChunk := maxOAEPPlainChunk(pubKey)
	if maxChunk <= 0 {
		return "", errors.New("rsaoaep: 公钥长度不足 OAEP-SHA256")
	}
	var (
		cipherBuf  bytes.Buffer
		start, end int
	)
	for start = 0; start < len(plainText); start += maxChunk {
		end = min(start+maxChunk, len(plainText))
		chunk, err := cryptorsa.EncryptOAEP(oaepHash(), rand.Reader, pubKey, plainText[start:end], nil)
		if err != nil {
			return "", err
		}
		cipherBuf.Write(chunk)
	}
	return base64.StdEncoding.EncodeToString(cipherBuf.Bytes()), nil
}

// Decrypt 私钥解密（RSA-OAEP-SHA256，label 为空），输入 base64 密文，自动分段。
func (my *RSAOAEP) Decrypt(cipherBase64 string) ([]byte, error) {
	priKey, ok := my.sem.GetPriKey().(*cryptorsa.PrivateKey)
	if !ok || priKey == nil {
		return nil, errors.New("rsaoaep: 私钥缺失或不是 RSA")
	}
	cipherText, err := base64.StdEncoding.DecodeString(secrets.PaddingBase64(cipherBase64))
	if err != nil {
		return nil, err
	}
	chunkSize := priKey.PublicKey.Size()
	if chunkSize <= 0 || len(cipherText)%chunkSize != 0 {
		return nil, errors.New("rsaoaep: 密文长度错误")
	}
	var plainBuf bytes.Buffer
	for start := 0; start < len(cipherText); start += chunkSize {
		chunk, err := cryptorsa.DecryptOAEP(oaepHash(), rand.Reader, priKey, cipherText[start:start+chunkSize], nil)
		if err != nil {
			return nil, err
		}
		plainBuf.Write(chunk)
	}
	return plainBuf.Bytes(), nil
}

// Sign 私钥签名（SHA-256 + PKCS1v15），返回 base64 签名。
func (my *RSAOAEP) Sign(data []byte) (string, error) {
	priKey, ok := my.sem.GetPriKey().(*cryptorsa.PrivateKey)
	if !ok || priKey == nil {
		return "", errors.New("rsaoaep: 私钥缺失或不是 RSA")
	}
	h := sha256.Sum256(data)
	sig, err := cryptorsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// Verify 公钥验签（SHA-256 + PKCS1v15），输入 base64 签名。
func (my *RSAOAEP) Verify(data []byte, sigBase64 string) (bool, error) {
	pubKey, ok := my.sem.GetPubKey().(*cryptorsa.PublicKey)
	if !ok || pubKey == nil {
		return false, errors.New("rsaoaep: 公钥缺失或不是 RSA")
	}
	h := sha256.Sum256(data)
	sig, err := base64.StdEncoding.DecodeString(sigBase64)
	if err != nil {
		return false, err
	}
	if err = cryptorsa.VerifyPKCS1v15(pubKey, crypto.SHA256, h[:], sig); err != nil {
		return false, nil
	}
	return true, nil
}
