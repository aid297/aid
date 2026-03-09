package simpleDBDriver

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aid297/aid/ptr"
	"github.com/google/uuid"
)

func TestSimpleDB_CRUDAndReload(t *testing.T) {
	database := "demo"
	table := "users"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("NewSimpleDB() error = %v", err)
	}

	if err = db.Put("user:1", []byte(`{"name":"alice"}`)); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	value, ok, err := db.Get("user:1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok || string(value) != `{"name":"alice"}` {
		t.Fatalf("Get() = (%q, %v), want alice payload", string(value), ok)
	}

	if err = db.Update("user:1", []byte(`{"name":"bob"}`)); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	filtered, err := db.Query("user:")
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got := string(filtered["user:1"]); got != `{"name":"bob"}` {
		t.Fatalf("Query() value = %q, want bob payload", got)
	}

	if err = db.Delete("user:1"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, ok, err = db.Get("user:1")
	if err != nil {
		t.Fatalf("Get() after delete error = %v", err)
	}
	if ok {
		t.Fatal("Get() after delete should not find key")
	}

	if err = db.Put("user:2", []byte(`{"name":"carol"}`)); err != nil {
		t.Fatalf("Put() second key error = %v", err)
	}

	if err = db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reloaded, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("reload NewSimpleDB() error = %v", err)
	}
	defer reloaded.Close()

	_, ok, err = reloaded.Get("user:1")
	if err != nil {
		t.Fatalf("reloaded Get(user:1) error = %v", err)
	}
	if ok {
		t.Fatal("deleted key should stay deleted after reload")
	}

	value, ok, err = reloaded.Get("user:2")
	if err != nil {
		t.Fatalf("reloaded Get(user:2) error = %v", err)
	}
	if !ok || string(value) != `{"name":"carol"}` {
		t.Fatalf("reloaded Get(user:2) = (%q, %v)", string(value), ok)
	}
}

func TestSimpleDB_Compact(t *testing.T) {
	database := "demo"
	table := "users"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("NewSimpleDB() error = %v", err)
	}
	defer db.Close()

	if err = db.Put("doc:1", []byte("v1")); err != nil {
		t.Fatalf("Put v1 error = %v", err)
	}
	if err = db.Update("doc:1", []byte("v2")); err != nil {
		t.Fatalf("Update v2 error = %v", err)
	}
	if err = db.Put("doc:2", []byte("trash")); err != nil {
		t.Fatalf("Put doc:2 error = %v", err)
	}
	if err = db.Delete("doc:2"); err != nil {
		t.Fatalf("Delete doc:2 error = %v", err)
	}

	before, err := os.Stat(filepath.Join(database, table, table+defaultDBFileEx))
	if err != nil {
		t.Fatalf("Stat before compact error = %v", err)
	}

	if err = db.Compact(); err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	after, err := os.Stat(filepath.Join(database, table, table+defaultDBFileEx))
	if err != nil {
		t.Fatalf("Stat after compact error = %v", err)
	}
	if after.Size() >= before.Size() {
		t.Fatalf("compact should shrink file, before=%d after=%d", before.Size(), after.Size())
	}

	value, ok, err := db.Get("doc:1")
	if err != nil {
		t.Fatalf("Get(doc:1) error = %v", err)
	}
	if !ok || string(value) != "v2" {
		t.Fatalf("Get(doc:1) = (%q, %v), want v2", string(value), ok)
	}

	_, ok, err = db.Get("doc:2")
	if err != nil {
		t.Fatalf("Get(doc:2) error = %v", err)
	}
	if ok {
		t.Fatal("deleted key should not reappear after compact")
	}
}

func TestSimpleDB_ErrorCases(t *testing.T) {
	database := "demo"
	table := "users"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("NewSimpleDB() error = %v", err)
	}

	if err = db.Update("missing", []byte("value")); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Update missing error = %v, want ErrKeyNotFound", err)
	}

	if err = db.Delete("missing"); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Delete missing error = %v, want ErrKeyNotFound", err)
	}

	if err = db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	if err = db.Put("closed", []byte("value")); !errors.Is(err, ErrDatabaseClosed) {
		t.Fatalf("Put after close error = %v, want ErrDatabaseClosed", err)
	}
}

func TestSimpleDB_FileLock(t *testing.T) {
	database := "demo"
	table := "users"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("first newSimpleDB() error = %v", err)
	}
	defer db.Close()

	second, err := newSimpleDB(database, table)
	if second != nil {
		_ = second.Close()
	}
	if !errors.Is(err, ErrDatabaseLocked) {
		t.Fatalf("second newSimpleDB() error = %v, want ErrDatabaseLocked", err)
	}
}

