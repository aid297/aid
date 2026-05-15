package ed25519

import (
	"bytes"
	"testing"
)

func TestNewSemGeneratesKeyPair(t *testing.T) {
	sem, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	pubB64, err := sem.GetPubKeyBase64()
	if err != nil || pubB64 == "" {
		t.Fatalf("公钥长度为空：base64，err=%v", err)
	}
	priB64, err := sem.GetPriKeyBase64()
	if err != nil || priB64 == "" {
		t.Fatalf("私钥长度为空：base64，err=%v", err)
	}
	t.Logf("公钥：%s", pubB64)
	t.Logf("私钥：%s", priB64)
}

func TestSemRoundTripKeys(t *testing.T) {
	sem0, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	pubBytes, err := sem0.GetPubKeyBytes()
	if err != nil {
		t.Fatalf("获取公钥(bytes)失败：%v", err)
	}
	priBytes, err := sem0.GetPriKeyBytes()
	if err != nil {
		t.Fatalf("获取私钥(bytes)失败：%v", err)
	}

	semPub, err := NewSem(PubKeyBytes(pubBytes))
	if err != nil {
		t.Fatalf("用公钥字节创建种子失败：%v", err)
	}
	if semPub.GetPriKey() != nil {
		t.Fatal("仅公钥种子不应持有私钥")
	}

	semPri, err := NewSem(PriKeyBytes(priBytes))
	if err != nil {
		t.Fatalf("用私钥字节创建种子失败：%v", err)
	}
	if semPri.GetPubKey() == nil {
		t.Fatal("仅私钥种子应能派生公钥")
	}
}

func TestSignVerifyRoundTrip(t *testing.T) {
	sem, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	asymm := New(sem)
	data := []byte("hello, Ed25519 数字签名测试")

	sig, err := asymm.Sign(data)
	if err != nil {
		t.Fatalf("签名失败：%v", err)
	}
	if sig == "" {
		t.Fatal("签名为空")
	}
	t.Logf("签名(base64)：%s", sig)

	ok, err := asymm.Verify(data, sig)
	if err != nil {
		t.Fatalf("验签失败：%v", err)
	}
	if !ok {
		t.Fatal("相同密钥与数据下验签应成功")
	}
	t.Logf("验签成功")
}

func TestVerifyWrongKey(t *testing.T) {
	semSign, err := NewSem()
	if err != nil {
		t.Fatalf("生成签名种子失败：%v", err)
	}
	semOther, err := NewSem()
	if err != nil {
		t.Fatalf("生成另一组密钥种子失败：%v", err)
	}

	data := []byte("payload")
	sig, err := New(semSign).Sign(data)
	if err != nil {
		t.Fatalf("签名失败：%v", err)
	}
	t.Logf("签名(base64)：%s", sig)

	ok, err := New(semOther).Verify(data, sig)
	if err != nil {
		t.Fatalf("验签失败：%v", err)
	}
	if ok {
		t.Fatal("异钥验签应失败")
	}
	t.Logf("异钥验签正确被拒绝")
}

func TestVerifyTamperedData(t *testing.T) {
	sem, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	asymm := New(sem)

	data := []byte("original")
	sig, err := asymm.Sign(data)
	if err != nil {
		t.Fatalf("签名失败：%v", err)
	}

	ok, err := asymm.Verify([]byte("tampered"), sig)
	if err != nil {
		t.Fatalf("验签失败：%v", err)
	}
	if ok {
		t.Fatal("篡改数据的验签应失败")
	}
	t.Logf("篡改数据正确被拒绝")
}

func TestVerifyInvalidBase64(t *testing.T) {
	sem, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	_, err = New(sem).Verify([]byte("x"), "not-valid-base64!!!")
	if err == nil {
		t.Fatal("非法 Base64 应返回解码错误")
	}
}

func TestVerifyBadSignatureLength(t *testing.T) {
	sem, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	// 解码后长度非 64 字节时 crypto/ed25519.Verify 返回 false，不返回 error
	ok, err := New(sem).Verify([]byte("msg"), "AAAA")
	if err != nil {
		t.Fatalf("验签不应因长度返回 err：%v", err)
	}
	if ok {
		t.Fatal("签名长度非法时验签应为 false")
	}
	t.Logf("过短签名正确被拒绝")
}

func TestEncryptDecryptUnsupported(t *testing.T) {
	sem, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	asymm := New(sem)

	if _, err := asymm.Encrypt([]byte("x")); err == nil {
		t.Fatal("Ed25519 不支持加密，应返回错误")
	}
	if _, err := asymm.Decrypt(""); err == nil {
		t.Fatal("Ed25519 不支持解密，应返回错误")
	}
}

