package sm4

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aid297/aid/v2/secrets"
	"github.com/tjfoc/gmsm/sm4"
)

// Mode 加密模式
type Mode int

const blockSize = 16 // blockSize SM4 块大小固定为 16 字节

var _ secrets.Symmetric = (*SM4)(nil)

type SM4 struct {
	key, iv   []byte
	algorithm string
}

// New 实例化
func New(attrs ...secrets.SymmetricAttr) (my secrets.Symmetric, err error) {
	my = &SM4{algorithm: "CBC"}
	err = my.SetAttrs(attrs...)
	return
}

// SetAttrs 设置属性
func (my *SM4) SetAttrs(attrs ...secrets.SymmetricAttr) (err error) {
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
func (my *SM4) GetKeyString() string { return string(my.key) }

// GetKeyBytes 获取 key：bytes
func (my *SM4) GetKeyBytes() []byte { return my.key }

// GetIVString 获取 iv：string
func (my *SM4) GetIVString() string { return string(my.iv) }

// GetIVBytes 获取 iv：bytes
func (my *SM4) GetIVBytes() []byte { return my.iv }

// SetKeyString 设置 key：string
func (my *SM4) SetKeyString(key string) { my.key = []byte(key) }

// SetKeyBytes 设置 key：bytes
func (my *SM4) SetKeyBytes(key []byte) { my.key = key }

// SetIVString 设置 iv：string
func (my *SM4) SetIVString(iv string) { my.iv = []byte(iv) }

// SetIVBytes 设置 iv：bytes
func (my *SM4) SetIVBytes(iv []byte) { my.iv = iv }

// SetAlgorithm 设置算法
func (my *SM4) SetAlgorithm(algorithm string) (err error) {
	switch algorithm {
	case "ECB", "CBC", "CTR", "GCM":
		my.algorithm = strings.ToUpper(algorithm)
	default:
		err = errors.New("对称加密算法目前只支持：ECB/CBC/CTR/GCM")
	}
	return
}

// Encrypt 加密：通过原始内容
func (my *SM4) Encrypt(plainText []byte) ([]byte, error) {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECB(plainText)
	case "CBC":
		return my.encryptCBC(plainText)
	case "CTR":
		return my.encryptCTR(plainText)
	case "GCM":
		return my.encryptGCM(plainText)
	default:
		return nil, errors.New("对称加密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// Decrypt 解密：通过密文
func (my *SM4) Decrypt(cipherText []byte) ([]byte, error) {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECB(cipherText)
	case "CBC":
		return my.decryptCBC(cipherText)
	case "CTR":
		return my.decryptCTR(cipherText)
	case "GCM":
		return my.decryptGCM(cipherText)
	default:
		return nil, errors.New("对称解密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// EcryptBase64 加密：通过原始内容，返回 base64 编码的密文
func (my SM4) EncryptBase64(plainText []byte) (string, error) {
	cipherText, err := my.Encrypt(plainText)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptBase64 解密：通过 base64 编码的密文
func (my *SM4) DecryptBase64(cipherBase64 string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(secrets.PaddingBase64(cipherBase64))
	if err != nil {
		return nil, fmt.Errorf("base64解码错误：%w", err)
	}
	return my.Decrypt(cipherText)
}

// encryptECB ECB 模式加密，返回原始字节
func (my *SM4) encryptECB(plainText []byte) ([]byte, error) {
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
func (my *SM4) decryptECB(cipherText []byte) ([]byte, error) {
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
func (my *SM4) encryptCBC(plainText []byte) ([]byte, error) {
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
func (my *SM4) decryptCBC(cipherText []byte) ([]byte, error) {
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

// validateNonce 校验 Nonce 长度（CTR/GCM 使用）
// CTR: 必须 16 字节（与 blockSize 相同）
// GCM: 12 字节是标准
func validateNonce(nonce []byte) error {
	if len(nonce) != blockSize && len(nonce) != 12 {
		return fmt.Errorf("错误的nonce长度 %d，CTR模式需要16字节，GCM模式需要12字节", len(nonce))
	}
	return nil
}

// encryptCTR CTR 模式加密，返回原始字节
// CTR 模式将分组密码转换为流密码，无需填充
// 输出格式：nonce(16) + 密文
func (my *SM4) encryptCTR(plainText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if err := validateNonce(my.iv); err != nil {
		return nil, err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}

	// 输出格式：nonce(16) + 密文
	cipherText := make([]byte, 16+len(plainText))
	copy(cipherText, my.iv)

	stream := cipher.NewCTR(block, my.iv)
	stream.XORKeyStream(cipherText[16:], plainText)

	return cipherText, nil
}

// decryptCTR CTR 模式解密
// 输入格式：nonce(16) + 密文
func (my *SM4) decryptCTR(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if len(cipherText) < 16 {
		return nil, errors.New("错误的密文长度：CTR 模式密文至少需要 16 字节（nonce）")
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}

	nonce := cipherText[:16]
	actualCipher := cipherText[16:]

	plainText := make([]byte, len(actualCipher))
	stream := cipher.NewCTR(block, nonce)
	stream.XORKeyStream(plainText, actualCipher)

	return plainText, nil
}

// encryptGCM GCM 模式加密（认证加密）
// GCM 同时提供机密性和完整性认证
// 输出格式：nonce(12) + 密文 + tag(16)
func (my *SM4) encryptGCM(plainText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if err := validateNonce(my.iv); err != nil {
		return nil, err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}

	// GCM 标准 nonce 长度为 12 字节
	nonce := make([]byte, 12)
	if len(my.iv) >= 12 {
		copy(nonce, my.iv[:12])
	} else {
		// 如果 IV 不足 12 字节，使用 IV 并填充
		copy(nonce, my.iv)
	}

	gcmCipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// GCM 会自动生成随机 nonce 扩展到 12 字节，但我们直接使用传入的 nonce
	cipherText := gcmCipher.Seal(nonce, nonce, plainText, nil)

	return cipherText, nil
}

// decryptGCM GCM 模式解密（认证解密）
// 输入格式：nonce(12) + 密文 + tag(16)
// 如果 tag 验证失败，返回错误
func (my *SM4) decryptGCM(cipherText []byte) ([]byte, error) {
	if err := validateKey(my.key); err != nil {
		return nil, err
	}
	if len(cipherText) < 28 { // 12 (nonce) + 16 (tag) minimum
		return nil, errors.New("错误的密文长度：GCM 模式密文至少需要 28 字节（12 字节 nonce + 16 字节 tag）")
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return nil, err
	}

	gcmCipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := cipherText[:12]
	// GCM 的 Open 方法会自动提取并验证 tag
	plainText, err := gcmCipher.Open(nil, nonce, cipherText[12:], nil)
	if err != nil {
		return nil, errors.New("GCM 认证失败：密文已被篡改或密钥错误")
	}

	return plainText, nil
}

// encryptCBCStream CBC 流式加密（适用于大文件）
func (my *SM4) encryptCBCStream(in io.Reader, out io.Writer) error {
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
func (my *SM4) decryptCBCStream(in io.Reader, out io.Writer) error {
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
func (my *SM4) encryptCBCFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
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
func (my *SM4) decryptCBCFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
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
func (my *SM4) encryptCBCLargeFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
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
func (my *SM4) decryptCBCLargeFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
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
func (my *SM4) encryptECBStream(in io.Reader, out io.Writer) error {
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
func (my *SM4) decryptECBStream(in io.Reader, out io.Writer) error {
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
func (my *SM4) encryptECBFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
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
func (my *SM4) decryptECBFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
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
func (my *SM4) encryptECBLargeFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
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
func (my *SM4) decryptECBLargeFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
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

// EncryptStream 流式加密（适用于大文件，根据 Algorithm 选择 ECB/CBC/CTR/GCM）
func (my *SM4) EncryptStream(in io.Reader, out io.Writer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECBStream(in, out)
	case "CBC":
		return my.encryptCBCStream(in, out)
	case "CTR":
		return my.encryptCTRStream(in, out)
	case "GCM":
		return my.encryptGCMStream(in, out)
	default:
		return errors.New("对称加密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// DecryptStream 流式解密（适用于大文件，根据 Algorithm 选择 ECB/CBC/CTR/GCM）
func (my *SM4) DecryptStream(in io.Reader, out io.Writer) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECBStream(in, out)
	case "CBC":
		return my.decryptCBCStream(in, out)
	case "CTR":
		return my.decryptCTRStream(in, out)
	case "GCM":
		return my.decryptGCMStream(in, out)
	default:
		return errors.New("对称解密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// EncryptFile 加密文件（根据 Algorithm 选择 ECB/CBC/CTR/GCM）
func (my *SM4) EncryptFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECBFile(plainFile, outFile, asymm)
	case "CBC":
		return my.encryptCBCFile(plainFile, outFile, asymm)
	case "CTR":
		return my.encryptCTRFile(plainFile, outFile, asymm)
	case "GCM":
		return my.encryptGCMFile(plainFile, outFile, asymm)
	default:
		return errors.New("对称加密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// DecryptFile 解密文件（根据 Algorithm 选择 ECB/CBC/CTR/GCM）
func (my *SM4) DecryptFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECBFile(cipherFile, outFile, asymm)
	case "CBC":
		return my.decryptCBCFile(cipherFile, outFile, asymm)
	case "CTR":
		return my.decryptCTRFile(cipherFile, outFile, asymm)
	case "GCM":
		return my.decryptGCMFile(cipherFile, outFile, asymm)
	default:
		return errors.New("对称解密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// EncryptLargeFile 加密大文件（根据 Algorithm 选择 ECB/CBC/CTR/GCM）
func (my *SM4) EncryptLargeFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.encryptECBLargeFile(plainFile, outFile, asymm)
	case "CBC":
		return my.encryptCBCLargeFile(plainFile, outFile, asymm)
	case "CTR":
		return my.encryptCTRLargeFile(plainFile, outFile, asymm)
	case "GCM":
		return my.encryptGCMLargeFile(plainFile, outFile, asymm)
	default:
		return errors.New("对称加密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// DecryptLargeFile 解密大文件（根据 Algorithm 选择 ECB/CBC/CTR/GCM）
func (my *SM4) DecryptLargeFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
	switch strings.ToUpper(my.algorithm) {
	case "ECB":
		return my.decryptECBLargeFile(cipherFile, outFile, asymm)
	case "CBC":
		return my.decryptCBCLargeFile(cipherFile, outFile, asymm)
	case "CTR":
		return my.decryptCTRLargeFile(cipherFile, outFile, asymm)
	case "GCM":
		return my.decryptGCMLargeFile(cipherFile, outFile, asymm)
	default:
		return errors.New("对称解密算法目前只支持：ECB/CBC/CTR/GCM")
	}
}

// encryptCTRStream CTR 流式加密（适用于大文件）
func (my *SM4) encryptCTRStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}
	if err := validateNonce(my.iv); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, my.iv)

	var (
		readBuf = make([]byte, 1024*1024)
	)

	for {
		n, readErr := in.Read(readBuf)
		if n > 0 {
			plainChunk := readBuf[:n]
			cipherChunk := make([]byte, n)
			stream.XORKeyStream(cipherChunk, plainChunk)
			if _, err = out.Write(cipherChunk); err != nil {
				return err
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return nil
}

// decryptCTRStream CTR 流式解密（适用于大文件）
func (my *SM4) decryptCTRStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}
	if err := validateNonce(my.iv); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, my.iv)

	var (
		readBuf = make([]byte, 1024*1024)
	)

	for {
		n, readErr := in.Read(readBuf)
		if n > 0 {
			cipherChunk := readBuf[:n]
			plainChunk := make([]byte, n)
			stream.XORKeyStream(plainChunk, cipherChunk)
			if _, err = out.Write(plainChunk); err != nil {
				return err
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return nil
}

// encryptGCMStream GCM 流式加密（适用于大文件）
// 注意：GCM 是认证加密，流式加密时无法在解密端验证完整性
// 如果需要验证完整性，请在完整数据流加密后单独验证，或使用分块 GCM
func (my *SM4) encryptGCMStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}
	if err := validateNonce(my.iv); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}

	nonce := make([]byte, 12)
	if len(my.iv) >= 12 {
		copy(nonce, my.iv[:12])
	} else {
		copy(nonce, my.iv)
	}

	// 生成随机 nonce 写入文件头部
	if _, err := out.Write(nonce); err != nil {
		return err
	}

	gcmCipher, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 流式加密：累积数据后一次性认证
	var (
		pending []byte
		readBuf = make([]byte, 1024*1024)
	)

	for {
		n, readErr := in.Read(readBuf)
		if n > 0 {
			pending = append(pending, readBuf[:n]...)
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	// 整体加密并认证
	cipherText := gcmCipher.Seal(nil, nonce, pending, nil)
	_, err = out.Write(cipherText)
	return err
}

// decryptGCMStream GCM 流式解密（适用于大文件）
func (my *SM4) decryptGCMStream(in io.Reader, out io.Writer) error {
	if err := validateKey(my.key); err != nil {
		return err
	}

	block, err := sm4.NewCipher(my.key)
	if err != nil {
		return err
	}

	// 读取 nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(in, nonce); err != nil {
		return errors.New("GCM 流式解密失败：无法读取 nonce")
	}

	gcmCipher, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 读取所有密文
	cipherData, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	// 解密并验证
	plainText, err := gcmCipher.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return errors.New("GCM 认证失败：密文已被篡改或密钥错误")
	}

	_, err = out.Write(plainText)
	return err
}

// encryptCTRFile CTR 加密文件
func (my *SM4) encryptCTRFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err               error
		plainData         []byte
		fileCipher        []byte
		sm4KeyAndNonce    []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
		keyLen            int
		out               []byte
	)

	if plainData, err = os.ReadFile(plainFile); err != nil {
		return err
	}

	if fileCipher, err = my.encryptCTR(plainData); err != nil {
		return err
	}

	// CTR 使用 16 字节 nonce
	nonce := fileCipher[:16]
	sm4KeyAndNonce = append(my.key, nonce...)
	if encryptedKeyStr, err = asymm.Encrypt(sm4KeyAndNonce); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr)

	if err = validateKey(my.key); err != nil {
		return err
	}

	keyLen = len(encryptedKeyBytes)
	out = make([]byte, 2+keyLen+len(fileCipher))
	out[0] = byte(keyLen >> 8)
	out[1] = byte(keyLen)
	copy(out[2:], encryptedKeyBytes)
	copy(out[2+keyLen:], fileCipher)

	return os.WriteFile(outFile, out, 0644)
}

// decryptCTRFile CTR 解密文件
func (my *SM4) decryptCTRFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err                error
		data               []byte
		keyLen             int
		encryptedKeyBase64 string
		fileCipher         []byte
		sm4KeyAndNonce     []byte
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

	if sm4KeyAndNonce, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(sm4KeyAndNonce) != 16+16 {
		return os.ErrInvalid
	}
	my.key = sm4KeyAndNonce[:16]
	my.iv = sm4KeyAndNonce[16:] // 16 字节是 nonce

	if plainData, err = my.decryptCTR(fileCipher); err != nil {
		return err
	}

	return os.WriteFile(outFile, plainData, 0644)
}

// encryptGCMFile GCM 加密文件
func (my *SM4) encryptGCMFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err               error
		plainData         []byte
		fileCipher        []byte
		sm4KeyAndNonce    []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
		keyLen            int
		out               []byte
	)

	if plainData, err = os.ReadFile(plainFile); err != nil {
		return err
	}

	if fileCipher, err = my.encryptGCM(plainData); err != nil {
		return err
	}

	// GCM 使用 12 字节 nonce
	nonce := fileCipher[:12]
	sm4KeyAndNonce = append(my.key, nonce...)
	if encryptedKeyStr, err = asymm.Encrypt(sm4KeyAndNonce); err != nil {
		return err
	}
	encryptedKeyBytes = []byte(encryptedKeyStr)

	if err = validateKey(my.key); err != nil {
		return err
	}

	keyLen = len(encryptedKeyBytes)
	out = make([]byte, 2+keyLen+len(fileCipher))
	out[0] = byte(keyLen >> 8)
	out[1] = byte(keyLen)
	copy(out[2:], encryptedKeyBytes)
	copy(out[2+keyLen:], fileCipher)

	return os.WriteFile(outFile, out, 0644)
}

// decryptGCMFile GCM 解密文件
func (my *SM4) decryptGCMFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err                error
		data               []byte
		keyLen             int
		encryptedKeyBase64 string
		fileCipher         []byte
		sm4KeyAndNonce     []byte
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

	if sm4KeyAndNonce, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(sm4KeyAndNonce) != 16+12 {
		return os.ErrInvalid
	}
	my.key = sm4KeyAndNonce[:16]
	my.iv = sm4KeyAndNonce[16:] // 前 12 字节是 nonce

	if plainData, err = my.decryptGCM(fileCipher); err != nil {
		return err
	}

	return os.WriteFile(outFile, plainData, 0644)
}

// encryptCTRLargeFile 用 SM2+SM4 CTR 流式加密大文件（TB级）
func (my *SM4) encryptCTRLargeFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err               error
		inF, outF         *os.File
		sm4KeyAndNonce    []byte
		encryptedKeyStr   string
		encryptedKeyBytes []byte
		nonce             []byte
	)

	if inF, err = os.Open(plainFile); err != nil {
		return err
	}
	defer func() { _ = inF.Close() }()

	if outF, err = os.Create(outFile); err != nil {
		return err
	}
	defer func() { _ = outF.Close() }()

	// 生成随机 16 字节 nonce
	nonce = make([]byte, 16)
	if _, err = rand.Read(nonce); err != nil {
		return err
	}

	sm4KeyAndNonce = append(my.key, nonce...)
	if encryptedKeyStr, err = asymm.Encrypt(sm4KeyAndNonce); err != nil {
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

	// 设置 nonce 到 iv
	my.iv = nonce

	return my.encryptCTRStream(inF, outF)
}

// decryptCTRLargeFile 用 SM2+SM4 CTR 流式解密大文件（TB级）
func (my *SM4) decryptCTRLargeFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err                error
		inF, outF          *os.File
		keyLenBuf          = make([]byte, 2)
		keyLen             int
		encryptedKeyBytes  []byte
		encryptedKeyBase64 string
		sm4KeyAndNonce     []byte
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

	if sm4KeyAndNonce, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(sm4KeyAndNonce) != 16+16 {
		return os.ErrInvalid
	}
	my.key = sm4KeyAndNonce[:16]
	my.iv = sm4KeyAndNonce[16:]

	return my.decryptCTRStream(inF, outF)
}

// encryptGCMLargeFile 用 SM2+SM4 GCM 流式加密大文件（TB级）
func (my *SM4) encryptGCMLargeFile(plainFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err               error
		inF, outF         *os.File
		sm4KeyAndNonce    []byte
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

	// 生成随机 nonce
	nonce := make([]byte, 12)
	if _, err = rand.Read(nonce); err != nil {
		return err
	}

	sm4KeyAndNonce = append(my.key, nonce...)
	if encryptedKeyStr, err = asymm.Encrypt(sm4KeyAndNonce); err != nil {
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

	// 设置 nonce 到 iv
	my.iv = nonce

	return my.encryptGCMStream(inF, outF)
}

// decryptGCMLargeFile 用 SM2+SM4 GCM 流式解密大文件（TB级）
func (my *SM4) decryptGCMLargeFile(cipherFile, outFile string, asymm secrets.Asymmetric) error {
	var (
		err                error
		inF, outF          *os.File
		keyLenBuf          = make([]byte, 2)
		keyLen             int
		encryptedKeyBytes  []byte
		encryptedKeyBase64 string
		sm4KeyAndNonce     []byte
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

	if sm4KeyAndNonce, err = asymm.Decrypt(encryptedKeyBase64); err != nil {
		return err
	}
	if len(sm4KeyAndNonce) != 16+12 {
		return os.ErrInvalid
	}
	my.key = sm4KeyAndNonce[:16]
	my.iv = sm4KeyAndNonce[16:]

	return my.decryptGCMStream(inF, outF)
}
