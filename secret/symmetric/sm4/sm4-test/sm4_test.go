package main

import (
	"bytes"
	"testing"

	"github.com/aid297/aid/v2/secret/symmetric/sm4"
)

var (
	testKey   = []byte("1234567890abcdef")
	testIV    = []byte("abcdef1234567890")
	testPlain = []byte("hello, SM4 encrypt test!")
)

func TestECB(t *testing.T) {
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.AlgorithmECB())
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	plain, err := sm4Helper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("ECB OK: %s", plain)
}

func TestECBBase64(t *testing.T) {
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.AlgorithmECB())
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}
	b64, err := sm4Helper.EncryptBase64(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	t.Logf("加密后：%s", b64)
	plain, err := sm4Helper.DecryptBase64(b64)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("OK: %s", plain)
}

func TestCBC(t *testing.T) {
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testIV), sm4.AlgorithmCBC())
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	plain, err := sm4Helper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("CBC OK: %s", plain)
}

func TestCBCBase64(t *testing.T) {
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testIV), sm4.AlgorithmCBC())
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}
	b64, err := sm4Helper.EncryptBase64(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	t.Logf("加密成功：%s", b64)
	plain, err := sm4Helper.DecryptBase64(b64)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("OK：%s", plain)
}

func TestInvalidKey(t *testing.T) {
	sm4Helper, err := sm4.New(sm4.KeyString("shortkey"), sm4.AlgorithmECB())
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}
	_, err = sm4Helper.Encrypt(testPlain)
	if err == nil {
		t.Fatal("验证key长度错误")
	}
	t.Logf("验证不通过：%v", err)
}

func TestCTR(t *testing.T) {
	testNonce := []byte("abcdef1234567890") // 16 字节 nonce
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmCTR())
	if err != nil {
		t.Fatalf("创建 CTR 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	plain, err := sm4Helper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("CTR OK: %s", plain)
}

func TestCTRBase64(t *testing.T) {
	testNonce := []byte("abcdef1234567890")
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmCTR())
	if err != nil {
		t.Fatalf("创建 CTR 对象失败：%v", err)
	}
	b64, err := sm4Helper.EncryptBase64(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	t.Logf("加密后：%s", b64)
	plain, err := sm4Helper.DecryptBase64(b64)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("OK: %s", plain)
}

func TestCTRWithRandNonce(t *testing.T) {
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.RandCTRNonce(), sm4.AlgorithmCTR())
	if err != nil {
		t.Fatalf("创建 CTR 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	plain, err := sm4Helper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("CTR RandNonce OK: %s", plain)
}

func TestGCM(t *testing.T) {
	testNonce := []byte("abcdef123456") // 12 字节 nonce
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmGCM())
	if err != nil {
		t.Fatalf("创建 GCM 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	t.Logf("GCM 密文长度：%d 字节", len(cipherText))
	plain, err := sm4Helper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("GCM OK: %s", plain)
}

func TestGCMBase64(t *testing.T) {
	testNonce := []byte("abcdef123456")
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmGCM())
	if err != nil {
		t.Fatalf("创建 GCM 对象失败：%v", err)
	}
	b64, err := sm4Helper.EncryptBase64(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	t.Logf("加密成功：%s", b64)
	plain, err := sm4Helper.DecryptBase64(b64)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("OK：%s", plain)
}

func TestGCMWithRandNonce(t *testing.T) {
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.RandGCMNonce(), sm4.AlgorithmGCM())
	if err != nil {
		t.Fatalf("创建 GCM 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}
	plain, err := sm4Helper.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("解密失败：%v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("比对失败：结果 %s，期望 %s", plain, testPlain)
	}
	t.Logf("GCM RandNonce OK: %s", plain)
}

func TestGCMTampered(t *testing.T) {
	testNonce := []byte("abcdef123456")
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmGCM())
	if err != nil {
		t.Fatalf("创建 GCM 对象失败：%v", err)
	}
	cipherText, err := sm4Helper.Encrypt(testPlain)
	if err != nil {
		t.Fatalf("加密失败：%v", err)
	}

	// 篡改密文
	if len(cipherText) > 20 {
		cipherText[20] ^= 0xFF
	}

	_, err = sm4Helper.Decrypt(cipherText)
	if err == nil {
		t.Fatal("GCM 认证应该失败：密文已被篡改")
	}
	t.Logf("GCM 篡改检测成功：%v", err)
}

func TestInvalidNonce(t *testing.T) {
	testNonce := []byte("short") // 不足 12 字节
	sm4Helper, err := sm4.New(sm4.KeyBytes(testKey), sm4.IVBytes(testNonce), sm4.AlgorithmCTR())
	if err != nil {
		t.Fatalf("创建 CTR 对象失败：%v", err)
	}
	_, err = sm4Helper.Encrypt(testPlain)
	if err == nil {
		t.Fatal("CTR 认证应该失败：nonce 长度不足")
	}
	t.Logf("Nonce 验证成功：%v", err)
}
