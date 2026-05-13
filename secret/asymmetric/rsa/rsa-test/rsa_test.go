package main

import (
	"bytes"
	"testing"

	"github.com/aid297/aid/v2/secret"
	"github.com/aid297/aid/v2/secret/asymmetric/rsa"
)

func TestGenerateKeyPair(t *testing.T) {
	var (
		err                        error
		sem                        secret.Semener
		pubKeyBase64, priKeyBase64 string
		pubKeyBytes, priKeyBytes   []byte
	)

	if sem, err = rsa.NewSem(); err != nil {
		t.Fatalf("生成种子(RSA)失败：%v", err)
	}

	if err = sem.GeneratePriKey(); err != nil {
		t.Fatalf("生成私钥失败：%v", err)
	}

	if pubKeyBase64, err = sem.GetPubKeyBase64(); pubKeyBase64 == "" {
		t.Fatalf("公钥长度为空：base64")
	}

	if pubKeyBytes, err = sem.GetPubKeyBytes(); len(pubKeyBytes) == 0 {
		t.Fatalf("公钥长度为空：bytes")
	}

	if priKeyBase64, err = sem.GetPriKeyBase64(); priKeyBase64 == "" {
		t.Fatalf("私钥长度为空：base64")
	}

	if priKeyBytes, err = sem.GetPriKeyBytes(); len(priKeyBytes) == 0 {
		t.Fatalf("私钥长度为空：bytes")
	}

	t.Logf("公钥：%s\n", pubKeyBase64)
	t.Logf("私钥：%s\n", priKeyBase64)
}

func TestEncryptDecrypt(t *testing.T) {
	var (
		err                error
		semA, semB         secret.Semener
		rsaA, rsaB         secret.Asymmetricer
		rsaSemAPriKeyBytes []byte
		plainText          = []byte("hello, RSA 非对称加密测试!")
		cipherBase64       string
		decrypted          []byte
	)

	if semA, err = rsa.NewSem(); err != nil {
		t.Fatalf("生成加密种子失败：%v", err)
	}
	rsaA = rsa.New(semA)

	if cipherBase64, err = rsaA.Encrypt(plainText); err != nil {
		t.Fatalf("加密错误：%v", err)
	}
	t.Logf("加密结果：%s\n", cipherBase64)

	if rsaSemAPriKeyBytes, err = semA.GetPriKeyBytes(); err != nil {
		t.Fatalf("获取私钥失败：%v", err)
	}

	if semB, err = rsa.NewSem(rsa.PriKeyBytes(rsaSemAPriKeyBytes)); err != nil {
		t.Fatalf("生成解密种子失败：%v", err)
	}
	a, _ := semB.GetPriKeyBase64()
	b, _ := semB.GetPubKeyBase64()
	t.Logf("解密种子私钥：%s, 公钥：%s", a, b)

	rsaB = rsa.New(semB)

	if decrypted, err = rsaB.Decrypt(cipherBase64); err != nil {
		t.Fatalf("解密错误：%v", err)
	}
	t.Logf("解密结果：%s", decrypted)

	if !bytes.Equal(plainText, decrypted) {
		t.Fatalf("比对结果不匹配")
	}
}

func TestSignVerify(t *testing.T) {
	var (
		err       error
		semA      secret.Semener
		rsaHelper secret.Asymmetricer
		data      = []byte("hello, RSA 数字签名测试!")
		sigHex    string
		ok        bool
	)

	if semA, err = rsa.NewSem(); err != nil {
		t.Fatalf("生成种子失败：%v", err)
	}

	rsaHelper = rsa.New(semA)

	if sigHex, err = rsaHelper.Sign(data); err != nil {
		t.Fatalf("签名失败：%v", err)
	}
	t.Logf("签名内容(hex): %s", sigHex)

	if ok, err = rsaHelper.Verify(data, sigHex); err != nil {
		t.Fatalf("验证失败：%v", err)
	}
	if !ok {
		t.Fatal("签名验证失败")
	}

	t.Logf("验证成功")
}

func TestVerifyWithWrongData(t *testing.T) {
	var (
		err        error
		semA, semB secret.Semener
		rsaA, rsaB secret.Asymmetricer
		data       = []byte("original data")
		sigHex     string
		ok         bool
	)
	if semA, err = rsa.NewSem(); err != nil {
		t.Fatalf("生成种子失败：%v", err)
	}

	if err = rsa.MustGeneratePriKey(semA); err != nil {
		t.Fatalf("%v", err)
	}

	rsaA = rsa.New(semA)
	if sigHex, err = rsaA.Sign(data); err != nil {
		t.Fatalf("签名失败：%v", err)
	}

	semB, err = rsa.NewSem(rsa.PriKey(semA.GetPriKey()))
	if err != nil {
		t.Fatalf("生成解密种子失败：%v", err)
	}

	rsaB = rsa.New(semB)
	if ok, err = rsaB.Verify([]byte("tampered data"), sigHex); err != nil {
		t.Fatalf("验证失败：%v", err)
	}
	if ok {
		t.Fatal("篡改数据的验证应失败")
	}
	t.Logf("篡改数据正确被拒绝")
}
