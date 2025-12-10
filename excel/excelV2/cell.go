package excelV2

type Cell struct {
	content                                                                                                              any
	contentType                                                                                                          CellContentType
	coordinate, fontRgb, patternRgb                                                                                      string
	fontBold, fontItalic                                                                                                 bool
	fontFamily                                                                                                           string
	fontSize                                                                                                             float64
	borderTopRgb, borderBottomRgb, borderLeftRgb, borderRightRgb, borderDiagonalUpRgb, borderDiagonalDownRgb             string
	borderTopStyle, borderBottomStyle, borderLeftStyle, borderRightStyle, borderDiagonalUpStyle, borderDiagonalDownStyle int
	wrapText                                                                                                             bool
}