func TestPubPriOnlySignVerify(t *testing.T) {
	semFull, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	pubBytes, _ := semFull.GetPubKeyBytes()
	priBytes, _ := semFull.GetPriKeyBytes()

	semPri, err := NewSem(PriKeyBytes(priBytes))
	if err != nil {
		t.Fatalf("用私钥字节创建种子失败：%v", err)
	}
	semPub, err := NewSem(PubKeyBytes(pubBytes))
	if err != nil {
		t.Fatalf("用公钥字节创建种子失败：%v", err)
	}

	data := []byte("split signer / verifier")
	sig, err := New(semPri).Sign(data)
	if err != nil {
		t.Fatalf("签名失败：%v", err)
	}
	t.Logf("签名(base64)：%s", sig)

	ok, err := New(semPub).Verify(data, sig)
	if err != nil {
		t.Fatalf("验签失败：%v", err)
	}
	if !ok {
		t.Fatal("仅公钥验签应接受仅私钥侧签名")
	}
	t.Logf("分角色验签成功")
}

func TestMultipleSignaturesDeterministic(t *testing.T) {
	sem, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	data := []byte("same message twice")
	a := New(sem)
	sig1, err := a.Sign(data)
	if err != nil {
		t.Fatalf("第一次签名失败：%v", err)
	}
	sig2, err := a.Sign(data)
	if err != nil {
		t.Fatalf("第二次签名失败：%v", err)
	}
	if sig1 != sig2 {
		t.Fatal("Ed25519 对同一消息签名应确定性相同")
	}
	ok, err := a.Verify(data, sig1)
	if err != nil || !ok {
		t.Fatalf("验签失败：ok=%v err=%v", ok, err)
	}
	t.Logf("确定性签名验签成功")
}

func TestMustGeneratePriKey(t *testing.T) {
	sem := &Ed25519Sem{}
	if err := MustGeneratePriKey(sem); err != nil {
		t.Fatalf("MustGeneratePriKey 失败：%v", err)
	}
	if sem.GetPriKey() == nil || sem.GetPubKey() == nil {
		t.Fatal("MustGeneratePriKey 后公私钥均不应为空")
	}
	t.Logf("MustGeneratePriKey 成功")
}

func TestSetPubKeyBytesInvalid(t *testing.T) {
	sem := &Ed25519Sem{}
	err := sem.SetPubKeyBytes([]byte("not-a-key"))
	if err == nil {
		t.Fatal("非法公钥字节应返回错误")
	}
}

func TestSetPriKeyBytesInvalid(t *testing.T) {
	sem := &Ed25519Sem{}
	err := sem.SetPriKeyBytes([]byte("not-a-key"))
	if err == nil {
		t.Fatal("非法私钥字节应返回错误")
	}
}

func TestGetPubKeyBytesWhenMissing(t *testing.T) {
	sem := &Ed25519Sem{}
	_, err := sem.GetPubKeyBytes()
	if err == nil {
		t.Fatal("未设置密钥时应无法获取公钥 bytes")
	}
}

func TestGetPriKeyBytesWhenMissing(t *testing.T) {
	sem := &Ed25519Sem{}
	_, err := sem.GetPriKeyBytes()
	if err == nil {
		t.Fatal("未设置私钥时应无法获取私钥 bytes")
	}
}

func TestNewSemWithBothKeysFromBytes(t *testing.T) {
	sem0, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	pubBytes, _ := sem0.GetPubKeyBytes()
	priBytes, _ := sem0.GetPriKeyBytes()

	sem, err := NewSem(PriKeyBytes(priBytes), PubKeyBytes(pubBytes))
	if err != nil {
		t.Fatalf("同时设置公私钥创建种子失败：%v", err)
	}
	pb, err := sem.GetPubKeyBytes()
	if err != nil {
		t.Fatalf("获取公钥(bytes)失败：%v", err)
	}
	if !bytes.Equal(pb, pubBytes) {
		t.Fatal("公钥 bytes 与原始不一致")
	}
	t.Logf("公私钥同时导入后公钥一致")
}

func TestPubKeyBase64RoundTrip(t *testing.T) {
	sem0, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	pubB64, err := sem0.GetPubKeyBase64()
	if err != nil {
		t.Fatalf("获取公钥 base64 失败：%v", err)
	}

	semPub, err := NewSem(PubKeyBase64(pubB64))
	if err != nil {
		t.Fatalf("用公钥 Base64 创建种子失败：%v", err)
	}
	got, err := semPub.GetPubKeyBase64()
	if err != nil {
		t.Fatalf("回读公钥 base64 失败：%v", err)
	}
	if got != pubB64 {
		t.Fatalf("公钥 Base64 不一致：期望 %s 实际 %s", pubB64, got)
	}
	t.Logf("公钥 Base64 往返一致")
}

func TestPriKeyBase64RoundTrip(t *testing.T) {
	sem0, err := NewSem()
	if err != nil {
		t.Fatalf("生成种子(Ed25519)失败：%v", err)
	}
	priB64, err := sem0.GetPriKeyBase64()
	if err != nil {
		t.Fatalf("获取私钥 base64 失败：%v", err)
	}

	semPri, err := NewSem(PriKeyBase64(priB64))
	if err != nil {
		t.Fatalf("用私钥 Base64 创建种子失败：%v", err)
	}
	got, err := semPri.GetPriKeyBase64()
	if err != nil {
		t.Fatalf("回读私钥 base64 失败：%v", err)
	}
	if got != priB64 {
		t.Fatalf("私钥 Base64 不一致")
	}
	t.Logf("私钥 Base64 往返一致")
}
