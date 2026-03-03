package symmetric

import "testing"

func TestAES1(t *testing.T) {
	sail := "tjp5OPIU1ETF5s33fsLWdA=="

	aes := NewAES(sail)
	aesEncrypt := aes.NewEncrypt().GetEncrypt()
	aesDecrypt := aes.NewDecrypt(aesEncrypt.GetOpenKey()).GetDecrypt()

	t.Logf("加密后：%v", aesEncrypt.GetOpenKey())
	t.Logf("还原后：%s", aesDecrypt.GetAesKey())
}
