package excelV2

import (
	"fmt"
	"sync"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/dict/anyDictV2"
	"github.com/xuri/excelize/v2"
)

type Reader struct {
	Error       error
	lock        *sync.RWMutex
	data        anyDictV2.AnyDict[uint64, anyArrayV2.AnyArray[string]]
	rawFile     *excelize.File
	filename    string
	sheetName   string
	originalRow int
	finishedRow int
	titleRow    int
}

func (Reader) New(attrs ...ReaderAttributer) Reader {
	return Reader{lock: &sync.RWMutex{}}.SetAttrs(APP.ReaderAttr.SheetName.New("Sheet 1"), APP.ReaderAttr.OriginalRow.New(1), APP.ReaderAttr.TitleRow.New(1)).SetAttrs(attrs...)
}

func (my Reader) SetAttrs(attrs ...ReaderAttributer) Reader {
	my.lock.Lock()
	for idx := range attrs {
		attrs[idx].Register(&my)
	}
	my.lock.Unlock()
	return my
}

func (my Reader) GetRawFile() *excelize.File {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.rawFile
}
func (my Reader) GetFilename() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.filename
}
func (my Reader) GetSheetName() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.sheetName
}
func (my Reader) GetOriginalRow() int {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.originalRow
}
func (my Reader) GetFinishedRow() int {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.finishedRow
}
func (my Reader) GetTitleRow() int {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.titleRow
}

func (my Reader) Read() Reader {
	my.lock.Lock()
	if my.rawFile, my.Error = excelize.OpenFile(my.filename); my.Error != nil {
		my.Error = fmt.Errorf("打开文件错误：%w", my.Error)
		return my
	}
	defer func(r *Reader) { _ = r.rawFile.Close() }(&my)

	rows, err := my.rawFile.GetRows(my.GetSheetName())
	if err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrRead, err)
		return my
	}

	my.data = anyDictV2.New[uint64, anyArrayV2.AnyArray[string]]()

	if my.finishedRow == 0 {
		for idx := range rows[my.GetOriginalRow():] {
			my.data = my.data.SetValue(uint64(idx), anyArrayV2.NewList(rows[uint64(my.GetOriginalRow()+idx)]))
		}
	} else {
		for idx := range rows[my.GetOriginalRow():my.GetFinishedRow()] {
			my.data = my.data.SetValue(uint64(idx), anyArrayV2.NewList(rows[uint64(my.GetOriginalRow()+idx)]))
		}
	}

	my.lock.Unlock()
	return my
}
