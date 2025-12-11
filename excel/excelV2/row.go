package excelV2

import "sync"

type Row struct {
	lock   sync.RWMutex
	cells  []*Cell
	number uint64
}

func (*Row) New(attrs ...RowAttributer) *Row {
	return (&Row{lock: sync.RWMutex{}}).setAttrs(attrs...)
}

func (my *Row) setAttrs(attrs ...RowAttributer) *Row {
	for i := range attrs {
		attrs[i].Register(my)
	}

	return my
}

func (my *Row) SetAttrs(attrs ...RowAttributer) *Row {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.setAttrs(attrs...)
}

func (my *Row) getCells() []*Cell { return my.cells }
func (my *Row) GetCells() []*Cell {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getCells()
}

func (my *Row) getNumber() uint64 { return my.number }
func (my *Row) GetNumber() uint64 {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.getNumber()
}
