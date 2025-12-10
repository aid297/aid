package excelV2

var APP struct {
	Reader     Reader
	ReaderAttr struct {
		Filename    AttrReaderFilename
		SheetName   AttrReaderSheetName
		OriginalRow AttrReaderOriginalRow
		FinishedRow AttrReaderFinishedRow
		TitleRow    AttrReaderTitleRow
	}
	Row     Row
	RowAttr struct {
		Cells  AttrRowCells
		Number AttrRowNumber
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
