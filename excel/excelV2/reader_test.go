package excelV2

import (
	"testing"
)

func TestToStruct_ExcelTag(t *testing.T) {
	r := (&Reader{}).New()
	// 原始数据：编号, 姓名, 邮箱
	r.originalData = [][]string{
		{"1", "Alice", "alice@example.com"},
		{"2", "Bob", "bob@example.com"},
	}

	titles := []string{"编号", "姓名", "邮箱"}

	type User struct {
		ID    int    `excel:"编号"`
		Name  string `excel:"姓名"`
		Email string `excel:"邮箱"`
	}

	var out []User
	r.ToStruct(titles, &out)
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}

	if len(out) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(out))
	}

	if out[0].ID != 1 || out[0].Name != "Alice" || out[0].Email != "alice@example.com" {
		t.Fatalf("row0 mismatch: %#v", out[0])
	}
}

func TestToStruct_JSONTagAndFieldName(t *testing.T) {
	r := (&Reader{}).New()
	r.originalData = [][]string{
		{"100", "Charlie", "true"},
	}

	titles := []string{"id", "name", "active"}

	type User2 struct {
		ID     int `json:"id"`
		Name   string
		Active bool `json:"active"`
	}

	var out []User2
	r.ToStruct(titles, &out)
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 row, got %d", len(out))
	}

	if out[0].ID != 100 || out[0].Name != "Charlie" || out[0].Active != true {
		t.Fatalf("mismatch: %#v", out[0])
	}
}

func TestToStruct_PointerElementsAndIgnore(t *testing.T) {
	r := (&Reader{}).New()
	r.originalData = [][]string{
		{"10", "Daisy", "p@ex.com", "skipme"},
	}

	titles := []string{"id", "name", "email", "ignore"}

	type User3 struct {
		ID    int     `json:"id"`
		Name  *string `json:"name"`
		Email *string `excel:"邮箱" json:"email"`
		Skip  string  `excel:"-"`
	}

	var out []*User3
	r.ToStruct(titles, &out)
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 row, got %d", len(out))
	}

	if out[0].ID != 10 {
		t.Fatalf("id mismatch: %d", out[0].ID)
	}

	if out[0].Name == nil || *out[0].Name != "Daisy" {
		t.Fatalf("name mismatch: %#v", out[0].Name)
	}

	if out[0].Email == nil || *out[0].Email != "p@ex.com" {
		t.Fatalf("email mismatch: %#v", out[0].Email)
	}

	// Skip should remain zero value
	if out[0].Skip != "" {
		t.Fatalf("expected Skip to be empty, got %q", out[0].Skip)
	}
}

func TestToStruct_ConvertError(t *testing.T) {
	r := (&Reader{}).New()
	r.originalData = [][]string{{"notanint", "X"}}
	titles := []string{"id", "name"}
	type T struct {
		ID int `json:"id"`
	}
	var out []T
	r.ToStruct(titles, &out)
	if r.Error == nil {
		t.Fatalf("expected conversion error, got nil")
	}
}

func TestToStruct_UnsupportedType(t *testing.T) {
	r := (&Reader{}).New()
	r.originalData = [][]string{{"data"}}
	titles := []string{"col"}
	type Bad struct {
		M map[string]string `json:"col"`
	}
	var out []Bad
	r.ToStruct(titles, &out)
	if r.Error == nil {
		t.Fatalf("expected unsupported type error, got nil")
	}
}

func TestToStruct_EmptyStringToZeroValue(t *testing.T) {
	r := (&Reader{}).New()
	r.originalData = [][]string{{"", ""}}
	titles := []string{"id", "name"}
	type S struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var out []S
	r.ToStruct(titles, &out)
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 row, got %d", len(out))
	}
	if out[0].ID != 0 || out[0].Name != "" {
		t.Fatalf("expected zero values, got %#v", out[0])
	}
}

func TestToStruct_FieldNameCaseInsensitive(t *testing.T) {
	r := (&Reader{}).New()
	r.originalData = [][]string{{"42"}}
	titles := []string{"Id"}
	type S struct{ ID int }
	var out []S
	r.ToStruct(titles, &out)
	if r.Error != nil {
		t.Fatalf("unexpected error: %v", r.Error)
	}
	if out[0].ID != 42 {
		t.Fatalf("expected 42, got %d", out[0].ID)
	}
}
