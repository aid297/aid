package validatorV3

import (
	"testing"

	"github.com/gin-gonic/gin"
)

type (
	UserRequest struct {
		Firstname string `v-rule:"ex:some-ex-check-fn" v-name:"姓"`
	}
)

func Test1(t *testing.T) {
	ctx := &gin.Context{}
	form, checker := WithGin[UserRequest](ctx, func(form any) (err error) {
		// 这里是一个示例的自定义验证函数，可以根据实际需求进行修改
		return
	})
	t.Logf("验证是否通过：%v\n", checker.OK())
	for _, wrong := range checker.Wrongs() {
		t.Logf("%v\n", wrong)
	}
	t.Logf("如果验证通过：%v\n", form)
}
