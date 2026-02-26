package validatorV3

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

type (
	UserRequest struct {
		Firstname string `v-rule:"(required)(min>10)" v-name:"姓"`
	}
)

func Test1(t *testing.T) {
	ctx := &gin.Context{}
	form, checker := WithGin[UserRequest](ctx, func(form any) (err error) {
		// 这里是一个示例的自定义验证函数，可以根据实际需求进行修改
		// 例如，检查 Firstname 是否等于 "John"
		if userForm, ok := form.(*UserRequest); ok {
			if userForm.Firstname != "John" {
				err = errors.New("firstname must be 'John'")
			}
		} else {
			err = errors.New("invalid form type")
		}
		return
	})

	t.Logf("验证是否通过：%v\n", checker.OK())

	for _, wrong := range checker.Wrongs() {
		t.Logf("%v\n", wrong)
	}

	t.Logf("如果验证通过：%v\n", form)
}
