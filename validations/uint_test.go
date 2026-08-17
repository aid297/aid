package validations_test

import (
	"testing"

	"github.com/aid297/aid/v2/validations"
)

type TestUint8 struct {
	Age uint8 `json:"age" v-rule:"(required)(min>=1)(max<=64)" v-name:"年龄"`
	Test2Uint8
}

type Test2Uint8 struct {
	Name string `json:"name" v-rule:"(required)(min>=2)" v-name:"姓名"`
}

func Test1_uint8(t *testing.T) {
	test := TestUint8{Age: 30, Test2Uint8: Test2Uint8{Name: "张三"}}

	checker := validations.Once().Checker(test).Validate()
	if checker.Invalid() {
		t.Errorf("错误：%v", checker.Error())
	}

	t.Logf("成功")
}
