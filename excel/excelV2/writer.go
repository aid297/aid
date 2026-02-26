package excelV2

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/aid297/aid/array/anySlice"
	"github.com/aid297/aid/operation/operationV2"
	"github.com/aid297/aid/str"
)

type (
	Writer interface {
		GetRawExcel() *excelize.File
		setFilename(filename string)
		SetFilename(attr ExcelAttributer) Writer
		setSheetByName(name string)
		setSheetByIndex(index int)
		SetSheet(attr SheetAttributer) Writer
		createSheet(name string)
		CreateSheet(attr SheetAttributer) Writer
		Save() (err error)
		Download(writer http.ResponseWriter) error
		Write(rows ...IRows) Writer
		setCellStyle(cell ICell)
	}

	Write struct {
		Error     error
		mu        sync.RWMutex
		filename  string
		excel     *excelize.File
		sheetName string
	}
)

func NewWriter() Writer { return &Write{mu: sync.RWMutex{}, excel: excelize.NewFile()} }

// GetRawExcel 获取原始 excelize.File 对象
func (my *Write) GetRawExcel() *excelize.File {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.excel
}

// setFilename 设置文件名
func (my *Write) setFilename(filename string) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.filename = filename
}

// SetFilename 设置文件名
func (my *Write) SetFilename(attr ExcelAttributer) Writer { attr.RegisterForWriter(my); return my }

// setSheetByName 通过名称设置工作表
func (my *Write) setSheetByName(name string) {
	var (
		err   error
		index int
	)

	my.mu.Lock()
	defer my.mu.Unlock()

	if index, err = my.excel.GetSheetIndex(name); err != nil {
		my.Error = fmt.Errorf("通过名称设置 Sheet 失败：%v", err)
		return
	}

	my.sheetName = name
	my.excel.SetActiveSheet(index)
}

// setSheetByIndex 通过索引设置工作表
func (my *Write) setSheetByIndex(index int) {
	my.mu.Lock()
	defer my.mu.Unlock()

	if index < 0 {
		my.Error = errors.New("工作表索引不能小于0")
		return
	}

	my.sheetName = my.excel.GetSheetName(index)
	my.excel.SetActiveSheet(index)
}

// SetSheet 设置当前工作表，参数为 SheetAttributer 接口类型，可以通过 SheetName 或 SheetIndex 来指定工作表
func (my *Write) SetSheet(attr SheetAttributer) Writer { attr.RegisterForWriter(my); return my }

// createSheet 创建工作表
func (my *Write) createSheet(name string) {
	var (
		err   error
		index int
	)

	my.mu.Lock()

	if index, err = my.excel.NewSheet(name); err != nil {
		my.Error = fmt.Errorf("设置工作表失败：%v", err)
		return
	}

	my.sheetName = name
	my.excel.SetActiveSheet(index)
	my.mu.Unlock()
}

// CreateSheet 创建工作表
func (my *Write) CreateSheet(attr SheetAttributer) Writer { attr.RegisterForWriter(my); return my }