func TestSimpleDB_TableConstraints(t *testing.T) {
	database := "demo"
	table := "accounts"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "email", Type: "string", Unique: true},
		{Name: "age", Type: "int", Indexed: true},
		{Name: "name", Type: "string"},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	first, err := db.InsertRow(Row{"email": "alice@example.com", "age": 18, "name": "alice"})
	if err != nil {
		t.Fatalf("InsertRow(first) error = %v", err)
	}
	if got := first["id"]; got != int64(1) {
		t.Fatalf("first inserted id = %#v, want 1", got)
	}

	second, err := db.InsertRow(Row{"email": "bob@example.com", "age": 18, "name": "bob"})
	if err != nil {
		t.Fatalf("InsertRow(second) error = %v", err)
	}
	if got := second["id"]; got != int64(2) {
		t.Fatalf("second inserted id = %#v, want 2", got)
	}

	if _, err = db.InsertRow(Row{"email": "bob@example.com", "age": 20, "name": "dup"}); !errors.Is(err, ErrUniqueConflict) {
		t.Fatalf("duplicate email error = %v, want ErrUniqueConflict", err)
	}

	row, found, err := db.FindByUnique("email", "bob@example.com")
	if err != nil {
		t.Fatalf("FindByUnique() error = %v", err)
	}
	if !found || row["name"] != "bob" {
		t.Fatalf("FindByUnique() = (%#v, %v), want bob", row, found)
	}

	rows, err := db.FindByIndex("age", 18)
	if err != nil {
		t.Fatalf("FindByIndex(age=18) error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("FindByIndex(age=18) len = %d, want 2", len(rows))
	}

	updated, err := db.UpdateRow(int64(2), Row{"email": "bob+new@example.com", "age": 21})
	if err != nil {
		t.Fatalf("UpdateRow() error = %v", err)
	}
	if updated["email"] != "bob+new@example.com" {
		t.Fatalf("updated email = %#v, want bob+new@example.com", updated["email"])
	}

	if _, found, err = db.FindByUnique("email", "bob@example.com"); err != nil {
		t.Fatalf("FindByUnique(old email) error = %v", err)
	} else if found {
		t.Fatal("old unique value should not be found after update")
	}

	rows, err = db.FindByIndex("age", 18)
	if err != nil {
		t.Fatalf("FindByIndex(age=18 after update) error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("FindByIndex(age=18 after update) len = %d, want 1", len(rows))
	}

	if err = db.DeleteRow(int64(1)); err != nil {
		t.Fatalf("DeleteRow() error = %v", err)
	}

	if _, found, err = db.FindByUnique("email", "alice@example.com"); err != nil {
		t.Fatalf("FindByUnique(alice after delete) error = %v", err)
	} else if found {
		t.Fatal("deleted row should be removed from unique index")
	}

	if err = db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reloaded, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("reload newSimpleDB() error = %v", err)
	}
	defer reloaded.Close()

	reloadedRow, found, err := reloaded.FindRow(int64(2))
	if err != nil {
		t.Fatalf("FindRow(reloaded) error = %v", err)
	}
	if !found || reloadedRow["email"] != "bob+new@example.com" {
		t.Fatalf("reloaded row = %#v, found=%v", reloadedRow, found)
	}
	if got := reloadedRow["age"]; got != int64(21) {
		t.Fatalf("reloaded age = %#v, want 21", got)
	}

	third, err := reloaded.InsertRow(Row{"email": "carol@example.com", "age": 30, "name": "carol"})
	if err != nil {
		t.Fatalf("InsertRow(reloaded) error = %v", err)
	}
	if got := third["id"]; got != int64(4) {
		t.Fatalf("third inserted id = %#v, want 4", got)
	}

	if err = reloaded.Compact(); err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	rows, err = reloaded.FindByIndex("age", 30)
	if err != nil {
		t.Fatalf("FindByIndex(age=30 after compact) error = %v", err)
	}
	if len(rows) != 1 || rows[0]["email"] != "carol@example.com" {
		t.Fatalf("FindByIndex(age=30 after compact) = %#v", rows)
	}
}

func TestSimpleDB_SchemaValidation(t *testing.T) {
	database := "demo"
	table := "schema_validation"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	if _, err = db.GetSchema(); !errors.Is(err, ErrSchemaNotConfigured) {
		t.Fatalf("GetSchema() error = %v, want ErrSchemaNotConfigured", err)
	}

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", AutoIncrement: true},
		{Name: "email"},
	}})
	if !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("invalid Configure() error = %v, want ErrInvalidSchema", err)
	}

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true},
		{Name: "email", Type: "string", Unique: true},
	}})
	if err != nil {
		t.Fatalf("valid Configure() error = %v", err)
	}

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true},
		{Name: "email", Type: "string", Indexed: true},
	}})
	if !errors.Is(err, ErrSchemaAlreadyExists) {
		t.Fatalf("reconfigure error = %v, want ErrSchemaAlreadyExists", err)
	}
}

