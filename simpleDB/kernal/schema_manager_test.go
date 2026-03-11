package kernal

import (
	"os"
	"path/filepath"
	"testing"
)

// ─── 测试辅助 ─────────────────────────────────────────────────────────────────

func newSchemaManagerDB(t *testing.T) (*SimpleDB, func()) {
	t.Helper()
	dir := t.TempDir()
	db, err := New.DB(dir, "sm_test")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db, func() { _ = db.Close() }
}

func userSchema() TableSchema {
	return TableSchema{
		PrimaryKey:    "id",
		AutoIncrement: true,
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "string", Required: true},
			{Name: "email", Type: "string", Unique: true},
		},
	}
}

// ─── HasSchema ────────────────────────────────────────────────────────────────

func TestSchemaManager_HasSchema(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if db.HasSchema() {
		t.Fatal("expected no schema initially")
	}

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	if !db.HasSchema() {
		t.Fatal("expected schema after CreateTable")
	}
}

// ─── CreateTable ──────────────────────────────────────────────────────────────

func TestSchemaManager_CreateTable_Basic(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	schema, err := db.GetSchema()
	if err != nil {
		t.Fatalf("GetSchema: %v", err)
	}
	if schema.PrimaryKey != "id" {
		t.Fatalf("expected PrimaryKey=id, got %s", schema.PrimaryKey)
	}
	if len(schema.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(schema.Columns))
	}
}

func TestSchemaManager_CreateTable_AlreadyExists(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("first CreateTable: %v", err)
	}

	// 即使相同 schema 也应返回 ErrSchemaAlreadyExists（严格 DDL 语义）
	err := db.CreateTable(userSchema())
	if err != ErrSchemaAlreadyExists {
		t.Fatalf("expected ErrSchemaAlreadyExists, got %v", err)
	}
}

// ─── DropTable ────────────────────────────────────────────────────────────────

func TestSchemaManager_DropTable(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// 插入几行
	for i := 0; i < 3; i++ {
		if _, err := db.InsertRow(Row{"name": "user", "email": "u@example.com" + string(rune('0'+i))}); err != nil {
			t.Fatalf("InsertRow: %v", err)
		}
	}

	if err := db.DropTable(); err != nil {
		t.Fatalf("DropTable: %v", err)
	}

	if db.HasSchema() {
		t.Fatal("schema should be gone after DropTable")
	}

	// 再次 CreateTable 应成功
	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable after DropTable: %v", err)
	}
	if !db.HasSchema() {
		t.Fatal("expected schema after re-create")
	}
}

func TestSchemaManager_DropTable_NoSchema(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	err := db.DropTable()
	if err != ErrSchemaNotConfigured {
		t.Fatalf("expected ErrSchemaNotConfigured, got %v", err)
	}
}

// ─── TruncateTable ────────────────────────────────────────────────────────────

func TestSchemaManager_TruncateTable(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// 插入 5 行
	for i := 0; i < 5; i++ {
		if _, err := db.InsertRow(Row{
			"name":  "user",
			"email": "u@example.com" + string(rune('a'+i)),
		}); err != nil {
			t.Fatalf("InsertRow: %v", err)
		}
	}

	rows, err := db.Find()
	if err != nil {
		t.Fatalf("Find before truncate: %v", err)
	}
	if len(rows) != 5 {
		t.Fatalf("expected 5 rows, got %d", len(rows))
	}

	if err := db.TruncateTable(); err != nil {
		t.Fatalf("TruncateTable: %v", err)
	}

	// Schema 仍在
	if !db.HasSchema() {
		t.Fatal("schema should survive TruncateTable")
	}

	// 行数据应为空
	rows, err = db.Find()
	if err != nil {
		t.Fatalf("Find after truncate: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected 0 rows after truncate, got %d", len(rows))
	}

	// 自增序列应重置：下一行 id 应从 1 开始
	row, err := db.InsertRow(Row{"name": "new", "email": "new@example.com"})
	if err != nil {
		t.Fatalf("InsertRow after truncate: %v", err)
	}
	if id, _ := row["id"].(int64); id != 1 {
		t.Fatalf("expected id=1 after truncate reset, got %v", row["id"])
	}
}

func TestSchemaManager_TruncateTable_NoSchema(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	err := db.TruncateTable()
	if err != ErrSchemaNotConfigured {
		t.Fatalf("expected ErrSchemaNotConfigured, got %v", err)
	}
}

// ─── AlterTable: AddColumn ────────────────────────────────────────────────────

