package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	cryptorsa "crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/aid297/aid/v2/secrets"
)

var _ secrets.Asymmetric = (*RSA)(nil)

type RSA struct{ sem secrets.Semen }

// New 实例化
func New(sem secrets.Semen) secrets.Asymmetric { return &RSA{sem: sem} }

// // NewByPEM 使用 PEM 公私钥创建 RSA 实例（可只传公钥或只传私钥）
// func NewByPEM(publicKeyPEM, privateKeyPEM []byte) (secrets.Asymmetric, error) {
// 	my := &RSAImpl{}
// 	if len(publicKeyPEM) > 0 {
// 		pub, err := parsePublicKeyPEM(publicKeyPEM)
// 		if err != nil {
// 			return nil, err
// 		}
// 		my.publicKey = pub
// 	}
// 	if len(privateKeyPEM) > 0 {
// 		pri, err := parsePrivateKeyPEM(privateKeyPEM)
// 		if err != nil {
// 			return nil, err
// 		}
// 		my.privateKey = pri
// 		if my.publicKey == nil {
// 			my.publicKey = &pri.PublicKey
// 		}
// 	}
// 	if my.publicKey == nil && my.privateKey == nil {
// 		return nil, errors.New("缺少公钥或私钥")
// 	}
// 	return my, nil
// }

// // NewByBase64 使用 base64(DER 或 PEM文本) 公私钥创建 RSA 实例
// func NewByBase64(base64PublicKey, base64PrivateKey string) (secrets.Asymmetric, error) {
// 	var (
// 		pubPEM []byte
// 		priPEM []byte
// 		err    error
// 	)
// 	if base64PublicKey != "" {
// 		pubPEM, err = decodeBase64KeyToPEM(base64PublicKey, true)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	if base64PrivateKey != "" {
// 		priPEM, err = decodeBase64KeyToPEM(base64PrivateKey, false)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return NewByPEM(pubPEM, priPEM)
// }

// // GenerateKeyPairBase64 生成 RSA 密钥对（PKIX公钥 + PKCS8私钥，base64 DER）
// func GenerateKeyPairBase64(bits int) (pubBase64, priBase64 string, err error) {
// 	if bits < 1024 {
// 		return "", "", errors.New("密钥长度过段")
// 	}
// 	pri, err := cryptorsa.GenerateKey(rand.Reader, bits)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	pubDer, err := x509.MarshalPKIXPublicKey(&pri.PublicKey)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	priDer, err := x509.MarshalPKCS8PrivateKey(pri)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return base64.StdEncoding.EncodeToString(pubDer), base64.StdEncoding.EncodeToString(priDer), nil
// }

// func GenerateKeyPairPEM(bits int) (pubPEM, priPEM []byte, err error) {
// 	var priKey *cryptorsa.PrivateKey
// 	if bits < 1024 {
// 		return nil, nil, errors.New("密钥长度过段")
// 	}

// 	if priKey, err = cryptorsa.GenerateKey(rand.Reader, bits); err != nil {
// 		return nil, nil, err
// 	}

// 	if pubPEM, err = x509.MarshalPKIXPublicKey(&priKey.PublicKey); err != nil {
// 		return nil, nil, err
// 	}

// 	if priPEM, err = x509.MarshalPKCS8PrivateKey(priKey); err != nil {
// 		return nil, nil, err
// 	}

// 	return pubDer, priDer, nil
// }

// Encrypt 公钥加密，返回 base64 密文（自动分段）
func (my *RSA) Encrypt(plainText []byte) (string, error) {
	var (
		err        error
		pubKey     = my.sem.GetPubKey().(*cryptorsa.PublicKey)
		maxChunk   = pubKey.Size() - 11
		cipherBuf  bytes.Buffer
		chunk      []byte
		start, end int
	)

	if maxChunk <= 0 {
		return "", errors.New("公钥长度错误")
	}

	for start = 0; start < len(plainText); start += maxChunk {
		end = min(start+maxChunk, len(plainText))
		if chunk, err = cryptorsa.EncryptPKCS1v15(rand.Reader, pubKey, plainText[start:end]); err != nil {
			return "", err
		}
		cipherBuf.Write(chunk)
	}

	return base64.StdEncoding.EncodeToString(cipherBuf.Bytes()), nil
}

