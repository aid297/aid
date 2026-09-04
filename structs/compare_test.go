package structs_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/aid297/aid/v2/structs"
)

type cmpSrc struct {
	Name     string     `name:"名称"`
	Age      *int       `name:"年龄"`
	Tag      *string    // 无name tag，用字段名
	Count    int        // 相同值，不应记录
	Inner    inner      // 结构体，应跳过
	InnerPtr *inner     // 结构体指针，应跳过
	Birthday *time.Time // time.Time特例，参与比较
}

type cmpDst struct {
	Name     string
	Age      int
	Tag      string
	Count    int
	Inner    inner
	InnerPtr *inner
	Birthday time.Time
}

func TestCompareBySkip(t *testing.T) {
	src := cmpSrc{Name: "src名字", Count: 1, Inner: inner{Name: "src内部"}}
	dst := cmpDst{
		Name:     "dst名字",
		Age:      20,
		Tag:      "dstTag",
		Count:    1,
		Inner:    inner{Name: "dst内部"},
		InnerPtr: &inner{Name: "dst指针"},
	}

	diffs := structs.NewStruct(src, dst).CompareBySkip()

	// Name有差异；Age为nil指针取零值0与20不同；Tag为nil指针取零值""与"dstTag"不同；Count相同不记录；Inner/InnerPtr跳过
	var expect = []string{"名称:src名字 -> dst名字", "年龄:0 -> 20", "Tag: -> dstTag"}
	if !slices.Equal(diffs, expect) {
		t.Errorf("差异记录不符合预期，期望 %v，实际 %v", expect, diffs)
	}
}

func TestCompareBySkipBlacklist(t *testing.T) {
	src := cmpSrc{Name: "src名字"}
	dst := cmpDst{Name: "dst名字", Age: 20}

	diffs := structs.NewStruct(src, dst).CompareBySkip("Name")

	// 黑名单中的Name被跳过，只记录Age（src为nil指针取零值0）
	var expect = []string{"年龄:0 -> 20"}
	if !slices.Equal(diffs, expect) {
		t.Errorf("黑名单中的字段不应比较，期望 %v，实际 %v", expect, diffs)
	}
}

func TestCompareByAssign(t *testing.T) {
	src := cmpSrc{Name: "src名字"}
	dst := cmpDst{Name: "dst名字", Age: 20}

	diffs := structs.NewStruct(src, dst).CompareByAssign("Name")

	// 白名单外的Age不比较，只记录Name
	var expect = []string{"名称:src名字 -> dst名字"}
	if !slices.Equal(diffs, expect) {
		t.Errorf("白名单外的字段不应比较，期望 %v，实际 %v", expect, diffs)
	}
}

func TestCompareAllEqual(t *testing.T) {
	t.Run("值全部相同返回nil", func(t *testing.T) {
		src := cmpSrc{Name: "相同", Count: 1}
		dst := cmpDst{Name: "相同", Count: 1}

		if diffs := structs.NewStruct(src, dst).CompareBySkip(); diffs != nil {
			t.Errorf("全部相同时应无差异记录，实际 %v", diffs)
		}
	})

	t.Run("nil指针与零值相同不记录", func(t *testing.T) {
		src := cmpSrc{}                // Age/Tag均为nil指针
		dst := cmpDst{Age: 0, Tag: ""} // 零值

		if diffs := structs.NewStruct(src, dst).CompareBySkip(); diffs != nil {
			t.Errorf("nil指针取零值后相同应无差异记录，实际 %v", diffs)
		}
	})

	t.Run("两侧都是nil指针不记录", func(t *testing.T) {
		src := cmpSrc{}
		var dst cmpDst

		if diffs := structs.NewStruct(src, dst).CompareBySkip(); diffs != nil {
			t.Errorf("两侧nil指针应无差异记录，实际 %v", diffs)
		}
	})
}

func TestComparePointerValue(t *testing.T) {
	src := cmpSrc{Age: ptr(10), Tag: ptr("srcTag")}
	dst := cmpDst{Age: 10, Tag: "dstTag"}

	diffs := structs.NewStruct(src, dst).CompareBySkip()

	// Age指针解引用后10与10相同不记录；Tag解引用后不同记录src与dst的值
	var expect = []string{"Tag:srcTag -> dstTag"}
	if !slices.Equal(diffs, expect) {
		t.Errorf("指针解引用比较不符合预期，期望 %v，实际 %v", expect, diffs)
	}
}

func TestCompareTime(t *testing.T) {
	now := time.Now()
	src := cmpSrc{Birthday: &now}
	dst := cmpDst{Birthday: now.Add(time.Hour)}

	diffs := structs.NewStruct(src, dst).CompareByAssign("Birthday")

	if len(diffs) != 1 || !strings.HasPrefix(diffs[0], "Birthday:") {
		t.Errorf("time.Time字段应参与比较并记录dst值，实际 %v", diffs)
	}
}

func TestCompareNotModifyInput(t *testing.T) {
	src := cmpSrc{Name: "src名字", Age: ptr(10)}
	dst := cmpDst{Name: "dst名字", Age: 20}

	structs.NewStruct(src, dst).CompareBySkip()

	if src.Name != "src名字" || *src.Age != 10 {
		t.Errorf("比较操作不应修改入参src，实际 %s %d", src.Name, *src.Age)
	}
}

func TestCompareTypeMismatch(t *testing.T) {
	type mismatchSrc struct {
		Num *int
	}
	type mismatchDst struct {
		Num *string // 解引用到底类型不同
	}

	src := mismatchSrc{Num: ptr(1)}
	dst := mismatchDst{Num: ptr("dst")}

	if diffs := structs.NewStruct(src, dst).CompareBySkip(); diffs != nil {
		t.Errorf("解引用到底类型不同时不应比较，实际 %v", diffs)
	}
}