func TestSchemaManager_AlterTable_AddColumn(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// 先插入一行
	row1, err := db.InsertRow(Row{"name": "alice", "email": "alice@example.com"})
	if err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	// 添加 age 列，默认值 0
	err = db.AlterTable(AlterTablePlan{
		AddColumns: []Column{
			{Name: "age", Type: "int", Default: float64(0)},
		},
	})
	if err != nil {
		t.Fatalf("AlterTable AddColumn: %v", err)
	}

	// 旧行应有 age=0
	existing, found, err := db.FindRow(row1["id"])
	if err != nil {
		t.Fatalf("FindRow: %v", err)
	}
	if !found {
		t.Fatal("row not found after AlterTable")
	}
	if existing["age"] == nil {
		t.Fatal("expected age field to be backfilled with default 0")
	}

	// 新行可以携带 age
	row2, err := db.InsertRow(Row{"name": "bob", "email": "bob@example.com", "age": float64(25)})
	if err != nil {
		t.Fatalf("InsertRow with new column: %v", err)
	}
	if row2["age"] == nil {
		t.Fatal("expected age in new row")
	}

	// Schema 应有 4 列
	schema, _ := db.GetSchema()
	if len(schema.Columns) != 4 {
		t.Fatalf("expected 4 columns after AddColumn, got %d", len(schema.Columns))
	}
}

func TestSchemaManager_AlterTable_AddColumn_AlreadyExists(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	err := db.AlterTable(AlterTablePlan{
		AddColumns: []Column{{Name: "name", Type: "string"}},
	})
	if err == nil {
		t.Fatal("expected error when adding existing column")
	}
}

// ─── AlterTable: DropColumn ───────────────────────────────────────────────────

func TestSchemaManager_AlterTable_DropColumn(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	row1, err := db.InsertRow(Row{"name": "alice", "email": "alice@example.com"})
	if err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	// 删除 email 列
	err = db.AlterTable(AlterTablePlan{DropColumns: []string{"email"}})
	if err != nil {
		t.Fatalf("AlterTable DropColumn: %v", err)
	}

	// 行中不再有 email 字段
	existing, found, err := db.FindRow(row1["id"])
	if err != nil {
		t.Fatalf("FindRow: %v", err)
	}
	if !found {
		t.Fatal("row not found after DropColumn")
	}
	if _, exists := existing["email"]; exists {
		t.Fatal("email field should be removed from row after DropColumn")
	}

	// Schema 应有 2 列
	schema, _ := db.GetSchema()
	if len(schema.Columns) != 2 {
		t.Fatalf("expected 2 columns after DropColumn, got %d", len(schema.Columns))
	}
}

func TestSchemaManager_AlterTable_DropPrimaryKey_Error(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	err := db.AlterTable(AlterTablePlan{DropColumns: []string{"id"}})
	if err == nil {
		t.Fatal("expected error when dropping primary key column")
	}
}

// ─── AlterTable: AddIndex / DropIndex / AddUnique / DropUnique ────────────────

func TestSchemaManager_AlterTable_IndexOps(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	schema := TableSchema{
		PrimaryKey:    "id",
		AutoIncrement: true,
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
		},
	}
	if err := db.CreateTable(schema); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// AddIndex on name
	if err := db.AlterTable(AlterTablePlan{AddIndexes: []string{"name"}}); err != nil {
		t.Fatalf("AddIndex: %v", err)
	}
	s, _ := db.GetSchema()
	col := findColumnByName(s, "name")
	if col == nil || !col.Indexed {
		t.Fatal("expected name to be indexed after AddIndex")
	}

	// AddUnique on email
	if err := db.AlterTable(AlterTablePlan{AddUniques: []string{"email"}}); err != nil {
		t.Fatalf("AddUnique: %v", err)
	}
	s, _ = db.GetSchema()
	col = findColumnByName(s, "email")
	if col == nil || !col.Unique {
		t.Fatal("expected email to be unique after AddUnique")
	}
	if !col.Indexed {
		t.Fatal("expected email to be indexed (implicit) after AddUnique")
	}

	// DropIndex on name
	if err := db.AlterTable(AlterTablePlan{DropIndexes: []string{"name"}}); err != nil {
		t.Fatalf("DropIndex: %v", err)
	}
	s, _ = db.GetSchema()
	col = findColumnByName(s, "name")
	if col == nil || col.Indexed {
		t.Fatal("expected name to be NOT indexed after DropIndex")
	}

	// DropUnique on email
	if err := db.AlterTable(AlterTablePlan{DropUniques: []string{"email"}}); err != nil {
		t.Fatalf("DropUnique: %v", err)
	}
	s, _ = db.GetSchema()
	col = findColumnByName(s, "email")
	if col == nil || col.Unique {
		t.Fatal("expected email to be NOT unique after DropUnique")
	}
}

