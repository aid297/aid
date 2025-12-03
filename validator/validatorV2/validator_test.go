package validatorV2

import (
	"testing"

	"github.com/aid297/aid/ptr"
	"github.com/aid297/aid/regexp"
)

type (
	Target struct {
		Name *string `json:"name" v-name:"名称" v-type:"string" v-rule:"required;max:10;in:abc,123;"`
		Age  *int32  `json:"age" v-name:"年龄" v-type:"int" v-rule:"required;no-zero;min:18;max:60;in:18,25,30,35,40;"`
	}
)

func Test1(t *testing.T) {
	t1 := &Target{
		Name: ptr.New(""),
		Age:  ptr.New(int32(119)),
	}
	v := APP.Validator.New(t1)

	for idx := range v.Errors {
		t.Logf("错误%d：%#v\n", idx+1, v.Errors[idx])
	}
}

func Test2(t *testing.T) {
	r := "required;zero;min:3;max:10;in:abc,123;"
	t.Logf("%v", regexp.APP.Regexp.New(`in:([^,]+)(?:,([^,]+))*`, regexp.TargetString(r)).MatchAll())
}
