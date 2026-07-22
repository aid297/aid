package validations_test

import (
	"testing"

	"github.com/aid297/aid/v2/validations"
)

// BinaryUUID 模拟 mysql.BinaryUUID 结构体
type BinaryUUID struct {
	ID   string `v-rule:"(required)" v-name:"ID"`
	Name string `v-rule:"(required)" v-name:"名称"`
}

// 测试：切片元素为指针类型且切片为 nil 时不应 panic
type RequestWithNilPtrSlice struct {
	UUIDs []*BinaryUUID `v-rule:"(required)" v-name:"UUID列表"`
}

// 测试：切片元素为指针类型且切片为空切片时不应 panic
type RequestWithEmptyPtrSlice struct {
	UUIDs []*BinaryUUID `v-rule:"(required)" v-name:"UUID列表"`
}

// 测试：切片元素为指针类型且部分元素为 nil 时不应 panic
type RequestWithPartialNilPtrSlice struct {
	UUIDs []*BinaryUUID `v-rule:"(required)" v-name:"UUID列表"`
}

// 测试：指向切片的 nil 指针不应 panic
type RequestWithNilSlicePtr struct {
	UUIDs *[]BinaryUUID `v-rule:"(required)" v-name:"UUID列表"`
}

// 测试：切片元素为结构体（非指针）的空切片
type RequestWithEmptyStructSlice struct {
	UUIDs []BinaryUUID `v-rule:"(required)" v-name:"UUID列表"`
}

func TestNilPointerSlice_NoPanic(t *testing.T) {
	// nil 切片，元素为指针类型
	req := &RequestWithNilPtrSlice{UUIDs: nil}
	checker := validations.OnceValidator().Checker(req)
	checker.Validate()

	t.Logf("验证是否通过：%v", checker.OK())
	for _, wrong := range checker.Errors() {
		t.Logf("  %v", wrong)
	}
}

func TestEmptyPointerSlice_NoPanic(t *testing.T) {
	// 空切片，元素为指针类型
	req := &RequestWithEmptyPtrSlice{UUIDs: []*BinaryUUID{}}
	checker := validations.OnceValidator().Checker(req)
	checker.Validate()

	t.Logf("验证是否通过：%v", checker.OK())
	for _, wrong := range checker.Errors() {
		t.Logf("  %v", wrong)
	}
}

func TestPartialNilPointerSlice_NoPanic(t *testing.T) {
	// 切片中有 nil 指针元素
	req := &RequestWithPartialNilPtrSlice{
		UUIDs: []*BinaryUUID{
			{ID: "abc", Name: "test"},
			nil, // nil 元素不应导致 panic
			{ID: "", Name: ""},
		},
	}
	checker := validations.OnceValidator().Checker(req)
	checker.Validate()

	t.Logf("验证是否通过：%v", checker.OK())
	for _, wrong := range checker.Errors() {
		t.Logf("  %v", wrong)
	}
}

func TestNilSlicePointer_NoPanic(t *testing.T) {
	// *[]BinaryUUID 为 nil
	req := &RequestWithNilSlicePtr{UUIDs: nil}
	checker := validations.OnceValidator().Checker(req)
	checker.Validate()

	t.Logf("验证是否通过：%v", checker.OK())
	for _, wrong := range checker.Errors() {
		t.Logf("  %v", wrong)
	}
}

func TestEmptyStructSlice_NoPanic(t *testing.T) {
	// 空切片，元素为结构体（非指针）
	req := &RequestWithEmptyStructSlice{UUIDs: []BinaryUUID{}}
	checker := validations.OnceValidator().Checker(req)
	checker.Validate()

	t.Logf("验证是否通过：%v", checker.OK())
	for _, wrong := range checker.Errors() {
		t.Logf("  %v", wrong)
	}
}
