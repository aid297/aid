package excelV2

import "testing"

func TestFromStruct_FontAndPatternTags(t *testing.T) {
	w := APP.Writer.New(APP.WriterAttr.SheetName.Set("Tags1"))

	type S struct {
		Name string `excel:"名称" excel-font-size:"14.5" excel-font-rgb:"FF0000" excel-pattern-rgb:"00FF00"`
	}

	data := []S{{Name: "X"}}
	titles := []string{"名称"}

	w.FromStruct(&data, titles, 1)
	if w.Error != nil {
		t.Fatalf("unexpected error: %v", w.Error)
	}

	if sid, err := w.rawFile.GetCellStyle(w.sheetName, "A2"); err != nil || sid == 0 {
		t.Fatalf("expected non-zero style id for font/pattern tags, got sid=%v err=%v", sid, err)
	}
}

func TestFromStruct_BorderAndDiagonalTags(t *testing.T) {
	w := APP.Writer.New(APP.WriterAttr.SheetName.Set("Tags2"))

	type S struct {
		Col string `excel:"列" excel-border-rgb:"111111,222222,333333,444444" excel-border-style:"1,2,3,4" excel-diagonal-rgb:"AAAAAA,BBBBBB" excel-diagonal-style:"1,2"`
	}

	data := []S{{Col: "v"}}
	titles := []string{"列"}

	w.FromStruct(&data, titles, 1)
	if w.Error != nil {
		t.Fatalf("unexpected error: %v", w.Error)
	}
	if sid, err := w.rawFile.GetCellStyle(w.sheetName, "A2"); err != nil || sid == 0 {
		t.Fatalf("expected non-zero style id for border/diagonal tags, got sid=%v err=%v", sid, err)
	}
}

func TestFromStruct_WrapAndBoolTags(t *testing.T) {
	w := APP.Writer.New(APP.WriterAttr.SheetName.Set("Tags3"))

	type S struct {
		Txt string `excel:"文本" excel-wrap-text:"true" excel-font-bold:"true" excel-font-italic:"false"`
	}

	data := []S{{Txt: "t"}}
	titles := []string{"文本"}

	w.FromStruct(&data, titles, 1)
	if w.Error != nil {
		t.Fatalf("unexpected error: %v", w.Error)
	}
	if sid, err := w.rawFile.GetCellStyle(w.sheetName, "A2"); err != nil || sid == 0 {
		t.Fatalf("expected non-zero style id for wrap/bool tags, got sid=%v err=%v", sid, err)
	}
}

func TestFromStruct_InvalidTagValues_DoNotPanic(t *testing.T) {
	w := APP.Writer.New(APP.WriterAttr.SheetName.Set("Tags4"))

	type S struct {
		X string `excel:"X" excel-font-size:"bad" excel-font-bold:"notabool" excel-border-style:"x,y"`
	}

	data := []S{{X: "x"}}
	titles := []string{"X"}

	w.FromStruct(&data, titles, 1)
	if w.Error != nil {
		t.Fatalf("unexpected error with invalid tags: %v", w.Error)
	}
}