func TestSimpleDB_FieldTypeValidation(t *testing.T) {
	database := "demo"
	table := "typed_rows"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string", Unique: true},
		{Name: "age", Type: "int", Indexed: true},
		{Name: "score", Type: "float"},
		{Name: "active", Type: "bool"},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	row, err := db.InsertRow(Row{"name": "alice", "age": 18, "score": 95.5, "active": true})
	if err != nil {
		t.Fatalf("InsertRow(valid) error = %v", err)
	}
	if age, ok := asInt64(row["age"]); !ok || age != 18 {
		t.Fatalf("normalized age = %#v, want integer 18", row["age"])
	}

	if _, err = db.InsertRow(Row{"name": "bob", "age": "18", "score": 88.5, "active": true}); !errors.Is(err, ErrFieldTypeMismatch) {
		t.Fatalf("InsertRow(age string) error = %v, want ErrFieldTypeMismatch", err)
	}

	if _, err = db.UpdateRow(int64(1), Row{"active": "yes"}); !errors.Is(err, ErrFieldTypeMismatch) {
		t.Fatalf("UpdateRow(active string) error = %v, want ErrFieldTypeMismatch", err)
	}

	if _, err = db.FindByIndex("age", "18"); !errors.Is(err, ErrFieldTypeMismatch) {
		t.Fatalf("FindByIndex(age string) error = %v, want ErrFieldTypeMismatch", err)
	}

	if _, found, err := db.FindByUnique("name", "alice"); err != nil || !found {
		t.Fatalf("FindByUnique(valid) = (%v, %v), want found", err, found)
	}
}

func TestSimpleDB_UUIDAutoIncrementPrimaryKey(t *testing.T) {
	database := "demo"
	table := "uuid_auto_pk"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string", Unique: true},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	generated, err := db.InsertRow(Row{"name": "generated"})
	if err != nil {
		t.Fatalf("InsertRow(generated uuid) error = %v", err)
	}
	generatedID, ok := generated["id"].(string)
	if !ok {
		t.Fatalf("generated id type = %#v, want string", generated["id"])
	}
	if _, err = uuid.Parse(generatedID); err != nil {
		t.Fatalf("generated id parse error = %v", err)
	}

	versionSamples := []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8", // v1
		"6fa459ea-ee8a-3ca4-894e-db77e160355e", // v3
		"550e8400-e29b-41d4-a716-446655440000", // v4
		"886313e1-3b8a-5372-9b90-0c9aee199e5d", // v5
		"1ec3f7d0-7c9b-6c9a-8b5d-9f8e7d6c5b4a", // v6
		"01890f3e-3f67-7b5a-bf7d-1f4b9f2a7c10", // v7
		"12345678-1234-8abc-8def-1234567890ab", // v8
	}

	for index, id := range versionSamples {
		name := "preset-" + id
		row, insertErr := db.InsertRow(Row{"id": id, "name": name})
		if insertErr != nil {
			t.Fatalf("InsertRow(uuid version sample %d) error = %v", index+1, insertErr)
		}
		if row["id"] != id {
			t.Fatalf("stored uuid = %#v, want %s", row["id"], id)
		}
	}

	if _, err = db.InsertRow(Row{"id": "not-a-uuid", "name": "invalid"}); !errors.Is(err, ErrFieldTypeMismatch) {
		t.Fatalf("InsertRow(invalid uuid) error = %v, want ErrFieldTypeMismatch", err)
	}
}

