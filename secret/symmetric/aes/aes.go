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
	"strings"

	"github.com/aid297/aid/v2/secret"
)

var _ secret.Symmetricer = (*AESImpl)(nil)

const (
	AES128 = 128
	AES192 = 192
	AES256 = 256
)

type AESImpl struct {
	key, iv   []byte
	keyBits   int
	algorithm string
}

// New 实例化
func New(attrs ...secret.SymmetricAttr) (my secret.Symmetricer, err error) {
	my = &AESImpl{keyBits: AES128, algorithm: "CBC"}
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

const (
	fileHeaderMagic   byte = 0xA5
	fileHeaderVersion byte = 0x01
)

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
		return nil, errors.New("错误的padding长度")
	}
	for _, b := range src[length-unPadding:] {
		if int(b) != unPadding {
			return nil, errors.New("错误的PKCS7 padding长度")
		}
	}
	return src[:length-unPadding], nil
}

func validateKeySize(keySize int) error {
	switch keySize {
	case 16, 24, 32:
		return nil
	default:
		return fmt.Errorf("错误的key长度 %d，必须16/24/32字节", keySize)
	}
}

// validateKey 校验密钥长度（支持 AES-128/192/256）
func validateKey(key []byte) error {
	return validateKeySize(len(key))
}

// validateIV 校验 IV 长度
func validateIV(iv []byte) error {
	if len(iv) != blockSize {
		return fmt.Errorf("错误的iv长度 %d，必须16字节", len(iv))
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

// SetAlgorithm 设置算法
func (my *AESImpl) SetAlgorithm(algorithm string) (err error) {
	switch algorithm {
	case "ECB", "CBC":
		my.algorithm = algorithm
	default:
		err = errors.New("对称加密算法目前只支持：ECB/CBC")
	}
	return
}

// Encrypt 加密：通过原始内容
func (my *AESImpl) Encrypt(plainText []byte) ([]byte, error) {
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
func (my *AESImpl) Decrypt(cipherText []byte) ([]byte, error) {
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
func (my *AESImpl) EncryptBase64(plainText []byte) (string, error) {
	cipherText, err := my.Encrypt(plainText)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptBase64 解密：通过 base64 编码的密文
func (my *AESImpl) DecryptBase64(cipherBase64 string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, fmt.Errorf("base64解码错误：%w", err)
	}
	return my.decryptECB(cipherText)
}

// encryptECB ECB 模式加密，返回原始字节
func (my *AESImpl) encryptECB(plainText []byte) ([]byte, error) {
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

// decryptECB ECB 模式解密
func (my *AESImpl) decryptECB(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%blockSize != 0 {
		return nil, errors.New("错误的密文长度")
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

// encryptCBC CBC 模式加密，返回原始字节
func (my *AESImpl) encryptCBC(plainText []byte) ([]byte, error) {
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

// decryptCBC CBC 模式解密
func (my *AESImpl) decryptCBC(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if err := validateIV(my.iv); err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%blockSize != 0 {
		return nil, errors.New("错误的密文长度")
	}
	block, err := stdaes.NewCipher(my.key)
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, my.iv).CryptBlocks(plainText, cipherText)
	return unPadPKCS7(plainText, blockSize)
}

// encryptCBCStream CBC 流式加密（适用于大文件）
func (my *AESImpl) encryptCBCStream(in io.Reader, out io.Writer) error {
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

// decryptCBCStream CBC 流式解密（适用于大文件）
func (my *AESImpl) decryptCBCStream(in io.Reader, out io.Writer) error {
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

// encryptCBCFile CBC 加密文件
func (my *AESImpl) encryptCBCFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
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

	if fileCipher, err = my.encryptCBC(plainData); err != nil {
		return err
	}

	aesKeyAndIV = append(my.key, my.iv...)
	if encryptedKeyStr, err = asymm.Encrypt(aesKeyAndIV); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr)

	if err = validateKey(my.key); err != nil {
		return err
	}

	keyLen = len(encryptedKeyBytes)
	out = make([]byte, 5+keyLen+len(fileCipher))
	out[0] = fileHeaderMagic
	out[1] = fileHeaderVersion
	out[2] = byte(len(my.key))
	out[3] = byte(keyLen >> 8)
	out[4] = byte(keyLen)
	copy(out[5:], encryptedKeyBytes)
	copy(out[5+keyLen:], fileCipher)

	return os.WriteFile(outFile, out, 0644)
}

// decryptCBCFile CBC 解密文件
func (my *AESImpl) decryptCBCFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
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

	keySize := 16 // 兼容旧格式默认 AES-128
	offset := 2
	if len(data) >= 5 && data[0] == fileHeaderMagic {
		if data[1] != fileHeaderVersion {
			return os.ErrInvalid
		}
		keySize = int(data[2])
		if err = validateKeySize(keySize); err != nil {
			return err
		}
		keyLen = int(data[3])<<8 | int(data[4])
		offset = 5
	} else {
		keyLen = int(data[0])<<8 | int(data[1])
	}
	if len(data) < offset+keyLen {
		return os.ErrInvalid
	}
	encryptedKeyBase64 = string(data[offset : offset+keyLen])
	fileCipher = data[offset+keyLen:]

	if aesKeyAndIV, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(aesKeyAndIV) != keySize+blockSize {
		return os.ErrInvalid
	}
	my.key = aesKeyAndIV[:keySize]
	my.iv = aesKeyAndIV[keySize:]

	if plainData, err = my.decryptCBC(fileCipher); err != nil {
		return err
	}

	return os.WriteFile(outFile, plainData, 0644)
}

// encryptCBCLargeFile 用 SM2+AES 流式加密大文件（TB级）
func (my *AESImpl) encryptCBCLargeFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
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
	if err = validateKey(my.key); err != nil {
		return err
	}
	if len(encryptedKeyBytes) > 0xFFFF {
		return errors.New("encrypted key too long")
	}

	if _, err = outF.Write([]byte{fileHeaderMagic, fileHeaderVersion, byte(len(my.key)), byte(len(encryptedKeyBytes) >> 8), byte(len(encryptedKeyBytes))}); err != nil {
		return err
	}
	if _, err = outF.Write(encryptedKeyBytes); err != nil {
		return err
	}

	return my.encryptCBCStream(inF, outF)
}

// decryptCBCLargeFile 用 SM2+AES 流式解密大文件（TB级）
func (my *AESImpl) decryptCBCLargeFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
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

	head := make([]byte, 1)
	if _, err = io.ReadFull(inF, head); err != nil {
		return err
	}

	keySize := 16 // 兼容旧格式默认 AES-128
	if head[0] == fileHeaderMagic {
		newHeader := make([]byte, 4)
		if _, err = io.ReadFull(inF, newHeader); err != nil {
			return err
		}
		if newHeader[0] != fileHeaderVersion {
			return os.ErrInvalid
		}
		keySize = int(newHeader[1])
		if err = validateKeySize(keySize); err != nil {
			return err
		}
		keyLen = int(newHeader[2])<<8 | int(newHeader[3])
	} else {
		keyLenBuf[0] = head[0]
		if _, err = io.ReadFull(inF, keyLenBuf[1:2]); err != nil {
			return err
		}
		keyLen = int(keyLenBuf[0])<<8 | int(keyLenBuf[1])
	}

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
	if len(aesKeyAndIV) != keySize+blockSize {
		return os.ErrInvalid
	}
	my.key = aesKeyAndIV[:keySize]
	my.iv = aesKeyAndIV[keySize:]

	return my.decryptCBCStream(inF, outF)
}

// encryptECBStream ECB 流式加密（适用于大文件）
func (my *AESImpl) encryptECBStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}

	block, err := stdaes.NewCipher(my.key)
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
func (my *AESImpl) decryptECBStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}

	block, err := stdaes.NewCipher(my.key)
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
func (my *AESImpl) encryptECBFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		plainData         []byte
		fileCipher        []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
		keyLen            int
		out               []byte
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

	keyLen = len(encryptedKeyBytes)
	out = make([]byte, 5+keyLen+len(fileCipher))
	out[0] = fileHeaderMagic
	out[1] = fileHeaderVersion
	out[2] = byte(len(my.key))
	out[3] = byte(keyLen >> 8)
	out[4] = byte(keyLen)
	copy(out[5:], encryptedKeyBytes)
	copy(out[5+keyLen:], fileCipher)

	return os.WriteFile(outFile, out, 0644)
}

// decryptECBFile ECB 解密文件
func (my *AESImpl) decryptECBFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
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

	keySize := 16 // 兼容旧格式默认 AES-128
	offset := 2
	if len(data) >= 5 && data[0] == fileHeaderMagic {
		if data[1] != fileHeaderVersion {
			return os.ErrInvalid
		}
		keySize = int(data[2])
		if err = validateKeySize(keySize); err != nil {
			return err
		}
		keyLen = int(data[3])<<8 | int(data[4])
		offset = 5
	} else {
		keyLen = int(data[0])<<8 | int(data[1])
	}
	if len(data) < offset+keyLen {
		return os.ErrInvalid
	}
	encryptedKeyBase64 = string(data[offset : offset+keyLen])
	fileCipher = data[offset+keyLen:]

	if my.key, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(my.key) != keySize {
		return os.ErrInvalid
	}

	if plainData, err = my.decryptECB(fileCipher); err != nil {
		return err
	}

	return os.WriteFile(outFile, plainData, 0644)
}

