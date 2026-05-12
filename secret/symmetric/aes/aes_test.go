package aes

import (
	"bytes"
	"testing"
)

var (
	testKey   = []byte("1234567890abcdef")
	testIV    = []byte("abcdef1234567890")
	testPlain = []byte("hello, AES encrypt test!")
)

/*






















































































}	}		t.Fatal("expected error for invalid key length")	if err == nil {	_, err = aesHelper.EncryptECB(testPlain)	}		t.Fatalf("创建对象失败：%v", err)	if err != nil {	aesHelper, err := New(KeyString("shortkey"))func TestInvalidKey(t *testing.T) {}	}		t.Fatalf("CBC base64 mismatch: got %s, want %s", plain, testPlain)	if !bytes.Equal(plain, testPlain) {	}		t.Fatalf("DecryptCBCBase64 failed: %v", err)	if err != nil {	plain, err := aesHelper.DecryptCBCBase64(b64)	}		t.Fatalf("EncryptCBCBase64 failed: %v", err)	if err != nil {	b64, err := aesHelper.EncryptCBCBase64(testPlain)	}		t.Fatalf("创建 CBC 对象失败：%v", err)	if err != nil {	aesHelper, err := New(KeyBytes(testKey), IVBytes(testIV))func TestCBCBase64(t *testing.T) {}	}		t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)	if !bytes.Equal(plain, testPlain) {	}		t.Fatalf("DecryptCBC failed: %v", err)	if err != nil {	plain, err := aesHelper.DecryptCBC(cipherText)	}		t.Fatalf("EncryptCBC failed: %v", err)	if err != nil {	cipherText, err := aesHelper.EncryptCBC(testPlain)	}		t.Fatalf("创建 CBC 对象失败：%v", err)	if err != nil {	aesHelper, err := New(KeyBytes(testKey), IVBytes(testIV))func TestCBC(t *testing.T) {}	}		t.Fatalf("ECB base64 mismatch: got %s, want %s", plain, testPlain)	if !bytes.Equal(plain, testPlain) {	}		t.Fatalf("DecryptECBBase64 failed: %v", err)	if err != nil {	plain, err := aesHelper.DecryptECBBase64(b64)	}		t.Fatalf("EncryptECBBase64 failed: %v", err)	if err != nil {	b64, err := aesHelper.EncryptECBBase64(testPlain)	}		t.Fatalf("创建 ECB 对象失败：%v", err)	if err != nil {	aesHelper, err := New(KeyBytes(testKey))func TestECBBase64(t *testing.T) {}	t.Logf("ECB OK: %s", plain)	}		t.Fatalf("ECB mismatch: got %s, want %s", plain, testPlain)	if !bytes.Equal(plain, testPlain) {	}		t.Fatalf("DecryptECB failed: %v", err)	if err != nil {	plain, err := aesHelper.DecryptECB(cipherText)	}		t.Fatalf("EncryptECB failed: %v", err)	if err != nil {	cipherText, err := aesHelper.EncryptECB(testPlain)	}		t.Fatalf("创建 ECB 对象失败：%v", err)	if err != nil {	aesHelper, err := New(KeyBytes(testKey))func TestECB(t *testing.T) {)	testPlain = []byte("hello, AES encrypt test!")
*/

func TestECB(t *testing.T) {
	aesHelper, err := New(KeyBytes(testKey))
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}
	cipherText, err := aesHelper.EncryptECB(testPlain)
	if err != nil {
		t.Fatalf("EncryptECB failed: %v", err)
	}
	plain, err := aesHelper.DecryptECB(cipherText)
	if err != nil {
		t.Fatalf("DecryptECB failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("ECB mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestECBBase64(t *testing.T) {
	aesHelper, err := New(KeyBytes(testKey))
	if err != nil {
		t.Fatalf("创建 ECB 对象失败：%v", err)
	}
	b64, err := aesHelper.EncryptECBBase64(testPlain)
	if err != nil {
		t.Fatalf("EncryptECBBase64 failed: %v", err)
	}
	plain, err := aesHelper.DecryptECBBase64(b64)
	if err != nil {
		t.Fatalf("DecryptECBBase64 failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("ECB base64 mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestCBC(t *testing.T) {
	aesHelper, err := New(KeyBytes(testKey), IVBytes(testIV))
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}
	cipherText, err := aesHelper.EncryptCBC(testPlain)
	if err != nil {
		t.Fatalf("EncryptCBC failed: %v", err)
	}
	plain, err := aesHelper.DecryptCBC(cipherText)
	if err != nil {
		t.Fatalf("DecryptCBC failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("CBC mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestCBCBase64(t *testing.T) {
	aesHelper, err := New(KeyBytes(testKey), IVBytes(testIV))
	if err != nil {
		t.Fatalf("创建 CBC 对象失败：%v", err)
	}
	b64, err := aesHelper.EncryptCBCBase64(testPlain)
	if err != nil {
		t.Fatalf("EncryptCBCBase64 failed: %v", err)
	}
	plain, err := aesHelper.DecryptCBCBase64(b64)
	if err != nil {
		t.Fatalf("DecryptCBCBase64 failed: %v", err)
	}
	if !bytes.Equal(plain, testPlain) {
		t.Fatalf("CBC base64 mismatch: got %s, want %s", plain, testPlain)
	}
}

func TestInvalidKey(t *testing.T) {
	aesHelper, err := New(KeyString("shortkey"))
	if err != nil {
		t.Fatalf("创建对象失败：%v", err)
	}
	_, err = aesHelper.EncryptECB(testPlain)
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}
