package excelV2

type (
	CellAttributer interface{ Register(cell *Cell) }

	AttrCellContent       struct{ content any }
	AttrCellContentType   struct{ contentType CellContentType }
	AttrCellCoordinate    struct{ coordinate string }
	AttrCellFontRGB       struct{ rgb string }
	AttrCellPatternRGB    struct{ rgb string }
	AttrCellFontBold      struct{ fontBold bool }
	AttrCellFontItalic    struct{ fontItalic bool }
	AttrCellFontSize      struct{ fontSize float32 }
	AttrCellBorderRGB     struct{ top, bottom, left, right string }
	AttrCellBorderStyle   struct{ top, bottom, left, right int }
	AttrCellDiagonalRGB   struct{ up, down string }
	AttrCellDiagonalStyle struct{ up, down int }
	AttrCellWrapText      struct{ wrapText bool }
)

func (AttrCellContent) New(val any) CellAttributer { return AttrCellContent{val} }
func (my AttrCellContent) Register(cell *Cell)     { cell.content = my.content }

func (AttrCellContentType) New(val CellContentType) CellAttributer { return AttrCellContentType{val} }
func (my AttrCellContentType) Register(cell *Cell)                 { cell.contentType = my.contentType }

func (AttrCellCoordinate) New(val string) CellAttributer { return AttrCellCoordinate{val} }
func (my AttrCellCoordinate) Register(cell *Cell)        { cell.coordinate = my.coordinate }

func (AttrCellFontRGB) New(val string) CellAttributer { return AttrCellFontRGB{val} }
func (my AttrCellFontRGB) Register(cell *Cell)        { cell.fontRgb = my.rgb }

func (AttrCellPatternRGB) New(val string) CellAttributer { return AttrCellPatternRGB{val} }
func (my AttrCellPatternRGB) Register(cell *Cell)        { cell.patternRgb = my.rgb }

func (AttrCellFontBold) New(val bool) CellAttributer { return AttrCellFontBold{val} }
func (my AttrCellFontBold) Register(cell *Cell)      { cell.fontBold = my.fontBold }

func (AttrCellFontItalic) New(val bool) CellAttributer { return AttrCellFontItalic{val} }
func (my AttrCellFontItalic) Register(cell *Cell)      { cell.fontItalic = my.fontItalic }

func (AttrCellFontSize) New(val float32) CellAttributer { return AttrCellFontSize{val} }
func (my AttrCellFontSize) Register(cell *Cell)         { cell.fontSize = my.fontSize }

func (AttrCellBorderRGB) New(top, bottom, left, right string) CellAttributer {
	return AttrCellBorderRGB{top, bottom, left, right}
}
func (my AttrCellBorderRGB) Register(cell *Cell) {
	cell.borderTopRgb = my.top
	cell.borderBottomRgb = my.bottom
	cell.borderLeftRgb = my.left
	cell.borderRightRgb = my.right
}

func (AttrCellBorderStyle) New(top, bottom, left, right int) CellAttributer {
	return AttrCellBorderStyle{top, bottom, left, right}
}
func (my AttrCellBorderStyle) Register(cell *Cell) {
	cell.borderTopStyle = my.top
	cell.borderBottomStyle = my.bottom
	cell.borderLeftStyle = my.left
	cell.borderRightStyle = my.right
}

func (AttrCellDiagonalRGB) New(up, down string) CellAttributer { return AttrCellDiagonalRGB{up, down} }
func (my AttrCellDiagonalRGB) Register(cell *Cell) {
	cell.borderDiagonalUpRgb = my.up
	cell.borderDiagonalDownRgb = my.down
}

func (AttrCellDiagonalStyle) New(up, down int) CellAttributer { return AttrCellDiagonalStyle{up, down} }
func (my AttrCellDiagonalStyle) Register(cell *Cell) {
	cell.borderDiagonalUpStyle = my.up
	cell.borderDiagonalDownStyle = my.down
}

func (AttrCellWrapText) New(val bool) CellAttributer { return AttrCellWrapText{val} }
func (my AttrCellWrapText) Register(cell *Cell)      { cell.wrapText = my.wrapText }
