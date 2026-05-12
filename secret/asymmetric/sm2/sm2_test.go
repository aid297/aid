package sm2

import (
	"bytes"
	"testing"

	"github.com/aid297/aid/v2/secret"
)

func TestGenerateKeyPair(t *testing.T) {
	var (
		err                        error
		sem                        secret.Semener
		pubKeyBase64, priKeyBase64 string
		pubKeyBytes, priKeyBytes   []byte
	)

	if sem, err = NewSem(); err != nil {
		t.Fatalf("生成种子失败：%v", err)
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
		sm2A, sm2B         secret.Asymmetricer
		sm2SemAPriKeyBytes []byte
		plainText          = []byte("hello, SM2 非对称加密测试!")
		cipherBase64       string
		decrypted          []byte
	)

	if semA, err = NewSem(); err != nil {
		t.Fatalf("生成加密种子失败：%v", err)
	}

	sm2A = New(semA)

	if cipherBase64, err = sm2A.Encrypt(plainText); err != nil {
		t.Fatalf("加密错误：%v", err)
	}
	t.Logf("加密结果：%s\n", cipherBase64)

	if sm2SemAPriKeyBytes, err = semA.GetPriKeyBytes(); err != nil {
		t.Fatalf("获取私钥失败：%v", err)
	}

	if semB, err = NewSem(PriKeyBytes(sm2SemAPriKeyBytes)); err != nil {
		t.Fatalf("生成解密种子失败：%v", err)
	}
	a, _ := semB.GetPriKeyBase64()
	b, _ := semB.GetPubKeyBase64()
	t.Logf("解密种子私钥：%s, 公钥：%s", a, b)

	sm2B = New(semB)

	if decrypted, err = sm2B.Decrypt(cipherBase64); err != nil {
		t.Fatalf("解密错误：%v", err)
	}
	t.Logf("解密结果：%s", decrypted)

	if !bytes.Equal(plainText, decrypted) {
		t.Fatalf("比对结果不匹配")
	}
}

func TestSignVerify(t *testing.T) {
	var (
		err    error
		semA   secret.Semener
		sm2    secret.Asymmetricer
		data   = []byte("hello, SM2 数字签名测试!")
		sigHex string
		ok     bool
	)

	if semA, err = NewSem(); err != nil {
		t.Fatalf("生成种子失败：%v", err)
	}

	sm2 = New(semA)

	if sigHex, err = sm2.Sign(data); err != nil {
		t.Fatalf("签名失败：%v", err)
	}
	t.Logf("签名内容(hex): %s", sigHex)

	if ok, err = sm2.Verify(data, sigHex); err != nil {
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
		sm2A, sm2B secret.Asymmetricer
		data       = []byte("original data")
		sigHex     string
		ok         bool
	)
	if semA, err = NewSem(); err != nil {
		t.Fatalf("生成种子失败：%v", err)
	}

	if err = MustGeneratePriKey(semA); err != nil {
		t.Fatalf("%v", err)
	}

	sm2A = New(semA)

	if sigHex, err = sm2A.Sign(data); err != nil {
		t.Fatalf("签名失败：%v", err)
	}

	semB, err = NewSem(PriKey(semA.GetPriKey()))
	if err != nil {
		t.Fatalf("生成解密种子失败：%v", err)
	}

	sm2B = New(semB)

	if ok, err = sm2B.Verify([]byte("tampered data"), sigHex); err != nil {
		t.Fatalf("验证失败：%v", err)
	}
	if ok {
		t.Fatal("篡改数据的验证应失败")
	}
	t.Logf("篡改数据正确被拒绝")
}
