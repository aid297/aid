package excelV2

import (
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type (
	IRows interface {
		GetError() error
		SetError(err error) IRows
		GetRows() []IRow
		SetRows(rows ...IRow) IRows
		AppendRows(rows ...IRow) IRows
	}

	Rows struct {
		Error       error
		mu          sync.RWMutex
		originalRow uint
		rows        []IRow
	}
)

func NewRows(originalRow uint, rows ...IRow) IRows {
	return (&Rows{
		originalRow: operationV2.
			NewTernary(
				operationV2.TrueValue(originalRow),
				operationV2.FalseValue[uint](1),
			).
			GetByValue(originalRow != 0),
	}).SetRows(rows...)
}

// GetError 获取错误
func (my *Rows) GetError() error {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.Error
}

// SetError 设置错误
func (my *Rows) SetError(err error) IRows {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.Error = err
	return my
}

// GetRows 获取 rows
func (my *Rows) GetRows() []IRow {
	my.mu.RLock()
	defer my.mu.RUnlock()

	return my.rows
}

// SetRows 设置 rows 自动赋值行号
func (my *Rows) SetRows(rows ...IRow) IRows {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.rows = rows

	for idx := range my.rows {
		my.rows[idx].SetRowNum(uint(idx) + my.originalRow)
	}

	return my
}

// AppendRows 追加 rows 自动赋值行号
func (my *Rows) AppendRows(rows ...IRow) IRows {
	my.mu.Lock()
	defer my.mu.Unlock()

	for idx := range rows {
		rows[idx].SetRowNum(uint(idx) + my.originalRow + uint(len(my.rows)))
	}

	my.rows = append(my.rows, rows...)

	return my
}