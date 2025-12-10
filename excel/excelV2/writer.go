package excelV2

import (
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

type Writer struct {
	Error     error
	rawFile   *excelize.File
	filename  string
	sheetName string
	lock      *sync.RWMutex
}

func (Writer) New(attrs ...WriterAttributer) Writer {
	return Writer{lock: &sync.RWMutex{}}.SetAttrs(attrs...)
}

func (my Writer) SetAttrs(attrs ...WriterAttributer) Writer {
	my.lock.Lock()

	for idx := range attrs {
		attrs[idx].Register(&my)
	}

	my.lock.Unlock()

	return my
}

func (my Writer) GetRawFile() *excelize.File {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.rawFile
}
func (my Writer) GetFilename() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.filename
}
func (my Writer) GetSheetName() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.sheetName
}

func (my Writer) Write() Writer {
	my.lock.Lock()
	if my.filename == "" {
		my.Error = ErrFilenameRequired
		return my
	}
	if my.sheetName == "" {
		my.Error = ErrSheetNameRequired
	}

	my.rawFile = excelize.NewFile()
	my.lock.Unlock()
	return my
}

// CreateSheet 创建工作表
func (my Writer) CreateSheet(sheetName string) Writer {
	my.lock.Lock()
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

	my.lock.Unlock()
	return my
}

// ActiveSheetByName 选择工作表（根据名称）
func (my Writer) ActiveSheetByName(sheetName string) Writer {
	my.lock.Lock()
	if sheetName == "" {
		my.Error = ErrSheetNameRequired
		return my
	}
	sheetIndex, err := my.rawFile.GetSheetIndex(sheetName)
	if err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrSheetNotFound, err)
		return my
	}

	my.rawFile.SetActiveSheet(sheetIndex)
	my.sheetName = my.rawFile.GetSheetName(sheetIndex)

	my.lock.Unlock()
	return my
}

// ActiveSheetByIndex 选择工作表（根据编号）
func (my *Writer) ActiveSheetByIndex(sheetIndex uint) *Writer {
	my.lock.Lock()
	my.rawFile.SetActiveSheet(int(sheetIndex))
	my.sheetName = my.rawFile.GetSheetName(int(sheetIndex))

	my.lock.Unlock()
	return my
}