func TestSimpleDB_ConditionalQueries(t *testing.T) {
	database := "demo"
	table := "conditional_queries"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string"},
		{Name: "age", Type: "int", Indexed: true},
		{Name: "status", Type: "string", Indexed: true},
		{Name: "created_at", Type: "timestamp", Indexed: true},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	rows := []Row{
		{"name": "alice", "age": 18, "status": "draft", "created_at": "2026-03-01T10:00:00Z"},
		{"name": "bob", "age": 25, "status": "active", "created_at": "2026-03-02T10:00:00Z"},
		{"name": "carol", "age": 31, "status": "done", "created_at": "2026-03-03T10:00:00Z"},
		{"name": "dave", "age": 40, "status": "active", "created_at": "2026-03-04T10:00:00Z"},
	}
	for _, row := range rows {
		if _, err = db.InsertRow(row); err != nil {
			t.Fatalf("InsertRow(%v) error = %v", row["name"], err)
		}
	}

	result, err := db.FindByConditions([]QueryCondition{{Field: "age", Operator: QueryOpGT, Value: 25}})
	if err != nil {
		t.Fatalf("FindByConditions(age > 25) error = %v", err)
	}
	if len(result) != 2 || result[0]["name"] != "carol" || result[1]["name"] != "dave" {
		t.Fatalf("FindByConditions(age > 25) = %#v", result)
	}

	result, err = db.FindByConditions([]QueryCondition{{Field: "status", Operator: QueryOpIn, Values: []any{"draft", "active"}}})
	if err != nil {
		t.Fatalf("FindByConditions(status in) error = %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("FindByConditions(status in) len = %d, want 3", len(result))
	}

	result, err = db.FindByConditions([]QueryCondition{{Field: "status", Operator: QueryOpNotIn, Values: []any{"draft", "done"}}})
	if err != nil {
		t.Fatalf("FindByConditions(status not in) error = %v", err)
	}
	if len(result) != 2 || result[0]["name"] != "bob" || result[1]["name"] != "dave" {
		t.Fatalf("FindByConditions(status not in) = %#v", result)
	}

	result, err = db.FindByConditions([]QueryCondition{{Field: "status", Operator: QueryOpNE, Value: "draft"}})
	if err != nil {
		t.Fatalf("FindByConditions(status != draft) error = %v", err)
	}
	if len(result) != 3 || result[0]["name"] != "bob" || result[1]["name"] != "carol" || result[2]["name"] != "dave" {
		t.Fatalf("FindByConditions(status != draft) = %#v", result)
	}

	result, err = db.FindByConditions([]QueryCondition{{Field: "age", Operator: QueryOpBetween, Lower: 20, Upper: 35}})
	if err != nil {
		t.Fatalf("FindByConditions(age between) error = %v", err)
	}
	if len(result) != 2 || result[0]["name"] != "bob" || result[1]["name"] != "carol" {
		t.Fatalf("FindByConditions(age between) = %#v", result)
	}

	result, err = db.FindByConditions([]QueryCondition{{Field: "age", Operator: QueryOpNotBetween, Lower: 20, Upper: 35}})
	if err != nil {
		t.Fatalf("FindByConditions(age not between) error = %v", err)
	}
	if len(result) != 2 || result[0]["name"] != "alice" || result[1]["name"] != "dave" {
		t.Fatalf("FindByConditions(age not between) = %#v", result)
	}

	result, err = db.FindByConditions([]QueryCondition{{Field: "age", Operator: QueryOpNotIn, Values: []any{18, 31}}})
	if err != nil {
		t.Fatalf("FindByConditions(age not in) error = %v", err)
	}
	if len(result) != 2 || result[0]["name"] != "bob" || result[1]["name"] != "dave" {
		t.Fatalf("FindByConditions(age not in) = %#v", result)
	}

	result, err = db.FindByConditions([]QueryCondition{
		{Field: "age", Operator: QueryOpGTE, Value: 25},
		{Field: "status", Operator: QueryOpEQ, Value: "active"},
		{Field: "created_at", Operator: QueryOpBetween, Lower: "2026-03-02T00:00:00Z", Upper: "2026-03-04T00:00:00Z"},
	})
	if err != nil {
		t.Fatalf("FindByConditions(combined) error = %v", err)
	}
	if len(result) != 1 || result[0]["name"] != "bob" {
		t.Fatalf("FindByConditions(combined) = %#v", result)
	}

	if _, err = db.FindByConditions([]QueryCondition{{Field: "age", Operator: QueryOpBetween, Lower: 30}}); !errors.Is(err, ErrInvalidQueryCondition) {
		t.Fatalf("FindByConditions(missing upper) error = %v, want ErrInvalidQueryCondition", err)
	}

	if _, err = db.FindByConditions([]QueryCondition{{Field: "age", Operator: QueryOpGT, Value: "18"}}); !errors.Is(err, ErrFieldTypeMismatch) {
		t.Fatalf("FindByConditions(type mismatch) error = %v, want ErrFieldTypeMismatch", err)
	}

	if _, err = db.FindByConditions([]QueryCondition{{Field: "status", Operator: "unknown", Value: "active"}}); !errors.Is(err, ErrInvalidQueryCondition) {
		t.Fatalf("FindByConditions(unknown operator) error = %v, want ErrInvalidQueryCondition", err)
	}

	result, err = db.FindByConditions([]QueryCondition{{Field: "name", Operator: QueryOpEQ, Value: "alice"}})
	if err != nil {
		t.Fatalf("FindByConditions(non-indexed field fallback) error = %v", err)
	}
	if len(result) != 1 || result[0]["name"] != "alice" {
		t.Fatalf("FindByConditions(non-indexed field fallback) = %#v", result)
	}
}

func TestSimpleDB_ForeignKeyCascadeJSON(t *testing.T) {
	database := "demo"
	for _, table := range []string{"users_rel", "orders_rel", "order_items_rel"} {
		_ = os.RemoveAll(filepath.Join(database, table))
	}
	t.Cleanup(func() {
		for _, table := range []string{"users_rel", "orders_rel", "order_items_rel"} {
			_ = os.RemoveAll(filepath.Join(database, table))
		}
	})

	users, err := newSimpleDB(database, "users_rel")
	if err != nil {
		t.Fatalf("newSimpleDB(users_rel) error = %v", err)
	}
	if err = users.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string"},
		{Name: "age", Type: "int", Indexed: true},
	}}); err != nil {
		t.Fatalf("Configure(users_rel) error = %v", err)
	}
	if _, err = users.InsertRow(Row{"name": "alice", "age": 30}); err != nil {
		t.Fatalf("InsertRow(users_rel alice) error = %v", err)
	}
	if _, err = users.InsertRow(Row{"name": "bob", "age": 22}); err != nil {
		t.Fatalf("InsertRow(users_rel bob) error = %v", err)
	}
	if err = users.Close(); err != nil {
		t.Fatalf("Close(users_rel) error = %v", err)
	}

	orders, err := newSimpleDB(database, "orders_rel")
	if err != nil {
		t.Fatalf("newSimpleDB(orders_rel) error = %v", err)
	}
	if err = orders.Configure(TableSchema{
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "user_id", Type: "int"},
			{Name: "amount", Type: "float"},
			{Name: "status", Type: "string"},
		},
		ForeignKeys: []ForeignKey{{Name: "fk_orders_user", Field: "user_id", RefTable: "users_rel", RefField: "id", Alias: "orders"}},
	}); err != nil {
		t.Fatalf("Configure(orders_rel) error = %v", err)
	}
	if _, err = orders.InsertRow(Row{"user_id": 1, "amount": 120.5, "status": "paid"}); err != nil {
		t.Fatalf("InsertRow(orders_rel 1) error = %v", err)
	}
	if _, err = orders.InsertRow(Row{"user_id": 1, "amount": 80.0, "status": "draft"}); err != nil {
		t.Fatalf("InsertRow(orders_rel 2) error = %v", err)
	}
	if _, err = orders.InsertRow(Row{"user_id": 2, "amount": 200.0, "status": "paid"}); err != nil {
		t.Fatalf("InsertRow(orders_rel 3) error = %v", err)
	}
	if err = orders.Close(); err != nil {
		t.Fatalf("Close(orders_rel) error = %v", err)
	}

	items, err := newSimpleDB(database, "order_items_rel")
	if err != nil {
		t.Fatalf("newSimpleDB(order_items_rel) error = %v", err)
	}
	if err = items.Configure(TableSchema{
		Columns: []Column{
			{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
			{Name: "order_id", Type: "int"},
			{Name: "sku", Type: "string"},
			{Name: "quantity", Type: "int"},
		},
		ForeignKeys: []ForeignKey{{Name: "fk_items_order", Field: "order_id", RefTable: "orders_rel", RefField: "id", Alias: "items"}},
	}); err != nil {
		t.Fatalf("Configure(order_items_rel) error = %v", err)
	}
	if _, err = items.InsertRow(Row{"order_id": 1, "sku": "keyboard", "quantity": 1}); err != nil {
		t.Fatalf("InsertRow(order_items_rel 1) error = %v", err)
	}
	if _, err = items.InsertRow(Row{"order_id": 1, "sku": "mouse", "quantity": 2}); err != nil {
		t.Fatalf("InsertRow(order_items_rel 2) error = %v", err)
	}
	if _, err = items.InsertRow(Row{"order_id": 3, "sku": "monitor", "quantity": 1}); err != nil {
		t.Fatalf("InsertRow(order_items_rel 3) error = %v", err)
	}
	if err = items.Close(); err != nil {
		t.Fatalf("Close(order_items_rel) error = %v", err)
	}

	users, err = newSimpleDB(database, "users_rel")
	if err != nil {
		t.Fatalf("reload users_rel error = %v", err)
	}
	defer users.Close()

	payload, err := users.QueryCascadeJSON(CascadeQuery{
		Conditions: []QueryCondition{{Field: "name", Operator: QueryOpEQ, Value: "alice"}},
		Includes: []CascadeInclude{{
			Table:      "orders_rel",
			Conditions: []QueryCondition{{Field: "amount", Operator: QueryOpGTE, Value: 100.0}},
			Includes: []CascadeInclude{{
				Table:      "order_items_rel",
				Conditions: []QueryCondition{{Field: "sku", Operator: QueryOpEQ, Value: "keyboard"}},
			}},
		}},
	})
	if err != nil {
		t.Fatalf("QueryCascadeJSON(users_rel) error = %v", err)
	}

	var userResults []map[string]any
	if err = json.Unmarshal(payload, &userResults); err != nil {
		t.Fatalf("Unmarshal(users_rel cascade) error = %v", err)
	}
	if len(userResults) != 1 || userResults[0]["name"] != "alice" {
		t.Fatalf("users_rel cascade root = %#v", userResults)
	}
	ordersValue, ok := userResults[0]["orders"].([]any)
	if !ok || len(ordersValue) != 1 {
		t.Fatalf("users_rel nested orders = %#v", userResults[0]["orders"])
	}
	firstOrder, ok := ordersValue[0].(map[string]any)
	if !ok || firstOrder["status"] != "paid" {
		t.Fatalf("users_rel first order = %#v", ordersValue[0])
	}
	itemsValue, ok := firstOrder["items"].([]any)
	if !ok || len(itemsValue) != 1 {
		t.Fatalf("users_rel order items = %#v", firstOrder["items"])
	}
	firstItem, ok := itemsValue[0].(map[string]any)
	if !ok || firstItem["sku"] != "keyboard" {
		t.Fatalf("users_rel first item = %#v", itemsValue[0])
	}
	if err = users.Close(); err != nil {
		t.Fatalf("Close(users_rel after cascade) error = %v", err)
	}

	orders, err = newSimpleDB(database, "orders_rel")
	if err != nil {
		t.Fatalf("reload orders_rel error = %v", err)
	}
	defer orders.Close()

	payload, err = orders.QueryCascadeJSON(CascadeQuery{
		Conditions: []QueryCondition{{Field: "status", Operator: QueryOpEQ, Value: "paid"}},
		Includes: []CascadeInclude{{
			Table: "users_rel",
			Alias: "user",
		}},
	})
	if err != nil {
		t.Fatalf("QueryCascadeJSON(orders_rel) error = %v", err)
	}

	var orderResults []map[string]any
	if err = json.Unmarshal(payload, &orderResults); err != nil {
		t.Fatalf("Unmarshal(orders_rel cascade) error = %v", err)
	}
	if len(orderResults) != 2 {
		t.Fatalf("orders_rel cascade len = %d, want 2", len(orderResults))
	}
	userValue, ok := orderResults[0]["user"].(map[string]any)
	if !ok || userValue["name"] == nil {
		t.Fatalf("orders_rel parent user = %#v", orderResults[0]["user"])
	}

	if _, err = orders.QueryCascadeJSON(CascadeQuery{
		Conditions: []QueryCondition{{Field: "status", Operator: QueryOpEQ, Value: "paid"}},
		MaxDepth:   7,
		Includes: []CascadeInclude{{
			Table: "users_rel",
		}},
	}); !errors.Is(err, ErrCascadeDepthExceeded) {
		t.Fatalf("QueryCascadeJSON(maxDepth>6) error = %v, want ErrCascadeDepthExceeded", err)
	}

	if _, err = orders.QueryCascadeJSON(CascadeQuery{
		Conditions: []QueryCondition{{Field: "status", Operator: QueryOpEQ, Value: "paid"}},
		Includes: []CascadeInclude{{
			Table: "users_rel",
			Includes: []CascadeInclude{{
				Table: "orders_rel",
			}},
		}},
	}); !errors.Is(err, ErrCascadeCycleNotAllow) {
		t.Fatalf("QueryCascadeJSON(cycle) error = %v, want ErrCascadeCycleNotAllow", err)
	}
}

