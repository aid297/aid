package reader

import (
	`errors`

	`github.com/xuri/excelize/v2`

	`github.com/aid297/aid/filesystem/filesystemV4`
)

type ReaderAttribute func(reader *Read) (err error)

func Filename(filename string) ReaderAttribute {
	return func(reader *Read) (err error) {
		var fs filesystemV4.IFilesystem
		if fs = filesystemV4.NewFile(filesystemV4.Abs(filename)); !fs.GetExist() {
			return errors.New("文件不存在")
		}

		reader.filename = fs.GetFullPath()

		return
	}
}

func Filesystem(fs filesystemV4.IFilesystem) ReaderAttribute {
	return func(reader *Read) (err error) {
		if !fs.GetExist() {
			return errors.New("文件不存在")
		}

		reader.filename = fs.GetFullPath()

		return
	}
}

func UnzipXMLSizeLimit(limit int64) ReaderAttribute {
	return func(reader *Read) (err error) { reader.unzipXMLSizeLimit = limit; return }
}

func UnzipSizeLimit(limit int64) ReaderAttribute {
	return func(reader *Read) (err error) { reader.unzipSizeLimit = limit; return }
}

func OriginalRow(row int) ReaderAttribute {
	return func(reader *Read) (err error) { reader.originalRow = max(row, 1); return }
}

func FinishedRow(row int) ReaderAttribute {
	return func(reader *Read) (err error) { reader.finishedRow = max(row, 1); return }
}

func OriginalColumn(column int) ReaderAttribute {
	return func(reader *Read) (err error) { reader.originalCol = max(column, 1); return }
}

func FinishedColumn(column int) ReaderAttribute {
	return func(reader *Read) (err error) { reader.finishedCol = max(column, 1); return }
}

func OriginalColumnText(column string) ReaderAttribute {
	return func(reader *Read) (err error) { reader.originalCol, err = excelize.ColumnNameToNumber(column); return }
}

func FinishedColumnText(column string) ReaderAttribute {
	return func(reader *Read) (err error) { reader.finishedCol, err = excelize.ColumnNameToNumber(column); return }
}
