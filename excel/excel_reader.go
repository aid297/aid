package excel

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/xuri/excelize/v2"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/dict"
)

// Reader Excel读取器
type Reader struct {
	Err         error
	data        *dict.AnyDict[uint64, *array.AnyArray[string]]
	excel       *excelize.File
	sheetName   string
	originalRow int
	finishedRow int
	titleRow    int
	titles      *array.AnyArray[string]
}

var ReaderApp Reader

func (*Reader) New() *Reader { return NewReader() }

// NewReader 构造函数
//
//go:fix 推荐使用New方法
func NewReader() *Reader {
	return &Reader{data: dict.Make[uint64, *array.AnyArray[string]]()}
}

// AutoRead 自动读取（默认第一行是表头，从第二行开始，默认Sheet名称为：Sheet1）
func (my *Reader) AutoRead(filename ...any) *Reader {
	return my.
		OpenFile(filename...).
		SetOriginalRow(2).
		SetTitleRow(1).
		SetSheetName("Sheet1").
		ReadTitle().
		Read()
}

// AutoReadBySheetName 自动读取（默认第一行是表头，从第二行开始）
func (my *Reader) AutoReadBySheetName(sheetName string, filename ...any) *Reader {
	return my.
		OpenFile(filename...).
		SetOriginalRow(2).
		SetTitleRow(1).
		SetSheetName(sheetName).
		ReadTitle().
		Read()
}

// Data 获取数据：有序字典
func (my *Reader) Data() *dict.AnyDict[uint64, *array.AnyArray[string]] { return my.data }

// DataWithTitle 获取数据：带有title的有序字典
func (my *Reader) DataWithTitle() (*dict.AnyDict[uint64, *dict.AnyDict[string, string]], error) {
	newDict := dict.Make[uint64, *dict.AnyDict[string, string]]()

	for idx, value := range my.data.ToMap() {
		newDict.Set(idx, dict.Zip(my.titles.ToSlice(), value.ToSlice()))
	}

	return newDict, nil
}

// SetDataByRow 设置单行数据
func (my *Reader) SetDataByRow(rowNumber uint64, data []string) *Reader {
	my.data.Set(rowNumber, array.New(data))

	return my
}

// GetSheetName 获取工作表名称
func (my *Reader) GetSheetName() string { return my.sheetName }

// SetSheetName 设置工作表名称
func (my *Reader) SetSheetName(sheetName string) *Reader {
	my.sheetName = sheetName

	return my
}

// GetOriginalRow 获取读取起始行
func (my *Reader) GetOriginalRow() int { return my.originalRow }

// SetOriginalRow 设置读取起始行
func (my *Reader) SetOriginalRow(originalRow int) *Reader {
	my.originalRow = originalRow - 1

	return my
}

// GetFinishedRow 获取读取终止行
func (my *Reader) GetFinishedRow() int { return my.finishedRow }

// SetFinishedRow 设置读取终止行
func (my *Reader) SetFinishedRow(finishedRow int) *Reader {
	my.finishedRow = finishedRow - 1

	return my
}

// GetTitleRow 获取表头行
func (my *Reader) GetTitleRow() int { return my.titleRow }

// SetTitleRow 设置表头行
func (my *Reader) SetTitleRow(titleRow int) *Reader {
	my.titleRow = titleRow - 1

	return my
}

// GetTitle 获取表头
func (my *Reader) GetTitle() *array.AnyArray[string] { return my.titles }

// SetTitle 设置表头
func (my *Reader) SetTitle(titles []string) *Reader {
	if len(titles) == 0 {
		my.Err = ReadErr.New("表头不能为空")
		return my
	}
	my.titles = array.New(titles)

	return my
}

