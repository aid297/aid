package excelV2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type (
	Reader struct {
		Error       error
		rawFile     *excelize.File
		filename    string
		sheetName   string
		originalRow int
		finishedRow int
		titleRow    int
	}
)

var (
	ExcelReaderApp Reader
)

func (*Reader) New(options ...IReaderOption) *Reader {
	return new(Reader).Set(options...)
}

func (*Reader) Set(options ...IReaderOption) *Reader {
	ins := new(Reader)

	ins.Set(
		ReaderOriginalRow(1),
		ReaderFinishedRow(0),
	)

	for _, option := range options {
		option.Register(ins)
	}

	return ins
}

func (my *Reader) Read() *Reader {
	if my.rawFile, my.Error = excelize.OpenFile(my.filename); my.Error != nil {
		my.Error = fmt.Errorf("打开文件错误：%w", my.Error)
		return my
	}
	defer func(r *Reader) { _ = r.rawFile.Close() }(my)

	return my
}
