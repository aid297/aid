package asymmetric

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/aid297/aid/compression"
	"github.com/aid297/aid/str"
)

type (
	Rsa       struct{}
	UnEncrypt struct {
		Username string `json:"username"`
		Password string `json:"password"`
		AesKey   string `json:"aes_key"`
	}
)

var RsaApp Rsa

func (*Rsa) New() *Rsa { return &Rsa{} }

// NewRsa 实例化：Rsa加密
//
//go:fix 推荐使用：New方法
func NewRsa() *Rsa { return &Rsa{} }

// EncryptByBase64 通过base64公钥加密
func (my *Rsa) EncryptByBase64(base64PublicKey string, plainText []byte) ([]byte, error) {
	var (
		pemBase64 *PemBase64
		err       error
	)

	pemBase64, err = NewPemBase64().
		SetBase64PublicKey(base64PublicKey).
		GeneratePemPublicKey()
	if err != nil {
		return nil, err
	}

	return my.EncryptByPem(pemBase64.GetPemPublicKey(), plainText)
}

// EncryptByPem 通过pem公钥加密
func (my *Rsa) EncryptByPem(pemPublicKey []byte, plainText []byte) ([]byte, error) {
	var (
		err                error
		block              *pem.Block
		publicKeyInterface any
		publicKey          *rsa.PublicKey
		cipherText         []byte
	)

	block, _ = pem.Decode(pemPublicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("公钥类型错误")
	}

	// x509解码
	publicKeyInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	publicKey = publicKeyInterface.(*rsa.PublicKey)

	if len(plainText) > publicKey.N.BitLen()/8-11 {
		// 密文长度超过密钥长度，需要分段加密
		cipherText, err = my.encryptWithTooLong(publicKey, plainText)
	} else {
		// 对明文进行加密
		cipherText, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	}
	if err != nil {
		return nil, err
	}

	// 返回密文
	return cipherText, nil
}

// encryptWithTooLong 分段加密处理
func (my *Rsa) encryptWithTooLong(publicKey *rsa.PublicKey, plainText []byte) ([]byte, error) {
	var (
		maxChunkSize = publicKey.N.BitLen()/8 - 11 // 计算每个分段的最大长度
		cipherTexts  [][]byte                      // 存储每个分段加密后的结果
	)

	// 分割明文
	for i := 0; i < len(plainText); i += maxChunkSize {
		end := i + maxChunkSize
		if end > len(plainText) {
			end = len(plainText)
		}

		chunk := plainText[i:end]

		// 加密当前分段
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, chunk)
		if err != nil {
			return nil, err
		}
		cipherTexts = append(cipherTexts, encryptedChunk)
	}

	// 合并所有加密后的分段
	finalCipherText := bytes.Join(cipherTexts, []byte{})

	return finalCipherText, nil
}

// DecryptByBase64 通过base64私钥解密
func (my *Rsa) DecryptByBase64(base64PrivateKey string, cipherText []byte) ([]byte, error) {
	var (
		pemBase64 *PemBase64
		err       error
	)

	pemBase64, err = NewPemBase64().SetBase64PrivateKye(base64PrivateKey).GeneratePemPrivateKey()
	if err != nil {
		return nil, err
	}

	return my.DecryptByPem(pemBase64.GetPemPrivateKey(), cipherText)
}

// DecryptByPem 使用PEM私钥进行RSA解密
func (my *Rsa) DecryptByPem(pemPrivateKey []byte, cipherText []byte) ([]byte, error) {
	var (
		err                       error
		block                     *pem.Block
		privateKey, rsaPrivateKey *rsa.PrivateKey
		privateKeyInterface       any
		ok                        bool
		plainText                 []byte
	)

	block, _ = pem.Decode(pemPrivateKey)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing the private key")
	}

	// 解析DER编码的私钥
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 尝试解析PKCS8格式的私钥
		privateKeyInterface, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		// 使用类型断言将interface{}转换为*rsa.PrivateKey
		rsaPrivateKey, ok = privateKeyInterface.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("failed to cast private key to *rsa.PrivateKey")
		}
		privateKey = rsaPrivateKey
	}

	if len(plainText) > privateKey.PublicKey.N.BitLen() {
		// 分段解密
		_, err2 := my.decryptWithTooLong(privateKey, cipherText)
		if err2 != nil {
			return nil, fmt.Errorf("分段解密错误：%v", err2)
		}
	} else {
		// 解密数据
		plainText, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt: %w", err)
		}
	}

	return plainText, nil
}

