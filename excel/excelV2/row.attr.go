package excelV2

type (
	RowAttributer interface{ Register(row *Row) }

	AttrCells  struct{ cells []*Cell }
	AttrNumber struct{ number uint64 }
)

func (AttrCells) Set(cells ...*Cell) RowAttributer { return AttrCells{cells: cells} }
func (my AttrCells) Register(row *Row)             { row.cells = my.cells }

func (AttrNumber) Set(number uint64) RowAttributer { return AttrNumber{number: number} }
func (my AttrNumber) Register(row *Row)            { row.number = my.number }
