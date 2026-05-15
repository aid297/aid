package secret

import (
	"crypto"
	"io"
)

type (
	Asymmetric interface {
		Encrypt(plainText []byte) (string, error)        // 非对称加密
		Decrypt(cipherBase64 string) ([]byte, error)     // 非对称解密
		Sign(data []byte) (string, error)                // 非对称签名
		Verify(data []byte, sigHex string) (bool, error) // 非对称验签
	}

	SemenPriKey interface{ Public() crypto.PublicKey } // 非对称加密种子私钥
	SemenPubKey crypto.PublicKey                       // 非对称加密种子公钥

	SemenAttr func(sem Semen) error // 非对称加密种子属性

	Semen interface {
		SetAttrs(attrs ...SemenAttr) error               // 设置属性
		GeneratePriKey() (err error)                     // 生成私钥
		SetPubKey(pubKey SemenPubKey) (err error)        // 设置公钥：crypto.PublicKey(secret.SemenPubKey)
		SetPubKeyBytes(pubKeyBytes []byte) (err error)   // 设置公钥：bytes
		SetPubKeyBase64(pubKeyBase64 string) (err error) // 设置公钥：base64
		GetPubKey() SemenPubKey                          // 获取公钥：crypto.PublicKey(secret.SemenPubKey) 如果公钥存在则返回公钥，如果公钥不存在则使用私钥返回公钥
		GetPubKeyBytes() ([]byte, error)                 // 获取公钥：bytes
		GetPubKeyBase64() (string, error)                // 获取公钥：base64
		SetPriKey(priKey SemenPriKey) (err error)        // 获取私钥：crypto.PrivateKey(secret.SemenPriKey)
		SetPriKeyBytes(priKeyBytes []byte) (err error)   // 设置私钥：bytes
		SetPriKeyBase64(priKeyBase64 string) (err error) // 设置私钥：base64
		GetPriKey() SemenPriKey                          // 获取私钥：crypto.PrivateKey(secret.SemenPriKey)
		GetPriKeyBytes() ([]byte, error)                 // 获取私钥：bytes
		GetPriKeyBase64() (string, error)                // 获取私钥：base64
	}

	SymmetricAttr func(s Symmetric) (err error) // 对称加密属性

	Symmetric interface {
		SetAttrs(attrs ...SymmetricAttr) (err error)                     // 设置属性
		GetKeyString() string                                            // 获取 key：string
		GetKeyBytes() []byte                                             // 获取 key：bytes
		GetIVString() string                                             // 获取 iv：string
		GetIVBytes() []byte                                              // 获取 iv：bytes
		SetKeyString(key string)                                         // 设置 key：string
		SetKeyBytes(key []byte)                                          // 设置 key：bytes
		SetIVString(iv string)                                           // 设置 iv：string
		SetIVBytes(iv []byte)                                            // 设置 iv：bytes
		SetAlgorithm(algorithm string) (err error)                       // 设置算法模式：ECB/CBC/CTR/GCM
		Encrypt(plainText []byte) ([]byte, error)                        // 加密：通过原始内容
		Decrypt(cipherText []byte) ([]byte, error)                       // 解密：通过密文
		EncryptBase64(plainText []byte) (string, error)                  // 加密：通过原始内容，返回 base64 编码的密文
		DecryptBase64(cipherBase64 string) ([]byte, error)               // 解密：通过 base64 编码的密文
		EncryptStream(in io.Reader, out io.Writer) error                 // 流式加密（适用于大文件，根据 Algorithm 选择 ECB/CBC）
		DecryptStream(in io.Reader, out io.Writer) error                 // 流式解密（适用于大文件，根据 Algorithm 选择 ECB/CBC）
		EncryptFile(plainFile, outFile string, a Asymmetric) error       // 加密文件（根据 Algorithm 选择 ECB/CBC）
		DecryptFile(cipherFile, outFile string, a Asymmetric) error      // 解密文件（根据 Algorithm 选择 ECB/CBC）
		EncryptLargeFile(plainFile, outFile string, a Asymmetric) error  // 加密大文件（根据 Algorithm 选择 ECB/CBC）
		DecryptLargeFile(cipherFile, outFile string, a Asymmetric) error // 解密大文件（根据 Algorithm 选择 ECB/CBC）
	}
)
