package excelV2

import (
	"fmt"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/dict/anyDictV2"
	"github.com/xuri/excelize/v2"
)

type Reader struct {
	Error       error
	data        anyDictV2.AnyDict[uint64, anyArrayV2.AnyArray[string]]
	rawFile     *excelize.File
	filename    string
	sheetName   string
	originalRow int
	finishedRow int
	titleRow    int
}

func (Reader) New(options ...ReaderAttributer) Reader {
	return new(Reader).SetAttrs(APP.ReaderAttr.SheetName("Sheet 1"), APP.ReaderAttr.OriginalRow(1), APP.ReaderAttr.TitleRow(1)).SetAttrs(options...)
}

func (my Reader) SetAttrs(options ...ReaderAttributer) Reader {
	for _, option := range options {
		option.Register(&my)
	}

	return my
}

func (my Reader) GetRawFile() *excelize.File { return my.rawFile }
func (my Reader) GetFilename() string        { return my.filename }
func (my Reader) GetSheetName() string       { return my.sheetName }
func (my Reader) GetOriginalRow() int        { return my.originalRow }
func (my Reader) GetFinishedRow() int        { return my.finishedRow }
func (my Reader) GetTitleRow() int           { return my.titleRow }

func (my Reader) Read() Reader {
	if my.rawFile, my.Error = excelize.OpenFile(my.filename); my.Error != nil {
		my.Error = fmt.Errorf("打开文件错误：%w", my.Error)
		return my
	}
	defer func(r *Reader) { _ = r.rawFile.Close() }(my)

	rows, err := my.excel.GetRows(my.GetSheetName())
	if err != nil {
		my.Err = ReadErr.Wrap(err)
		return my
	}

	my.data = anyDictV2.New[uint64, anyArrayV2.AnyArray[string]]()

	if my.finishedRow == 0 {
		for idx := range rows[my.GetOriginalRow():] {
			my.data = my.data.SetValue(rowNumber, anyArrayV2.NewList(rows[my.GetOriginalRow+idx]))
		}
	} else {
		for idx := range rows[my.GetOriginalRow():my.GetFinishedRow()] {
			my.data = my.data.SetValue(rowNumber, anyArrayV2.NewList(rows[my.GetOriginalRow+idx]))
		}
	}

	return my
}
