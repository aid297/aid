package excelV2

import (
	"regexp"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"

	"github.com/aid297/aid/array/anySlice"
)

type (
	CellContentType string

	ICell interface {
		SetAttrs(attrs ...CellAttributer) ICell
		GetCoordinate() string
		SetCoordinate(coordinate string) ICell
		SetRowNum(rowNum uint) ICell
		GetContent() any
		GetContentType() CellContentType
		GetFont() CellFontOpt
		GetBorder() anySlice.AnySlicer[excelize.Border]
		GetAlignment() CellAlignmentOpt
	}

	Cell struct {
		mu          sync.RWMutex
		content     any
		contentType CellContentType
		coordinate  string
		font        CellFontOpt
		borderRGB   CellBorderRGBOpt
		borderStyle CellBorderStyleOpt
		alignment   CellAlignmentOpt
	}

	CellAlignmentOpt struct {
		Horizontal, Vertical string
		WrapText             bool
	}
	CellBorderOpt struct {
		Type  string
		RGB   string
		Style int
	}
	CellBorderRGBOpt   struct{ Top, Bottom, Left, Right, DiagonalUp, DiagonalDown string }
	CellBorderStyleOpt struct{ Top, Bottom, Left, Right, DiagonalUp, DiagonalDown int }
	CellFontOpt        struct {
		Family          string
		Bold, Italic    bool
		RGB, PatternRGB string
		Size            float64
	}
)

const (
	CellContentTypeAny     CellContentType = "ANY"
	CellContentTypeFormula CellContentType = "FORMULA"
	CellContentTypeInt     CellContentType = "INT"
	CellContentTypeFloat64 CellContentType = "FLOAT"
	CellContentTypeBool    CellContentType = "BOOL"
	CellContentTypeTime    CellContentType = "TIME"
)

func NewCell(content any, attrs ...CellAttributer) ICell {
	return (&Cell{content: content, contentType: CellContentTypeAny, mu: sync.RWMutex{}}).SetAttrs(
		attrs...)
}

func NewCellFormula(content string, attrs ...CellAttributer) ICell {
	return (&Cell{content: content, contentType: CellContentTypeFormula, mu: sync.RWMutex{}}).SetAttrs(
		attrs...)
}

func NewCellInt(content int, attrs ...CellAttributer) ICell {
	return (&Cell{content: content, contentType: CellContentTypeInt, mu: sync.RWMutex{}}).SetAttrs(
		attrs...)
}

func NewCellFloat64(content float64, attrs ...CellAttributer) ICell {
	return (&Cell{content: content, contentType: CellContentTypeFloat64, mu: sync.RWMutex{}}).SetAttrs(
		attrs...)
}

func NewCellBool(content bool, attrs ...CellAttributer) ICell {
	return (&Cell{content: content, contentType: CellContentTypeBool, mu: sync.RWMutex{}}).SetAttrs(
		attrs...)
}

func NewCellTime(content time.Time, attrs ...CellAttributer) ICell {
	return (&Cell{content: content, contentType: CellContentTypeTime, mu: sync.RWMutex{}}).SetAttrs(
		attrs...)
}

// SetAttrs 设置属性
func (my *Cell) SetAttrs(attrs ...CellAttributer) ICell {
	my.mu.Lock()
	defer my.mu.Unlock()

	for _, attr := range attrs {
		attr.Register(my)
	}
	return my
}

// GetCoordinate 获取坐标
func (my *Cell) GetCoordinate() string {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.coordinate
}

// SetCoordinate 设置坐标
func (my *Cell) SetCoordinate(coordinate string) ICell {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.coordinate = coordinate
	return my
}

// SetRowNum 设置行号
func (my *Cell) SetRowNum(rowNum uint) ICell {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.coordinate = regexp.MustCompile(`[A-Za-z]+`).
		FindString(my.coordinate) +
		cast.ToString(
			rowNum,
		)
	return my
}

// GetContent 获取内容
func (my *Cell) GetContent() any {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.content
}

// GetContentType 获取内容类型
func (my *Cell) GetContentType() CellContentType {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.contentType
}

// GetFont 获取字体属性
func (my *Cell) GetFont() CellFontOpt {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.font
}

// GetBorder 获取边框属性
func (my *Cell) GetBorder() anySlice.AnySlicer[excelize.Border] {
	my.mu.RLock()
	defer my.mu.RUnlock()

	borders := anySlice.New(anySlice.Cap[excelize.Border](6))

	if my.borderRGB.Top != "" && my.borderStyle.Top > 0 {
		borders.Append(
			excelize.Border{Type: "top", Color: my.borderRGB.Top, Style: my.borderStyle.Top},
		)
	}

	if my.borderRGB.Bottom != "" && my.borderStyle.Bottom > 0 {
		borders.Append(
			excelize.Border{
				Type:  "bottom",
				Color: my.borderRGB.Bottom,
				Style: my.borderStyle.Bottom,
			},
		)
	}

	if my.borderRGB.Left != "" && my.borderStyle.Left > 0 {
		borders.Append(
			excelize.Border{Type: "left", Color: my.borderRGB.Left, Style: my.borderStyle.Left},
		)
	}

	if my.borderRGB.Right != "" && my.borderStyle.Right > 0 {
		borders.Append(
			excelize.Border{Type: "right", Color: my.borderRGB.Right, Style: my.borderStyle.Right},
		)
	}

	if my.borderRGB.DiagonalUp != "" && my.borderStyle.DiagonalUp > 0 {
		borders.Append(
			excelize.Border{
				Type:  "diagonalUp",
				Color: my.borderRGB.DiagonalUp,
				Style: my.borderStyle.DiagonalUp,
			},
		)
	}

	if my.borderRGB.DiagonalDown != "" && my.borderStyle.DiagonalDown > 0 {
		borders.Append(
			excelize.Border{
				Type:  "diagonalDown",
				Color: my.borderRGB.DiagonalDown,
				Style: my.borderStyle.DiagonalDown,
			},
		)
	}

	return borders
}

// GetAlignment 获取字体对齐
func (my *Cell) GetAlignment() CellAlignmentOpt {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.alignment
}