// ─── AlterTable: 重启后 Schema 持久化 ────────────────────────────────────────

func TestSchemaManager_AlterTable_Persistence(t *testing.T) {
	dir := t.TempDir()

	// 第一次打开：建表、改表
	db1, err := New.DB(dir, "sm_persist")
	if err != nil {
		t.Fatalf("open db1: %v", err)
	}
	if err = db1.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}
	if err = db1.AlterTable(AlterTablePlan{
		AddColumns: []Column{{Name: "age", Type: "int", Default: float64(18)}},
	}); err != nil {
		t.Fatalf("AlterTable: %v", err)
	}
	if _, err = db1.InsertRow(Row{"name": "carol", "email": "carol@example.com"}); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}
	_ = db1.Close()

	// 第二次打开：验证 schema 持久化
	db2, err := New.DB(dir, "sm_persist")
	if err != nil {
		t.Fatalf("open db2: %v", err)
	}
	defer func() { _ = db2.Close() }()

	s, err := db2.GetSchema()
	if err != nil {
		t.Fatalf("GetSchema after reopen: %v", err)
	}
	if findColumnByName(s, "age") == nil {
		t.Fatal("age column should persist after reopen")
	}

	rows, err := db2.Find()
	if err != nil {
		t.Fatalf("Find after reopen: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
}

// ─── AlterTable: EmptyPlan ────────────────────────────────────────────────────

func TestSchemaManager_AlterTable_EmptyPlan(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	err := db.AlterTable(AlterTablePlan{})
	if err == nil {
		t.Fatal("expected error for empty AlterTablePlan")
	}
}

// ─── DropTable + CreateTable 生命周期 ─────────────────────────────────────────

func TestSchemaManager_FullLifecycle(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	// Create
	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// Insert
	if _, err := db.InsertRow(Row{"name": "test", "email": "test@test.com"}); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	// Truncate
	if err := db.TruncateTable(); err != nil {
		t.Fatalf("TruncateTable: %v", err)
	}
	rows, _ := db.Find()
	if len(rows) != 0 {
		t.Fatalf("expected empty after truncate, got %d rows", len(rows))
	}

	// Insert again after truncate
	row, err := db.InsertRow(Row{"name": "v2", "email": "v2@v2.com"})
	if err != nil {
		t.Fatalf("InsertRow after truncate: %v", err)
	}
	if id, _ := row["id"].(int64); id != 1 {
		t.Fatalf("expected id reset to 1, got %v", row["id"])
	}

	// Drop
	if err := db.DropTable(); err != nil {
		t.Fatalf("DropTable: %v", err)
	}
	if db.HasSchema() {
		t.Fatal("HasSchema should be false after DropTable")
	}

	// Re-create with different schema
	newSchema := TableSchema{
		PrimaryKey:    "uid",
		AutoIncrement: true,
		Columns: []Column{
			{Name: "uid", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "title", Type: "string"},
		},
	}
	if err := db.CreateTable(newSchema); err != nil {
		t.Fatalf("CreateTable with new schema: %v", err)
	}
	s, _ := db.GetSchema()
	if s.PrimaryKey != "uid" {
		t.Fatalf("expected PrimaryKey=uid in new schema, got %s", s.PrimaryKey)
	}
}

// ─── 辅助 ────────────────────────────────────────────────────────────────────

func findColumnByName(schema *TableSchema, name string) *Column {
	if schema == nil {
		return nil
	}
	for i := range schema.Columns {
		if schema.Columns[i].Name == name {
			return &schema.Columns[i]
		}
	}
	return nil
}

// 确保测试不依赖文件系统残留（已通过 t.TempDir() 保证隔离）
var _ = filepath.Join
var _ = os.TempDir

// ─── SchemaDiff ───────────────────────────────────────────────────────────────

func TestSchemaManager_SchemaDiff_NoSchema(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	plan, exists, err := db.SchemaDiff(userSchema())
	if err != nil {
		t.Fatalf("SchemaDiff: %v", err)
	}
	if exists {
		t.Fatal("expected exists=false when no schema")
	}
	if plan != nil {
		t.Fatal("expected nil plan when no schema")
	}
}

func TestSchemaManager_SchemaDiff_Identical(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	plan, exists, err := db.SchemaDiff(userSchema())
	if err != nil {
		t.Fatalf("SchemaDiff: %v", err)
	}
	if !exists {
		t.Fatal("expected exists=true")
	}
	if plan != nil {
		t.Fatalf("expected nil plan for identical schema, got %+v", plan)
	}
}

func TestSchemaManager_SchemaDiff_AddColumn(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	target := userSchema()
	target.Columns = append(target.Columns, Column{Name: "age", Type: "int"})

	plan, exists, err := db.SchemaDiff(target)
	if err != nil {
		t.Fatalf("SchemaDiff: %v", err)
	}
	if !exists {
		t.Fatal("expected exists=true")
	}
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if len(plan.AddColumns) != 1 || plan.AddColumns[0].Name != "age" {
		t.Fatalf("expected AddColumns=[age], got %+v", plan.AddColumns)
	}
	if len(plan.DropColumns) != 0 {
		t.Fatalf("expected no DropColumns, got %+v", plan.DropColumns)
	}
}

func TestSchemaManager_SchemaDiff_DropColumn(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.CreateTable(userSchema()); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// target 移除了 email 列
	target := TableSchema{
		PrimaryKey:    "id",
		AutoIncrement: true,
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "string", Required: true},
		},
	}

	plan, _, err := db.SchemaDiff(target)
	if err != nil {
		t.Fatalf("SchemaDiff: %v", err)
	}
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if len(plan.DropColumns) != 1 || plan.DropColumns[0] != "email" {
		t.Fatalf("expected DropColumns=[email], got %+v", plan.DropColumns)
	}
}

