package excelV2

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/xuri/excelize/v2"
)

type (
	Reader interface {
		GetRawExcel() *excelize.File
		setFilename(filename string)
		SetFilename(attr ExcelAttributer) Reader
		setOriginalRow(row int)
		setFinishedRow(row int)
		setUnzipXMLSizeLimit(size int64)
		setUnzipSizeLimit(size int64)
		SetOpenFile(attrs ...OpenFileAttributer) Reader
		Read(sheetName string, callback func(rowNum int, rows *excelize.Rows) (err error), attrs ...ReadRangeAttributer) Reader
	}

	Read struct {
		Error                             error
		mu                                sync.RWMutex
		filename                          string
		file                              *os.File
		excel                             *excelize.File
		originalRow                       int
		finishedRow                       int
		unzipXMLSizeLimit, unzipSizeLimit int64
	}
)

func NewReader(attrs ...ExcelAttributer) Reader {
	return &Read{
		mu:                sync.RWMutex{},
		unzipXMLSizeLimit: 10 * 1024 * 1024,
		unzipSizeLimit:    1 << 30,
	}
}

// GetRawExcel 获取原始 excelize.File 对象
func (my *Read) GetRawExcel() *excelize.File { return my.excel }

// setFilename 设置文件名
func (my *Read) setFilename(filename string) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.filename = filename
}

// SetFilename 设置文件名
func (my *Read) SetFilename(attr ExcelAttributer) Reader { attr.RegisterForReader(my); return my }

// 设置起始行
func (my *Read) setOriginalRow(row int) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.originalRow = row
}

// 设置终止行
func (my *Read) setFinishedRow(row int) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.finishedRow = row
}

// setUnzipXMLSizeLimit 设置解压缩XML大小限制，超过限制则写入临时文件，降低内存占用
func (my *Read) setUnzipXMLSizeLimit(size int64) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.unzipXMLSizeLimit = size
}

// setUnzipSizeLimit 设置解压缩大小限制，超过限制则写入临时文件，降低内存占用
func (my *Read) setUnzipSizeLimit(size int64) {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.unzipSizeLimit = size
}

// SetOpenFile 设置打开文件的属性，例如解压缩大小限制等
func (my *Read) SetOpenFile(attrs ...OpenFileAttributer) Reader {
	for _, attr := range attrs {
		attr.RegisterForReader(my)
	}
	return my
}

// Read 读取数据，参数为可变参数 ReadRangeAttributer 接口类型，可以通过 OriginalRow 和 FinishedRow 来指定读取范围
func (my *Read) Read(
	sheetName string,
	callback func(rowNum int, rows *excelize.Rows) (err error),
	attrs ...ReadRangeAttributer,
) Reader {
	var (
		errOpen = errors.New("打开文件失败")
		errRead = errors.New("读取文件失败")
		err     error
		rows    *excelize.Rows
		rowNum  int
	)

	if my.filename == "" {
		my.mu.RUnlock()
		my.Error = fmt.Errorf("%w：文件名不能为空", errOpen)
		return my
	}

	if sheetName == "" {
		my.Error = fmt.Errorf("%w-未设置工作表", errRead)
		return my
	}

	// 1. 初始化io.Reader（示例：本地文件流，可替换为HTTP流、bytes.Buffer等）
	if my.file, err = os.Open(my.filename); err != nil {
		my.mu.RUnlock()
		my.Error = fmt.Errorf("%w-打开文件错误：%w", errOpen, err)
		return my
	}

	// 2. 从io.Reader初始化File，可配置内存优化参数
	if my.excel, err = excelize.OpenReader(my.file, excelize.Options{
		UnzipXMLSizeLimit: my.unzipXMLSizeLimit, // 超过10MB则将XML写入临时文件，降低内存
		UnzipSizeLimit:    my.unzipSizeLimit,    // 解压缩大小限制1GB
	}); err != nil {
		my.mu.RUnlock()
		my.Error = fmt.Errorf("%w-解析文件错误：%w", errOpen, err)
		return my
	}
	defer my.file.Close()
	defer my.excel.Close()

	for _, attr := range attrs {
		attr.RegisterForReader(my)
	}

	// 3. 行迭代器流式读取数据
	if rows, err = my.excel.Rows(sheetName); err != nil {
		my.Error = fmt.Errorf("%w-未找到工作表：%w", errRead, err)
		return my
	}
	defer rows.Close()

	// 4. 逐行遍历，仅加载当前行数据到内存
	for rows.Next() {
		rowNum++

		if (my.originalRow != 0 && rowNum < my.originalRow) || (my.finishedRow != 0 && rowNum > my.finishedRow) {
			continue
		}

		if err = callback(rowNum, rows); err != nil {
			my.Error = fmt.Errorf("%w-回调函数执行失败：%w [行号：%d]", errRead, err, rowNum)
			return my
		}
	}

	return my
}
