package excelV2

type (
	ReaderAttributer interface{ Register(excelReader *Reader) }

	AttrReaderFilename    struct{ filename string }
	AttrReaderSheetName   struct{ sheetName string }
	AttrReaderOriginalRow struct{ originalRow int }
	AttrReaderFinishedRow struct{ finishedRow int }
	AttrReaderTitleRow    struct{ titleRow int }
)

func (AttrReaderFilename) New(val string) ReaderAttributer { return AttrReaderFilename{val} }
func (my AttrReaderFilename) Register(reader *Reader)      { reader.filename = my.filename }

func (AttrReaderSheetName) New(val string)             { return AttrReaderSheetName{val} }
func (my AttrReaderSheetName) Register(reader *Reader) { reader.sheetName = my.sheetName }

func (AttrReaderOriginalRow) New(val int)                { return AttrReaderOriginalRow{val} }
func (my AttrReaderOriginalRow) Register(reader *Reader) { reader.originalRow = my.originalRow }

func (AttrReaderFinishedRow) New(val int)                { return AttrReaderFinishedRow{val} }
func (my AttrReaderFinishedRow) Register(reader *Reader) { reader.finishedRow = my.finishedRow }

func (AttrReaderTitleRow) New(val int)                { return AttrReaderTitleRow{val} }
func (my AttrReaderTitleRow) Register(reader *Reader) { reader.titleRow = my.titleRow }
