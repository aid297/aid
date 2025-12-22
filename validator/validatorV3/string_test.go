package validatorV3

import (
	"testing"

	"github.com/aid297/aid/ptr"
)

type (
	StringTest struct {
		Name1 string  `v-rule:"required;min>2;max<=10;not-in:张三" v-name:"姓名"`
		Name2 *string `v-rule:"required" v-name:"姓名1"`
		Name3 *string `v-rule:"required;min>0;in:王五,赵六"`
	}

	IntTest struct {
		Age1 int   `v-rule:"required"`
		Age2 *int  `v-rule:"required"`
		Age3 *int8 `v-rule:"required;min>=5"`
	}
)

func Test1(t *testing.T) {
	st := StringTest{Name1: "张三", Name2: nil, Name3: ptr.New("")}
	t.Logf("%v", APP.Validator.New(st).Validate().Wrongs())
}

func Test2(t *testing.T) {
	it := &IntTest{0, nil, ptr.New(int8(5))}

	t.Logf("%v", APP.Validator.New(it).Validate().Wrongs())
}

func Test3(t *testing.T) {
	it := &IntTest{0, ptr.New(1), ptr.New(int8(5))}

	t.Logf("%v", APP.Validator.New(it).Validate(func(data any) error {
		data.(*IntTest).Age2 = ptr.New(111)
		return nil
	}).Wrongs())

	t.Logf("%v", *it.Age2)
}
