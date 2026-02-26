package excelV2

import (
	"github.com/aid297/aid/filesystem/filesystemV4"
)

type (
	BaseAttributer interface {
		RegisterForWriter(writer Writer)
		RegisterForReader(reader Reader)
	}

	ExcelAttributer interface {
		FilenameAttributer
		SheetAttributer
		ReadRangeAttributer
		OpenFileAttributer
	}
)

type (
	FilenameAttributer BaseAttributer

	AttrFilename struct{ filename string }
)

func Filename(filename string) FilenameAttributer { return &AttrFilename{filename: filename} }
func File(file filesystemV4.Filesystemer) FilenameAttributer {
	return &AttrFilename{filename: file.GetFullPath()}
}
func (my *AttrFilename) RegisterForWriter(writer Writer) { writer.setFilename(my.filename) }
func (my *AttrFilename) RegisterForReader(reader Reader) { reader.setFilename(my.filename) }

type (
	SheetAttributer BaseAttributer

	AttrSheet struct {
		kind  string
		name  string
		index int
	}

	AttrSheetName  struct{ name string }
	AttrSheetIndex struct{ index int }
)

func SheetName(name string) SheetAttributer   { return &AttrSheet{name: name, kind: "name"} }
func SheetIndex(index int) SheetAttributer    { return &AttrSheet{index: index, kind: "index"} }
func CreateSheet(name string) SheetAttributer { return &AttrSheet{name: name, kind: "create"} }
func (my *AttrSheet) RegisterForWriter(writer Writer) {
	switch my.kind {
	case "name":
		writer.setSheetByName(my.name)
	case "index":
		writer.setSheetByIndex(my.index)
	case "create":
		writer.createSheet(my.name)
	}
}
func (my *AttrSheet) RegisterForReader(reader Reader) {}

type (
	ReadRangeAttributer BaseAttributer

	AttrReadRangeRow struct {
		row  int
		kind string
	}
)

func OriginalRow(row int) ReadRangeAttributer           { return &AttrReadRangeRow{row: row, kind: "ORIGINAL"} }
func FinishedRow(row int) ReadRangeAttributer           { return &AttrReadRangeRow{row: row, kind: "FINISHED"} }
func (my *AttrReadRangeRow) RegisterForWriter(_ Writer) {}
func (my *AttrReadRangeRow) RegisterForReader(reader Reader) {
	switch my.kind {
	case "ORIGINAL":
		reader.setOriginalRow(my.row)
	case "FINISHED":
		reader.setFinishedRow(my.row)
	}
}

type (
	OpenFileAttributer BaseAttributer

	AttrUnzipXMLSizeLimit struct{ size int64 }
	AttrUnzipSizeLimit    struct{ size int64 }
	AttrOpenFileSize      struct{ size int64 }
)

func UnzipXMLSizeLimit(size int64) OpenFileAttributer        { return &AttrUnzipXMLSizeLimit{size: size} }
func (my *AttrUnzipXMLSizeLimit) RegisterForWriter(_ Writer) {}
func (my *AttrUnzipXMLSizeLimit) RegisterForReader(reader Reader) {
	reader.setUnzipXMLSizeLimit(my.size)
}

func UnzipSizeLimit(size int64) OpenFileAttributer             { return &AttrUnzipSizeLimit{size: size} }
func (my *AttrUnzipSizeLimit) RegisterForWriter(_ Writer)      {}
func (my *AttrUnzipSizeLimit) RegisterForReader(reader Reader) { reader.setUnzipSizeLimit(my.size) }
