package validations_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/aid297/aid/v2/validations"
)

type (
	UserRequest struct {
		Firstname string `v-rule:"(required)(min>10)" v-name:"姓"`
		Lastname  string `v-rule:"in:,三,四" v-name:"名"`
	}
)

func Test1(t *testing.T) {
	ctx := &gin.Context{}
	form, checker := validations.WithGin[UserRequest](ctx, func(form any) (err error) {
		// 这里是一个示例的自定义验证函数，可以根据实际需求进行修改
		// 例如，检查 Firstname 是否等于 "John"
		if userForm, ok := form.(*UserRequest); ok {
			if userForm.Firstname != "John" {
				err = errors.New(`firstname must be "John"`)
			}
		} else {
			err = errors.New("invalid form type")
		}
		return
	})

	t.Logf("验证是否通过：%v\n", checker.OK())

	for _, wrong := range checker.Errors() {
		t.Logf("%v\n", wrong)
	}

	t.Logf("如果验证通过：%v\n", form)
}

func Test2(t *testing.T) {
	ur := &UserRequest{
		Firstname: "张",
		Lastname:  "",
	}

	valid := validations.Once()
	checker := valid.Checker(ur)
	checker.Validate()

	t.Logf("验证是否通过：%v\n", checker.OK())
	for _, wrong := range checker.Errors() {
		t.Logf("%v\n", wrong)
	}
}

func Test3(t *testing.T) {
	type T struct {
		Time1 string  `v-rule:"(!)(str-timers>3s)(str-timers<10m)" v-name:"时间1"`
		Time2 *string `v-rule:"(!)(str-timers>=10m)(str-timers<=1h)" v-name:"时间2"`
	}

	t1 := &T{Time1: "", Time2: nil}
	if errs := validations.
		Once().
		DefaultErrorSplitChar("\n").
		Checker(t1).
		Validate().
		Errors(); errs != nil {
		t.Errorf("验证不通过：%v", errs)
	}

	t.Logf("完成")
}