func TestSchemaManager_SchemaDiff_IndexChange(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	base := TableSchema{
		PrimaryKey: "id",
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true},
			{Name: "name", Type: "string", Indexed: true},
			{Name: "email", Type: "string", Unique: true},
		},
	}
	if err := db.CreateTable(base); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// target: name 去掉索引，email 去掉唯一，加上 phone 唯一索引
	target := TableSchema{
		PrimaryKey: "id",
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true},
			{Name: "name", Type: "string"},                // 去掉 indexed
			{Name: "email", Type: "string"},               // 去掉 unique
			{Name: "phone", Type: "string", Unique: true}, // 新列+唯一
		},
	}

	plan, _, err := db.SchemaDiff(target)
	if err != nil {
		t.Fatalf("SchemaDiff: %v", err)
	}
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if len(plan.DropIndexes) != 1 || plan.DropIndexes[0] != "name" {
		t.Fatalf("expected DropIndexes=[name], got %v", plan.DropIndexes)
	}
	if len(plan.DropUniques) != 1 || plan.DropUniques[0] != "email" {
		t.Fatalf("expected DropUniques=[email], got %v", plan.DropUniques)
	}
	if len(plan.AddColumns) != 1 || plan.AddColumns[0].Name != "phone" {
		t.Fatalf("expected AddColumns=[phone], got %v", plan.AddColumns)
	}
}

// ─── AutoMigrate ──────────────────────────────────────────────────────────────

func TestSchemaManager_AutoMigrate_CreateNew(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	// 无 Schema 时 AutoMigrate 等同 CreateTable
	if err := db.AutoMigrate(userSchema()); err != nil {
		t.Fatalf("AutoMigrate (create): %v", err)
	}
	if !db.HasSchema() {
		t.Fatal("expected schema after AutoMigrate")
	}
}

func TestSchemaManager_AutoMigrate_Idempotent(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.AutoMigrate(userSchema()); err != nil {
		t.Fatalf("first AutoMigrate: %v", err)
	}
	// 多次调用相同 schema 应幂等
	if err := db.AutoMigrate(userSchema()); err != nil {
		t.Fatalf("second AutoMigrate (idempotent): %v", err)
	}
	s, _ := db.GetSchema()
	if len(s.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(s.Columns))
	}
}

