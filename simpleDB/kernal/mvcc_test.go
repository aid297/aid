package kernal

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSimpleDB_MVCCSnapshotRead(t *testing.T) {
	database := "demo"
	table := "mvcc_snapshot"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	if err = db.Put("k", []byte("v1")); err != nil {
		t.Fatalf("Put(v1) error = %v", err)
	}

	tx, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}

	if err = db.Put("k", []byte("v2")); err != nil {
		t.Fatalf("Put(v2) error = %v", err)
	}

	value, exists, err := tx.Get("k")
	if err != nil {
		t.Fatalf("tx.Get() error = %v", err)
	}
	if !exists || string(value) != "v1" {
		t.Fatalf("tx snapshot value = (%s, %v), want (v1, true)", string(value), exists)
	}

	if err = tx.Rollback(); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}
}

func TestSimpleDB_MVCCConflict(t *testing.T) {
	database := "demo"
	table := "mvcc_conflict"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	if err = db.Put("k", []byte("v1")); err != nil {
		t.Fatalf("Put(v1) error = %v", err)
	}

	tx, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}

	if _, _, err = tx.Get("k"); err != nil {
		t.Fatalf("tx.Get() error = %v", err)
	}

	if err = db.Put("k", []byte("v2")); err != nil {
		t.Fatalf("Put(v2) error = %v", err)
	}

	if err = tx.Put("k", []byte("v3")); err != nil {
		t.Fatalf("tx.Put() error = %v", err)
	}

	err = tx.Commit()
	if !errors.Is(err, ErrTxConflict) {
		t.Fatalf("Commit() error = %v, want ErrTxConflict", err)
	}
}

func TestSimpleDB_MVCCWriteWriteConflict(t *testing.T) {
	database := "demo"
	table := "mvcc_ww_conflict"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	if err = db.Put("k", []byte("v1")); err != nil {
		t.Fatalf("Put(v1) error = %v", err)
	}

	tx1, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx(tx1) error = %v", err)
	}
	tx2, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx(tx2) error = %v", err)
	}

	if err = tx1.Put("k", []byte("v2")); err != nil {
		t.Fatalf("tx1.Put() error = %v", err)
	}
	if err = tx2.Put("k", []byte("v3")); err != nil {
		t.Fatalf("tx2.Put() error = %v", err)
	}

	if err = tx1.Commit(); err != nil {
		t.Fatalf("tx1.Commit() error = %v", err)
	}
	if err = tx2.Commit(); !errors.Is(err, ErrTxConflict) {
		t.Fatalf("tx2.Commit() error = %v, want ErrTxConflict", err)
	}
}

func TestSimpleDB_MVCCReadOnlyAndSnapshotQuery(t *testing.T) {
	database := "demo"
	table := "mvcc_readonly_query"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	if err = db.Put("a:1", []byte("v1")); err != nil {
		t.Fatalf("Put(a:1) error = %v", err)
	}
	if err = db.Put("a:2", []byte("v2")); err != nil {
		t.Fatalf("Put(a:2) error = %v", err)
	}

	tx, err := db.BeginReadOnlyTx()
	if err != nil {
		t.Fatalf("BeginReadOnlyTx() error = %v", err)
	}

	if err = db.Put("a:3", []byte("v3")); err != nil {
		t.Fatalf("Put(a:3) error = %v", err)
	}

	rows, err := tx.Query("a:")
	if err != nil {
		t.Fatalf("tx.Query() error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("snapshot rows len = %d, want 2", len(rows))
	}
	if _, exists := rows["a:3"]; exists {
		t.Fatalf("snapshot should not contain a:3")
	}

	keys, err := tx.Keys()
	if err != nil {
		t.Fatalf("tx.Keys() error = %v", err)
	}
	if !reflect.DeepEqual(keys, []string{"a:1", "a:2"}) {
		t.Fatalf("tx.Keys() = %#v, want [a:1 a:2]", keys)
	}

	if err = tx.Put("a:4", []byte("v4")); !errors.Is(err, ErrTxReadOnly) {
		t.Fatalf("tx.Put(readonly) error = %v, want ErrTxReadOnly", err)
	}

	if err = tx.Rollback(); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}
}

func TestSimpleDB_MVCCRowTransaction(t *testing.T) {
	database := "demo"
	table := "mvcc_row_tx"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
		{Name: "age", Type: "int", Required: true, Default: int64(0)},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	tx, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}

	row, err := tx.InsertRow(Row{"username": "alice", "age": 20})
	if err != nil {
		t.Fatalf("tx.InsertRow() error = %v", err)
	}

	rows, err := tx.Find(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: "alice"})
	if err != nil {
		t.Fatalf("tx.Find() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("tx.Find() len = %d, want 1", len(rows))
	}

	pk := row["id"]
	updated, err := tx.UpdateRow(pk, Row{"age": 21})
	if err != nil {
		t.Fatalf("tx.UpdateRow() error = %v", err)
	}
	if updated["age"].(int64) != 21 {
		t.Fatalf("updated age = %#v, want 21", updated["age"])
	}

	if err = tx.Commit(); err != nil {
		t.Fatalf("tx.Commit() error = %v", err)
	}

	stored, found, err := db.FindRow(pk)
	if err != nil {
		t.Fatalf("db.FindRow() error = %v", err)
	}
	if !found {
		t.Fatalf("db.FindRow() found = false, want true")
	}
	if stored["age"].(int64) != 21 {
		t.Fatalf("stored age = %#v, want 21", stored["age"])
	}

	tx2, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx(tx2) error = %v", err)
	}
	if err = tx2.DeleteRow(pk); err != nil {
		t.Fatalf("tx2.DeleteRow() error = %v", err)
	}
	if err = tx2.Commit(); err != nil {
		t.Fatalf("tx2.Commit() error = %v", err)
	}

	_, found, err = db.FindRow(pk)
	if err != nil {
		t.Fatalf("db.FindRow(after delete) error = %v", err)
	}
	if found {
		t.Fatalf("db.FindRow(after delete) found = true, want false")
	}
}

