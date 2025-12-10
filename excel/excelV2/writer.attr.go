package excelV2

type (
	WriterAttributer interface{ Register(writer *Writer) }

	AttrWriterFilename struct{ filename string }
	AttrSheetName      struct{ sheetName string }
)

func (AttrSheetName) New(val string) WriterAttributer { return AttrSheetName{val} }
func (my AttrSheetName) Register(writer *Writer)      { writer.sheetName = my.sheetName }

func (AttrWriterFilename) New(val string) WriterAttributer { return AttrWriterFilename{val} }
func (my AttrWriterFilename) Register(writer *Writer)      { writer.filename = my.filename }
