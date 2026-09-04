package structs_test

import (
	"testing"
	"time"

	"github.com/aid297/aid/v2/structs"
)

type ptrSrcStruct struct {
	Tag      *string
	TagNil   *string
	Num      *int
	NumPtr   **int
	In       *inner
	InNil    *inner
	Name     string
	Age      int
	Birthday *time.Time
}

type ptrDstStruct struct {
	Tag      string
	TagNil   string
	Num      int
	NumPtr   int
	In       inner
	InNil    inner
	Name     *string
	Age      *int
	Birthday time.Time
}

func TestCoverByDeref(t *testing.T) {
	dstName := "dst名字"
	dstAge := 99

	src := ptrSrcStruct{
		Tag:  ptr("srcTag"),
		Num:  ptr(1),
		Name: "src名字",
		Age:  20,
	}
	dst := ptrDstStruct{
		Tag:      "dstTag",
		TagNil:   "dstTagNil",
		Num:      100,
		NumPtr:   200,
		In:       inner{Name: "dst内部", Age: 30},
		InNil:    inner{Name: "dst内部nil"},
		Name:     &dstName,
		Age:      &dstAge,
		Birthday: time.Now(),
	}

	src = structs.NewStruct(src, dst).CoverBySkip().Value()

	t.Run("src指针_dst值_解引用覆盖", func(t *testing.T) {
		if *src.Tag != "dstTag" {
			t.Errorf("*string vs string 应解引用覆盖，期望 dstTag，实际 %s", *src.Tag)
		}
		if *src.Num != 100 {
			t.Errorf("*int vs int 应解引用覆盖，期望 100，实际 %d", *src.Num)
		}
	})

	t.Run("src为nil指针时新建对象覆盖", func(t *testing.T) {
		if src.TagNil == nil || *src.TagNil != "dstTagNil" {
			t.Errorf("src为nil的指针字段应新建对象覆盖，实际 %v", src.TagNil)
		}
	})

	t.Run("src多级指针_dst值_解引用覆盖", func(t *testing.T) {
		if src.NumPtr == nil || **src.NumPtr != 200 {
			t.Errorf("**int vs int 应解引用覆盖，期望 200，实际 %v", src.NumPtr)
		}
	})

	t.Run("src值_dst指针_解引用覆盖", func(t *testing.T) {
		if src.Name != "dst名字" {
			t.Errorf("string vs *string 应解引用覆盖，期望 dst名字，实际 %s", src.Name)
		}
		if src.Age != 99 {
			t.Errorf("int vs *int 应解引用覆盖，期望 99，实际 %d", src.Age)
		}
	})

	t.Run("src指针嵌套结构体_dst值_解引用值拷贝覆盖", func(t *testing.T) {
		if src.In == nil || src.In.Name != "dst内部" || src.In.Age != 30 {
			t.Errorf("*inner vs inner 应解引用值拷贝覆盖，实际 %+v", src.In)
		}
		if src.InNil == nil || src.InNil.Name != "dst内部nil" {
			t.Errorf("src为nil的嵌套指针字段应新建对象覆盖，实际 %+v", src.InNil)
		}
	})

	t.Run("src指针时间_dst值时间_解引用覆盖", func(t *testing.T) {
		if src.Birthday == nil || !src.Birthday.Equal(dst.Birthday) {
			t.Errorf("*time.Time vs time.Time 应解引用覆盖，实际 %v", src.Birthday)
		}
	})
}

func TestCoverByDerefDstNil(t *testing.T) {
	t.Run("dst为nil指针时用零值覆盖", func(t *testing.T) {
		src := ptrSrcStruct{Name: "src名字", Age: 20}
		dst := ptrDstStruct{} // Name和Age都是nil指针

		src = structs.NewStruct(src, dst).CoverBySkip().Value()

		if src.Name != "" {
			t.Errorf("dst为nil指针时应用零值覆盖，期望空字符串，实际 %s", src.Name)
		}
		if src.Age != 0 {
			t.Errorf("dst为nil指针时应用零值覆盖，期望0，实际 %d", src.Age)
		}
	})
}

func TestCoverByDerefTypeMismatch(t *testing.T) {
	type mismatchSrc struct {
		Num *int
	}
	type mismatchDst struct {
		Num *string // 解引用到底类型不同
	}

	src := mismatchSrc{Num: ptr(1)}
	dst := mismatchDst{Num: ptr("dst")}

	src = structs.NewStruct(src, dst).CoverBySkip().Value()

	if *src.Num != 1 {
		t.Errorf("解引用到底类型不同时不应覆盖，期望 1，实际 %d", *src.Num)
	}
}

func ptr[T any](v T) *T { return &v }