func TestSimpleDB_MVCCBatchRowTransaction(t *testing.T) {
	database := "demo"
	table := "mvcc_batch_row_tx"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
		{Name: "age", Type: "int", Required: true, Default: int64(0)},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	tx, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}

	insertedRows, err := tx.InsertRows([]Row{
		{"username": "u1", "age": 18},
		{"username": "u2", "age": 19},
	})
	if err != nil {
		t.Fatalf("tx.InsertRows() error = %v", err)
	}
	if len(insertedRows) != 2 {
		t.Fatalf("tx.InsertRows() len = %d, want 2", len(insertedRows))
	}

	updates := []RowUpdate{
		{PrimaryKey: insertedRows[0]["id"], Updates: Row{"age": 20}},
		{PrimaryKey: insertedRows[1]["id"], Updates: Row{"age": 21}},
	}
	updatedRows, err := tx.UpdateRows(updates)
	if err != nil {
		t.Fatalf("tx.UpdateRows() error = %v", err)
	}
	if len(updatedRows) != 2 {
		t.Fatalf("tx.UpdateRows() len = %d, want 2", len(updatedRows))
	}

	if err = tx.DeleteRows([]any{insertedRows[0]["id"]}); err != nil {
		t.Fatalf("tx.DeleteRows() error = %v", err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatalf("tx.Commit() error = %v", err)
	}

	rows, err := db.Find()
	if err != nil {
		t.Fatalf("db.Find() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("db.Find() len = %d, want 1", len(rows))
	}
	if rows[0]["username"].(string) != "u2" {
		t.Fatalf("remaining username = %v, want u2", rows[0]["username"])
	}
}

func TestSimpleDB_MVCCUniqueConflictAcrossTransactions(t *testing.T) {
	database := "demo"
	table := "mvcc_unique_conflict"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	tx1, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx(tx1) error = %v", err)
	}
	tx2, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx(tx2) error = %v", err)
	}

	if _, err = tx1.InsertRow(Row{"username": "duplicate"}); err != nil {
		t.Fatalf("tx1.InsertRow() error = %v", err)
	}
	if _, err = tx2.InsertRow(Row{"username": "duplicate"}); err != nil {
		t.Fatalf("tx2.InsertRow() error = %v", err)
	}

	if err = tx1.Commit(); err != nil {
		t.Fatalf("tx1.Commit() error = %v", err)
	}
	if err = tx2.Commit(); !errors.Is(err, ErrTxConflict) {
		t.Fatalf("tx2.Commit() error = %v, want ErrTxConflict", err)
	}
}

func TestSimpleDB_MVCCRemoveByCondition(t *testing.T) {
	database := "demo"
	table := "mvcc_remove_condition"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
		{Name: "age", Type: "int", Required: true, Default: int64(0)},
	}})
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	_, err = db.InsertRows([]Row{
		{"username": "u1", "age": 20},
		{"username": "u2", "age": 25},
		{"username": "u3", "age": 30},
	})
	if err != nil {
		t.Fatalf("db.InsertRows() error = %v", err)
	}

	tx, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}

	deleted, err := tx.RemoveByCondition(QueryCondition{Field: "age", Operator: QueryOpGTE, Value: 25})
	if err != nil {
		t.Fatalf("tx.RemoveByCondition() error = %v", err)
	}
	if deleted != 2 {
		t.Fatalf("tx.RemoveByCondition() deleted = %d, want 2", deleted)
	}

	if err = tx.Commit(); err != nil {
		t.Fatalf("tx.Commit() error = %v", err)
	}

	rows, err := db.Find()
	if err != nil {
		t.Fatalf("db.Find() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("db.Find() len = %d, want 1", len(rows))
	}
	if rows[0]["username"].(string) != "u1" {
		t.Fatalf("remaining username = %v, want u1", rows[0]["username"])
	}
}

func TestSimpleDB_MVCCReadCommitted(t *testing.T) {
	database := "demo"
	table := "mvcc_read_committed"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := newSimpleDB(database, table)
	if err != nil {
		t.Fatalf("newSimpleDB() error = %v", err)
	}
	defer db.Close()

	if err = db.Put("k", []byte("v1")); err != nil {
		t.Fatalf("Put(v1) error = %v", err)
	}

	tx, err := db.BeginTxWithOptions(TxOptions{Isolation: TxIsolationReadCommitted})
	if err != nil {
		t.Fatalf("BeginTxWithOptions() error = %v", err)
	}

	value, exists, err := tx.Get("k")
	if err != nil {
		t.Fatalf("tx.Get(first) error = %v", err)
	}
	if !exists || string(value) != "v1" {
		t.Fatalf("tx.Get(first) = (%s, %v), want (v1, true)", string(value), exists)
	}

	if err = db.Put("k", []byte("v2")); err != nil {
		t.Fatalf("Put(v2) error = %v", err)
	}

	value, exists, err = tx.Get("k")
	if err != nil {
		t.Fatalf("tx.Get(second) error = %v", err)
	}
	if !exists || string(value) != "v2" {
		t.Fatalf("tx.Get(second) = (%s, %v), want (v2, true)", string(value), exists)
	}

	if err = tx.Rollback(); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}
}
