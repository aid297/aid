package excelV2

type Cell struct {
	content                                                                                                              any
	contentType                                                                                                          CellContentType
	coordinate, fontRgb, patternRgb                                                                                      string
	fontBold, fontItalic                                                                                                 bool
	fontFamily                                                                                                           string
	fontSize                                                                                                             float32
	borderTopRgb, borderBottomRgb, borderLeftRgb, borderRightRgb, borderDiagonalUpRgb, borderDiagonalDownRgb             string
	borderTopStyle, borderBottomStyle, borderLeftStyle, borderRightStyle, borderDiagonalUpStyle, borderDiagonalDownStyle int
	wrapText                                                                                                             bool
}

func (Cell) New(attrs ...CellAttributer) Cell {
	return Cell{}.SetAttrs(attrs...)
}

func (my Cell) SetAttrs(attrs ...CellAttributer) Cell {
	for idx := range attrs {
		attrs[idx].Register(&my)
	}
	return my
}
