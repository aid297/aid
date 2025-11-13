package operationV2

import (
	"fmt"
	"testing"
)

type T struct {
	Name string
}

func TestMatch1(t *testing.T) {
	var (
		err error
		t1  T
	)

	t1 = T{Name: "李四"}

	APP.Match.New(
		MatchItem("张三", func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(4, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(nil, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(err != nil, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(T{Name: "李四"}, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(&T{Name: "王五"}, func(val any) { fmt.Printf("执行：%v\n", val) }),
	).SetDefault(func() { fmt.Printf("执行：default") }).Do(t1)
}

func TestMatch2(t *testing.T) {
	var (
		err error
		t2  *T
	)

	t2 = &T{Name: "王五"}

	APP.Match.New(
		MatchItem("张三", func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(4, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(nil, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(err != nil, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(&T{Name: "王五"}, func(val any) { fmt.Printf("执行：%v\n", val) }),
	).SetDefault(func() { fmt.Printf("执行：default") }).Do(t2)
}

func TestMatch3(t *testing.T) {
	var err error

	APP.Match.New(
		MatchItem("张三", func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(4, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(nil, func(val any) { fmt.Printf("执行：%v\n", val) }),
		MatchItem(err != nil, func(val any) { fmt.Printf("执行：%v\n", val) }),
	).SetDefault(func() { fmt.Printf("执行：default") }).Do("")
}
