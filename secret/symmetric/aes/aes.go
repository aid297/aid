package aes

import (
	"bytes"
	stdaes "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aid297/aid/secret"
)

var _ secret.Symmetricer = (*AESImpl)(nil)

type AESImpl struct{ key, iv []byte }

// New 实例化 AESHelper
func New(attrs ...secret.SymmetricAttr) (my secret.Symmetricer, err error) {
	my = &AESImpl{}
	err = my.SetAttrs(attrs...)
	return
}

// SetAttrs 设置属性
func (my *AESImpl) SetAttrs(attrs ...secret.SymmetricAttr) (err error) {
	for idx := range attrs {
		if err = attrs[idx](my); err != nil {
			return
		}
	}
	return
}

// blockSize AES 块大小固定为 16 字节
const blockSize = 16

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
		return nil, fmt.Errorf("invalid blockSize: %d", size)
	}
	if length == 0 || length%size != 0 {
		return nil, errors.New("invalid data length")
	}
	unPadding := int(src[length-1])
	if unPadding == 0 || unPadding > size {
		return nil, errors.New("invalid padding size")
	}
	for _, b := range src[length-unPadding:] {
		if int(b) != unPadding {
			return nil, errors.New("invalid PKCS7 padding")
		}
	}
	return src[:length-unPadding], nil
}

// validateKey 校验密钥长度（当前对齐 SM4，仅支持 AES-128 的 16 字节）
func validateKey(key []byte) error {
	if len(key) != blockSize {
		return fmt.Errorf("aes: invalid key size %d, must be 16 bytes", len(key))
	}
	return nil
}

// validateIV 校验 IV 长度
func validateIV(iv []byte) error {
	if len(iv) != blockSize {
		return fmt.Errorf("aes: invalid iv size %d, must be 16 bytes", len(iv))
	}
	return nil
}

// GetKeyString 获取 key：string
func (my *AESImpl) GetKeyString() string { return string(my.key) }

// GetKeyBytes 获取 key：bytes
func (my *AESImpl) GetKeyBytes() []byte { return my.key }

// GetIVString 获取 iv：string
func (my *AESImpl) GetIVString() string { return string(my.iv) }

// GetIVBytes 获取 iv：bytes
func (my *AESImpl) GetIVBytes() []byte { return my.iv }

// SetKeyString 设置 key：string
func (my *AESImpl) SetKeyString(key string) { my.key = []byte(key) }

// SetKeyBytes 设置 key：bytes
func (my *AESImpl) SetKeyBytes(key []byte) { my.key = key }

// SetIVString 设置 iv：string
func (my *AESImpl) SetIVString(iv string) { my.iv = []byte(iv) }

// SetIVBytes 设置 iv：bytes
func (my *AESImpl) SetIVBytes(iv []byte) { my.iv = iv }

// EncryptECB ECB 模式加密，返回原始字节
func (my *AESImpl) EncryptECB(plainText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	block, err := stdaes.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	padded := padPKCS7(plainText, blockSize)
	cipherText := make([]byte, len(padded))
	for i := 0; i < len(padded); i += blockSize {
		block.Encrypt(cipherText[i:i+blockSize], padded[i:i+blockSize])
	}
	return cipherText, nil
}

// DecryptECB ECB 模式解密
func (my *AESImpl) DecryptECB(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%blockSize != 0 {
		return nil, errors.New("aes: invalid cipherText length")
	}
	block, err := stdaes.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(cipherText))
	for i := 0; i < len(cipherText); i += blockSize {
		block.Decrypt(plainText[i:i+blockSize], cipherText[i:i+blockSize])
	}
	return unPadPKCS7(plainText, blockSize)
}

// EncryptECBBase64 ECB 模式加密，返回 base64 字符串
func (my *AESImpl) EncryptECBBase64(plainText []byte) (string, error) {
	cipherText, err := my.EncryptECB(plainText)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptECBBase64 ECB 模式解密，输入 base64 字符串
func (my *AESImpl) DecryptECBBase64(cipherBase64 string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, fmt.Errorf("aes: base64 decode error: %w", err)
	}
	return my.DecryptECB(cipherText)
}

// EncryptCBC CBC 模式加密，返回原始字节
func (my *AESImpl) EncryptCBC(plainText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if err := validateIV(my.iv); err != nil {
		return nil, err
	}
	block, err := stdaes.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	padded := padPKCS7(plainText, blockSize)
	cipherText := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, my.iv).CryptBlocks(cipherText, padded)
	return cipherText, nil
}

// DecryptCBC CBC 模式解密
func (my *AESImpl) DecryptCBC(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if err := validateIV(my.iv); err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%blockSize != 0 {
		return nil, errors.New("aes: invalid cipherText length")
	}
	block, err := stdaes.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, my.iv).CryptBlocks(plainText, cipherText)
	return unPadPKCS7(plainText, blockSize)
}