// decryptWithTooLong 分段解密
func (my *Rsa) decryptWithTooLong(privateKey *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	var (
		maxChunkSize int
		plainTexts   [][]byte
	)

	// 计算每个分段的最大长度
	maxChunkSize = privateKey.PublicKey.N.BitLen() / 8

	// 分割密文
	for i := 0; i < len(cipherText); i += maxChunkSize {
		end := i + maxChunkSize
		if end > len(cipherText) {
			end = len(cipherText)
		}

		chunk := cipherText[i:end]

		// 解密当前分段
		decryptedChunk, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt segment: %w", err)
		}
		plainTexts = append(plainTexts, decryptedChunk)
	}

	// 合并所有解密后的分段
	finalPlainText := bytes.Join(plainTexts, []byte{})

	return finalPlainText, nil
}

func (*Rsa) DemoEncryptRsa(unEncrypt []byte) string {
	var (
		base64PublicKey                     = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCFbbjNGuqhF3HhmvnZxjG6mS6Q3OmD/vh9voriZTyNCVLJ7y2r0bHZZ7brWwkgtGPQXosZ0IzUZAvlMuZ0m11DiuXZzlCnRz1owwMXKalJeeKQwA8CoJBSy99zCo9fxIErqTMhGwPFCKUaByt8TEIkNq8fUsmqjqqshRLKSazWuwIDAQAB"
		encrypted                           []byte
		base64Encrypted                     string
		pemPublicKey                        []byte
		pemBase64                           *PemBase64
		generatePemPublicKeyErr, encryptErr error
	)

	pemBase64, generatePemPublicKeyErr = PemBase64App.New().SetBase64PublicKey(base64PublicKey).GeneratePemPublicKey()
	if generatePemPublicKeyErr != nil {
		str.TerminalLogApp.New("[RSA] generate public key: %v").Error(generatePemPublicKeyErr)
	}

	pemPublicKey = pemBase64.GetPemPublicKey()
	str.TerminalLogApp.New("[RSA] generate public key: \n%s").Info(pemPublicKey)

	encrypted, encryptErr = RsaApp.New().EncryptByPem(pemPublicKey, unEncrypt)
	if encryptErr != nil {
		str.TerminalLogApp.New("[RSA] encrypt: %v").Error(encryptErr)
	}
	base64Encrypted = base64.StdEncoding.EncodeToString(encrypted)

	return base64Encrypted
}