func TestSimpleDB_BatchRowOperations(t *testing.T) {
	database := "demo"
	table := "batch_rows"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string"},
		{Name: "age", Type: "int", Indexed: true},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	inserted, err := db.InsertRows([]Row{
		{"name": "alice", "age": 20},
		{"name": "bob", "age": 21},
		{"name": "carol", "age": 22},
	})
	if err != nil {
		t.Fatalf("InsertRows() error = %v", err)
	}
	if len(inserted) != 3 || inserted[0]["id"] != int64(1) || inserted[2]["id"] != int64(3) {
		t.Fatalf("InsertRows() result = %#v", inserted)
	}

	updated, err := db.UpdateRows([]RowUpdate{
		{PrimaryKey: int64(1), Updates: Row{"age": 30}},
		{PrimaryKey: int64(3), Updates: Row{"name": "carol-updated", "age": 32}},
	})
	if err != nil {
		t.Fatalf("UpdateRows() error = %v", err)
	}
	if len(updated) != 2 || updated[0]["age"] != int64(30) || updated[1]["name"] != "carol-updated" {
		t.Fatalf("UpdateRows() result = %#v", updated)
	}

	if err = db.DeleteRows([]any{int64(2)}); err != nil {
		t.Fatalf("DeleteRows() error = %v", err)
	}

	if _, found, err := db.FindRow(int64(2)); err != nil || found {
		t.Fatalf("FindRow(deleted) = (%v, %v), want not found", err, found)
	}

	rows, err := db.FindByConditions([]QueryCondition{{Field: "age", Operator: QueryOpGTE, Value: 30}})
	if err != nil {
		t.Fatalf("FindByConditions(after batch) error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("FindByConditions(after batch) len = %d, want 2", len(rows))
	}

	if _, err = db.InsertRows([]Row{}); !errors.Is(err, ErrBatchEmpty) {
		t.Fatalf("InsertRows(empty) error = %v, want ErrBatchEmpty", err)
	}
	if _, err = db.UpdateRows([]RowUpdate{}); !errors.Is(err, ErrBatchEmpty) {
		t.Fatalf("UpdateRows(empty) error = %v, want ErrBatchEmpty", err)
	}
	if err = db.DeleteRows([]any{}); !errors.Is(err, ErrBatchEmpty) {
		t.Fatalf("DeleteRows(empty) error = %v, want ErrBatchEmpty", err)
	}
}

func TestSimpleDB_DefaultNullableRequired(t *testing.T) {
	database := "demo"
	table := "defaults"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string", Required: true},
		{Name: "status", Type: "string", Default: "pending"},
		{Name: "note", Type: "string", Nullable: ptr.New(true)},
		{Name: "enabled", Type: "bool", Default: true, Nullable: ptr.New(false)},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	row, err := db.InsertRow(Row{"name": "alice"})
	if err != nil {
		t.Fatalf("InsertRow(defaults) error = %v", err)
	}
	if row["status"] != "pending" {
		t.Fatalf("default status = %#v, want pending", row["status"])
	}
	if row["enabled"] != true {
		t.Fatalf("default enabled = %#v, want true", row["enabled"])
	}
	if _, exists := row["note"]; exists {
		t.Fatalf("nullable note should stay absent when omitted, got %#v", row["note"])
	}

	if _, err = db.InsertRow(Row{"status": "draft"}); !errors.Is(err, ErrFieldRequired) {
		t.Fatalf("InsertRow(missing required) error = %v, want ErrFieldRequired", err)
	}

	row, err = db.InsertRow(Row{"name": "bob", "note": nil})
	if err != nil {
		t.Fatalf("InsertRow(nullable nil) error = %v", err)
	}
	if value, exists := row["note"]; !exists || value != nil {
		t.Fatalf("nullable note = %#v exists=%v, want nil true", value, exists)
	}

	if _, err = db.UpdateRow(int64(1), Row{"name": nil}); !errors.Is(err, ErrFieldNotNullable) {
		t.Fatalf("UpdateRow(non-nullable nil) error = %v, want ErrFieldNotNullable", err)
	}

	if err = db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reloaded, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("reload newSimpleDB() error = %v", err)
	}
	defer reloaded.Close()

	reloadedRow, found, err := reloaded.FindRow(int64(1))
	if err != nil || !found {
		t.Fatalf("FindRow(reloaded) = (%v, %v), want found", err, found)
	}
	if reloadedRow["status"] != "pending" || reloadedRow["enabled"] != true {
		t.Fatalf("reloaded defaults = %#v, want pending/true", reloadedRow)
	}
}

