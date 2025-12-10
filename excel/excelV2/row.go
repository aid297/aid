package excelV2

import (
	"github.com/aid297/aid/array"
	"github.com/aid297/aid/array/anyArrayV2"
)

type Row struct {
	Error  error
	cells  *array.AnyArray[*Cell]
	number uint64
}

type (
	RowAttributer interface{ Register(row *Row) }

	AttrRowCells  struct{ cells anyArrayV2.AnyArray[Cell] }
	AttrRowNumber struct{ number uint64 }
)

func (AttrRowCells) New(values ...Cell) RowAttributer {
	return AttrRowCells{cells: anyArrayV2.NewItems(values...)}
}
func (my AttrRowCells) Register(row *Row) { row.cells = my.cells }

func (AttrRowNumber) New(val uint64) RowAttributer { return AttrRowNumber{val} }
func (my AttrRowNumber) Register(row *Row)         { row.number = my.number }