// Write 写入数据
func (my *Write) Write(rows ...IRows) Writer {
	my.mu.Lock()
	defer my.mu.Unlock()

	for _, row := range rows {
		for _, cells := range row.GetRows() {
			for _, cell := range cells.GetCells() {
				switch cell.GetContentType() {
				case CellContentTypeFormula:
					if err := my.excel.SetCellFormula(
						my.sheetName,
						cell.GetCoordinate(),
						cell.GetContent().(string),
					); err != nil {
						my.Error = fmt.Errorf(
							"写入数据错误（公式）%s %s：%v",
							cell.GetCoordinate(),
							cell.GetContent(),
							err.Error(),
						)
						return my
					}
				case CellContentTypeAny:
					if err := my.excel.SetCellValue(
						my.sheetName,
						cell.GetCoordinate(),
						cell.GetContent(),
					); err != nil {
						my.Error = fmt.Errorf(
							"写入ExcelCell（通用） %s %s：%v",
							cell.GetCoordinate(),
							cell.GetContent(),
							err.Error(),
						)
						return my
					}
				case CellContentTypeInt:
					if err := my.excel.SetCellInt(
						my.sheetName,
						cell.GetCoordinate(),
						cell.GetContent().(int),
					); err != nil {
						my.Error = fmt.Errorf(
							"写入ExcelCell（整数） %s %s：%v",
							cell.GetCoordinate(),
							cell.GetContent(),
							err.Error(),
						)
						return my
					}
				case CellContentTypeFloat64:
					if err := my.excel.SetCellFloat(
						my.sheetName,
						cell.GetCoordinate(),
						cell.GetContent().(float64),
						2,
						64,
					); err != nil {
						my.Error = fmt.Errorf(
							"写入ExcelCell（浮点） %s %s：%v",
							cell.GetCoordinate(),
							cell.GetContent(),
							err.Error(),
						)
						return my
					}
				case CellContentTypeBool:
					if err := my.excel.SetCellBool(
						my.sheetName,
						cell.GetCoordinate(),
						cell.GetContent().(bool),
					); err != nil {
						my.Error = fmt.Errorf(
							"写入ExcelCell（布尔） %s %s：%v",
							cell.GetCoordinate(),
							cell.GetContent(),
							err.Error(),
						)
						return my
					}
				case CellContentTypeTime:
					if err := my.excel.SetCellValue(
						my.sheetName,
						cell.GetCoordinate(),
						cell.GetContent().(time.Time),
					); err != nil {
						my.Error = fmt.Errorf(
							"写入ExcelCell（时间） %s %s：%v",
							cell.GetCoordinate(),
							cell.GetContent(),
							err.Error(),
						)
					}
				}

				my.setCellStyle(cell)
			}
		}
	}

	return my
}

// setStyle 设置单元格样式
func (my *Write) setCellStyle(cell ICell) {
	var (
		err           error
		style         int
		cellBorders   anySlice.AnySlicer[excelize.Border]
		fill          excelize.Fill
		cellFont      CellFontOpt
		cellAlignment CellAlignmentOpt
	)

	// 设置填充
	fill = excelize.Fill{Type: "pattern", Pattern: 0, Color: []string{""}}
	cellFont = cell.GetFont()
	if cellFont.PatternRGB != "" {
		fill.Pattern = 1
		fill.Color[0] = cellFont.PatternRGB
	}

	cellBorders = cell.GetBorder()
	cellAlignment = cell.GetAlignment()
	if style, err = my.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   cellFont.Bold,
			Italic: cellFont.Italic,
			Family: cellFont.Family,
			Size: operationV2.NewTernary(
				operationV2.TrueValue(cellFont.Size),
				operationV2.FalseValue[float64](9),
			).GetByValue(cellFont.Size > 0),
			Color: cellFont.RGB,
		},
		Alignment: &excelize.Alignment{
			Horizontal: cellAlignment.Horizontal,
			Vertical:   cellAlignment.Vertical,
			WrapText:   cellAlignment.WrapText,
		},
		Fill:   fill,
		Border: cellBorders.ToSlice(),
	}); err != nil {
		my.Error = fmt.Errorf("设置字体错误：%s", cell.GetCoordinate())
		return
	}

	my.Error = my.excel.SetCellStyle(
		my.sheetName,
		cell.GetCoordinate(),
		cell.GetCoordinate(),
		style,
	)
}

// Save 保存文件
func (my *Write) Save() (err error) {
	my.mu.RLock()
	defer my.mu.RUnlock()

	if my.Error != nil {
		return my.Error
	}

	if my.filename == "" {
		return errors.New("保存文件失败：未设置文件名")
	}

	if err = my.excel.SaveAs(my.filename); err != nil {
		return fmt.Errorf("保存文件失败：%w", err)
	}

	return
}

// Download 下载Excel
func (my *Write) Download(writer http.ResponseWriter) error {
	{
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().
			Set("Content-Disposition", str.APP.Buffer.JoinString("attachment; filename=", url.QueryEscape(my.filename)))
		writer.Header().Set("Content-Transfer-Encoding", "binary")
		writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	}

	return my.excel.Write(writer)
}