func TestSchemaManager_AutoMigrate_AddOnly(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.AutoMigrate(userSchema()); err != nil {
		t.Fatalf("first AutoMigrate: %v", err)
	}

	// 插入一行
	if _, err := db.InsertRow(Row{"name": "alice", "email": "alice@example.com"}); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	// 新版 schema：增加 age、phone 列，但移除 email（AutoMigrate 不应删 email）
	newSchema := TableSchema{
		PrimaryKey:    "id",
		AutoIncrement: true,
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "string", Required: true},
			{Name: "age", Type: "int", Default: float64(0)},
			{Name: "phone", Type: "string"},
		},
	}
	if err := db.AutoMigrate(newSchema); err != nil {
		t.Fatalf("AutoMigrate (add only): %v", err)
	}

	s, _ := db.GetSchema()
	// email 应被保留（AutoMigrate 不删列），加上 age、phone = 5 列
	if len(s.Columns) != 5 {
		t.Fatalf("expected 5 columns (email preserved + age + phone added), got %d", len(s.Columns))
	}
	if findColumnByName(s, "email") == nil {
		t.Fatal("email column should NOT be dropped by AutoMigrate")
	}
	if findColumnByName(s, "age") == nil {
		t.Fatal("age column should be added")
	}
}

// ─── SyncSchema ───────────────────────────────────────────────────────────────

func TestSchemaManager_SyncSchema_CreateNew(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.SyncSchema(userSchema()); err != nil {
		t.Fatalf("SyncSchema (create): %v", err)
	}
	if !db.HasSchema() {
		t.Fatal("expected schema after SyncSchema")
	}
}

func TestSchemaManager_SyncSchema_Idempotent(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.SyncSchema(userSchema()); err != nil {
		t.Fatalf("first SyncSchema: %v", err)
	}
	if err := db.SyncSchema(userSchema()); err != nil {
		t.Fatalf("second SyncSchema (idempotent): %v", err)
	}
}

func TestSchemaManager_SyncSchema_AddAndDrop(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	if err := db.SyncSchema(userSchema()); err != nil {
		t.Fatalf("SyncSchema initial: %v", err)
	}

	// 插入一行
	row, err := db.InsertRow(Row{"name": "bob", "email": "bob@example.com"})
	if err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	// 新版 schema：移除 email，增加 age
	newSchema := TableSchema{
		PrimaryKey:    "id",
		AutoIncrement: true,
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "name", Type: "string", Required: true},
			{Name: "age", Type: "int", Default: float64(18)},
		},
	}
	if err := db.SyncSchema(newSchema); err != nil {
		t.Fatalf("SyncSchema update: %v", err)
	}

	s, _ := db.GetSchema()
	if len(s.Columns) != 3 {
		t.Fatalf("expected 3 columns after SyncSchema, got %d", len(s.Columns))
	}
	if findColumnByName(s, "email") != nil {
		t.Fatal("email should be dropped by SyncSchema")
	}
	if findColumnByName(s, "age") == nil {
		t.Fatal("age should be added by SyncSchema")
	}

	// 旧行的 email 字段应已删除，age 应回填默认值
	existing, found, err := db.FindRow(row["id"])
	if err != nil || !found {
		t.Fatalf("FindRow: %v, found=%v", err, found)
	}
	if _, emailExists := existing["email"]; emailExists {
		t.Fatal("email field should be removed from existing row")
	}
	if existing["age"] == nil {
		t.Fatal("age should be backfilled with default value")
	}
}

func TestSchemaManager_SyncSchema_Persistence(t *testing.T) {
	dir := t.TempDir()

	db1, err := New.DB(dir, "sync_persist")
	if err != nil {
		t.Fatalf("open db1: %v", err)
	}
	if err = db1.SyncSchema(userSchema()); err != nil {
		t.Fatalf("SyncSchema: %v", err)
	}
	// 追加一列
	evolved := userSchema()
	evolved.Columns = append(evolved.Columns, Column{Name: "score", Type: "float", Default: float64(0)})
	if err = db1.SyncSchema(evolved); err != nil {
		t.Fatalf("SyncSchema evolved: %v", err)
	}
	_ = db1.Close()

	// 重开：验证 score 列持久化
	db2, err := New.DB(dir, "sync_persist")
	if err != nil {
		t.Fatalf("open db2: %v", err)
	}
	defer func() { _ = db2.Close() }()

	s, err := db2.GetSchema()
	if err != nil {
		t.Fatalf("GetSchema: %v", err)
	}
	if findColumnByName(s, "score") == nil {
		t.Fatal("score column should persist after reopen")
	}
}

