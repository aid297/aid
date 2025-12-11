package excelV2

type (
	ReaderAttributer interface{ Register(reader *Reader) }

	AttrReaderFilename    struct{ filename string }
	AttrReaderSheetName   struct{ sheetName string }
	AttrReaderOriginalRow struct{ originalRow uint }
	AttrReaderFinishedRow struct{ finishedRow uint }
	AttrReaderOriginalCol struct{ originalCol uint }
	AttrReaderFinishedCol struct{ finishedCol uint }
)

func (AttrReaderFilename) Set(val string) ReaderAttributer { return AttrReaderFilename{val} }
func (my AttrReaderFilename) Register(reader *Reader)      { reader.filename = my.filename }

func (AttrReaderSheetName) Set(val string) ReaderAttributer { return AttrReaderSheetName{val} }
func (my AttrReaderSheetName) Register(reader *Reader)      { reader.sheetName = my.sheetName }

func (AttrReaderOriginalRow) Set(val uint) ReaderAttributer { return AttrReaderOriginalRow{val} }
func (my AttrReaderOriginalRow) Register(reader *Reader)    { reader.originalRow = my.originalRow }

func (AttrReaderFinishedRow) Set(val uint) ReaderAttributer { return AttrReaderFinishedRow{val} }
func (my AttrReaderFinishedRow) Register(reader *Reader)    { reader.finishedRow = my.finishedRow }

func (AttrReaderOriginalCol) Set(val uint) ReaderAttributer { return AttrReaderOriginalCol{val} }
func (my AttrReaderOriginalCol) Register(reader *Reader) {
	reader.originalColNo = my.originalCol
	reader.originalColTxt, reader.Error = ColumnNumberToText(int(my.originalCol))
}

func (AttrReaderFinishedCol) Set(val uint) ReaderAttributer { return AttrReaderFinishedCol{val} }
func (my AttrReaderFinishedCol) Register(reader *Reader) {
	reader.finishedColNo = my.finishedCol
	reader.finishedColTxt, reader.Error = ColumnNumberToText(int(my.finishedCol))
}
