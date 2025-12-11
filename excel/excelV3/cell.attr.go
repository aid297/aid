package excelV3

type (
	CellAttributer interface{ Register(cell *Cell) }

	AttrCellContent       struct{ content any }
	AttrCellContentType   struct{ contentType CellContentType }
	AttrCellCoordinate    struct{ coordinate string }
	AttrCellFontRGB       struct{ rgb string }
	AttrCellPatternRGB    struct{ rgb string }
	AttrCellFontBold      struct{ fontBold bool }
	AttrCellFontItalic    struct{ fontItalic bool }
	AttrCellFontSize      struct{ fontSize float64 }
	AttrCellBorderRGB     struct{ top, bottom, left, right string }
	AttrCellBorderStyle   struct{ top, bottom, left, right int }
	AttrCellDiagonalRGB   struct{ up, down string }
	AttrCellDiagonalStyle struct{ up, down int }
	AttrCellWrapText      struct{ wrapText bool }
)

func (AttrCellContent) Set(val any) CellAttributer { return AttrCellContent{val} }
func (my AttrCellContent) Register(cell *Cell)     { cell.content = my.content }

func (AttrCellContentType) Set(val CellContentType) CellAttributer { return AttrCellContentType{val} }
func (my AttrCellContentType) Register(cell *Cell)                 { cell.contentType = my.contentType }

func (AttrCellCoordinate) Set(val string) CellAttributer { return AttrCellCoordinate{val} }
func (my AttrCellCoordinate) Register(cell *Cell)        { cell.coordinate = my.coordinate }

func (AttrCellFontRGB) Set(val string) CellAttributer { return AttrCellFontRGB{val} }
func (my AttrCellFontRGB) Register(cell *Cell)        { cell.fontRGB = my.rgb }

func (AttrCellPatternRGB) Set(val string) CellAttributer { return AttrCellPatternRGB{val} }
func (my AttrCellPatternRGB) Register(cell *Cell)        { cell.patternRGB = my.rgb }

func (AttrCellFontBold) Set(val bool) CellAttributer { return AttrCellFontBold{val} }
func (AttrCellFontBold) SetTrue() CellAttributer     { return AttrCellFontBold{true} }
func (AttrCellFontBold) SetFalse() CellAttributer    { return AttrCellFontBold{false} }
func (my AttrCellFontBold) Register(cell *Cell)      { cell.fontBold = my.fontBold }

func (AttrCellFontItalic) Set(val bool) CellAttributer { return AttrCellFontItalic{val} }
func (AttrCellFontItalic) SetTrue() CellAttributer     { return AttrCellFontItalic{true} }
func (AttrCellFontItalic) SetFalse() CellAttributer    { return AttrCellFontItalic{false} }
func (my AttrCellFontItalic) Register(cell *Cell)      { cell.fontItalic = my.fontItalic }

func (AttrCellFontSize) Set(val float64) CellAttributer { return AttrCellFontSize{val} }
func (my AttrCellFontSize) Register(cell *Cell)         { cell.fontSize = my.fontSize }

func (AttrCellBorderRGB) Set(top, bottom, left, right string) CellAttributer {
	return AttrCellBorderRGB{top, bottom, left, right}
}
func (my AttrCellBorderRGB) Register(cell *Cell) {
	cell.borderTopRGB = my.top
	cell.borderBottomRGB = my.bottom
	cell.borderLeftRGB = my.left
	cell.borderRightRGB = my.right
}

func (AttrCellBorderStyle) Set(top, bottom, left, right int) CellAttributer {
	return AttrCellBorderStyle{top, bottom, left, right}
}
func (my AttrCellBorderStyle) Register(cell *Cell) {
	cell.borderTopStyle = my.top
	cell.borderBottomStyle = my.bottom
	cell.borderLeftStyle = my.left
	cell.borderRightStyle = my.right
}

func (AttrCellDiagonalRGB) Set(up, down string) CellAttributer { return AttrCellDiagonalRGB{up, down} }
func (my AttrCellDiagonalRGB) Register(cell *Cell) {
	cell.borderDiagonalUpRGB = my.up
	cell.borderDiagonalDownRGB = my.down
}

func (AttrCellDiagonalStyle) Set(up, down int) CellAttributer { return AttrCellDiagonalStyle{up, down} }
func (my AttrCellDiagonalStyle) Register(cell *Cell) {
	cell.borderDiagonalUpStyle = my.up
	cell.borderDiagonalDownStyle = my.down
}

func (AttrCellWrapText) Set(val bool) CellAttributer { return AttrCellWrapText{val} }
func (AttrCellWrapText) SetTrue() CellAttributer     { return AttrCellWrapText{true} }
func (AttrCellWrapText) SetFalse() CellAttributer    { return AttrCellWrapText{false} }
func (my AttrCellWrapText) Register(cell *Cell)      { cell.wrapText = my.wrapText }
