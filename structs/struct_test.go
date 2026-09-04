package structs_test

import (
	"testing"
	"time"

	"github.com/aid297/aid/v2/structs"
)

type inner struct {
	Name  string
	Age   int
	Slice []int
}

type srcStruct struct {
	Name     string
	Age      int
	Inner    inner
	Ptr      *inner
	PtrPtr   **inner
	Birthday time.Time
	Tags     []string
	Secret   string // dst中没有的字段，用于验证保留
}

type dstStruct struct {
	Name     string
	Age      string  // 类型不一致，不应覆盖
	Inner    inner   // 嵌套结构体，应跳过
	Ptr      *inner  // 嵌套指针结构体，应跳过
	PtrPtr   **inner // 多级指针，应跳过
	Birthday time.Time
	Tags     []string
	Extra    int // src没有的字段，应忽略
}

func TestCover(t *testing.T) {
	src := srcStruct{Name: "src名字", Age: 20, Tags: []string{"a"}, Secret: "秘密"}
	dst := dstStruct{Name: "dst名字", Age: "30", Extra: 99}

	src = structs.NewStruct(src, dst).Cover()

	if src.Name != "dst名字" {
		t.Errorf("同名同类型字段应被覆盖，期望 dst名字，实际 %s", src.Name)
	}
	if src.Age != 20 {
		t.Errorf("同名但类型不一致的字段不应覆盖，期望 20，实际 %d", src.Age)
	}
	if src.Secret != "秘密" {
		t.Errorf("dst中没有的同名字段应保留，实际 %s", src.Secret)
	}
	if src.Tags != nil {
		t.Errorf("dst中同名同类型字段即使为零值也应覆盖，期望 nil，实际 %v", src.Tags)
	}
}

func TestCoverNotModifyInput(t *testing.T) {
	src := srcStruct{Name: "src名字", Tags: []string{"old"}}
	dst := dstStruct{Name: "dst名字", Tags: []string{"new"}}

	structs.NewStruct(src, dst).Cover()

	if src.Name != "src名字" || src.Tags[0] != "old" {
		t.Errorf("入参src不应被修改，实际 %s %v", src.Name, src.Tags)
	}
}

func TestCoverNestedStruct(t *testing.T) {
	src := srcStruct{Inner: inner{Name: "src内部", Age: 10, Slice: []int{1}}}
	dst := dstStruct{Inner: inner{Name: "dst内部", Age: 20, Slice: []int{9, 9}}}

	src = structs.NewStruct(src, dst).Cover()

	// 嵌套结构体字段（非time.Time）直接跳过，src保持原值
	if src.Inner.Name != "src内部" || src.Inner.Age != 10 || len(src.Inner.Slice) != 1 {
		t.Errorf("嵌套结构体字段应直接跳过，期望保持原值，实际 %+v", src.Inner)
	}
}

func TestCoverNestedPtr(t *testing.T) {
	t.Run("非nil时直接跳过", func(t *testing.T) {
		src := srcStruct{Ptr: &inner{Name: "src指针", Age: 1}}
		dst := dstStruct{Ptr: &inner{Name: "dst指针", Age: 2}}

		src = structs.NewStruct(src, dst).Cover()

		if src.Ptr.Name != "src指针" {
			t.Errorf("嵌套指针字段应直接跳过，期望保持原值，实际 %+v", src.Ptr)
		}
	})

	t.Run("src为nil时也跳过", func(t *testing.T) {
		src := srcStruct{}
		dstInner := inner{Name: "dst指针", Age: 2}
		dst := dstStruct{Ptr: &dstInner}

		src = structs.NewStruct(src, dst).Cover()

		if src.Ptr != nil {
			t.Errorf("src为nil的嵌套指针字段也应跳过，实际 %+v", src.Ptr)
		}
	})

	t.Run("dst为nil时也跳过", func(t *testing.T) {
		src := srcStruct{Ptr: &inner{Name: "src指针"}}
		dst := dstStruct{}

		src = structs.NewStruct(src, dst).Cover()

		if src.Ptr == nil || src.Ptr.Name != "src指针" {
			t.Errorf("dst为nil的嵌套指针字段也应跳过，期望保持原值，实际 %+v", src.Ptr)
		}
	})
}

