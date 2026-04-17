package reader

import (
	`errors`
	`fmt`
	`os`
	`sync`

	`github.com/xuri/excelize/v2`
)

type (
	Reader interface {
		GetRawExcel() *excelize.File
		GetError() error
		Read(sheetName string, callback func(rowNum int, rows *excelize.Rows) (err error), attrs ...ReaderAttribute) Reader
	}

	Read struct {
		Error                             error
		mu                                sync.RWMutex
		filename                          string
		file                              *os.File
		excel                             *excelize.File
		originalRow, finishedRow          int
		originalCol, finishedCol          int
		unzipXMLSizeLimit, unzipSizeLimit int64
	}
)

func NewReader(attrs ...ReaderAttribute) Reader {
	return (&Read{
		mu:                sync.RWMutex{},
		unzipXMLSizeLimit: 10 * 1024 * 1024,
		unzipSizeLimit:    1 << 30,
	}).setAttrs(attrs...)
}

func (my *Read) setAttrs(attrs ...ReaderAttribute) Reader {
	if len(attrs) != 0 {
		for i := range attrs {
			if my.Error = attrs[i](my); my.Error != nil {
				return my
			}
		}
	}

	return my
}

// GetRawExcel 获取原始 excelize.File 对象
func (my *Read) GetRawExcel() *excelize.File { return my.excel }

// GetError 获取错误信息
func (my *Read) GetError() error { return my.Error }

// Read 读取数据，参数为可变参数 ReaderAttribute 接口类型，可以通过 OriginalRow 和 FinishedRow 来指定读取范围
func (my *Read) Read(
	sheetName string,
	callback func(rowNum int, rows *excelize.Rows) (err error),
	attrs ...ReaderAttribute,
) Reader {
	var (
		errOpen = errors.New("打开文件失败")
		errRead = errors.New("读取文件失败")
		err     error
		rows    *excelize.Rows
		rowNum  int
	)

	if my.setAttrs(attrs...); my.Error != nil {
		return my
	}

	if my.filename == "" {
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
	defer func() { _ = my.file.Close() }()
	defer func() { _ = my.excel.Close() }()

	// 3. 行迭代器流式读取数据
	if rows, err = my.excel.Rows(sheetName); err != nil {
		my.Error = fmt.Errorf("%w-未找到工作表：%w", errRead, err)
		return my
	}
	defer func() { _ = rows.Close() }()

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
