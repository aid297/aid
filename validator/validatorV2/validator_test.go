package validatorV2

import (
	"testing"

	"github.com/aid297/aid/ptr"
	"github.com/aid297/aid/regexp"
)

type (
	Target struct {
		Name *string `json:"name" v-name:"名称" v-type:"string" v-rule:"required;not-zero;min:3;max:10;in:abc,123;"`
		Age  *int32  `json:"age" v-name:"年龄" v-type:"int" v-rule:"required;no-zero;min:18;max:60;in:18,25,30,35,40;"`
	}
)

func Test1(t *testing.T) {
	t1 := &Target{
		Name: ptr.New("Alice"),
		Age:  ptr.New(int32(119)),
	}
	v := APP.Validator.New(t1)
	t.Log(v.Errors)
}

func Test2(t *testing.T) {
	r := "required;zero;min:3;max:10;in:abc,123;"
	t.Logf("%v", regexp.APP.Regexp.New(`in:([^,]+)(?:,([^,]+))*`, regexp.TargetString(r)).MatchAll())
}
