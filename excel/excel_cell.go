package excel

import (
	"time"

	"github.com/aid297/aid/array"
)

type (
	// CellContentType 单元格内容类型
	CellContentType string

	// Cell Excel单元格
	Cell struct {
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

	// border 单元格边框
	border struct {
		Type  string
		Rgb   string
		Style int
	}
)

var CellApp Cell

func (*Cell) NewByAny(content any) *Cell        { return NewCellAny(content) }
func (*Cell) NewByInt(content any) *Cell        { return NewCellInt(content) }
func (*Cell) NewByFloat64(content any) *Cell    { return NewCellFloat64(content) }
func (*Cell) NewByBool(content any) *Cell       { return NewCellBool(content) }
func (*Cell) NewByTime(content time.Time) *Cell { return NewCellTime(content) }
func (*Cell) NewByFormula(content string) *Cell { return NewCellFormula(content) }

const (
	CellContentTypeAny     CellContentType = "any"
	CellContentTypeFormula CellContentType = "formula"
	CellContentTypeInt     CellContentType = "int"
	CellContentTypeFloat64 CellContentType = "float64"
	CellContentTypeBool    CellContentType = "bool"
	CellContentTypeTime    CellContentType = "time"
)

// NewCellAny 实例化：任意值
//
//go:fix 推荐使用：NewByAny方法
func NewCellAny(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeAny}
}

// NewCellInt 实例化：整数
//
//go:fix 推荐使用：NewByInt方法
func NewCellInt(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeInt}
}

// NewCellFloat64 实例化：浮点
//
//go:fix 推荐使用：NewByFloat64方法
func NewCellFloat64(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeFloat64}
}

// NewCellBool 实例化：布尔
//
//go:fix 推荐使用：NewByBool方法
func NewCellBool(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeBool}
}

// NewCellTime 实例化：时间
//
//go:fix 推荐使用：NewByTime方法
func NewCellTime(content time.Time) *Cell {
	return &Cell{content: content, contentType: CellContentTypeTime}
}

// NewCellFormula 实例化：公式
//
//go:fix 推荐使用：NewByFormula方法
func NewCellFormula(content string) *Cell {
	return &Cell{content: content, contentType: CellContentTypeFormula}
}

// GetBorder 获取边框
func (my *Cell) GetBorder() *array.AnyArray[border] {
	borders := array.Make[border](0)

	if my.borderTopRgb != "" {
		borders.Append(border{Type: "top", Rgb: my.borderTopRgb, Style: my.borderTopStyle})
	}

	if my.borderBottomRgb != "" {
		borders.Append(border{Type: "bottom", Rgb: my.borderBottomRgb, Style: my.borderBottomStyle})
	}

	if my.borderLeftRgb != "" {
		borders.Append(border{Type: "left", Rgb: my.borderLeftRgb, Style: my.borderLeftStyle})
	}

	if my.borderRightRgb != "" {
		borders.Append(border{Type: "right", Rgb: my.borderRightRgb, Style: my.borderRightStyle})
	}

	if my.borderDiagonalUpRgb != "" {
		borders.Append(border{Type: "diagonalUp", Rgb: my.borderDiagonalUpRgb, Style: my.borderDiagonalUpStyle})
	}

	if my.borderDiagonalDownRgb != "" {
		borders.Append(border{Type: "diagonalDown", Rgb: my.borderDiagonalDownRgb, Style: my.borderDiagonalDownStyle})
	}

	return borders
}

// SetWrapText 设置自动换行
func (my *Cell) SetWrapText(wrapText bool) *Cell {
	my.wrapText = wrapText

	return my
}

// GetWrapText 获取自动换行
func (my *Cell) GetWrapText() bool { return my.wrapText }

// SetBorderSurrounding 设置四周边框
func (my *Cell) SetBorderSurrounding(borderRgb string, borderStyle int, condition bool) *Cell {
	if condition {
		my.borderTopRgb = borderRgb
		my.borderBottomRgb = borderRgb
		my.borderLeftRgb = borderRgb
		my.borderRightRgb = borderRgb
		my.borderTopStyle = borderStyle
		my.borderBottomStyle = borderStyle
		my.borderLeftStyle = borderStyle
		my.borderRightStyle = borderStyle
	}

	return my
}

// SetBorderSurroundingFunc 设置四周边框 函数
func (my *Cell) SetBorderSurroundingFunc(condition func() (string, int, bool)) *Cell {
	if condition != nil {
		my.SetBorderSurrounding(condition())
	}

	return my
}