func TestCoverMultiLevelPtr(t *testing.T) {
	srcInner := inner{Name: "src多级", Age: 1}
	srcPtr := &srcInner
	src := srcStruct{PtrPtr: &srcPtr}

	dstInner := inner{Name: "dst多级", Age: 2}
	dstPtr := &dstInner
	dst := dstStruct{PtrPtr: &dstPtr}

	src = structs.NewStruct(src, dst).Cover()

	if src.PtrPtr == dst.PtrPtr || (*src.PtrPtr).Name != "src多级" {
		t.Errorf("多级指针字段应直接跳过，期望保持原值，实际 %v", src.PtrPtr)
	}
}

func TestCoverTime(t *testing.T) {
	now := time.Now()
	src := srcStruct{Birthday: time.Time{}}
	dst := dstStruct{Birthday: now}

	src = structs.NewStruct(src, dst).Cover()

	if !src.Birthday.Equal(now) {
		t.Errorf("time.Time字段应整体覆盖，实际 %v", src.Birthday)
	}
}

func TestCoverSlice(t *testing.T) {
	src := srcStruct{Tags: []string{"old"}}
	dst := dstStruct{Tags: []string{"new1", "new2"}}

	src = structs.NewStruct(src, dst).Cover()

	if len(src.Tags) != 2 || src.Tags[0] != "new1" || src.Tags[1] != "new2" {
		t.Errorf("slice字段应整体覆盖，实际 %v", src.Tags)
	}
}

func TestCoverBySkip(t *testing.T) {
	src := srcStruct{Name: "src名字", Age: 20, Tags: []string{"old"}}
	dst := dstStruct{Name: "dst名字", Age: "30", Tags: []string{"new"}}

	src = structs.NewStruct(src, dst).CoverBySkip("Name", "Tags")

	if src.Name != "src名字" {
		t.Errorf("黑名单中的字段不应被覆盖，期望 src名字，实际 %s", src.Name)
	}
	if src.Tags == nil || src.Tags[0] != "old" {
		t.Errorf("黑名单中的字段不应被覆盖，期望 [old]，实际 %v", src.Tags)
	}

	t.Run("不传黑名单时等同Cover", func(t *testing.T) {
		src := srcStruct{Name: "src名字", Tags: []string{"old"}}
		dst := dstStruct{Name: "dst名字", Tags: []string{"new"}}

		src = structs.NewStruct(src, dst).CoverBySkip()

		if src.Name != "dst名字" || src.Tags[0] != "new" {
			t.Errorf("不传黑名单时应正常覆盖，实际 %s %v", src.Name, src.Tags)
		}
	})
}

func TestCoverByAssign(t *testing.T) {
	src := srcStruct{Name: "src名字", Age: 20, Birthday: time.Time{}, Tags: []string{"old"}}
	dst := dstStruct{Name: "dst名字", Age: "30", Birthday: time.Now(), Tags: []string{"new"}}

	src = structs.NewStruct(src, dst).CoverByAssign("Name")

	if src.Name != "dst名字" {
		t.Errorf("白名单中的字段应被覆盖，期望 dst名字，实际 %s", src.Name)
	}
	if src.Tags[0] != "old" {
		t.Errorf("不在白名单中的字段不应被覆盖，期望 [old]，实际 %v", src.Tags)
	}
	if !src.Birthday.IsZero() {
		t.Errorf("不在白名单中的字段不应被覆盖，期望零值，实际 %v", src.Birthday)
	}

	t.Run("空白名单时什么都不覆盖", func(t *testing.T) {
		src := srcStruct{Name: "src名字"}
		dst := dstStruct{Name: "dst名字"}

		src = structs.NewStruct(src, dst).CoverByAssign()

		if src.Name != "src名字" {
			t.Errorf("空白名单时不应覆盖任何字段，实际 %s", src.Name)
		}
	})
}

func TestNewStructPanic(t *testing.T) {
	t.Run("src传指针时panic", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("src传指针时应panic")
			}
		}()

		structs.NewStruct(&srcStruct{}, dstStruct{})
	})

	t.Run("dst传指针时panic", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("dst传指针时应panic")
			}
		}()

		structs.NewStruct(srcStruct{}, &dstStruct{})
	})

	t.Run("src非结构体时panic", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("src非结构体时应panic")
			}
		}()

		structs.NewStruct(123, dstStruct{})
	})

	t.Run("dst非结构体时panic", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("dst非结构体时应panic")
			}
		}()

		structs.NewStruct(srcStruct{}, 123)
	})
}
