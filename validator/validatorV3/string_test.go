package validatorV3

import (
	"testing"
)

func Test1(t *testing.T) {
	type UserRequest struct {
		Firstname string `json:"firstname" v-rule:"(min>0)" v-name:"firstname"`
	}

	ur := UserRequest{Firstname: ""}
	checker := APP.Validator.Once().Checker(ur)
	checker.Validate()
	if !checker.OK() {
		for _, wrong := range checker.Wrongs() {
			t.Errorf("%v\n", wrong)
		}
	}
}
