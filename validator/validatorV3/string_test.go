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
	wrongs := checker.Validate().Wrongs()
	for _, wrong := range wrongs {
		t.Errorf("%v\n", wrong)
	}
}
