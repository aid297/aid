package sm4

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tjfoc/gmsm/sm4"

	"github.com/aid297/aid/v2/secret"
)

// Mode 加密模式
type Mode int

const blockSize = 16 // blockSize SM4 块大小固定为 16 字节

var _ secret.Symmetricer = (*SM4Impl)(nil)

type SM4Impl struct {
	key, iv   []byte
	algorithm string
}

// New 实例化
func New(attrs ...secret.SymmetricAttr) (my secret.Symmetricer, err error) {
	my = &SM4Impl{algorithm: "CBC"}
	err = my.SetAttrs(attrs...)
	return
}

// SetAttrs 设置属性
func (my *SM4Impl) SetAttrs(attrs ...secret.SymmetricAttr) (err error) {
	for idx := range attrs {
		if err = attrs[idx](my); err != nil {
			return
		}
	}
	return
}

// padPKCS7 PKCS7 填充
func padPKCS7(src []byte, size int) []byte {
	padding := size - len(src)%size
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// unPadPKCS7 移除 PKCS7 填充
func unPadPKCS7(src []byte, size int) ([]byte, error) {
	length := len(src)
	if size <= 0 {
		return nil, fmt.Errorf("错误的填充块长度: %d", size)
	}
	if length == 0 || length%size != 0 {
		return nil, errors.New("错误的填充内容长度")
	}
	unPadding := int(src[length-1])
	if unPadding == 0 || unPadding > size {
		return nil, errors.New("错误的padding大小")
	}
	for _, b := range src[length-unPadding:] {
		if int(b) != unPadding {
			return nil, errors.New("错误的PKCS7 padding大小")
		}
	}
	return src[:length-unPadding], nil
}

// validateKey 校验密钥长度（SM4 要求 16 字节）
func validateKey(key []byte) error {
	if len(key) != blockSize {
		return fmt.Errorf("错误的key长度 %d，必须16字节", len(key))
	}
	return nil
}

// validateIV 校验 IV 长度
func validateIV(iv []byte) error {
	if len(iv) != blockSize {
		return fmt.Errorf("错误的iv长度 %d，必须16字节", len(iv))
	}
	return nil
}

// GetKeyString 获取 key：string
func (my *SM4Impl) GetKeyString() string { return string(my.key) }

// GetKeyBytes 获取 key：bytes
func (my *SM4Impl) GetKeyBytes() []byte { return my.key }

// GetIVString 获取 iv：string
func (my *SM4Impl) GetIVString() string { return string(my.iv) }

// GetIVBytes 获取 iv：bytes
func (my *SM4Impl) GetIVBytes() []byte { return my.iv }

// SetKeyString 设置 key：string
func (my *SM4Impl) SetKeyString(key string) { my.key = []byte(key) }

// SetKeyBytes 设置 key：bytes
func (my *SM4Impl) SetKeyBytes(key []byte) { my.key = key }

// SetIVString 设置 iv：string
func (my *SM4Impl) SetIVString(iv string) { my.iv = []byte(iv) }

// SetIVBytes 设置 iv：bytes
func (my *SM4Impl) SetIVBytes(iv []byte) { my.iv = iv }

// SetAlgorithm 设置算法
func (my *SM4Impl) SetAlgorithm(algorithm string) (err error) {
	switch algorithm {
	case "ECB", "CBC":
		my.algorithm = algorithm
	default:
		err = errors.New("对称加密算法目前只支持：ECB/CBC")
	}
	return
}

// Encrypt 加密：通过原始内容
func (my *SM4Impl) Encrypt(plainText []byte) ([]byte, error) {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECB(plainText)
	case "CBC":
		return my.encryptCBC(plainText)
	default:
		return nil, errors.New("对称加密算法目前只支持：ECB/CBC")
	}
}

// Decrypt 解密：通过密文
func (my *SM4Impl) Decrypt(cipherText []byte) ([]byte, error) {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECB(cipherText)
	case "CBC":
		return my.decryptCBC(cipherText)
	default:
		return nil, errors.New("对称解密算法目前只支持：ECB/CBC")
	}
}

