package excel

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/xuri/excelize/v2"
)

// Writer Excel写入器
type Writer struct {
	Err       error
	filename  string
	excel     *excelize.File
	sheetName string
}

var WriterApp Writer

func (*Writer) New(filename ...any) *Writer { return NewWriter(filename...) }

// NewWriter 初始化
//
//go:fix 推荐使用：New方法
func NewWriter(filename ...any) *Writer {
	ins := &Writer{}
	if filename[0].(string) == "" {
		ins.Err = WriteErr.New("文件名不能为空")
		return ins
	}
	ins.filename = fmt.Sprintf(filename[0].(string), filename[1:]...)
	ins.excel = excelize.NewFile()

	return ins
}

// GetFilename 获取文件名
func (my *Writer) GetFilename() string { return my.filename }

// SetFilename 设置文件名
func (my *Writer) SetFilename(filename string) *Writer {
	my.filename = filename

	return my
}

// CreateSheet 创建工作表
func (my *Writer) CreateSheet(sheetName string) *Writer {
	if sheetName == "" {
		my.Err = WriteErr.New("工作表名称不能为空")
		return my
	}
	sheetIndex, err := my.excel.NewSheet(sheetName)
	if err != nil {
		my.Err = WriteErr.Wrap(err)
		return my
	}

	my.excel.SetActiveSheet(sheetIndex)
	my.sheetName = my.excel.GetSheetName(sheetIndex)

	return my
}

// ActiveSheetByName 选择工作表（根据名称）
func (my *Writer) ActiveSheetByName(sheetName string) *Writer {
	if sheetName == "" {
		my.Err = WriteErr.New("工作表名称不能为空")
		return my
	}
	sheetIndex, err := my.excel.GetSheetIndex(sheetName)
	if err != nil {
		my.Err = WriteErr.Wrap(err)
		return my
	}

	my.excel.SetActiveSheet(sheetIndex)
	my.sheetName = sheetName

	return my
}

// ActiveSheetByIndex 选择工作表（根据编号）
func (my *Writer) ActiveSheetByIndex(sheetIndex int) *Writer {
	if sheetIndex < 0 {
		my.Err = WriteErr.New("工作表索引不能小于0")
		return my
	}
	my.excel.SetActiveSheet(sheetIndex)
	my.sheetName = my.excel.GetSheetName(sheetIndex)

	return my
}

// SetSheetName 设置sheet名称
func (my *Writer) SetSheetName(sheetName string) *Writer {
	my.excel.SetSheetName(my.sheetName, sheetName)
	my.sheetName = sheetName

	return my
}

// setStyleFont 设置字体
func (my *Writer) setStyleFont(cell *Cell) {
	fill := excelize.Fill{Type: "pattern", Pattern: 0, Color: []string{""}}
	if cell.GetPatternRgb() != "" {
		fill.Pattern = 1
		fill.Color[0] = cell.GetPatternRgb()
	}

	var borders = make([]excelize.Border, 0)
	if cell.GetBorder().Len() > 0 {
		for _, border := range cell.GetBorder().ToSlice() {
			borders = append(borders, excelize.Border{
				Type:  border.Type,
				Color: border.Rgb,
				Style: border.Style,
			})
		}
	}

	if style, err := my.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   cell.GetFontBold(),
			Italic: cell.GetFontItalic(),
			Family: cell.GetFontFamily(),
			Size:   cell.GetFontSize(),
			Color:  cell.GetFontRgb(),
		},
		Alignment: &excelize.Alignment{
			WrapText: cell.GetWrapText(),
		},
		Fill:   fill,
		Border: borders,
	}); err != nil {
		my.Err = WriteErr.Wrap(fmt.Errorf("设置字体错误：%s", cell.GetCoordinate()))
	} else {
		my.Err = my.excel.SetCellStyle(my.sheetName, cell.GetCoordinate(), cell.GetCoordinate(), style)
	}
}

// SetColumnWidthByIndex 设置单列宽：通过列索引
func (my *Writer) SetColumnWidthByIndex(col int, width float64) *Writer {
	my.SetColumnsWidthByIndex(col, col, width)

	return my
}

// SetColumnWidthByText 设置单列宽：通过列名称
func (my *Writer) SetColumnWidthByText(col string, width float64) *Writer {
	my.SetColumnsWidthByText(col, col, width)

	return my
}

