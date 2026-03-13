package api

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
	Rows  []kernal.Row
}

func (s InsertStmt) Type() StatementType { return StmtInsert }
func (s InsertStmt) TableName() string   { return s.Table }

type UpdateStmt struct {
	Table       string
	PrimaryKey  any
	PrimaryKeys []any
	Updates     kernal.Row
}

func (s UpdateStmt) Type() StatementType { return StmtUpdate }
func (s UpdateStmt) TableName() string   { return s.Table }

type DeleteStmt struct {
	Table      string
	Conditions []kernal.QueryCondition
}

func (s DeleteStmt) Type() StatementType { return StmtDelete }
func (s DeleteStmt) TableName() string   { return s.Table }

// SubqueryCondition 表示 IN/NOT IN 子查询条件，由 engine 展开为具体值列表后传入 kernal
type SubqueryCondition struct {
	Field   string
	NotIn   bool
	SubStmt SelectStmt
}

type JoinType string

const (
	JoinInner JoinType = "INNER"
	JoinLeft  JoinType = "LEFT"
)

type JoinClause struct {
	Type       JoinType // INNER | LEFT
	Table      string   // 从表名
	LeftAlias  string   // 主表字段（可含 "alias.field"）
	RightAlias string   // 从表字段（可含 "alias.field"）
}

type AggFunc string

const (
	AggNone  AggFunc = ""
	AggCount AggFunc = "COUNT"
	AggSum   AggFunc = "SUM"
	AggAvg   AggFunc = "AVG"
	AggMin   AggFunc = "MIN"
	AggMax   AggFunc = "MAX"
)

type SelectField struct {
	Star     bool    // SELECT *
	Field    string  // plain column name
	Agg      AggFunc // COUNT / SUM / AVG / MIN / MAX
	AggField string  // field inside aggregate ("*" for COUNT(*))
	Alias    string  // AS alias
}

type SelectStmt struct {
	Table      string
	Fields     []SelectField
	Joins      []JoinClause
	Conditions []kernal.QueryCondition
	SubConds   []SubqueryCondition
	GroupBy    []string
	OrderBy    string
	OrderDesc  bool
	Limit      int
	Offset     int
	Page       int
	PageSize   int
}

func (s SelectStmt) Type() StatementType { return StmtSelect }
func (s SelectStmt) TableName() string   { return s.Table }

type DropTableStmt struct{ Table string }

func (s DropTableStmt) Type() StatementType { return StmtDropTable }
func (s DropTableStmt) TableName() string   { return s.Table }

type TruncateTableStmt struct{ Table string }

func (s TruncateTableStmt) Type() StatementType { return StmtTruncate }
func (s TruncateTableStmt) TableName() string   { return s.Table }

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
	TotalItems  int `json:"totalItems"`
	PageSize    int `json:"pageSize"`
}

type ExecResult struct {
	Rows         []kernal.Row  `json:"rows,omitempty"`
	Affected     int           `json:"affected"`
	Inserted     kernal.Row    `json:"inserted,omitempty"`
	InsertedRows []kernal.Row  `json:"insertedRows,omitempty"`
	Updated      kernal.Row    `json:"updated,omitempty"`
	UpdatedRows  []kernal.Row  `json:"updatedRows,omitempty"`
	Statement    StatementType `json:"statement"`
	Pagination   *Pagination   `json:"pagination,omitempty"`
}
