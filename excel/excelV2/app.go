package excelV2

var APP struct {
	Writer     Writer
	WriterAttr struct {
		Filename  AttrWriterFilename
		SheetName AttrWriterSheetName
		Cells     AttrWriterCells
		Rows      AttrWriterRows
		Offset    AttrWriterOffset
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
	Reader     Reader
	ReaderAttr struct {
		Filename    AttrReaderFilename
		SheetName   AttrReaderSheetName
		OriginalRow AttrReaderOriginalRow
		FinishedRow AttrReaderFinishedRow
		OriginalCol AttrReaderOriginalCol
		FinishedCol AttrReaderFinishedCol
	}
}
