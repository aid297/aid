package validatorV2

import "testing"

type (
	T1 struct {
		Name string `json:"name" v-rule:"required" v-name:"名称"`
		Age  int    `json:"age" v-rule:"min=18;max=60" v-name:"年龄"`
	}
)

func Test1(t *testing.T) {
	t1 := T1{
		Name: "Alice",
		Age:  30,
	}
	od, err := NewValidator(t1)
	t.Logf("%v, %v", od, err)
}