// Decrypt 私钥解密，输入 base64 密文（自动分段）
func (my *RSA) Decrypt(cipherBase64 string) ([]byte, error) {
	var (
		err        error
		priKey     = my.sem.GetPriKey().(*cryptorsa.PrivateKey)
		cipherText []byte
		chunkSize  = priKey.PublicKey.Size()
		plainBuf   bytes.Buffer
		start      int
		chunk      []byte
	)

	if cipherText, err = base64.StdEncoding.DecodeString(secrets.PaddingBase64(cipherBase64)); err != nil {
		return nil, err
	}

	if chunkSize <= 0 || len(cipherText)%chunkSize != 0 {
		return nil, errors.New("密文长度错误")
	}

	for start = 0; start < len(cipherText); start += chunkSize {
		if chunk, err = cryptorsa.DecryptPKCS1v15(rand.Reader, priKey, cipherText[start:start+chunkSize]); err != nil {
			return nil, err
		}

		plainBuf.Write(chunk)
	}

	return plainBuf.Bytes(), nil
}

// Sign 私钥签名（SHA-256 + PKCS1v15），返回 base64 签名
func (my *RSA) Sign(data []byte) (string, error) {
	var (
		err    error
		h      [32]byte = sha256.Sum256(data)
		priKey          = my.sem.GetPriKey().(*cryptorsa.PrivateKey)
		sig    []byte
	)

	if sig, err = cryptorsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, h[:]); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

// Verify 公钥验签（SHA-256 + PKCS1v15），输入 base64 签名
func (my *RSA) Verify(data []byte, sigBase64 string) (bool, error) {
	var (
		err    error
		h      [32]byte = sha256.Sum256(data)
		pubKey          = my.sem.GetPubKey().(*cryptorsa.PublicKey)
		sig    []byte
	)

	if sig, err = base64.StdEncoding.DecodeString(sigBase64); err != nil {
		return false, err
	}

	if err = cryptorsa.VerifyPKCS1v15(pubKey, crypto.SHA256, h[:], sig); err != nil {
		return false, nil
	}

	return true, nil
}

// func decodeBase64KeyToPEM(base64Key string, isPublic bool) ([]byte, error) {
// 	derOrPem, err := base64.StdEncoding.DecodeString(base64Key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if block, _ := pem.Decode(derOrPem); block != nil {
// 		return derOrPem, nil
// 	}
// 	if isPublic {
// 		if _, err = x509.ParsePKIXPublicKey(derOrPem); err != nil {
// 			return nil, fmt.Errorf("rsa: invalid public key der: %w", err)
// 		}
// 		return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: derOrPem}), nil
// 	}
// 	if _, err = x509.ParsePKCS8PrivateKey(derOrPem); err == nil {
// 		return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: derOrPem}), nil
// 	}
// 	if _, err = x509.ParsePKCS1PrivateKey(derOrPem); err == nil {
// 		return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: derOrPem}), nil
// 	}
// 	return nil, errors.New("rsa: invalid private key der")
// }

// func parsePublicKeyPEM(publicKeyPEM []byte) (*cryptorsa.PublicKey, error) {
// 	block, _ := pem.Decode(publicKeyPEM)
// 	if block == nil {
// 		return nil, errors.New("rsa: invalid public key pem")
// 	}
// 	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pub, ok := pubAny.(*cryptorsa.PublicKey)
// 	if !ok {
// 		return nil, errors.New("rsa: not an rsa public key")
// 	}
// 	return pub, nil
// }

// func parsePrivateKeyPEM(privateKeyPEM []byte) (*cryptorsa.PrivateKey, error) {
// 	block, _ := pem.Decode(privateKeyPEM)
// 	if block == nil {
// 		return nil, errors.New("rsa: invalid private key pem")
// 	}
// 	if pri, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
// 		return pri, nil
// 	}
// 	priAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pri, ok := priAny.(*cryptorsa.PrivateKey)
// 	if !ok {
// 		return nil, errors.New("rsa: not an rsa private key")
// 	}
// 	return pri, nil
// }
