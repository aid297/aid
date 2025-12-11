package excelV3

type (
	WriterAttributer interface{ Register(writer *Writer) }

	AttrWriterFilename  struct{ filename string }
	AttrWriterSheetName struct{ sheetName string }
	AttrWriterCell      struct{ cell *Cell }
)

func (AttrWriterSheetName) Set(val string) WriterAttributer { return AttrWriterSheetName{val} }
func (my AttrWriterSheetName) Register(writer *Writer)      { writer.sheetName = my.sheetName }

func (AttrWriterFilename) Set(val string) WriterAttributer { return AttrWriterFilename{val} }
func (my AttrWriterFilename) Register(writer *Writer)      { writer.filename = my.filename }

func (AttrWriterCell) Set(val *Cell) WriterAttributer { return AttrWriterCell{val} }
func (my AttrWriterCell) Register(writer *Writer)     { writer.setCell(my.cell) }