func TestSimpleDB_AutoTimestamps(t *testing.T) {
	database := "demo"
	table := "timestamps"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string", Required: true},
		{Name: "created_at", Type: "timestamp", DefaultExpr: "current_timestamp"},
		{Name: "updated_at", Type: "timestamp", DefaultExpr: "current_timestamp", OnUpdateExpr: "current_timestamp"},
		{Name: "login_time", Type: "time", DefaultExpr: "current_time"},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	row, err := db.InsertRow(Row{"name": "alice"})
	if err != nil {
		t.Fatalf("InsertRow() error = %v", err)
	}

	createdAt, ok := row["created_at"].(string)
	if !ok {
		t.Fatalf("created_at type = %#v, want string", row["created_at"])
	}
	updatedAt, ok := row["updated_at"].(string)
	if !ok {
		t.Fatalf("updated_at type = %#v, want string", row["updated_at"])
	}
	loginTime, ok := row["login_time"].(string)
	if !ok {
		t.Fatalf("login_time type = %#v, want string", row["login_time"])
	}

	createdTS, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		t.Fatalf("parse created_at error = %v", err)
	}
	updatedTS, err := time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		t.Fatalf("parse updated_at error = %v", err)
	}
	if _, err = time.Parse("15:04:05.999999999Z07:00", loginTime); err != nil {
		t.Fatalf("parse login_time error = %v", err)
	}
	if updatedTS.Before(createdTS) {
		t.Fatalf("updated_at should be >= created_at, got created=%s updated=%s", createdAt, updatedAt)
	}

	time.Sleep(10 * time.Millisecond)
	updatedRow, err := db.UpdateRow(int64(1), Row{"name": "alice-updated"})
	if err != nil {
		t.Fatalf("UpdateRow() error = %v", err)
	}
	newUpdatedAt := updatedRow["updated_at"].(string)
	newUpdatedTS, err := time.Parse(time.RFC3339Nano, newUpdatedAt)
	if err != nil {
		t.Fatalf("parse new updated_at error = %v", err)
	}
	if !newUpdatedTS.After(updatedTS) {
		t.Fatalf("updated_at should move forward, old=%s new=%s", updatedAt, newUpdatedAt)
	}
	if updatedRow["created_at"] != createdAt {
		t.Fatalf("created_at should stay unchanged, old=%s new=%v", createdAt, updatedRow["created_at"])
	}

	if err = db.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reloaded, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("reload newSimpleDB() error = %v", err)
	}
	defer reloaded.Close()

	reloadedRow, found, err := reloaded.FindRow(int64(1))
	if err != nil || !found {
		t.Fatalf("FindRow(reloaded) = (%v, %v), want found", err, found)
	}
	if reloadedRow["created_at"] != createdAt {
		t.Fatalf("reloaded created_at = %v, want %s", reloadedRow["created_at"], createdAt)
	}
	if reloadedRow["updated_at"] != newUpdatedAt {
		t.Fatalf("reloaded updated_at = %v, want %s", reloadedRow["updated_at"], newUpdatedAt)
	}
}

