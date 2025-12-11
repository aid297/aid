package excelV2

import (
	"sync"

	"github.com/aid297/aid/array/anyArrayV2"
)

type (
	Cell struct {
		lock                                                                                                                 *sync.RWMutex
		content                                                                                                              any
		contentType                                                                                                          CellContentType
		coordinate, fontRGB, patternRGB                                                                                      string
		fontBold, fontItalic                                                                                                 bool
		fontFamily                                                                                                           string
		fontSize                                                                                                             float64
		borderTopRGB, borderBottomRGB, borderLeftRGB, borderRightRGB, borderDiagonalUpRGB, borderDiagonalDownRGB             string
		borderTopStyle, borderBottomStyle, borderLeftStyle, borderRightStyle, borderDiagonalUpStyle, borderDiagonalDownStyle int
		wrapText                                                                                                             bool
	}

	// border 单元格边框
	border struct {
		Type  string
		RGB   string
		Style int
	}
)

const (
	CellContentTypeAny     CellContentType = "any"
	CellContentTypeFormula CellContentType = "formula"
	CellContentTypeInt     CellContentType = "int"
	CellContentTypeFloat   CellContentType = "float64"
	CellContentTypeBool    CellContentType = "bool"
	CellContentTypeTime    CellContentType = "time"
)

func (*Cell) New(attrs ...CellAttributer) *Cell {
	return (&Cell{lock: &sync.RWMutex{}, contentType: CellContentTypeAny}).setAttrs(attrs...)
}

func (*Cell) NewAny(content any, attrs ...CellAttributer) *Cell {
	return APP.Cell.New(APP.CellAttr.Content.Set(content), APP.CellAttr.ContentType.Set(CellContentTypeAny)).setAttrs(attrs...)
}

func (*Cell) NewFormula(content any, attrs ...CellAttributer) *Cell {
	return APP.Cell.New(APP.CellAttr.Content.Set(content), APP.CellAttr.ContentType.Set(CellContentTypeFormula)).setAttrs(attrs...)
}

func (*Cell) NewInt(content any, attrs ...CellAttributer) *Cell {
	return APP.Cell.New(APP.CellAttr.Content.Set(content), APP.CellAttr.ContentType.Set(CellContentTypeInt)).setAttrs(attrs...)
}

func (*Cell) NewFloat(content any, attrs ...CellAttributer) *Cell {
	return APP.Cell.New(APP.CellAttr.Content.Set(content), APP.CellAttr.ContentType.Set(CellContentTypeFloat)).setAttrs(attrs...)
}

func (*Cell) NewBool(content any, attrs ...CellAttributer) *Cell {
	return APP.Cell.New(APP.CellAttr.Content.Set(content), APP.CellAttr.ContentType.Set(CellContentTypeBool)).setAttrs(attrs...)
}

func (*Cell) NewTime(content any, attrs ...CellAttributer) *Cell {
	return APP.Cell.New(APP.CellAttr.Content.Set(content), APP.CellAttr.ContentType.Set(CellContentTypeTime)).setAttrs(attrs...)
}

func (my *Cell) setAttrs(attrs ...CellAttributer) *Cell {
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}
func (my *Cell) SetAttrs(attrs ...CellAttributer) *Cell {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.setAttrs(attrs...)
}

func (my *Cell) getContent() any { return my.content }
func (my *Cell) GetContent() any {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getContent()
}
func (my *Cell) getContentType() CellContentType { return my.contentType }
func (my *Cell) GetContentType() CellContentType {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getContentType()
}
func (my *Cell) getCoordinate() string { return my.coordinate }
func (my *Cell) GetCoordinate() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getCoordinate()
}
func (my *Cell) getFontRGB() string { return my.fontRGB }
func (my *Cell) GetFontRGB() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getFontRGB()
}
func (my *Cell) getPatternRGB() string { return my.patternRGB }
func (my *Cell) GetPatternRGB() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getPatternRGB()
}
func (my *Cell) getFontBold() bool { return my.fontBold }
func (my *Cell) GetFontBold() bool {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getFontBold()
}
func (my *Cell) getFontItalic() bool { return my.fontItalic }
func (my *Cell) GetFontItalic() bool {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getFontItalic()
}
func (my *Cell) getFontFamily() string { return my.fontFamily }
func (my *Cell) GetFontFamily() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getFontFamily()
}
func (my *Cell) getFontSize() float64 { return my.fontSize }
func (my *Cell) GetFontSize() float64 {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getFontSize()
}
func (my *Cell) getBorder() anyArrayV2.AnyArray[border] {
	borders := anyArrayV2.New[border]()

	if my.borderTopRGB != "" {
		borders = borders.Append(border{Type: "top", RGB: my.borderTopRGB, Style: my.borderTopStyle})
	}

	if my.borderBottomRGB != "" {
		borders = borders.Append(border{Type: "bottom", RGB: my.borderBottomRGB, Style: my.borderBottomStyle})
	}

	if my.borderLeftRGB != "" {
		borders = borders.Append(border{Type: "left", RGB: my.borderLeftRGB, Style: my.borderLeftStyle})
	}

	if my.borderRightRGB != "" {
		borders = borders.Append(border{Type: "right", RGB: my.borderRightRGB, Style: my.borderRightStyle})
	}

	if my.borderDiagonalUpRGB != "" {
		borders = borders.Append(border{Type: "diagonalUp", RGB: my.borderDiagonalUpRGB, Style: my.borderDiagonalUpStyle})
	}

	if my.borderDiagonalDownRGB != "" {
		borders = borders.Append(border{Type: "diagonalDown", RGB: my.borderDiagonalDownRGB, Style: my.borderDiagonalDownStyle})
	}

	return borders
}
func (my *Cell) GetBorder() anyArrayV2.AnyArray[border] {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getBorder()
}
func (my *Cell) getWrapText() bool { return my.wrapText }
func (my *Cell) GetWrapText() bool {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getWrapText()
}
