package rsaOAEP_test

import (
	"bytes"
	"testing"

	myrsa "github.com/aid297/aid/v2/secrets/asymmetric/rsa"
	"github.com/aid297/aid/v2/secrets/asymmetric/rsaOAEP"
)

func TestRSAOAEP_EncryptDecrypt_RoundTrip(t *testing.T) {
	sem, err := rsaOAEP.NewSem()
	if err != nil {
		t.Fatal(err)
	}
	pubSem, err := myrsa.NewSem(myrsa.PubKey(sem.GetPubKey()))
	if err != nil {
		t.Fatal(err)
	}
	enc := rsaOAEP.New(pubSem)
	dec := rsaOAEP.New(sem)

	plain := []byte("hello OAEP")
	b64, err := enc.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	got, err := dec.Decrypt(b64)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("round-trip: got %q want %q", got, plain)
	}
}

func TestRSAOAEP_EncryptDecrypt_MultiChunk(t *testing.T) {
	sem, err := rsaOAEP.NewSem()
	if err != nil {
		t.Fatal(err)
	}
	pubSem, err := myrsa.NewSem(myrsa.PubKey(sem.GetPubKey()))
	if err != nil {
		t.Fatal(err)
	}
	enc := rsaOAEP.New(pubSem)
	dec := rsaOAEP.New(sem)

	// 2048-bit RSA OAEP-SHA256 单段约 190 字节；构造更长明文触发分段
	plain := bytes.Repeat([]byte("A"), 500)
	b64, err := enc.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	got, err := dec.Decrypt(b64)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatal("multi-chunk round-trip mismatch")
	}
}

func TestRSAOAEP_SignVerify(t *testing.T) {
	signSem, err := rsaOAEP.NewSem()
	if err != nil {
		t.Fatal(err)
	}
	verifySem, err := myrsa.NewSem(myrsa.PubKey(signSem.GetPubKey()))
	if err != nil {
		t.Fatal(err)
	}
	signer := rsaOAEP.New(signSem)
	verifier := rsaOAEP.New(verifySem)

	msg := []byte("payload")
	sig, err := signer.Sign(msg)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	ok, err := verifier.Verify(msg, sig)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !ok {
		t.Fatal("Verify expected true")
	}
}
