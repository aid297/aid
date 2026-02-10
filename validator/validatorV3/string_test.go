package validatorV3

import (
	"testing"

	`github.com/aid297/aid/ptr`
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

func Test2(t *testing.T) {
	type UserRequest struct {
		Firstname string  `json:"firstname" v-rule:"(required)" v-name:"firstname"`
		Lastname  *string `json:"lastname" v-rule:"(required)" v-name:"lastname"`
	}

	ur := UserRequest{Firstname: "123", Lastname: ptr.New("")}
	checker := APP.Validator.Once().Checker(ur)
	checker.Validate()
	if !checker.OK() {
		for _, wrong := range checker.Wrongs() {
			t.Errorf("%v\n", wrong)
		}
	}
}
