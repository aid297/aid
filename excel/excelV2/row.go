package excelV2

import (
	"github.com/aid297/aid/array/anyArrayV2"
)

type Row struct {
	Error  error
	cells  anyArrayV2.AnyArray[Cell]
	number uint64
}

func (Row) New(attrs ...RowAttributer) Row {
	return Row{}.SetAttrs(attrs...)
}

func (my Row) SetAttrs(attrs ...RowAttributer) Row {
	for idx := range attrs {
		attrs[idx].Register(&my)
	}
	return my
}