// EncryptCBCBase64 CBC 模式加密，返回 base64 字符串
func (my *AESImpl) EncryptCBCBase64(plainText []byte) (string, error) {
	cipherText, err := my.EncryptCBC(plainText)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptCBCBase64 CBC 模式解密，输入 base64 字符串
func (my *AESImpl) DecryptCBCBase64(cipherBase64 string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, fmt.Errorf("aes: base64 decode error: %w", err)
	}
	return my.DecryptCBC(cipherText)
}

// EncryptCBCStream CBC 流式加密（适用于大文件）
func (my *AESImpl) EncryptCBCStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}
	if err := validateIV(my.iv); err != nil {
		return err
	}

	block, err := stdaes.NewCipher(my.key)
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

// DecryptCBCStream CBC 流式解密（适用于大文件）
func (my *AESImpl) DecryptCBCStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}
	if err := validateIV(my.iv); err != nil {
		return err
	}

	block, err := stdaes.NewCipher(my.key)
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
				processLen := (fullBlocks - 1) * blockSize
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
		return errors.New("aes: invalid cipherText length")
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

// EncryptCBCFile 加密文件
func (my *AESImpl) EncryptCBCFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		plainData         []byte
		fileCipher        []byte
		aesKeyAndIV       []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
		keyLen            int
		out               []byte
	)

	if plainData, err = os.ReadFile(plainFile); err != nil {
		return err
	}

	if fileCipher, err = my.EncryptCBC(plainData); err != nil {
		return err
	}

	aesKeyAndIV = append(my.key, my.iv...)
	if encryptedKeyStr, err = asymm.Encrypt(aesKeyAndIV); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr)

	keyLen = len(encryptedKeyBytes)
	out = make([]byte, 2+keyLen+len(fileCipher))
	out[0] = byte(keyLen >> 8)
	out[1] = byte(keyLen)
	copy(out[2:], encryptedKeyBytes)
	copy(out[2+keyLen:], fileCipher)

	return os.WriteFile(outFile, out, 0644)
}

// DecryptCBCFile 解密文件
func (my *AESImpl) DecryptCBCFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err                error
		data               []byte
		keyLen             int
		encryptedKeyBase64 string
		fileCipher         []byte
		aesKeyAndIV        []byte
		plainData          []byte
	)

	if data, err = os.ReadFile(cipherFile); err != nil {
		return err
	}
	if len(data) < 2 {
		return os.ErrInvalid
	}

	keyLen = int(data[0])<<8 | int(data[1])
	if len(data) < 2+keyLen {
		return os.ErrInvalid
	}
	encryptedKeyBase64 = string(data[2 : 2+keyLen])
	fileCipher = data[2+keyLen:]

	if aesKeyAndIV, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(aesKeyAndIV) != 32 {
		return os.ErrInvalid
	}
	my.key = aesKeyAndIV[:16]
	my.iv = aesKeyAndIV[16:]

	if plainData, err = my.DecryptCBC(fileCipher); err != nil {
		return err
	}

	return os.WriteFile(outFile, plainData, 0644)
}

// EncryptCBCLargeFile 用 SM2+AES 流式加密大文件（TB级）
func (my *AESImpl) EncryptCBCLargeFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		inF, outF         *os.File
		aesKeyAndIV       []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
	)

	if inF, err = os.Open(plainFile); err != nil {
		return err
	}
	defer inF.Close()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer outF.Close()

	aesKeyAndIV = append(my.key, my.iv...)
	if encryptedKeyStr, err = asymm.Encrypt(aesKeyAndIV); err != nil {
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

	return my.EncryptCBCStream(inF, outF)
}

// DecryptCBCLargeFile 用 SM2+AES 流式解密大文件（TB级）
func (my *AESImpl) DecryptCBCLargeFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err                error
		inF, outF          *os.File
		keyLenBuf          = make([]byte, 2)
		keyLen             int
		encryptedKeyBytes  []byte
		encryptedKeyBase64 string
		aesKeyAndIV        []byte
	)

	if inF, err = os.Open(cipherFile); err != nil {
		return err
	}
	defer inF.Close()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer outF.Close()

	if _, err = io.ReadFull(inF, keyLenBuf); err != nil {
		return err
	}
	keyLen = int(keyLenBuf[0])<<8 | int(keyLenBuf[1])
	if keyLen <= 0 {
		return os.ErrInvalid
	}

	encryptedKeyBytes = make([]byte, keyLen)
	if _, err = io.ReadFull(inF, encryptedKeyBytes); err != nil {
		return err
	}
	encryptedKeyBase64 = string(encryptedKeyBytes)

	if aesKeyAndIV, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(aesKeyAndIV) != 32 {
		return os.ErrInvalid
	}
	my.key = aesKeyAndIV[:16]
	my.iv = aesKeyAndIV[16:]

	return my.DecryptCBCStream(inF, outF)
}
