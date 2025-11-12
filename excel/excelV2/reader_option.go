package excelV2

type (
	IReaderOption interface {
		Register(excelReader *Reader)
	}

	ReaderFilenameOption struct{ filename string }

	ReaderSheetNameOption   struct{ sheetName string }
	ReaderOriginalRowOption struct{ originalRow int }
	ReaderFinishedRowOption struct{ finishedRow int }
	ReaderTitleRowOption    struct{ titleRow int }
)

func ReaderFilename(filename string) IReaderOption {
	return &ReaderFilenameOption{filename: filename}
}

func (my *ReaderFilenameOption) Register(excelReader *Reader) { excelReader.filename = my.filename }

func ReaderSheetName(sheetName string) IReaderOption {
	return &ReaderSheetNameOption{sheetName: sheetName}
}

func (my *ReaderSheetNameOption) Register(excelReader *Reader) {
	excelReader.sheetName = my.sheetName
}

func ReaderOriginalRow(originalRow int) IReaderOption {
	return &ReaderOriginalRowOption{originalRow: originalRow}
}

func (my *ReaderOriginalRowOption) Register(excelReader *Reader) {
	excelReader.originalRow = my.originalRow
}

func ReaderFinishedRow(finishedRow int) IReaderOption {
	return &ReaderFinishedRowOption{finishedRow: finishedRow}
}

func (my *ReaderFinishedRowOption) Register(excelReader *Reader) {
	excelReader.finishedRow = my.finishedRow
}

func ReaderTitleRow(titleRow int) IReaderOption {
	return &ReaderTitleRowOption{titleRow: titleRow}
}

func (my *ReaderTitleRowOption) Register(excelReader *Reader) {
	excelReader.titleRow = my.titleRow
}
