package excelV2

import (
	"testing"
	"time"
)

func TestFromStruct_BasicTypesAndHeader(t *testing.T) {
	w := APP.Writer.New(APP.WriterAttr.SheetName.Set("Sheet1"))

	type User struct {
		ID     int     `excel:"编号"`
		Name   string  `excel:"姓名"`
		Active bool    `excel:"激活"`
		Score  float64 `excel:"得分"`
	}

	data := []User{{ID: 1, Name: "Alice", Active: true, Score: 3.5}, {ID: 2, Name: "Bob", Active: false, Score: 4.0}}
	titles := []string{"编号", "姓名", "激活", "得分"}

	w.FromStruct(&data, titles, 1)
	if w.Error != nil {
		t.Fatalf("unexpected error: %v", w.Error)
	}

	v, _ := w.rawFile.GetCellValue(w.sheetName, "A1")
	if v != "编号" {
		t.Fatalf("expected header A1=编号 got %q", v)
	}

	v, _ = w.rawFile.GetCellValue(w.sheetName, "B1")
	if v != "姓名" {
		t.Fatalf("expected header B1=姓名 got %q", v)
	}

	// first data row
	v, _ = w.rawFile.GetCellValue(w.sheetName, "A2")
	if v != "1" {
		t.Fatalf("expected A2=1 got %q", v)
	}
	v, _ = w.rawFile.GetCellValue(w.sheetName, "B2")
	if v != "Alice" {
		t.Fatalf("expected B2=Alice got %q", v)
	}
}

func TestFromStruct_PointerElementsAndIgnore(t *testing.T) {
	w := APP.Writer.New(APP.WriterAttr.SheetName.Set("Sheet2"))

	type U struct {
		ID    int     `json:"id"`
		Name  *string `json:"name"`
		Email *string `excel:"邮箱" json:"email"`
		Skip  string  `excel:"-"`
	}

	n := "Nina"
	e := "nina@ex.com"
	data := []*U{{ID: 10, Name: &n, Email: &e, Skip: "shouldskip"}, nil}
	titles := []string{"id", "name", "邮箱", "skip"}

	w.FromStruct(&data, titles, 2)
	if w.Error != nil {
		t.Fatalf("unexpected error: %v", w.Error)
	}

	v, _ := w.rawFile.GetCellValue(w.sheetName, "A2")
	if v != "id" {
		t.Fatalf("expected header A2=id got %q", v)
	}

	v, _ = w.rawFile.GetCellValue(w.sheetName, "A3")
	if v != "10" {
		t.Fatalf("expected A3=10 got %q", v)
	}
	v, _ = w.rawFile.GetCellValue(w.sheetName, "B3")
	if v != "Nina" {
		t.Fatalf("expected B3=Nina got %q", v)
	}
	v, _ = w.rawFile.GetCellValue(w.sheetName, "C3")
	if v != "nina@ex.com" {
		t.Fatalf("expected C3=nina@ex.com got %q", v)
	}

	v, _ = w.rawFile.GetCellValue(w.sheetName, "D3")
	if v != "" {
		t.Fatalf("expected D3 empty got %q", v)
	}
}

func TestFromStruct_TimeField(t *testing.T) {
	w := APP.Writer.New(APP.WriterAttr.SheetName.Set("Sheet3"))

	type T struct {
		When time.Time `excel:"时间"`
	}

	now := time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)
	data := []T{{When: now}}
	titles := []string{"时间"}

	w.FromStruct(&data, titles, 1)
	if w.Error != nil {
		t.Fatalf("unexpected error: %v", w.Error)
	}

	v, _ := w.rawFile.GetCellValue(w.sheetName, "A2")
	if v == "" {
		t.Fatalf("expected non-empty time cell, got empty")
	}
}
