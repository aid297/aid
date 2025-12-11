package excelV2

var APP struct {
	Writer     Writer
	WriterAttr struct {
		Filename  AttrWriterFilename
		SheetName AttrWriterSheetName
		Cells     AttrWriterCells
		Rows      AttrWriterRows
	}
	Row     Row
	RowAttr struct {
		Cells  AttrCells
		Number AttrNumber
	}
	Cell     Cell
	CellAttr struct {
		Content       AttrCellContent
		ContentType   AttrCellContentType
		Coordinate    AttrCellCoordinate
		FontRGB       AttrCellFontRGB
		PatternRGB    AttrCellPatternRGB
		FontBold      AttrCellFontBold
		FontItalic    AttrCellFontItalic
		FontSize      AttrCellFontSize
		BorderRGB     AttrCellBorderRGB
		BorderStyle   AttrCellBorderStyle
		DiagonalRGB   AttrCellDiagonalRGB
		DiagonalStyle AttrCellDiagonalStyle
		WrapText      AttrCellWrapText
	}
}
