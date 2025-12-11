package excelV2

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"

	"github.com/aid297/aid/str"
)

type Writer struct {
	Error      error
	lock       sync.RWMutex
	rawFile    *excelize.File
	filename   string
	sheetName  string
	isSheetSet bool
}

func (*Writer) New(attrs ...WriterAttributer) *Writer {
	return (&Writer{lock: sync.RWMutex{}, rawFile: excelize.NewFile()}).setAttrs(attrs...)
}

func (my *Writer) setAttrs(attrs ...WriterAttributer) *Writer {
	for i := range attrs {
		attrs[i].Register(my)
	}
	return my
}

func (my *Writer) SetAttrs(attrs ...WriterAttributer) *Writer {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.setAttrs(attrs...)
}

func (my *Writer) getFilename() string { return my.filename }
func (my *Writer) GetFilename() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getFilename()
}
func (my *Writer) getSheetName() string { return my.sheetName }
func (my *Writer) GetSheetName() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getSheetName()
}

// setCell 设置 cells 值和格式
func (my *Writer) setCell(cell *Cell) *Writer {
	return my.setSheet().setCellValue(cell).setCellStyle(cell)
}

// setCellValue 设置 cells 值
func (my *Writer) setCellValue(cell *Cell) *Writer {
	var err error

	switch cell.GetContentType() {
	case CellContentTypeFormula:
		if err = my.rawFile.SetCellFormula(my.sheetName, cell.GetCoordinate(), cast.ToString(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellFormula, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeInt:
		if err = my.rawFile.SetCellInt(my.sheetName, cell.GetCoordinate(), cast.ToInt(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellInt, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeFloat:
		if err = my.rawFile.SetCellFloat(my.sheetName, cell.GetCoordinate(), cast.ToFloat64(cell.GetContent()), 2, 64); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellFloat, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeBool:
		if err = my.rawFile.SetCellBool(my.sheetName, cell.GetCoordinate(), cast.ToBool(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellBool, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	case CellContentTypeTime:
		if err = my.rawFile.SetCellValue(my.sheetName, cell.GetCoordinate(), cast.ToTime(cell.GetContent())); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellTime, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	default:
		if err = my.rawFile.SetCellValue(my.sheetName, cell.GetCoordinate(), cell.GetContent()); err != nil {
			my.Error = fmt.Errorf("%w：%s %s %w", ErrWriteCellAny, cell.GetCoordinate(), cell.GetContent(), err)
			return my
		}
	}

	return my
}

// setCellStyle 设置 cells 格式
func (my *Writer) setCellStyle(cell *Cell) *Writer {
	fill := excelize.Fill{Type: "pattern", Pattern: 0, Color: []string{""}}
	if cell.GetPatternRGB() != "" {
		fill.Pattern = 1
		fill.Color[0] = cell.GetPatternRGB()
	}

	var borders = make([]excelize.Border, 0)
	if cell.GetBorder().LengthNotEmpty() > 0 {
		cell.GetBorder().Each(func(_ int, item border) {
			borders = append(borders, excelize.Border{
				Type:  item.Type,
				Color: item.RGB,
				Style: item.Style,
			})
		})
	}

	if style, err := my.rawFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   cell.GetFontBold(),
			Italic: cell.GetFontItalic(),
			Family: cell.GetFontFamily(),
			Size:   cell.GetFontSize(),
			Color:  cell.GetFontRGB(),
		},
		Alignment: &excelize.Alignment{WrapText: cell.GetWrapText()},
		Fill:      fill,
		Border:    borders,
	}); err != nil {
		my.Error = fmt.Errorf("%w：%s", ErrSetFont, cell.GetCoordinate())
	} else {
		my.Error = my.rawFile.SetCellStyle(my.sheetName, cell.GetCoordinate(), cell.GetCoordinate(), style)
	}

	return my
}

func (my *Writer) createSheet(sheetName string) *Writer {
	if sheetName == "" {
		my.Error = ErrSheetNameRequired
		return my
	}
	sheetIndex, err := my.rawFile.NewSheet(sheetName)
	if err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrCreateSheet, err)
		return my
	}

	my.rawFile.SetActiveSheet(sheetIndex)
	my.sheetName = my.rawFile.GetSheetName(sheetIndex)

	return my
}

// CreateSheet 创建工作表
func (my *Writer) CreateSheet(sheetName string) *Writer {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.createSheet(sheetName)
}

func (my *Writer) setSheet() *Writer {
	var (
		sheetIndex int
		err        error
	)

	if my.isSheetSet {
		return my
	}

	if my.sheetName == "" {
		my.Error = ErrSheetNameRequired
		return my
	}

	if sheetIndex, err = my.rawFile.GetSheetIndex(my.sheetName); err != nil {
		my.Error = fmt.Errorf("%w：%w %s", ErrSetSheet, err, my.sheetName)
		return my
	}

	if sheetIndex == -1 {
		// sheet 不存在，创建sheet
		my.createSheet(my.sheetName)
	} else {
		// sheet 存在，设置为活动sheet
		my.rawFile.SetActiveSheet(sheetIndex)
	}

	my.isSheetSet = true
	return my
}

// Save 保存文件
func (my *Writer) Save() *Writer {
	var err error

	my.lock.Lock()

	if my.filename == "" {
		my.Error = ErrFilenameRequired
		return my
	}

	if err = my.rawFile.SaveAs(my.filename); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrSave, err)
		return my
	}

	my.lock.Unlock()
	return my
}

// Download 下载文件
func (my *Writer) Download(w http.ResponseWriter) *Writer {
	my.lock.Lock()
	defer my.lock.Unlock()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", str.APP.Buffer.JoinString("attachment; filename=%s", url.QueryEscape(my.filename)))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	my.Error = my.rawFile.Write(w)
	return my
}