// SetColumnsWidthByIndex 设置多列宽：通过列索引
func (my *Writer) SetColumnsWidthByIndex(startCol, endCol int, width float64) *Writer {
	startColText, err := ColumnNumberToText(startCol)
	if err != nil {
		my.Err = WriteErr.Wrap(fmt.Errorf("设置列宽错误：%s", err))
	}

	endColText, err := ColumnNumberToText(endCol)
	if err != nil {
		my.Err = WriteErr.Wrap(fmt.Errorf("设置列宽错误：%s", err))
	}

	if err = my.excel.SetColWidth(my.sheetName, startColText, endColText, width); err != nil {
		my.Err = WriteErr.Wrap(fmt.Errorf("设置列宽错误：%s", err))
	}

	return my
}

// SetColumnsWidthByText 设置多列宽：通过列名称
func (my *Writer) SetColumnsWidthByText(startCol, endCol string, width float64) *Writer {
	if err := my.excel.SetColWidth(my.sheetName, startCol, endCol, width); err != nil {
		my.Err = WriteErr.Wrap(fmt.Errorf("设置列宽错误：%s", err))
	}

	return my
}

// SetRows 设置行数据
func (my *Writer) SetRows(excelRows []*Row) *Writer {
	for _, row := range excelRows {
		my.AddRow(row)
	}

	return my
}

// AddRow 增加一行行数据
func (my *Writer) AddRow(excelRow *Row) *Writer {
	for _, cell := range excelRow.GetCells().ToSlice() {
		my.Err = my.excel.SetCellValue(my.sheetName, cell.GetCoordinate(), cell.GetContent())
		switch cell.GetContentType() {
		case CellContentTypeFormula:
			if err := my.excel.SetCellFormula(my.sheetName, cell.GetCoordinate(), cell.GetContent().(string)); err != nil {
				my.Err = WriteErr.Wrap(fmt.Errorf("写入数据错误（公式）%s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error()))
				return my
			}
		case CellContentTypeAny:
			if err := my.excel.SetCellValue(my.sheetName, cell.GetCoordinate(), cell.GetContent()); err != nil {
				my.Err = WriteErr.Wrap(fmt.Errorf("写入ExcelCell（任意） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error()))
				return my
			}
		case CellContentTypeInt:
			if err := my.excel.SetCellInt(my.sheetName, cell.GetCoordinate(), cell.GetContent().(int)); err != nil {
				my.Err = WriteErr.Wrap(fmt.Errorf("写入ExcelCell（整数） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error()))
				return my
			}
		case CellContentTypeFloat64:
			if err := my.excel.SetCellFloat(my.sheetName, cell.GetCoordinate(), cell.GetContent().(float64), 2, 64); err != nil {
				my.Err = WriteErr.Wrap(fmt.Errorf("写入ExcelCell（浮点数） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error()))
				return my
			}
		case CellContentTypeBool:
			if err := my.excel.SetCellBool(my.sheetName, cell.GetCoordinate(), cell.GetContent().(bool)); err != nil {
				my.Err = WriteErr.Wrap(fmt.Errorf("写入ExcelCell（布尔） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error()))
				return my
			}
		case CellContentTypeTime:
			if err := my.excel.SetCellValue(my.sheetName, cell.GetCoordinate(), cell.GetContent().(time.Time)); err != nil {
				my.Err = WriteErr.Wrap(fmt.Errorf("写入ExcelCell（时间） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error()))
			}
		}
		my.setStyleFont(cell)
	}

	return my
}

// SetTitleRow 设置标题行
func (my *Writer) SetTitleRow(titles []string, rowNumber uint64) *Writer {
	var (
		titleRow   *Row
		titleCells = make([]*Cell, len(titles))
	)

	if len(titles) > 0 {
		for idx, title := range titles {
			titleCells[idx] = NewCellAny(title)
		}

		titleRow = NewRow().SetRowNumber(rowNumber).SetCells(titleCells)

		my.AddRow(titleRow)
	}

	return my
}

// Save 保存文件
func (my *Writer) Save() error {
	if my.filename == "" {
		return WriteErr.New("未设置文件名")
	}

	return my.excel.SaveAs(my.filename)
}

// Download 下载Excel
func (my *Writer) Download(writer http.ResponseWriter) error {
	{
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(my.filename)))
		writer.Header().Set("Content-Transfer-Encoding", "binary")
		writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	}

	return my.excel.Write(writer)
}

// GetExcelizeFile 获取excelize文件对象
func (my *Writer) GetExcelizeFile() *excelize.File { return my.excel }
