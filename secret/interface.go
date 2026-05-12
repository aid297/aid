package secret

import (
	"io"

	"github.com/tjfoc/gmsm/sm2"
)

type (
	Asymmetricer interface {
		Encrypt(plainText []byte) (string, error)
		Decrypt(cipherBase64 string) ([]byte, error)
		Sign(data []byte) (string, error)
		Verify(data []byte, sigHex string) (bool, error)
	}

	SemenerAttr func(sm2Sem Semener) error

	Semener interface {
		SetAttrs(attrs ...SemenerAttr) error
		GeneratePriKey() (err error)
		GetPriKey() *sm2.PrivateKey
		GetPriKeyBytes() ([]byte, error)
		GetPriKeyBase64() (string, error)
		GetPubKey() *sm2.PublicKey
		GetPubKeyBytes() ([]byte, error)
		GetPubKeyBase64() (string, error)
		SetPriKey(priKey *sm2.PrivateKey) (err error)
		SetPriKeyBytes(priKeyBytes []byte) (err error)
		SetPriKeyBase64(priKeyBase64 string) (err error)
	}

	SymmetricAttr func(symm Symmetricer) (err error)

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
