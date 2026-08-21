package main

import (
	"bytes"
	"testing"

	"github.com/aid297/aid/v3/secrets"
	"github.com/aid297/aid/v3/secrets/asymmetric/sm2"
)

func TestGenerateKeyPair(t *testing.T) {
	var (
		err                        error
		sem                        secrets.Semen
		pubKeyBase64, priKeyBase64 string
		pubKeyBytes, priKeyBytes   []byte
	)

	if sem, err = sm2.NewSem(); err != nil {
		t.Fatalf("生成种子(SM2)失败：%v", err)
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
		err                        error
		semEncrypter, semDecrypter secrets.Semen
		sm2Encrypter, sm2Decrypter secrets.Asymmetric
		sm2SemAPriKeyBytes         []byte
		plainText                  = []byte("hello, SM2 非对称加密测试!")
		cipherBase64               string
		decrypted                  []byte
	)

	if semEncrypter, err = sm2.NewSem(); err != nil {
		t.Fatalf("生成加密种子失败：%v", err)
	}

	sm2Encrypter = sm2.New(semEncrypter)

	if cipherBase64, err = sm2Encrypter.Encrypt(plainText); err != nil {
		t.Fatalf("加密错误：%v", err)
	}
	t.Logf("加密结果：%s\n", cipherBase64)

	if sm2SemAPriKeyBytes, err = semEncrypter.GetPriKeyBytes(); err != nil {
		t.Fatalf("获取私钥失败：%v", err)
	}

	if semDecrypter, err = sm2.NewSem(sm2.PriKeyBytes(sm2SemAPriKeyBytes)); err != nil {
		t.Fatalf("生成解密种子失败：%v", err)
	}
	a, _ := semDecrypter.GetPriKeyBase64()
	b, _ := semDecrypter.GetPubKeyBase64()
	t.Logf("解密种子私钥：%s, 公钥：%s", a, b)

	sm2Decrypter = sm2.New(semDecrypter)

	if decrypted, err = sm2Decrypter.Decrypt(cipherBase64); err != nil {
		t.Fatalf("解密错误：%v", err)
	}
	t.Logf("解密结果：%s", decrypted)

	if !bytes.Equal(plainText, decrypted) {
		t.Fatalf("比对结果不匹配")
	}
}

func TestSignVerify(t *testing.T) {
	var (
		err                error
		semSign, semVerify secrets.Semen
		sm2Sign, sm2Verify secrets.Asymmetric
		data               = []byte("hello, SM2 数字签名测试!")
		sigHex             string
		ok                 bool
	)

	if semSign, err = sm2.NewSem(); err != nil {
		t.Fatalf("生成签名种子失败：%v", err)
	}
	sm2Sign = sm2.New(semSign)

	if sigHex, err = sm2Sign.Sign(data); err != nil {
		t.Fatalf("签名失败：%v", err)
	}
	t.Logf("签名内容(hex): %s", sigHex)

	if semVerify, err = sm2.NewSem(sm2.PriKey(semSign.GetPriKey())); err != nil {
		t.Fatalf("生成验证种子失败：%v", err)
	}
	sm2Verify = sm2.New(semVerify)

	if ok, err = sm2Verify.Verify(data, sigHex); err != nil {
		t.Fatalf("验证失败：%v", err)
	}
	if !ok {
		t.Fatal("签名验证失败")
	}

	t.Logf("验证成功")
}

func TestVerifyWithWrongData(t *testing.T) {
	var (
		err                error
		semSign, semVerify secrets.Semen
		sm2Sign, sm2Verify secrets.Asymmetric
		data               = []byte("original data")
		sigHex             string
		ok                 bool
	)

	if semSign, err = sm2.NewSem(); err != nil {
		t.Fatalf("生成种子失败：%v", err)
	}
	sm2Sign = sm2.New(semSign)

	if sigHex, err = sm2Sign.Sign(data); err != nil {
		t.Fatalf("签名失败：%v", err)
	}

	if semVerify, err = sm2.NewSem(sm2.PriKey(semSign.GetPriKey())); err != nil {
		t.Fatalf("生成解密种子失败：%v", err)
	}
	sm2Verify = sm2.New(semVerify)

	if ok, err = sm2Verify.Verify([]byte("tampered data"), sigHex); err != nil {
		t.Fatalf("验证失败：%v", err)
	}
	if ok {
		t.Fatal("篡改数据的验证应失败")
	}
	t.Logf("篡改数据正确被拒绝")
}
