package validatorV3

import (
	"fmt"
	"testing"
)

type (
	UserRequest struct {
		Firstname string `v-rule:"ex:some-ex-check-fn" v-name:"姓"`
	}
)

func Test1(t *testing.T) {
	ur := &UserRequest{
		Firstname: "张三",
	}

	validator := APP.Validator.Once()

	checker := validator.Checker(ur)
	checker.Validate(func(form any) (err error) {
		// 这里是一次性自定义验证（模拟去数据库中进行验证）
		if form.(*UserRequest).Firstname != "王五" {
			err = fmt.Errorf("名字必须是：王五")
		}

		return
	})
	t.Logf("验证是否通过：%v\n", checker.OK())
	for _, wrong := range checker.Wrongs() {
		t.Logf("%v\n", wrong)
	}

	// 验证是否通过：false
	// 名字必须是：王五
}