// EcryptBase64 加密：通过原始内容，返回 base64 编码的密文
func (my SM4Impl) EncryptBase64(plainText []byte) (string, error) {
	cipherText, err := my.Encrypt(plainText)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptBase64 解密：通过 base64 编码的密文
func (my *SM4Impl) DecryptBase64(cipherBase64 string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, fmt.Errorf("base64解码错误：%w", err)
	}
	return my.decryptECB(cipherText)
}

// encryptECB ECB 模式加密，返回原始字节
func (my *SM4Impl) encryptECB(plainText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	padded := padPKCS7(plainText, blockSize)
	cipherText := make([]byte, len(padded))
	// ECB 逐块加密
	for i := 0; i < len(padded); i += blockSize {
		block.Encrypt(cipherText[i:i+blockSize], padded[i:i+blockSize])
	}
	return cipherText, nil
}

// decryptECB ECB 模式解密
func (my *SM4Impl) decryptECB(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%blockSize != 0 {
		return nil, errors.New("错误的密文长度")
	}
	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(cipherText))
	for i := 0; i < len(cipherText); i += blockSize {
		block.Decrypt(plainText[i:i+blockSize], cipherText[i:i+blockSize])
	}
	return unPadPKCS7(plainText, blockSize)
}

// encryptCBC CBC 模式加密，返回原始字节
func (my *SM4Impl) encryptCBC(plainText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if err := validateIV(my.iv); err != nil {
		return nil, err
	}
	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	padded := padPKCS7(plainText, blockSize)
	cipherText := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, my.iv).CryptBlocks(cipherText, padded)
	return cipherText, nil
}

// decryptCBC CBC 模式解密
func (my *SM4Impl) decryptCBC(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if err := validateIV(my.iv); err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%blockSize != 0 {
		return nil, errors.New("错误的密文长度")
	}
	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, my.iv).CryptBlocks(plainText, cipherText)
	return unPadPKCS7(plainText, blockSize)
}

// encryptCBCStream CBC 流式加密（适用于大文件）
func (my *SM4Impl) encryptCBCStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}
	if err := validateIV(my.iv); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}
	encrypter := cipher.NewCBCEncrypter(block, my.iv)

	var (
		pending []byte
		readBuf = make([]byte, 1024*1024)
	)

	for {
		n, readErr := in.Read(readBuf)
		if n > 0 {
			pending = append(pending, readBuf[:n]...)
			full := len(pending) / blockSize * blockSize
			if full > 0 {
				chunk := make([]byte, full)
				encrypter.CryptBlocks(chunk, pending[:full])
				if _, err = out.Write(chunk); err != nil {
					return err
				}
				pending = pending[full:]
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	final := padPKCS7(pending, blockSize)
	finalCipher := make([]byte, len(final))
	encrypter.CryptBlocks(finalCipher, final)
	if _, err = out.Write(finalCipher); err != nil {
		return err
	}

	return nil
}

// decryptCBCStream CBC 流式解密（适用于大文件）
func (my *SM4Impl) decryptCBCStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}
	if err := validateIV(my.iv); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}
	decrypter := cipher.NewCBCDecrypter(block, my.iv)

	var (
		pending []byte
		readBuf = make([]byte, 1024*1024)
	)

	for {
		n, readErr := in.Read(readBuf)
		if n > 0 {
			pending = append(pending, readBuf[:n]...)
			fullBlocks := len(pending) / blockSize
			if fullBlocks >= 2 {
				processLen := (fullBlocks - 1) * blockSize // 留最后一块去掉 padding
				chunk := make([]byte, processLen)
				decrypter.CryptBlocks(chunk, pending[:processLen])
				if _, err = out.Write(chunk); err != nil {
					return err
				}
				pending = pending[processLen:]
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	if len(pending) == 0 || len(pending)%blockSize != 0 {
		return errors.New("错误的密文长度")
	}

	finalPlain := make([]byte, len(pending))
	decrypter.CryptBlocks(finalPlain, pending)
	finalPlain, err = unPadPKCS7(finalPlain, blockSize)
	if err != nil {
		return err
	}

	if _, err = out.Write(finalPlain); err != nil {
		return err
	}

	return nil
}

// encryptCBCFile 加密文件
func (my *SM4Impl) encryptCBCFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		plainData         []byte
		fileCipher        []byte
		sm4KeyAndIV       []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
		keyLen            int
		out               []byte
	)

	// 1. 读取明文文件
	if plainData, err = os.ReadFile(plainFile); err != nil {
		return err
	}

	// 2. 加密
	if fileCipher, err = my.encryptCBC(plainData); err != nil {
		return err
	}

	// 3. 用 非对称 公钥加密 对称 密钥和 IV（共 32 字节）
	sm4KeyAndIV = append(my.key, my.iv...)
	if encryptedKeyStr, err = asymm.Encrypt(sm4KeyAndIV); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr) // base64 字符串转字节

	// 5. 拼装输出：[2字节：密钥段长度][SM2加密的SM4密钥(base64)][SM4密文]
	keyLen = len(encryptedKeyBytes)
	out = make([]byte, 2+keyLen+len(fileCipher))
	out[0] = byte(keyLen >> 8)
	out[1] = byte(keyLen)
	copy(out[2:], encryptedKeyBytes)
	copy(out[2+keyLen:], fileCipher)

	return os.WriteFile(outFile, out, 0644)
}

