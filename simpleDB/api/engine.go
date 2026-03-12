package api

import (
	"fmt"

	"github.com/aid297/aid/simpleDB/driver"
	"github.com/aid297/aid/simpleDB/kernal"
)

type Backend string

const (
	BackendDriver Backend = "driver"
	BackendKernal Backend = "kernal"
)

type Engine struct {
	Database string
	Backend  Backend
}

func NewEngine(database string, backend Backend) *Engine {
	if backend == "" {
		backend = BackendDriver
	}
	return &Engine{Database: database, Backend: backend}
}

func (e *Engine) Parse(sql string) (Statement, error) { return Parse(sql) }

func (e *Engine) Execute(sql string) (ExecResult, error) {
	stmt, err := Parse(sql)
	if err != nil {
		return ExecResult{}, err
	}
	return e.ExecuteStatement(stmt)
}

func (e *Engine) ExecuteStatement(stmt Statement) (ExecResult, error) {
	switch s := stmt.(type) {
	case CreateTableStmt:
		return e.execCreate(s)
	case AlterTableStmt:
		return e.execAlter(s)
	case InsertStmt:
		return e.execInsert(s)
	case UpdateStmt:
		return e.execUpdate(s)
	case DeleteStmt:
		return e.execDelete(s)
	case SelectStmt:
		return e.execSelect(s)
	case DropTableStmt:
		return e.execDrop(s)
	case TruncateTableStmt:
		return e.execTruncate(s)
	default:
		return ExecResult{}, fmt.Errorf("unsupported statement")
	}
}

type tableDB interface {
	Close() error
	CreateTable(schema kernal.TableSchema) error
	AlterTable(plan kernal.AlterTablePlan) error
	InsertRow(values kernal.Row) (kernal.Row, error)
	UpdateRow(primaryKey any, updates kernal.Row) (kernal.Row, error)
	RemoveByCondition(conditions ...kernal.QueryCondition) (int, error)
	Find(conditions ...kernal.QueryCondition) ([]kernal.Row, error)
	DropTable() error
	TruncateTable() error
}

func (e *Engine) open(table string) (tableDB, error) {
	if e.Backend == BackendKernal {
		return kernal.New.DB(e.Database, table)
	}
	return driver.New.DB(e.Database, table)
}

func (e *Engine) execCreate(s CreateTableStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	if err = db.CreateTable(s.Schema); err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtCreateTable, Affected: 1}, nil
}

func (e *Engine) execAlter(s AlterTableStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	if err = db.AlterTable(s.Plan); err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtAlterTable, Affected: 1}, nil
}

func (e *Engine) execInsert(s InsertStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	row, err := db.InsertRow(s.Row)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtInsert, Affected: 1, Inserted: row}, nil
}

func (e *Engine) execUpdate(s UpdateStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	row, err := db.UpdateRow(s.PrimaryKey, s.Updates)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtUpdate, Affected: 1, Updated: row}, nil
}

func (e *Engine) execDelete(s DeleteStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	count, err := db.RemoveByCondition(s.Conditions...)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtDelete, Affected: count}, nil
}

func (e *Engine) execSelect(s SelectStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	rows, err := db.Find(s.Conditions...)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtSelect, Rows: rows, Affected: len(rows)}, nil
}

func (e *Engine) execDrop(s DropTableStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	if err = db.DropTable(); err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtDropTable, Affected: 1}, nil
}

func (e *Engine) execTruncate(s TruncateTableStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	defer db.Close()
	if err = db.TruncateTable(); err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtTruncate, Affected: 1}, nil
}
