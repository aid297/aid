package excelV2

import "testing"

func Test1NewCell(t *testing.T) {
	var cell = APP.Cell.New(
		APP.CellAttr.Content.Set(123),
		APP.CellAttr.ContentType.Set(CellContentTypeInt),
		APP.CellAttr.Coordinate.Set("A1"),
		APP.CellAttr.FontRGB.Set("FF0000"),
		APP.CellAttr.PatternRGB.Set("00FF00"),
		APP.CellAttr.FontBold.SetTrue(),
		APP.CellAttr.FontItalic.SetFalse(),
		APP.CellAttr.FontSize.Set(12.5),
		APP.CellAttr.BorderRGB.Set("0000FF", "0000FF", "0000FF", "0000FF"),
		APP.CellAttr.BorderStyle.Set(1, 1, 1, 1),
		APP.CellAttr.DiagonalRGB.Set("FFFF00", "FFFF00"),
		APP.CellAttr.DiagonalStyle.Set(1, 1),
		APP.CellAttr.WrapText.SetTrue(),
	)

	t.Logf("%+v", *cell)
}
