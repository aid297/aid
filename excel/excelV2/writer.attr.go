package excelV2

import (
	"fmt"

	"github.com/aid297/aid/operation/operationV2"
)

type (
	WriterAttributer interface{ Register(writer *Writer) }

	AttrWriterFilename  struct{ filename string }
	AttrWriterSheetName struct{ sheetName string }
	AttrWriterCells     struct{ cells []*Cell }
	AttrWriterRows      struct{ rows []*Row }
)

func (AttrWriterSheetName) Set(val string) WriterAttributer { return AttrWriterSheetName{val} }
func (my AttrWriterSheetName) Register(writer *Writer)      { writer.sheetName = my.sheetName }

func (AttrWriterFilename) Set(val string) WriterAttributer { return AttrWriterFilename{val} }
func (my AttrWriterFilename) Register(writer *Writer)      { writer.filename = my.filename }

func (AttrWriterCells) Set(vals ...*Cell) WriterAttributer { return AttrWriterCells{vals} }
func (my AttrWriterCells) Register(writer *Writer) {
	for _, cell := range my.cells {
		writer.setCell(cell)
	}
}

func (AttrWriterRows) Set(vals ...*Row) WriterAttributer { return AttrWriterRows{vals} }
func (my AttrWriterRows) Register(writer *Writer) {
	for rn, row := range my.rows {
		rn = operationV2.NewTernary(operationV2.TrueValue[int](int(row.getNumber())), operationV2.FalseValue(rn)).GetByValue(row.getNumber() > 0)
		for cn, cell := range row.cells {
			var (
				err error
				col string
			)

			if col, err = ColumnNumberToText(cn + 1); err != nil {
				writer.Error = fmt.Errorf("%w：行 %d 列索引 %d 转换为列名称错误", ErrColumnNumber, rn+1, cn+1)
				return
			}

			col = operationV2.NewTernary(operationV2.TrueValue(cell.GetCoordinate()), operationV2.FalseValue(col)).GetByValue(cell.getCoordinate() != "")
			cell.setAttrs(APP.CellAttr.Coordinate.Set(fmt.Sprintf("%s%d", col, rn+1)))
			writer.setCell(cell)
		}
	}
}