// encryptECBLargeFile 用 SM2+AES ECB 流式加密大文件（TB级）
func (my *AESImpl) encryptECBLargeFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
	var (
		err               error
		inF, outF         *os.File
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

	if _, err = outF.Write([]byte{fileHeaderMagic, fileHeaderVersion, byte(len(my.key)), byte(len(encryptedKeyBytes) >> 8), byte(len(encryptedKeyBytes))}); err != nil {
		return err
	}
	if _, err = outF.Write(encryptedKeyBytes); err != nil {
		return err
	}

	return my.encryptECBStream(inF, outF)
}

// decryptECBLargeFile 用 SM2+AES ECB 流式解密大文件（TB级）
func (my *AESImpl) decryptECBLargeFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
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
	defer inF.Close()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer outF.Close()

	head := make([]byte, 1)
	if _, err = io.ReadFull(inF, head); err != nil {
		return err
	}

	keySize := 16 // 兼容旧格式默认 AES-128
	if head[0] == fileHeaderMagic {
		newHeader := make([]byte, 4)
		if _, err = io.ReadFull(inF, newHeader); err != nil {
			return err
		}
		if newHeader[0] != fileHeaderVersion {
			return os.ErrInvalid
		}
		keySize = int(newHeader[1])
		if err = validateKeySize(keySize); err != nil {
			return err
		}
		keyLen = int(newHeader[2])<<8 | int(newHeader[3])
	} else {
		keyLenBuf[0] = head[0]
		if _, err = io.ReadFull(inF, keyLenBuf[1:2]); err != nil {
			return err
		}
		keyLen = int(keyLenBuf[0])<<8 | int(keyLenBuf[1])
	}

	if keyLen <= 0 {
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
	if len(my.key) != keySize {
		return os.ErrInvalid
	}

	return my.decryptECBStream(inF, outF)
}

// EncryptStream 流式加密（适用于大文件，根据 Algorithm 选择 ECB/CBC）
func (my *AESImpl) EncryptStream(in io.Reader, out io.Writer) error {
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
func (my *AESImpl) DecryptStream(in io.Reader, out io.Writer) error {
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
func (my *AESImpl) EncryptFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
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
func (my *AESImpl) DecryptFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
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
func (my *AESImpl) EncryptLargeFile(plainFile, outFile string, asymm secret.Asymmetricer) error {
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
func (my *AESImpl) DecryptLargeFile(cipherFile, outFile string, asymm secret.Asymmetricer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECBLargeFile(cipherFile, outFile, asymm)
	case "CBC":
		return my.decryptCBCLargeFile(cipherFile, outFile, asymm)
	default:
		return errors.New("对称解密算法目前只支持：ECB/CBC")
	}
}