func TestSimpleDB_TimeExpressionValidation(t *testing.T) {
	database := "demo"
	table := "time_validation"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true},
		{Name: "created_at", Type: "string", DefaultExpr: "current_timestamp"},
	}})
	if !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("defaultExpr type mismatch error = %v, want ErrInvalidSchema", err)
	}

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true},
		{Name: "updated_at", Type: "timestamp", OnUpdateExpr: "current_time"},
	}})
	if !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("onUpdateExpr type mismatch error = %v, want ErrInvalidSchema", err)
	}
}

func TestSimpleDB_LengthEnumAndChecks(t *testing.T) {
	database := "demo"
	table := "column_checks"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
		{Name: "code", Type: "string", MaxLength: 5, Checks: []ColumnCheck{{Operator: "regex", Value: "^[A-Z0-9]+$"}}},
		{Name: "status", Type: "string", Enum: []any{"draft", "active", "done"}},
		{Name: "age", Type: "int", Checks: []ColumnCheck{{Operator: "gte", Value: 18}, {Operator: "lte", Value: 60}}},
		{Name: "tags", Type: "array", Checks: []ColumnCheck{{Operator: "len_lte", Value: 3}}},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	row, err := db.InsertRow(Row{"code": "A001", "status": "active", "age": 20, "tags": []any{"go", "db"}})
	if err != nil {
		t.Fatalf("InsertRow(valid) error = %v", err)
	}
	if row["status"] != "active" {
		t.Fatalf("status = %#v, want active", row["status"])
	}

	if _, err = db.InsertRow(Row{"code": "TOO-LONG", "status": "draft", "age": 20}); !errors.Is(err, ErrFieldLengthViolation) {
		t.Fatalf("InsertRow(length violation) error = %v, want ErrFieldLengthViolation", err)
	}

	if _, err = db.InsertRow(Row{"code": "B002", "status": "archived", "age": 20}); !errors.Is(err, ErrFieldEnumViolation) {
		t.Fatalf("InsertRow(enum violation) error = %v, want ErrFieldEnumViolation", err)
	}

	if _, err = db.InsertRow(Row{"code": "B002", "status": "draft", "age": 12}); !errors.Is(err, ErrFieldCheckViolation) {
		t.Fatalf("InsertRow(check violation) error = %v, want ErrFieldCheckViolation", err)
	}

	if _, err = db.UpdateRow(int64(1), Row{"code": "bad-1"}); !errors.Is(err, ErrFieldCheckViolation) {
		t.Fatalf("UpdateRow(regex violation) error = %v, want ErrFieldCheckViolation", err)
	}

	if _, err = db.UpdateRow(int64(1), Row{"tags": []any{"a", "b", "c", "d"}}); !errors.Is(err, ErrFieldCheckViolation) {
		t.Fatalf("UpdateRow(array length violation) error = %v, want ErrFieldCheckViolation", err)
	}
}

