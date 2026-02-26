package excelV2

type (
	CellAttributer interface{ Register(cell *Cell) }

	AttrContent     struct{ content any }
	AttrContentType struct{ contentType CellContentType }
	AttrFont        struct{ font CellFontOpt }
	AttrBorder      struct {
		borderRGB   CellBorderRGBOpt
		borderStyle CellBorderStyleOpt
	}
	AttrAlignment  struct{ alignment CellAlignmentOpt }
	AttrCoordinate struct{ coordinate string }
)

func Content(content any) CellAttributer    { return &AttrContent{content: content} }
func (my *AttrContent) Register(cell *Cell) { cell.content = my.content }

func ContentType(contentType CellContentType) CellAttributer {
	return &AttrContentType{contentType: contentType}
}
func (my *AttrContentType) Register(cell *Cell) { cell.contentType = my.contentType }

func Font(font CellFontOpt) CellAttributer { return &AttrFont{font: font} }
func (my *AttrFont) Register(cell *Cell)   { cell.font = my.font }

func Border(borderRGB CellBorderRGBOpt, borderStyle CellBorderStyleOpt) CellAttributer {
	return &AttrBorder{borderRGB: borderRGB, borderStyle: borderStyle}
}

func (my *AttrBorder) Register(cell *Cell) {
	cell.borderRGB = my.borderRGB
	cell.borderStyle = my.borderStyle
}

func Alignment(alignment CellAlignmentOpt) CellAttributer {
	return &AttrAlignment{alignment: alignment}
}
func (my *AttrAlignment) Register(cell *Cell) { cell.alignment = my.alignment }

func Coordinate(coordinate string) CellAttributer { return &AttrCoordinate{coordinate: coordinate} }
func (my *AttrCoordinate) Register(cell *Cell)    { cell.coordinate = my.coordinate }
