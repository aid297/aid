package validatorV3

import (
	"testing"
	`time`
)

func Test1(t *testing.T) {
	type UserRequest struct {
		Birthday time.Time `v-rule:"required" v-name:"生日"`
	}

	ur := &UserRequest{Birthday: time.Time{}}

	checker := APP.Validator.Once().Checker(ur).Validate()
	t.Logf("%v\n", checker.OK())
	for _, wrong := range checker.Wrongs() {
		t.Logf("%v\n", wrong)
	}
}
