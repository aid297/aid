package secret

import (
	"crypto"
	"io"
)

type (
	Asymmetricer interface {
		Encrypt(plainText []byte) (string, error)        // 非对称加密
		Decrypt(cipherBase64 string) ([]byte, error)     // 非对称解密
		Sign(data []byte) (string, error)                // 非对称签名
		Verify(data []byte, sigHex string) (bool, error) // 非对称验签
	}

	SemenerPriKey interface{ Public() crypto.PublicKey } // 非对称加密种子私钥
	SemenerPubKey crypto.PublicKey                       // 非对称加密种子公钥

	SemenerAttr func(sm2Sem Semener) error // 非对称加密种子属性

	Semener interface {
		SetAttrs(attrs ...SemenerAttr) error             // 设置属性
		GeneratePriKey() (err error)                     // 生成私钥
		SetPubKey(pubKey SemenerPubKey) (err error)      // 设置公钥：crypto.PublicKey(secret.SemenerPubKey)
		SetPubKeyBytes(pubKeyBytes []byte) (err error)   // 设置公钥：bytes
		SetPubKeyBase64(pubKeyBase64 string) (err error) // 设置公钥：base64
		GetPubKey() SemenerPubKey                        // 获取公钥：crypto.PublicKey(secret.SemenerPubKey) 如果公钥存在则返回公钥，如果公钥不存在则使用私钥返回公钥
		GetPubKeyBytes() ([]byte, error)                 // 获取公钥：bytes
		GetPubKeyBase64() (string, error)                // 获取公钥：base64
		SetPriKey(priKey SemenerPriKey) (err error)      // 获取私钥：crypto.PrivateKey(secret.SemenerPriKey)
		SetPriKeyBytes(priKeyBytes []byte) (err error)   // 设置私钥：bytes
		SetPriKeyBase64(priKeyBase64 string) (err error) // 设置私钥：base64
		GetPriKey() SemenerPriKey                        // 获取私钥：crypto.PrivateKey(secret.SemenerPriKey)
		GetPriKeyBytes() ([]byte, error)                 // 获取私钥：bytes
		GetPriKeyBase64() (string, error)                // 获取私钥：base64
	}

	SymmetricAttr func(symm Symmetricer) (err error) // 对称加密属性

	Symmetricer interface {
		SetAttrs(attrs ...SymmetricAttr) (err error)                              // 设置属性
		GetKeyString() string                                                     // 获取 key：string
		GetKeyBytes() []byte                                                      // 获取 key：bytes
		GetIVString() string                                                      // 获取 iv：string
		GetIVBytes() []byte                                                       // 获取 iv：bytes
		SetKeyString(key string)                                                  // 设置 key：string
		SetKeyBytes(key []byte)                                                   // 设置 key：bytes
		SetIVString(iv string)                                                    // 设置 iv：string
		SetIVBytes(iv []byte)                                                     // 设置 iv：bytes
		EncryptECBBase64(plainText []byte) (string, error)                        // ECB 模式加密，返回 base64 字符串
		DecryptECBBase64(cipherBase64 string) ([]byte, error)                     // ECB 模式解密，输入 base64 字符串
		EncryptECB(plainText []byte) ([]byte, error)                              // ECB 模式加密，返回原始字节
		DecryptECB(cipherText []byte) ([]byte, error)                             // ECB 模式解密
		EncryptCBCBase64(plainText []byte) (string, error)                        // CBC 模式加密，返回 base64 字符串
		DecryptCBCBase64(cipherBase64 string) ([]byte, error)                     // CBC 模式解密，输入 base64 字符串
		EncryptCBC(plainText []byte) ([]byte, error)                              // CBC 模式加密，返回原始字节
		DecryptCBC(cipherText []byte) ([]byte, error)                             // CBC 模式解密
		EncryptCBCStream(in io.Reader, out io.Writer) error                       // CBC 流式加密（适用于大文件）
		DecryptCBCStream(in io.Reader, out io.Writer) error                       // CBC 流式解密（适用于大文件）
		EncryptCBCFile(plainFile, outFile string, asymm Asymmetricer) error       // CBC 加密文件
		DecryptCBCFile(cipherFile, outFile string, asymm Asymmetricer) error      // CBC 解密文件
		EncryptCBCLargeFile(plainFile, outFile string, asymm Asymmetricer) error  // CBC 加密大文件
		DecryptCBCLargeFile(cipherFile, outFile string, asymm Asymmetricer) error // CBC 解密大文件
	}
)
