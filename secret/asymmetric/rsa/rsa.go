package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	cryptorsa "crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/aid297/aid/secret"
)

var _ secret.Asymmetricer = (*RSAImpl)(nil)

type RSAImpl struct {
	publicKey  *cryptorsa.PublicKey
	privateKey *cryptorsa.PrivateKey
}

// NewByPEM 使用 PEM 公私钥创建 RSA 实例（可只传公钥或只传私钥）
func NewByPEM(publicKeyPEM, privateKeyPEM []byte) (secret.Asymmetricer, error) {
	my := &RSAImpl{}
	if len(publicKeyPEM) > 0 {
		pub, err := parsePublicKeyPEM(publicKeyPEM)
		if err != nil {
			return nil, err
		}
		my.publicKey = pub
	}
	if len(privateKeyPEM) > 0 {
		pri, err := parsePrivateKeyPEM(privateKeyPEM)
		if err != nil {
			return nil, err
		}
		my.privateKey = pri
		if my.publicKey == nil {
			my.publicKey = &pri.PublicKey
		}
	}
	if my.publicKey == nil && my.privateKey == nil {
		return nil, errors.New("rsa: at least one key is required")
	}
	return my, nil
}

// NewByBase64 使用 base64(DER 或 PEM文本) 公私钥创建 RSA 实例
func NewByBase64(base64PublicKey, base64PrivateKey string) (secret.Asymmetricer, error) {
	var (
		pubPEM []byte
		priPEM []byte
		err    error
	)
	if base64PublicKey != "" {
		pubPEM, err = decodeBase64KeyToPEM(base64PublicKey, true)
		if err != nil {
			return nil, err
		}
	}
	if base64PrivateKey != "" {
		priPEM, err = decodeBase64KeyToPEM(base64PrivateKey, false)
		if err != nil {
			return nil, err
		}
	}
	return NewByPEM(pubPEM, priPEM)
}

// GenerateKeyPairBase64 生成 RSA 密钥对（PKIX公钥 + PKCS8私钥，base64 DER）
func GenerateKeyPairBase64(bits int) (pubBase64, priBase64 string, err error) {
	if bits < 1024 {
		return "", "", errors.New("rsa: bits too small")
	}
	pri, err := cryptorsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	pubDer, err := x509.MarshalPKIXPublicKey(&pri.PublicKey)
	if err != nil {
		return "", "", err
	}
	priDer, err := x509.MarshalPKCS8PrivateKey(pri)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pubDer), base64.StdEncoding.EncodeToString(priDer), nil
}

// Encrypt 公钥加密，返回 base64 密文（自动分段）
func (my *RSAImpl) Encrypt(plainText []byte) (string, error) {
	if my.publicKey == nil {
		return "", errors.New("rsa: public key is required")
	}
	maxChunk := my.publicKey.Size() - 11
	if maxChunk <= 0 {
		return "", errors.New("rsa: invalid public key size")
	}
	var cipherBuf bytes.Buffer
	for start := 0; start < len(plainText); start += maxChunk {
		end := start + maxChunk
		if end > len(plainText) {
			end = len(plainText)
		}
		chunk, err := cryptorsa.EncryptPKCS1v15(rand.Reader, my.publicKey, plainText[start:end])
		if err != nil {
			return "", err
		}
		cipherBuf.Write(chunk)
	}
	return base64.StdEncoding.EncodeToString(cipherBuf.Bytes()), nil
}

// Decrypt 私钥解密，输入 base64 密文（自动分段）
func (my *RSAImpl) Decrypt(cipherBase64 string) ([]byte, error) {
	if my.privateKey == nil {
		return nil, errors.New("rsa: private key is required")
	}
	cipherText, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, err
	}
	chunkSize := my.privateKey.PublicKey.Size()
	if chunkSize <= 0 || len(cipherText)%chunkSize != 0 {
		return nil, errors.New("rsa: invalid ciphertext length")
	}
	var plainBuf bytes.Buffer
	for start := 0; start < len(cipherText); start += chunkSize {
		chunk, err := cryptorsa.DecryptPKCS1v15(rand.Reader, my.privateKey, cipherText[start:start+chunkSize])
		if err != nil {
			return nil, err
		}
		plainBuf.Write(chunk)
	}
	return plainBuf.Bytes(), nil
}

// Sign 私钥签名（SHA-256 + PKCS1v15），返回 base64 签名
func (my *RSAImpl) Sign(data []byte) (string, error) {
	if my.privateKey == nil {
		return "", errors.New("rsa: private key is required")
	}
	h := sha256.Sum256(data)
	sig, err := cryptorsa.SignPKCS1v15(rand.Reader, my.privateKey, crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// Verify 公钥验签（SHA-256 + PKCS1v15），输入 base64 签名
func (my *RSAImpl) Verify(data []byte, sigBase64 string) (bool, error) {
	if my.publicKey == nil {
		return false, errors.New("rsa: public key is required")
	}
	sig, err := base64.StdEncoding.DecodeString(sigBase64)
	if err != nil {
		return false, err
	}
	h := sha256.Sum256(data)
	if err = cryptorsa.VerifyPKCS1v15(my.publicKey, crypto.SHA256, h[:], sig); err != nil {
		return false, nil
	}
	return true, nil
}

func decodeBase64KeyToPEM(base64Key string, isPublic bool) ([]byte, error) {
	derOrPem, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	if block, _ := pem.Decode(derOrPem); block != nil {
		return derOrPem, nil
	}
	if isPublic {
		if _, err = x509.ParsePKIXPublicKey(derOrPem); err != nil {
			return nil, fmt.Errorf("rsa: invalid public key der: %w", err)
		}
		return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: derOrPem}), nil
	}
	if _, err = x509.ParsePKCS8PrivateKey(derOrPem); err == nil {
		return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: derOrPem}), nil
	}
	if _, err = x509.ParsePKCS1PrivateKey(derOrPem); err == nil {
		return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: derOrPem}), nil
	}
	return nil, errors.New("rsa: invalid private key der")
}

func parsePublicKeyPEM(publicKeyPEM []byte) (*cryptorsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, errors.New("rsa: invalid public key pem")
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pubAny.(*cryptorsa.PublicKey)
	if !ok {
		return nil, errors.New("rsa: not an rsa public key")
	}
	return pub, nil
}

func parsePrivateKeyPEM(privateKeyPEM []byte) (*cryptorsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("rsa: invalid private key pem")
	}
	if pri, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return pri, nil
	}
	priAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pri, ok := priAny.(*cryptorsa.PrivateKey)
	if !ok {
		return nil, errors.New("rsa: not an rsa private key")
	}
	return pri, nil
}