// OpenFile 打开文件
func (my *Reader) OpenFile(filename ...any) *Reader {
	if filename[0].(string) == "" {
		my.Err = ReadErr.New("文件名不能为空")
		return my
	}
	f, err := excelize.OpenFile(fmt.Sprintf(filename[0].(string), filename[1:]...))
	if err != nil {
		my.Err = ReadErr.Wrap(fmt.Errorf("打开文件错误：%w", err))
		return my
	}
	my.excel = f

	defer func(r *Reader) {
		if err = r.excel.Close(); err != nil {
			r.Err = ReadErr.New("文件关闭错误")
		}
	}(my)

	my.SetTitleRow(1)
	my.SetOriginalRow(2)
	my.data = dict.Make[uint64, *array.AnyArray[string]]()

	return my
}

// ReadTitle 读取表头
func (my *Reader) ReadTitle() *Reader {
	if my.GetSheetName() == "" {
		my.Err = ReadErr.New("未设置工作表名称")
		return my
	}

	rows, err := my.excel.GetRows(my.GetSheetName())
	if err != nil {
		my.Err = ReadErr.New("读取表头错误")
		return my
	}
	my.SetTitle(rows[my.GetTitleRow()])

	return my
}

// Read 读取Excel
func (my *Reader) Read() *Reader {
	if my.GetSheetName() == "" {
		my.Err = ReadErr.New("未设置工作表名称")
		return my
	}

	rows, err := my.excel.GetRows(my.GetSheetName())
	if err != nil {
		my.Err = ReadErr.Wrap(err)
		return my
	}

	if my.finishedRow == 0 {
		for rowNumber, values := range rows[my.GetOriginalRow():] {
			my.SetDataByRow(uint64(rowNumber), values)
		}
	} else {
		for rowNumber, values := range rows[my.GetOriginalRow():my.GetFinishedRow()] {
			my.SetDataByRow(uint64(rowNumber), values)
		}
	}

	return my
}

// ToDataFrameDefaultType 获取DataFrame类型数据 通过Excel表头自定义数据类型
func (my *Reader) ToDataFrameDefaultType() dataframe.DataFrame {
	titleWithType := make(map[string]series.Type)
	for _, title := range my.GetTitle().ToSlice() {
		titleWithType[title] = series.String
	}

	return my.ToDataFrame(titleWithType)
}

// ToDataFrame 获取DataFrame类型数据
func (my *Reader) ToDataFrame(titleWithType map[string]series.Type) dataframe.DataFrame {
	if my.GetSheetName() == "" {
		my.Err = ReadErr.New("未设置工作表名称")
		return dataframe.DataFrame{}
	}

	var _content [][]string

	rows, err := my.excel.GetRows(my.GetSheetName())
	if err != nil {
		my.Err = ReadErr.Wrap(err)
		return dataframe.DataFrame{}
	}

	if my.finishedRow == 0 {
		_content = rows[my.GetTitleRow():]
	} else {
		_content = rows[my.GetTitleRow():my.GetFinishedRow()]
	}

	return dataframe.LoadRecords(
		_content,
		dataframe.DetectTypes(false),
		dataframe.DefaultType(series.String),
		dataframe.WithTypes(titleWithType),
	)
}

// ToDataFrameDetectType 获取DataFrame类型数据 通过自动探寻数据类型
func (my *Reader) ToDataFrameDetectType() dataframe.DataFrame {
	if my.GetSheetName() == "" {
		my.Err = ReadErr.New("未设置工作表名称")
		return dataframe.DataFrame{}
	}

	var _content [][]string

	rows, err := my.excel.GetRows(my.GetSheetName())
	if err != nil {
		my.Err = ReadErr.Wrap(err)
		return dataframe.DataFrame{}
	}

	if my.finishedRow == 0 {
		_content = rows[my.GetTitleRow():]
	} else {
		_content = rows[my.GetTitleRow():my.GetFinishedRow()]
	}

	return dataframe.LoadRecords(
		_content,
		dataframe.DetectTypes(true),
		dataframe.DefaultType(series.String),
	)
}
