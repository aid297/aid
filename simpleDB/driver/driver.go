package driver

import "github.com/aid297/aid/simpleDB/kernal"

var New app

type app struct{}

type (
	Row            = kernal.Row
	RowUpdate      = kernal.RowUpdate
	Column         = kernal.Column
	ForeignKey     = kernal.ForeignKey
	TableSchema    = kernal.TableSchema
	QueryCondition = kernal.QueryCondition
	CascadeQuery   = kernal.CascadeQuery
	CascadeInclude = kernal.CascadeInclude
	DatabaseConfig = kernal.DatabaseConfig
	AlterTablePlan = kernal.AlterTablePlan
)

const (
	QueryOpEQ         = kernal.QueryOpEQ
	QueryOpNE         = kernal.QueryOpNE
	QueryOpGT         = kernal.QueryOpGT
	QueryOpGTE        = kernal.QueryOpGTE
	QueryOpLT         = kernal.QueryOpLT
	QueryOpLTE        = kernal.QueryOpLTE
	QueryOpIn         = kernal.QueryOpIn
	QueryOpNotIn      = kernal.QueryOpNotIn
	QueryOpBetween    = kernal.QueryOpBetween
	QueryOpNotBetween = kernal.QueryOpNotBetween
)

type DB struct {
	core *kernal.SimpleDB
}

type TxIsolation string

const (
	TxIsolationSnapshot      TxIsolation = "snapshot"
	TxIsolationReadCommitted TxIsolation = "read_committed"
)

type TxOptions struct {
	ReadOnly  bool
	Isolation TxIsolation
}

type Tx struct {
	core *kernal.Tx
}

// DB opens a database/table and returns a driver-level DB wrapper.
func (*app) DB(database, table string, attrs ...kernal.SchemaAttributer) (*DB, error) {
	core, err := kernal.New.DB(database, table)
	if err != nil {
		return nil, wrapError(err)
	}
	if len(attrs) > 0 {
		core.SetAttrs(attrs...)
	}
	return &DB{core: core}, nil
}

func (*app) EnsureSystemTables(database string) error {
	return wrapError(kernal.New.EnsureSystemTables(database))
}

func (db *DB) Close() error { return wrapError(db.core.Close()) }

func (db *DB) SetAttrs(attrs ...kernal.SchemaAttributer) *DB {
	db.core.SetAttrs(attrs...)
	return db
}

func (db *DB) GetConfig() DatabaseConfig { return db.core.GetConfig() }

func (db *DB) Compact() error { return wrapError(db.core.Compact()) }

func (db *DB) Configure(schema TableSchema) error { return wrapError(db.core.Configure(schema)) }

func (db *DB) GetSchema() (*TableSchema, error) {
	schema, err := db.core.GetSchema()
	if err != nil {
		return nil, wrapError(err)
	}
	return schema, nil
}

// ─── DDL: Schema Manager ─────────────────────────────────────────────────────

// HasSchema 返回表是否已有 Schema 定义。
func (db *DB) HasSchema() bool { return db.core.HasSchema() }

// CreateTable 以严格 DDL 语义创建表结构，若已存在 Schema 则返回错误。
func (db *DB) CreateTable(schema TableSchema) error { return wrapError(db.core.CreateTable(schema)) }

// DropTable 删除所有行数据及 Schema，但不关闭数据库文件。
func (db *DB) DropTable() error { return wrapError(db.core.DropTable()) }

// TruncateTable 清空所有行数据，保留 Schema 定义。
func (db *DB) TruncateTable() error { return wrapError(db.core.TruncateTable()) }

// AlterTable 按 AlterTablePlan 修改表结构（AddColumn/DropColumn/AddIndex 等）。
func (db *DB) AlterTable(plan AlterTablePlan) error { return wrapError(db.core.AlterTable(plan)) }

// SchemaDiff 计算从当前 Schema 到目标 Schema 的迁移计划，不执行任何变更。
// 返回值：(plan, schemaExists, error)
//   - plan == nil && !schemaExists → 当前无 Schema，需调用 CreateTable
//   - plan == nil && schemaExists  → Schema 已与目标一致，无需变更
//   - plan != nil                  → 有差异，plan 描述所需变更
func (db *DB) SchemaDiff(target TableSchema) (*AlterTablePlan, bool, error) {
	plan, exists, err := db.core.SchemaDiff(target)
	if err != nil {
		return nil, exists, wrapError(err)
	}
	return plan, exists, nil
}

// AutoMigrate 保守迁移：无表则建，有表则只**增**列/索引，绝不删除（幂等安全）。
func (db *DB) AutoMigrate(schema TableSchema) error { return wrapError(db.core.AutoMigrate(schema)) }