// SetBorderTopRgb 设置边框颜色：上
func (my *Cell) SetBorderTopRgb(borderTopRgb string, condition bool) *Cell {
	if condition {
		my.borderTopRgb = borderTopRgb
	}

	return my
}

// SetBorderTopRbgFunc 设置边框颜色：上 函数
func (my *Cell) SetBorderTopRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetBorderTopRgb(condition())
	}

	return my
}

// SetBorderTopStyle 设置边框样式：上
func (my *Cell) SetBorderTopStyle(borderTopStyle int, condition bool) *Cell {
	if condition {
		my.borderTopStyle = borderTopStyle
	}

	return my
}

// SetBorderTopStyleFunc 设置边框样式：上 函数
func (my *Cell) SetBorderTopStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		my.SetBorderTopStyle(condition())
	}

	return my
}

// SetBorderBottomRgb 设置边框颜色：下
func (my *Cell) SetBorderBottomRgb(borderBottomRgb string, condition bool) *Cell {
	if condition {
		my.borderBottomRgb = borderBottomRgb
	}

	return my
}

// SetBorderBottomRbgFunc 设置边框颜色：下 函数
func (my *Cell) SetBorderBottomRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetBorderBottomRgb(condition())
	}

	return my
}

// SetBorderBottomStyle 设置边框样式：下
func (my *Cell) SetBorderBottomStyle(borderBottomStyle int, condition bool) *Cell {
	if condition {
		my.borderBottomStyle = borderBottomStyle
	}

	return my
}

// SetBorderBottomStyleFunc 设置边框样式：下 函数
func (my *Cell) SetBorderBottomStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		my.SetBorderBottomStyle(condition())
	}

	return my
}

// SetBorderLeftRgb 设置边框颜色：左
func (my *Cell) SetBorderLeftRgb(borderLeftRgb string, condition bool) *Cell {
	if condition {
		my.borderLeftRgb = borderLeftRgb
	}

	return my
}

// SetBorderLeftRbgFunc 设置边框颜色：左 函数
func (my *Cell) SetBorderLeftRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetBorderLeftRgb(condition())
	}

	return my
}

// SetBorderLeftStyle 设置边框样式：左
func (my *Cell) SetBorderLeftStyle(borderLeftStyle int, condition bool) *Cell {
	if condition {
		my.borderLeftStyle = borderLeftStyle
	}

	return my
}

// SetBorderLeftStyleFunc 设置边框样式：左 函数
func (my *Cell) SetBorderLeftStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		my.SetBorderLeftStyle(condition())
	}

	return my
}

// SetBorderRightRgb 设置边框颜色：右
func (my *Cell) SetBorderRightRgb(borderRightRgb string, condition bool) *Cell {
	if condition {
		my.borderRightRgb = borderRightRgb
	}

	return my
}

// SetBorderRightRbgFunc 设置边框颜色：右 函数
func (my *Cell) SetBorderRightRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetBorderRightRgb(condition())
	}

	return my
}

// SetBorderRightStyle 设置边框样式：右
func (my *Cell) SetBorderRightStyle(borderRightStyle int, condition bool) *Cell {
	if condition {
		my.borderRightStyle = borderRightStyle
	}

	return my
}

// SetBorderRightStyleFunc 设置边框样式：右 函数
func (my *Cell) SetBorderRightStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		my.SetBorderRightStyle(condition())
	}

	return my
}

// SetBorderDiagonalUpRgb 设置边框颜色：对角线上
func (my *Cell) SetBorderDiagonalUpRgb(borderDiagonalUpRgb string, condition bool) *Cell {
	if condition {
		my.borderDiagonalUpRgb = borderDiagonalUpRgb
	}

	return my
}

// SetBorderDiagonalUpRbgFunc 设置边框颜色：对角线上 函数
func (my *Cell) SetBorderDiagonalUpRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetBorderDiagonalUpRgb(condition())
	}

	return my
}

// SetBorderDiagonalUpStyle 设置边框样式：对角线上
func (my *Cell) SetBorderDiagonalUpStyle(borderDiagonalUpStyle int, condition bool) *Cell {
	if condition {
		my.borderDiagonalUpStyle = borderDiagonalUpStyle
	}

	return my
}

// SetBorderDiagonalUpStyleFunc 设置边框样式：对角线上 函数
func (my *Cell) SetBorderDiagonalUpStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		my.SetBorderDiagonalUpStyle(condition())
	}

	return my
}

// SetBorderDiagonalDownRgb 设置边框颜色：对角线下
func (my *Cell) SetBorderDiagonalDownRgb(borderDiagonalDownRgb string, condition bool) *Cell {
	if condition {
		my.borderDiagonalDownRgb = borderDiagonalDownRgb
	}

	return my
}

