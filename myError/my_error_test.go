package myError

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type (
	MyError1 struct{ MyError }
	MyError2 struct{ MyError }
)

var (
	MyErr1 MyError1
	MyErr2 MyError2
)

func (*MyError1) New(msg string) IMyError { return &MyError1{MyError{Msg: msg}} }

func (*MyError1) Wrap(err error) IMyError {
	return &MyError1{MyError{Msg: fmt.Errorf("%w", err).Error()}}
}

func (*MyError1) Panic() IMyError {
	return &MyError1{MyError{Msg: "Some error occurred"}}
}

func (my *MyError1) Error() string { return my.Msg }

func (my *MyError1) Is(target error) bool { return reflect.DeepEqual(target, &MyError1{}) }

func (*MyError2) New(msg string) IMyError { return &MyError2{MyError{Msg: msg}} }

func (*MyError2) Wrap(err error) IMyError {
	return &MyError2{MyError{Msg: fmt.Errorf("%w", err).Error()}}
}

func (my *MyError2) Panic() IMyError {
	return &MyError2{MyError{Msg: "Some error occurred"}}
}

func (my *MyError2) Error() string { return my.Msg }

func (my *MyError2) Is(target error) bool { return reflect.DeepEqual(target, &MyError2{}) }

func Test1(t *testing.T) {
	t.Run("自定义错误", func(t *testing.T) {
		err1 := MyErr1.New("Some error occurred")
		err2 := MyErr2.New("Some error occurred2")

		// 使用 errors.Is 来判断错误是否是 ErrMyError
		if errors.Is(err1, &MyError1{}) {
			t.Logf("Is OK1: %s\n", err1)
		} else {
			t.Error("Is NO1")
		}

		if errors.Is(err2, &MyError2{}) {
			t.Logf("Is OK2: %s\n", err2)
		} else {
			t.Error("Is NO2")
		}

		var (
			as1 *MyError1
			as2 *MyError2
		)

		if errors.As(err1, &as1) {
			t.Logf("As OK1: %s\n", err1)
		} else {
			t.Error("As NO1")
		}

		if errors.As(err2, &as2) {
			t.Logf("As OK2: %s\n", err2)
		} else {
			t.Error("As NO2")
		}
	})
}