// SyncSchema 完全同步：无表则建，有差异则全量同步（含删列，⚠️ 不可逆）。
func (db *DB) SyncSchema(schema TableSchema) error { return wrapError(db.core.SyncSchema(schema)) }

func (db *DB) InsertRow(values Row) (Row, error) {
	row, err := db.core.InsertRow(values)
	return row, wrapError(err)
}

func (db *DB) InsertRows(values []Row) ([]Row, error) {
	rows, err := db.core.InsertRows(values)
	return rows, wrapError(err)
}

func (db *DB) UpdateRow(primaryKey any, updates Row) (Row, error) {
	row, err := db.core.UpdateRow(primaryKey, updates)
	return row, wrapError(err)
}

func (db *DB) UpdateRows(updates []RowUpdate) ([]Row, error) {
	rows, err := db.core.UpdateRows(updates)
	return rows, wrapError(err)
}

func (db *DB) DeleteRow(primaryKey any) error { return wrapError(db.core.DeleteRow(primaryKey)) }

func (db *DB) DeleteRows(primaryKeys []any) error { return wrapError(db.core.DeleteRows(primaryKeys)) }

func (db *DB) FindRow(primaryKey any) (Row, bool, error) {
	row, ok, err := db.core.FindRow(primaryKey)
	return row, ok, wrapError(err)
}

func (db *DB) FindByConditions(conditions []QueryCondition) ([]Row, error) {
	rows, err := db.core.FindByConditions(conditions)
	return rows, wrapError(err)
}

func (db *DB) Find(conditions ...QueryCondition) ([]Row, error) {
	rows, err := db.core.Find(conditions...)
	return rows, wrapError(err)
}

func (db *DB) FindOne(conditions ...QueryCondition) (Row, bool, error) {
	row, ok, err := db.core.FindOne(conditions...)
	return row, ok, wrapError(err)
}

func (db *DB) RemoveByCondition(conditions ...QueryCondition) (int, error) {
	count, err := db.core.RemoveByCondition(conditions...)
	return count, wrapError(err)
}

func (db *DB) RemoveOneByCondition(conditions ...QueryCondition) (bool, error) {
	ok, err := db.core.RemoveOneByCondition(conditions...)
	return ok, wrapError(err)
}

func (db *DB) Put(key string, value []byte) error { return wrapError(db.core.Put(key, value)) }

func (db *DB) Get(key string) ([]byte, bool, error) {
	value, ok, err := db.core.Get(key)
	return value, ok, wrapError(err)
}

func (db *DB) Update(key string, value []byte) error { return wrapError(db.core.Update(key, value)) }

func (db *DB) Delete(key string) error { return wrapError(db.core.Delete(key)) }

func (db *DB) Query(prefix string) (map[string][]byte, error) {
	rows, err := db.core.Query(prefix)
	return rows, wrapError(err)
}

func (db *DB) Keys() ([]string, error) {
	keys, err := db.core.Keys()
	return keys, wrapError(err)
}

func (db *DB) QueryCascadeJSON(query CascadeQuery) ([]byte, error) {
	data, err := db.core.QueryCascadeJSON(query)
	return data, wrapError(err)
}

func (db *DB) QueryCascade(query CascadeQuery) ([]map[string]any, error) {
	rows, err := db.core.QueryCascade(query)
	return rows, wrapError(err)
}

func (db *DB) BeginTx() (*Tx, error) {
	tx, err := db.core.BeginTx()
	if err != nil {
		return nil, wrapError(err)
	}
	return &Tx{core: tx}, nil
}

func (db *DB) BeginReadOnlyTx() (*Tx, error) {
	tx, err := db.core.BeginReadOnlyTx()
	if err != nil {
		return nil, wrapError(err)
	}
	return &Tx{core: tx}, nil
}

func (db *DB) BeginTxWithOptions(options TxOptions) (*Tx, error) {
	isolation := kernal.TxIsolationSnapshot
	if options.Isolation == TxIsolationReadCommitted {
		isolation = kernal.TxIsolationReadCommitted
	}
	tx, err := db.core.BeginTxWithOptions(kernal.TxOptions{ReadOnly: options.ReadOnly, Isolation: isolation})
	if err != nil {
		return nil, wrapError(err)
	}
	return &Tx{core: tx}, nil
}

func (db *DB) WithTx(fn func(tx *Tx) error) error {
	return wrapError(db.core.WithTx(func(coreTx *kernal.Tx) error {
		if err := fn(&Tx{core: coreTx}); err != nil {
			return wrapError(err)
		}
		return nil
	}))
}
