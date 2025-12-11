package excelV2

import (
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

type Reader struct {
	Error          error
	lock           sync.RWMutex
	rawFile        *excelize.File
	filename       string
	sheetName      string
	originalRow    uint
	finishedRow    uint
	originalColNo  uint
	originalColTxt string
	finishedColNo  uint
	finishedColTxt string
	originalData   [][]string
	dataMap        map[string]string
	dataMaps       []map[string]string
}

func (*Reader) New(attrs ...ReaderAttributer) *Reader {
	return (&Reader{lock: sync.RWMutex{}, sheetName: "Sheet 1", originalRow: 1, originalColNo: 1, originalData: make([][]string, 0), dataMap: make(map[string]string), dataMaps: make([]map[string]string, 0)}).setAttrs(attrs...)
}

func (my *Reader) setAttrs(attrs ...ReaderAttributer) *Reader {
	for i := range attrs {
		attrs[i].Register(my)
	}
	return my
}

func (my *Reader) SetAttrs(attrs ...ReaderAttributer) *Reader {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.setAttrs(attrs...)
}

func (my *Reader) read() *Reader {
	var (
		err  error
		rows = make([][]string, 0)
	)

	if my.filename == "" {
		my.Error = ErrFilenameRequired
		return my
	}

	if my.sheetName == "" {
		my.Error = ErrSheetNameRequired
		return my
	}

	if my.rawFile, err = excelize.OpenFile(fmt.Sprintf(my.filename)); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrOpen, err)
		return my
	}

	defer func(r *Reader) {
		if err = r.rawFile.Close(); err != nil {
			r.Error = fmt.Errorf("%w：%w", ErrClose, err)
		}
	}(my)

	if rows, err = my.rawFile.GetRows(my.sheetName); err != nil {
		my.Error = fmt.Errorf("%w：%w", ErrRead, err)
		return my
	}

	if len(rows) > 0 {
		if my.finishedRow == 0 {
			for rowNo := range rows[my.originalRow:] {
				if my.finishedColNo == 0 {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:]
				} else {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:my.finishedColNo]
				}
			}
		} else {
			for rowNo := range rows[my.originalRow:my.finishedRow] {
				if my.finishedColNo == 0 {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:]
				} else {
					my.originalData[rowNo] = rows[rowNo][my.originalColNo:my.finishedColNo]
				}
			}
		}
	}

	return my
}

func (my *Reader) getOriginalData() [][]string {
	if len(my.originalData) == 0 && my.Error == nil {
		my.read()
	}

	return my.originalData
}

func (my *Reader) GetOriginalData() [][]string {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.getOriginalData()
}

func (my *Reader) GetMaps() []map[string]string {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.getOriginalData()

	for row := range my.originalData {
		my.dataMaps[row] = make(map[string]string)
		for col := range my.originalData[row] {
			var colTxt string
			if colTxt, my.Error = ColumnNumberToText(col); my.Error != nil {
				my.Error = fmt.Errorf("%w：行 %d 列索引 %d 转换为列名称错误", ErrColumnNumber, row+1, col+1)
				return nil
			}

			my.dataMaps[row][colTxt] = my.originalData[row][col]
		}
	}

	return my.dataMaps
}

func (my *Reader) GetMap() map[string]string {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.getOriginalData()

	for row := range my.originalData {
		for col := range my.originalData[row] {
			var colTxt string
			if colTxt, my.Error = ColumnNumberToText(col); my.Error != nil {
				my.Error = fmt.Errorf("%w：行 %d 列索引 %d 转换为列名称错误", ErrColumnNumber, row+1, col+1)
				return nil
			}

			my.dataMap[fmt.Sprintf("%s%d", colTxt, row)] = my.originalData[row][col]
		}
	}

	return my.dataMap
}
