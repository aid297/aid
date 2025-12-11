package excelV3

import "testing"

func Test1NewWriter(t *testing.T) {
	var writer *Writer

	writer = APP.Writer.New(
		APP.WriterAttr.Filename.Set("表格1.xlsx"),
		APP.WriterAttr.SheetName.Set("图1"),
	)

	t.Logf("filename: %v", writer.GetFilename())
	t.Logf("sheet name: %v", writer.GetSheetName())
}

func Test1SaveWriterToFile(t *testing.T) {
	var writer *Writer

	writer = APP.Writer.New(
		APP.WriterAttr.Filename.Set("表格1.xlsx"),
		APP.WriterAttr.SheetName.Set("图1"),
		APP.WriterAttr.Cell.Set(
			APP.Cell.New(
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
			),
		),
	)

	writer.Save()
}