func TestSimpleDB_ConstraintSchemaValidation(t *testing.T) {
	database := "demo"
	table := "constraint_schema_validation"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true},
		{Name: "title", Type: "int", MaxLength: 10},
	}})
	if !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("length on int schema error = %v, want ErrInvalidSchema", err)
	}

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true},
		{Name: "status", Type: "string", Enum: []any{"draft", 1}},
	}})
	if !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("mixed enum schema error = %v, want ErrInvalidSchema", err)
	}

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "int", PrimaryKey: true},
		{Name: "code", Type: "string", Checks: []ColumnCheck{{Operator: "regex", Value: "["}}},
	}})
	if !errors.Is(err, ErrInvalidSchema) {
		t.Fatalf("invalid regex schema error = %v, want ErrInvalidSchema", err)
	}
}

func TestSimpleDB_ConfigAndUUIDVersions(t *testing.T) {
	database := "demo"
	table := "config_test"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(database, table))
	})

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	// Test default config
	defaultConfig := db.GetConfig()
	if defaultConfig.DefaultUUIDVersion != DefaultUUIDVersion {
		t.Fatalf("default UUID version = %d, want %d", defaultConfig.DefaultUUIDVersion, DefaultUUIDVersion)
	}
	if defaultConfig.DefaultCascadeMaxDepth != DefaultCascadeMaxDepth {
		t.Fatalf("default cascade depth = %d, want %d", defaultConfig.DefaultCascadeMaxDepth, DefaultCascadeMaxDepth)
	}

	// Test SetConfig
	newConfig := DatabaseConfig{
		DefaultUUIDVersion:     7,
		DefaultCascadeMaxDepth: 4,
	}
	db.SetConfig(newConfig)
	retrievedConfig := db.GetConfig()
	if retrievedConfig.DefaultUUIDVersion != 7 {
		t.Fatalf("set UUID version = %d, want 7", retrievedConfig.DefaultUUIDVersion)
	}
	if retrievedConfig.DefaultCascadeMaxDepth != 4 {
		t.Fatalf("set cascade depth = %d, want 4", retrievedConfig.DefaultCascadeMaxDepth)
	}

	// Test UUID type with version format (uuid:v6)
	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid:v6", PrimaryKey: true, AutoIncrement: true},
		{Name: "name", Type: "string"},
	}})
	if err != nil {
		t.Fatalf("Configure with uuid:v6 error = %v", err)
	}

	row, err := db.InsertRow(Row{"name": "test"})
	if err != nil {
		t.Fatalf("InsertRow error = %v", err)
	}

	id, ok := row["id"].(string)
	if !ok {
		t.Fatalf("id type = %#v, want string", row["id"])
	}
	if _, err = uuid.Parse(id); err != nil {
		t.Fatalf("generated id parse error = %v", err)
	}
}
