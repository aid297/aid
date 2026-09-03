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
	Name    string
	Age     int
	Inner   inner
	Ptr     *inner
	PtrPtr  **inner
	Birthday time.Time
	Tags    []string
	Secret  string // 用于验证未导出字段场景的对照（导出但配合嵌入未导出测试）
}

type dstStruct struct {
	Name    string
	Age     string   // 类型不一致，不应覆盖
	Inner   inner    // 嵌套结构体，应递归替换
	Ptr     *inner   // 嵌套指针结构体
	PtrPtr  **inner  // 多级指针
	Birthday time.Time
	Tags    []string
	Extra   int      // src没有的字段，应忽略
}

func TestCopyBasic(t *testing.T) {
	src := &srcStruct{Name: "src名字", Age: 20, Tags: []string{"a"}, Secret: "秘密"}
	dst := dstStruct{Name: "dst名字", Age: "30", Extra: 99}

	structs.Copy(src, dst)

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

func TestCopyNestedStruct(t *testing.T) {
	src := &srcStruct{Inner: inner{Name: "src内部", Age: 10, Slice: []int{1}}}
	dst := dstStruct{Inner: inner{Name: "dst内部", Age: 20, Slice: []int{9, 9}}}

	structs.Copy(src, dst)

	// 不递归：嵌套结构体整体覆盖，src原有的Slice一并被替换
	if src.Inner.Name != "dst内部" || src.Inner.Age != 20 {
		t.Errorf("嵌套结构体应整体覆盖，实际 %+v", src.Inner)
	}
	if len(src.Inner.Slice) != 2 || src.Inner.Slice[0] != 9 || src.Inner.Slice[1] != 9 {
		t.Errorf("嵌套结构体整体覆盖时slice应一并替换，实际 %v", src.Inner.Slice)
	}
}

func TestCopyNestedPtr(t *testing.T) {
	t.Run("非nil时整体覆盖（浅拷贝，指向同一地址）", func(t *testing.T) {
		src := &srcStruct{Ptr: &inner{Name: "src指针", Age: 1}}
		dst := dstStruct{Ptr: &inner{Name: "dst指针", Age: 2}}

		structs.Copy(src, dst)

		if src.Ptr != dst.Ptr {
			t.Errorf("不递归：指针字段应整体覆盖为dst的指针（同一地址），实际 %+v", src.Ptr)
		}
	})

	t.Run("src为nil时整体覆盖", func(t *testing.T) {
		src := &srcStruct{}
		dstInner := inner{Name: "dst指针", Age: 2}
		dst := dstStruct{Ptr: &dstInner}

		structs.Copy(src, dst)

		if src.Ptr == nil || src.Ptr.Name != "dst指针" {
			t.Errorf("src为nil时应整体覆盖，实际 %+v", src.Ptr)
		}
	})

	t.Run("dst为nil时覆盖为nil", func(t *testing.T) {
		src := &srcStruct{Ptr: &inner{Name: "src指针"}}
		dst := dstStruct{}

		structs.Copy(src, dst)

		if src.Ptr != nil {
			t.Errorf("dst为nil时应覆盖为nil，实际 %+v", src.Ptr)
		}
	})
}

func TestCopyMultiLevelPtr(t *testing.T) {
	srcInner := inner{Name: "src多级", Age: 1}
	srcPtr := &srcInner
	src := &srcStruct{PtrPtr: &srcPtr}

	dstInner := inner{Name: "dst多级", Age: 2}
	dstPtr := &dstInner
	dst := dstStruct{PtrPtr: &dstPtr}

	structs.Copy(src, dst)

	if src.PtrPtr != dst.PtrPtr {
		t.Errorf("不递归：多级指针字段应整体覆盖（同一地址），实际 %v", src.PtrPtr)
	}
}

func TestCopyTime(t *testing.T) {
	now := time.Now()
	src := &srcStruct{Birthday: time.Time{}}
	dst := dstStruct{Birthday: now}

	structs.Copy(src, dst)

	if !src.Birthday.Equal(now) {
		t.Errorf("time.Time字段应整体覆盖，实际 %v", src.Birthday)
	}
}

func TestCopySlice(t *testing.T) {
	src := &srcStruct{Tags: []string{"old"}}
	dst := dstStruct{Tags: []string{"new1", "new2"}}

	structs.Copy(src, dst)

	if len(src.Tags) != 2 || src.Tags[0] != "new1" || src.Tags[1] != "new2" {
		t.Errorf("slice字段应整体覆盖，实际 %v", src.Tags)
	}
}

func TestCopyDstPtr(t *testing.T) {
	src := &srcStruct{Name: "src名字"}
	dst := &dstStruct{Name: "dst指针名字"}

	structs.Copy(src, dst)

	if src.Name != "dst指针名字" {
		t.Errorf("dst传指针时应正常覆盖，实际 %s", src.Name)
	}
}

func TestCopyPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("src非指针时应panic")
		}
	}()

	structs.Copy(srcStruct{}, dstStruct{})
}

func TestCopyNoAction(t *testing.T) {
	t.Run("dst为nil指针时不做任何事", func(t *testing.T) {
		src := &srcStruct{Name: "src名字"}
		var dst *dstStruct

		structs.Copy(src, dst)

		if src.Name != "src名字" {
			t.Errorf("dst为nil指针时src应保持不变，实际 %s", src.Name)
		}
	})

	t.Run("dst非结构体时不做任何事", func(t *testing.T) {
		src := &srcStruct{Name: "src名字"}

		structs.Copy(src, 123)

		if src.Name != "src名字" {
			t.Errorf("dst非结构体时src应保持不变，实际 %s", src.Name)
		}
	})
}