// SetBorderDiagonalDownRbgFunc 设置边框颜色：对角线下 函数
func (my *Cell) SetBorderDiagonalDownRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetBorderDiagonalDownRgb(condition())
	}

	return my
}

// SetBorderDiagonalDownStyle 设置边框样式：对角线下
func (my *Cell) SetBorderDiagonalDownStyle(borderDiagonalDownStyle int, condition bool) *Cell {
	if condition {
		my.borderDiagonalDownStyle = borderDiagonalDownStyle
	}

	return my
}

// SetBorderDiagonalDownStyleFunc 设置边框样式：对角线下 函数
func (my *Cell) SetBorderDiagonalDownStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		my.SetBorderDiagonalDownStyle(condition())
	}

	return my
}

// GetFontRgb 获取字体颜色
func (my *Cell) GetFontRgb() string { return my.fontRgb }

// SetFontRgb 设置字体颜色
func (my *Cell) SetFontRgb(fontRgb string, condition bool) *Cell {
	if condition {
		my.fontRgb = fontRgb
	}

	return my
}

// SetFontRgbFunc 设置字体颜色：函数
func (my *Cell) SetFontRgbFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetFontRgb(condition())
	}

	return my
}

// GetPatternRgb 获取填充色
func (my *Cell) GetPatternRgb() string { return my.patternRgb }

// SetPatternRgb 设置填充色
func (my *Cell) SetPatternRgb(patternRgb string, condition bool) *Cell {
	if condition {
		my.patternRgb = patternRgb
	}

	return my
}

// SetPatternRgbFunc 设置填充色：函数
func (my *Cell) SetPatternRgbFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetPatternRgb(condition())
	}

	return my
}

// GetFontBold 获取字体粗体
func (my *Cell) GetFontBold() bool { return my.fontBold }

// SetFontBold 设置字体粗体
func (my *Cell) SetFontBold(fontBold bool, condition bool) *Cell {
	if condition {
		my.fontBold = fontBold
	}

	return my
}

// SetFontBoldFunc 设置字体粗体：函数
func (my *Cell) SetFontBoldFunc(condition func() (bool, bool)) *Cell {
	if condition != nil {
		my.SetFontBold(condition())
	}

	return my
}

// GetFontItalic 获取字体斜体
func (my *Cell) GetFontItalic() bool { return my.fontItalic }

// SetFontItalic 设置字体斜体
func (my *Cell) SetFontItalic(fontItalic bool, condition bool) *Cell {
	if condition {
		my.fontItalic = fontItalic
	}

	return my
}

// SetFontItalicFunc 设置字体斜体：函数
func (my *Cell) SetFontItalicFunc(condition func() (bool, bool)) *Cell {
	if condition != nil {
		my.SetFontItalic(condition())
	}

	return my
}

// GetFontFamily 获取字体
func (my *Cell) GetFontFamily() string { return my.fontFamily }

// SetFontFamily 设置字体
func (my *Cell) SetFontFamily(fontFamily string, condition bool) *Cell {
	if condition {
		my.fontFamily = fontFamily
	}

	return my
}

// SetFontFamilyFunc 设置字体：函数
func (my *Cell) SetFontFamilyFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		my.SetFontFamily(condition())
	}

	return my
}

// GetFontSize 获取字体字号
func (my *Cell) GetFontSize() float64 { return my.fontSize }

// SetFontSize 设置字体字号
func (my *Cell) SetFontSize(fontSize float64, condition bool) *Cell {
	if condition {
		my.fontSize = fontSize
	}

	return my
}

// SetFontSizeFunc 设置字体字号：函数
func (my *Cell) SetFontSizeFunc(condition func() (float64, bool)) *Cell {
	if condition != nil {
		my.SetFontSize(condition())
	}

	return my
}

// Init 初始化
func (my *Cell) Init(content any) *Cell {
	my.content = content

	return my
}

// GetContent 获取内容
func (my *Cell) GetContent() any { return my.content }

// SetContent 设置内容
func (my *Cell) SetContent(content any) *Cell {
	my.content = content

	return my
}

// GetCoordinate 获取单元格坐标
func (my *Cell) GetCoordinate() string { return my.coordinate }

// SetCoordinate 设置单元格坐标
func (my *Cell) SetCoordinate(coordinate string) *Cell {
	my.coordinate = coordinate

	return my
}

// GetContentType 获取单元格类型
func (my *Cell) GetContentType() CellContentType { return my.contentType }

// SetContentType 设置单元格类型
func (my *Cell) SetContentType(contentType CellContentType) *Cell {
	my.contentType = contentType

	return my
}