// decryptCBCFile 解密文件
func (my *SM4Impl) decryptCBCFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err                error
		data               []byte
		keyLen             int
		encryptedKeyBase64 string
		fileCipher         []byte
		sm4KeyAndIV        []byte
		plainData          []byte
	)

	// 1. 读取密文文件
	if data, err = os.ReadFile(cipherFile); err != nil {
		return err
	}

	if len(data) < 2 {
		return os.ErrInvalid
	}

	// 2. 解析头部：取出 非对称 加密的 对称加密 密钥段
	keyLen = int(data[0])<<8 | int(data[1])
	if len(data) < 2+keyLen {
		return os.ErrInvalid
	}
	encryptedKeyBase64 = string(data[2 : 2+keyLen])
	fileCipher = data[2+keyLen:]

	// 3. 非对称私钥解密 对称 key 和 IV
	if sm4KeyAndIV, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(sm4KeyAndIV) != 32 {
		return os.ErrInvalid
	}
	my.key = sm4KeyAndIV[:16]
	my.iv = sm4KeyAndIV[16:]

	// 4. SM4-CBC 解密文件内容
	if plainData, err = my.decryptCBC(fileCipher); err != nil {
		return err
	}

	return os.WriteFile(outFile, plainData, 0644)
}

// encryptCBCLargeFile 用 非对称+对称 流式加密大文件（TB级）
func (my *SM4Impl) encryptCBCLargeFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		inF, outF         *os.File
		sm4KeyAndIV       []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
	)

	if inF, err = os.Open(plainFile); err != nil {
		return err
	}
	defer func() { _ = inF.Close() }()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer func() { _ = outF.Close() }()

	sm4KeyAndIV = append(my.key, my.iv...)
	if encryptedKeyStr, err = asymm.Encrypt(sm4KeyAndIV); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr)
	if len(encryptedKeyBytes) > 0xFFFF {
		return errors.New("encrypted key too long")
	}

	if _, err = outF.Write([]byte{byte(len(encryptedKeyBytes) >> 8), byte(len(encryptedKeyBytes))}); err != nil {
		return err
	}
	if _, err = outF.Write(encryptedKeyBytes); err != nil {
		return err
	}

	return my.encryptCBCStream(inF, outF)
}

// decryptCBCLargeFile 用 非对称+对称 流式解密大文件（TB级）
func (my *SM4Impl) decryptCBCLargeFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err                error
		inF, outF          *os.File
		keyLenBuf          = make([]byte, 2)
		keyLen             int
		encryptedKeyBytes  []byte
		encryptedKeyBase64 string
		sm4KeyAndIV        []byte
	)

	if inF, err = os.Open(cipherFile); err != nil {
		return err
	}
	defer func() { _ = inF.Close() }()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer func() { _ = outF.Close() }()

	if _, err = io.ReadFull(inF, keyLenBuf); err != nil {
		return err
	}
	if keyLen = int(keyLenBuf[0])<<8 | int(keyLenBuf[1]); keyLen <= 0 {
		return os.ErrInvalid
	}

	encryptedKeyBytes = make([]byte, keyLen)
	if _, err = io.ReadFull(inF, encryptedKeyBytes); err != nil {
		return err
	}
	encryptedKeyBase64 = string(encryptedKeyBytes)

	if sm4KeyAndIV, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(sm4KeyAndIV) != 32 {
		return os.ErrInvalid
	}
	my.key = sm4KeyAndIV[:16]
	my.iv = sm4KeyAndIV[16:]

	return my.decryptCBCStream(inF, outF)
}

