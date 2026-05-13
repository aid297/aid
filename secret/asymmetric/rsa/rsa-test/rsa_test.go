package main

import (
	"bytes"
	"testing"

	"github.com/aid297/aid/v2/secret/asymmetric/rsa"
)

func TestEncryptDecrypt(t *testing.T) {
	pub, pri, err := rsa.GenerateKeyPairBase64(2048)
	if err != nil {
		t.Fatalf("generate key failed: %v", err)
	}

	helper, err := rsa.NewByBase64(pub, pri)
	if err != nil {
		t.Fatalf("new helper failed: %v", err)
	}

	plain := []byte("hello, RSA interface test!")
	cipherBase64, err := helper.Encrypt(plain)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := helper.Decrypt(cipherBase64)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if !bytes.Equal(plain, decrypted) {
		t.Fatalf("decrypt mismatch: got %s, want %s", decrypted, plain)
	}
}

func TestSignVerify(t *testing.T) {
	pub, pri, err := rsa.GenerateKeyPairBase64(2048)
	if err != nil {
		t.Fatalf("generate key failed: %v", err)
	}

	helper, err := rsa.NewByBase64(pub, pri)
	if err != nil {
		t.Fatalf("new helper failed: %v", err)
	}

	data := []byte("hello, rsa sign")
	sig, err := helper.Sign(data)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	ok, err := helper.Verify(data, sig)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if !ok {
		t.Fatalf("verify should be true")
	}

	ok, err = helper.Verify([]byte("tampered"), sig)
	if err != nil {
		t.Fatalf("verify tampered failed: %v", err)
	}
	if ok {
		t.Fatalf("verify should be false for tampered data")
	}
}