func (*Rsa) DemoDecryptRsa(base64Encrypted string) string {
	var (
		base64PrivateKey                                     = "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAIVtuM0a6qEXceGa+dnGMbqZLpDc6YP++H2+iuJlPI0JUsnvLavRsdlntutbCSC0Y9BeixnQjNRkC+Uy5nSbXUOK5dnOUKdHPWjDAxcpqUl54pDADwKgkFLL33MKj1/EgSupMyEbA8UIpRoHK3xMQiQ2rx9SyaqOqqyFEspJrNa7AgMBAAECgYATaA4E5vFRVNOfeKb2YblB5p27PCZKqH8D6v7QRuEzsjN0Y3FFGE7BzC/ys170fsg1ukqJCqgxDAwe3fRe6Wn6/Y5IEF/wRYODQn6yAXhCUepheaRl9zK+P+XXbGWENdL2N/KchNZrKUF97Eu00OhBI7uEKpUrhPuzaYDPiHujQQJBAOvc+Xwz3j/srv26bk5UJOAJtU096pNseEeVzFqSTU903NdgFUQupTsPeokUtMBeMihAYlfDZypIK0kvBoymTNkCQQCQ0e/vEGnqh9C0y340HUlIZe0Q5mAJ5e+3a7lR21LS9ki5vQLUf2Wjxw/QVbPDZthGK33BusrobyuwcVOMmROzAkEAz9lefeZTb6/Kkcvtktcx28CSZawvgJTw9dx7RkFxIZkRWDbS5s/YSdCdIhn+IxufRbtfLooC6s7IXmizc9TFGQJAZP1hum7RzbFkg4+ctK7vmcMqbKyasIxefKRsmX6+5UrGMHB0dsdYk7uPdZMuRseDbnuJuP2P3kMYTnTY9KUTLQJANq7Cy5OjtHiJ5EsRBePfGm9Qvs3mwJZAKDpZsmTRSyaQCTCpL6RQ+7gVFIEmiEU4REjag9/aq8C1G0MyvwxkiA=="
		pemPrivateKey, encrypted                             []byte
		base64DecodeErr, generatePemPublicKeyErr, decryptErr error
		pemBase64                                            *PemBase64
		decrypted                                            []byte
	)

	pemBase64, generatePemPublicKeyErr = PemBase64App.New().SetBase64PrivateKye(base64PrivateKey).GeneratePemPrivateKey()
	if generatePemPublicKeyErr != nil {
		str.TerminalLogApp.New("[RSA] generate private key: %v").Info(generatePemPublicKeyErr)
	}

	pemPrivateKey = pemBase64.GetPemPrivateKey()
	str.TerminalLogApp.New("[RSA] generate private key: %s").Info(pemPrivateKey)

	encrypted, base64DecodeErr = base64.StdEncoding.DecodeString(base64Encrypted)
	if base64DecodeErr != nil {
		str.TerminalLogApp.New("[RSA] base64 decode: %v").Error(base64DecodeErr)
	}

	decrypted, decryptErr = RsaApp.New().DecryptByPem(pemPrivateKey, encrypted)
	if decryptErr != nil {
		str.TerminalLogApp.New("[RSA] decrypt: %v").Error(decryptErr)
	}

	return string(decrypted)
}

func (my *Rsa) Demo() {
	var (
		unEncrypt = UnEncrypt{
			Username: "cbit",
			Password: "cbit-pwd",
			// AesKey:   "tjp5OPIU1ETF5s33fsLWdA==",
			AesKey: "87dwQRkoNFNoIcq1A+zFHA==",
		}
		jsonByte, zipByte, unzipByte  []byte
		jsonErr, zipErr, unzipByteErr error
	)

	// json序列化
	jsonByte, jsonErr = json.Marshal(unEncrypt)
	if jsonErr != nil {
		str.TerminalLogApp.New("[RSA] json marshal failed: %v").Error(zipErr)
	}

	// zip压缩
	zipByte, zipErr = compression.NewZlib().Compress(jsonByte)
	if zipErr != nil {
		str.TerminalLogApp.New("[RSA] zip failed: %v").Error(zipErr)
	}

	base64Encrypted := my.DemoEncryptRsa(zipByte) // 加密
	str.TerminalLogApp.New("[RSA] encrypting: %s").Success(base64Encrypted)

	decrypted := my.DemoDecryptRsa(base64Encrypted) // 解密
	str.TerminalLogApp.New("[RSA] decrypted").Info()

	// 解密后解压缩
	unzipByte, unzipByteErr = compression.NewZlib().Decompress([]byte(decrypted))
	if unzipByteErr != nil {
		str.TerminalLogApp.New("[RSA] unzipped failed: %v").Error(unzipByteErr)
	}
	str.TerminalLogApp.New("[RSA] unzipped").Info()

	str.TerminalLogApp.New("[RSA] decrypted: %s").Success(string(unzipByte))
}