// encryptECBStream ECB 流式加密（适用于大文件）
func (my *SM4Impl) encryptECBStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}

	var (
		pending []byte
		readBuf = make([]byte, 1024*1024)
	)

	for {
		n, readErr := in.Read(readBuf)
		if n > 0 {
			pending = append(pending, readBuf[:n]...)
			full := len(pending) / blockSize * blockSize
			if full > 0 {
				chunk := make([]byte, full)
				for i := 0; i < full; i += blockSize {
					block.Encrypt(chunk[i:i+blockSize], pending[i:i+blockSize])
				}
				if _, err = out.Write(chunk); err != nil {
					return err
				}
				pending = pending[full:]
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	final := padPKCS7(pending, blockSize)
	finalCipher := make([]byte, len(final))
	for i := 0; i < len(final); i += blockSize {
		block.Encrypt(finalCipher[i:i+blockSize], final[i:i+blockSize])
	}
	if _, err = out.Write(finalCipher); err != nil {
		return err
	}

	return nil
}

// decryptECBStream ECB 流式解密（适用于大文件）
func (my *SM4Impl) decryptECBStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}

	var (
		pending []byte
		readBuf = make([]byte, 1024*1024)
	)

	for {
		n, readErr := in.Read(readBuf)
		if n > 0 {
			pending = append(pending, readBuf[:n]...)
			full := len(pending) / blockSize * blockSize
			if full > 0 {
				chunk := make([]byte, full)
				for i := 0; i < full; i += blockSize {
					block.Decrypt(chunk[i:i+blockSize], pending[i:i+blockSize])
				}
				if _, err = out.Write(chunk); err != nil {
					return err
				}
				pending = pending[full:]
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	finalPlain, err := unPadPKCS7(pending, blockSize)
	if err != nil {
		return err
	}

	if _, err = out.Write(finalPlain); err != nil {
		return err
	}

	return nil
}

// encryptECBFile ECB 加密文件
func (my *SM4Impl) encryptECBFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		plainData         []byte
		fileCipher        []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
	)

	if plainData, err = os.ReadFile(plainFile); err != nil {
		return err
	}

	if fileCipher, err = my.encryptECB(plainData); err != nil {
		return err
	}

	if encryptedKeyStr, err = asymm.Encrypt(my.key); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr)

	if err = validateKey(my.key); err != nil {
		return err
	}

	if len(encryptedKeyBytes) > 0xFFFF {
		return errors.New("encrypted key too long")
	}

	out := make([]byte, 2+len(encryptedKeyBytes)+len(fileCipher))
	out[0] = byte(len(encryptedKeyBytes) >> 8)
	out[1] = byte(len(encryptedKeyBytes))
	copy(out[2:], encryptedKeyBytes)
	copy(out[2+len(encryptedKeyBytes):], fileCipher)

	return os.WriteFile(outFile, out, 0644)
}

// decryptECBFile ECB 解密文件
func (my *SM4Impl) decryptECBFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err                error
		data               []byte
		keyLen             int
		encryptedKeyBase64 string
		fileCipher         []byte
		plainData          []byte
	)

	if data, err = os.ReadFile(cipherFile); err != nil {
		return err
	}
	if len(data) < 2 {
		return os.ErrInvalid
	}

	keyLen = int(data[0])<<8 | int(data[1])
	if keyLen <= 0 || len(data) < 2+keyLen {
		return os.ErrInvalid
	}
	encryptedKeyBase64 = string(data[2 : 2+keyLen])
	fileCipher = data[2+keyLen:]

	if my.key, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(my.key) != 16 {
		return os.ErrInvalid
	}

	if plainData, err = my.decryptECB(fileCipher); err != nil {
		return err
	}

	return os.WriteFile(outFile, plainData, 0644)
}

