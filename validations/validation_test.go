package validations_test

import (
	"testing"

	"github.com/aid297/aid/v2/validations"
)

func TestCustomType(t *testing.T) {
	type A string
	var (
		Custom1 A = "CUSTOM-A"
		Custom2 A = "CUSTOM-B"
	)
	type S struct {
		AField A `json:"a_field" v-rule:"(required)(min>0)(in==CUSTOM-A,CUSTOM-B)" v-name:"a_field"`
	}

	_ = Custom2
	var s = S{AField: Custom1}

	checker := validations.Once().Checker(s).Validate()
	if checker.Invalid() {
		t.Errorf("验证不通过：%v", checker.Error())
	}

	t.Logf("OK")
}
