package sqlapi

import "github.com/aid297/aid/simpleDB/kernal"

type StatementType string

const (
	StmtCreateTable StatementType = "create_table"
	StmtAlterTable  StatementType = "alter_table"
	StmtInsert      StatementType = "insert"
	StmtUpdate      StatementType = "update"
	StmtDelete      StatementType = "delete"
	StmtSelect      StatementType = "select"
	StmtDropTable   StatementType = "drop_table"
	StmtTruncate    StatementType = "truncate_table"
)

type Statement interface {
	Type() StatementType
	TableName() string
}

type CreateTableStmt struct {
	Table  string
	Schema kernal.TableSchema
}

func (s CreateTableStmt) Type() StatementType { return StmtCreateTable }
func (s CreateTableStmt) TableName() string   { return s.Table }

type AlterTableStmt struct {
	Table string
	Plan  kernal.AlterTablePlan
}

func (s AlterTableStmt) Type() StatementType { return StmtAlterTable }
func (s AlterTableStmt) TableName() string   { return s.Table }

type InsertStmt struct {
	Table string
	Row   kernal.Row
}

func (s InsertStmt) Type() StatementType { return StmtInsert }
func (s InsertStmt) TableName() string   { return s.Table }

type UpdateStmt struct {
	Table      string
	PrimaryKey any
	Updates    kernal.Row
}

func (s UpdateStmt) Type() StatementType { return StmtUpdate }
func (s UpdateStmt) TableName() string   { return s.Table }

type DeleteStmt struct {
	Table      string
	Conditions []kernal.QueryCondition
}

func (s DeleteStmt) Type() StatementType { return StmtDelete }
func (s DeleteStmt) TableName() string   { return s.Table }

type SelectStmt struct {
	Table      string
	Conditions []kernal.QueryCondition
}

func (s SelectStmt) Type() StatementType { return StmtSelect }
func (s SelectStmt) TableName() string   { return s.Table }

type DropTableStmt struct{ Table string }

func (s DropTableStmt) Type() StatementType { return StmtDropTable }
func (s DropTableStmt) TableName() string   { return s.Table }

type TruncateTableStmt struct{ Table string }

func (s TruncateTableStmt) Type() StatementType { return StmtTruncate }
func (s TruncateTableStmt) TableName() string   { return s.Table }

type ExecResult struct {
	Rows      []kernal.Row
	Affected  int
	Inserted  kernal.Row
	Updated   kernal.Row
	Statement StatementType
}
