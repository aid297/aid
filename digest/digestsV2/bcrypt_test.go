package digestsV2_test

import (
	"testing"

	"github.com/aid297/aid/v2/digest/digestsV2"
)

func Test1(t *testing.T) {
	pwd := "123456"
	hashed := string(digestsV2.NewBcrypt(pwd).Hash())

	res := digestsV2.NewBcrypt(pwd).Valid(hashed)
	t.Logf("结果：%v", res)
}