func TestSchemaManager_SchemaDiff_ForeignKeys(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	base := TableSchema{
		PrimaryKey: "id",
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true},
			{Name: "user_id", Type: "int"},
		},
	}
	if err := db.CreateTable(base); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	target := cloneSchema(base)
	target.ForeignKeys = []ForeignKey{{Name: "fk_user", Field: "user_id", RefTable: "users", RefField: "id", Alias: "user"}}

	plan, exists, err := db.SchemaDiff(target)
	if err != nil {
		t.Fatalf("SchemaDiff add fk: %v", err)
	}
	if !exists || plan == nil {
		t.Fatalf("expected non-nil plan with schema exists, got exists=%v plan=%+v", exists, plan)
	}
	if len(plan.AddForeignKeys) != 1 || plan.AddForeignKeys[0].Name != "fk_user" {
		t.Fatalf("expected AddForeignKeys=[fk_user], got %+v", plan.AddForeignKeys)
	}

	if err = db.SyncSchema(target); err != nil {
		t.Fatalf("SyncSchema target: %v", err)
	}

	plan, _, err = db.SchemaDiff(base)
	if err != nil {
		t.Fatalf("SchemaDiff drop fk: %v", err)
	}
	if plan == nil || len(plan.DropForeignKeys) != 1 {
		t.Fatalf("expected DropForeignKeys size=1, got %+v", plan)
	}
}

func TestSchemaManager_AlterTable_ForeignKeys(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	schema := TableSchema{
		PrimaryKey: "id",
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true},
			{Name: "user_id", Type: "int"},
		},
	}
	if err := db.CreateTable(schema); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	if err := db.AlterTable(AlterTablePlan{
		AddForeignKeys: []ForeignKey{{Name: "fk_user", Field: "user_id", RefTable: "users", RefField: "id", Alias: "user"}},
	}); err != nil {
		t.Fatalf("AlterTable AddForeignKeys: %v", err)
	}

	current, err := db.GetSchema()
	if err != nil {
		t.Fatalf("GetSchema: %v", err)
	}
	if len(current.ForeignKeys) != 1 || current.ForeignKeys[0].Name != "fk_user" {
		t.Fatalf("expected 1 foreign key fk_user, got %+v", current.ForeignKeys)
	}

	if err := db.AlterTable(AlterTablePlan{DropForeignKeys: []string{"fk_user"}}); err != nil {
		t.Fatalf("AlterTable DropForeignKeys: %v", err)
	}

	current, err = db.GetSchema()
	if err != nil {
		t.Fatalf("GetSchema: %v", err)
	}
	if len(current.ForeignKeys) != 0 {
		t.Fatalf("expected foreign keys cleared, got %+v", current.ForeignKeys)
	}
}

func TestSchemaManager_AutoMigrate_KeepForeignKeys(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	withFK := TableSchema{
		PrimaryKey: "id",
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true},
			{Name: "user_id", Type: "int"},
		},
		ForeignKeys: []ForeignKey{{Name: "fk_user", Field: "user_id", RefTable: "users", RefField: "id"}},
	}
	if err := db.AutoMigrate(withFK); err != nil {
		t.Fatalf("AutoMigrate create with FK: %v", err)
	}

	withoutFK := cloneSchema(withFK)
	withoutFK.ForeignKeys = nil
	if err := db.AutoMigrate(withoutFK); err != nil {
		t.Fatalf("AutoMigrate with target no FK: %v", err)
	}

	current, err := db.GetSchema()
	if err != nil {
		t.Fatalf("GetSchema: %v", err)
	}
	if len(current.ForeignKeys) != 1 {
		t.Fatalf("expected AutoMigrate keep foreign key, got %+v", current.ForeignKeys)
	}
}

func TestSchemaManager_SyncSchema_DropForeignKeys(t *testing.T) {
	db, teardown := newSchemaManagerDB(t)
	defer teardown()

	withFK := TableSchema{
		PrimaryKey: "id",
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true},
			{Name: "user_id", Type: "int"},
		},
		ForeignKeys: []ForeignKey{{Name: "fk_user", Field: "user_id", RefTable: "users", RefField: "id"}},
	}
	if err := db.SyncSchema(withFK); err != nil {
		t.Fatalf("SyncSchema with FK: %v", err)
	}

	withoutFK := cloneSchema(withFK)
	withoutFK.ForeignKeys = nil
	if err := db.SyncSchema(withoutFK); err != nil {
		t.Fatalf("SyncSchema drop FK: %v", err)
	}

	current, err := db.GetSchema()
	if err != nil {
		t.Fatalf("GetSchema: %v", err)
	}
	if len(current.ForeignKeys) != 0 {
		t.Fatalf("expected SyncSchema to drop foreign keys, got %+v", current.ForeignKeys)
	}
}