// encryptECBLargeFile 用 SM2+SM4 ECB 流式加密大文件（TB级）
func (my *SM4Impl) encryptECBLargeFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		inF, outF         *os.File
		encryptedKeyStr   string
		encryptedKeyBytes []byte
	)

	if inF, err = os.Open(plainFile); err != nil {
		return err
	}
	defer func() { _ = inF.Close() }()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer func() { _ = outF.Close() }()

	if encryptedKeyStr, err = asymm.Encrypt(my.key); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr)
	if err = validateKey(my.key); err != nil {
		return err
	}
	if len(encryptedKeyBytes) > 0xFFFF {
		return errors.New("encrypted key too long")
	}

	if _, err = outF.Write([]byte{byte(len(encryptedKeyBytes) >> 8), byte(len(encryptedKeyBytes))}); err != nil {
		return err
	}
	if _, err = outF.Write(encryptedKeyBytes); err != nil {
		return err
	}

	return my.encryptECBStream(inF, outF)
}

// decryptECBLargeFile 用 SM2+SM4 ECB 流式解密大文件（TB级）
func (my *SM4Impl) decryptECBLargeFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err                error
		inF, outF          *os.File
		keyLenBuf          = make([]byte, 2)
		keyLen             int
		encryptedKeyBytes  []byte
		encryptedKeyBase64 string
	)

	if inF, err = os.Open(cipherFile); err != nil {
		return err
	}
	defer func() { _ = inF.Close() }()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer func() { _ = outF.Close() }()

	if _, err = io.ReadFull(inF, keyLenBuf); err != nil {
		return err
	}
	if keyLen = int(keyLenBuf[0])<<8 | int(keyLenBuf[1]); keyLen <= 0 {
		return os.ErrInvalid
	}

	encryptedKeyBytes = make([]byte, keyLen)
	if _, err = io.ReadFull(inF, encryptedKeyBytes); err != nil {
		return err
	}
	encryptedKeyBase64 = string(encryptedKeyBytes)

	if my.key, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(my.key) != 16 {
		return os.ErrInvalid
	}

	return my.decryptECBStream(inF, outF)
}

// EncryptStream 流式加密（适用于大文件，根据 Algorithm 选择 ECB/CBC）
func (my *SM4Impl) EncryptStream(in io.Reader, out io.Writer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECBStream(in, out)
	case "CBC":
		return my.encryptCBCStream(in, out)
	default:
		return errors.New("对称加密算法目前只支持：ECB/CBC")
	}
}

// DecryptStream 流式解密（适用于大文件，根据 Algorithm 选择 ECB/CBC）
func (my *SM4Impl) DecryptStream(in io.Reader, out io.Writer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECBStream(in, out)
	case "CBC":
		return my.decryptCBCStream(in, out)
	default:
		return errors.New("对称解密算法目前只支持：ECB/CBC")
	}
}

// EncryptFile 加密文件（根据 Algorithm 选择 ECB/CBC）
func (my *SM4Impl) EncryptFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECBFile(plainFile, outFile, asymm)
	case "CBC":
		return my.encryptCBCFile(plainFile, outFile, asymm)
	default:
		return errors.New("对称加密算法目前只支持：ECB/CBC")
	}
}

// DecryptFile 解密文件（根据 Algorithm 选择 ECB/CBC）
func (my *SM4Impl) DecryptFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECBFile(cipherFile, outFile, asymm)
	case "CBC":
		return my.decryptCBCFile(cipherFile, outFile, asymm)
	default:
		return errors.New("对称解密算法目前只支持：ECB/CBC")
	}
}

// EncryptLargeFile 加密大文件（根据 Algorithm 选择 ECB/CBC）
func (my *SM4Impl) EncryptLargeFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECBLargeFile(plainFile, outFile, asymm)
	case "CBC":
		return my.encryptCBCLargeFile(plainFile, outFile, asymm)
	default:
		return errors.New("对称加密算法目前只支持：ECB/CBC")
	}
}

// DecryptLargeFile 解密大文件（根据 Algorithm 选择 ECB/CBC）
func (my *SM4Impl) DecryptLargeFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECBLargeFile(cipherFile, outFile, asymm)
	case "CBC":
		return my.decryptCBCLargeFile(cipherFile, outFile, asymm)
	default:
		return errors.New("对称解密算法目前只支持：ECB/CBC")
	}
}
