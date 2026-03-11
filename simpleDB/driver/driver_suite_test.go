package driver

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aid297/aid/simpleDB/kernal"
)

func TestDriver_CRUDAndTx(t *testing.T) {
	database := "demo"
	table := "driver_crud_tx"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := New.DB(database, table)
	if err != nil {
		t.Fatalf("Open DB error = %v", err)
	}
	defer db.Close()

	err = db.Configure(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid:v7", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Unique: true, Required: true},
		{Name: "age", Type: "int", Required: true, Default: int64(0)},
	}})
	if err != nil {
		t.Fatalf("Configure error = %v", err)
	}

	tx, err := db.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx error = %v", err)
	}

	row, err := tx.InsertRow(Row{"username": "driver_alice", "age": 22})
	if err != nil {
		t.Fatalf("InsertRow(tx) error = %v", err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatalf("Commit error = %v", err)
	}

	found, ok, err := db.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: "driver_alice"})
	if err != nil {
		t.Fatalf("FindOne error = %v", err)
	}
	if !ok {
		t.Fatalf("FindOne ok=false, want true")
	}
	if found["id"] != row["id"] {
		t.Fatalf("id mismatch = %v vs %v", found["id"], row["id"])
	}
}

func TestDriver_ErrorModel(t *testing.T) {
	database := "demo"
	table := "driver_error_model"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := New.DB(database, table)
	if err != nil {
		t.Fatalf("Open DB error = %v", err)
	}
	defer db.Close()

	err = db.Delete("missing-key")
	if err == nil {
		t.Fatalf("Delete(missing) error=nil, want not nil")
	}

	var dErr *DriverError
	if !errors.As(err, &dErr) {
		t.Fatalf("error type = %T, want *DriverError", err)
	}
	if dErr.Code != ErrorCodeNotFound {
		t.Fatalf("error code = %s, want %s", dErr.Code, ErrorCodeNotFound)
	}
}

func TestDriver_TxOptionsReadCommitted(t *testing.T) {
	database := "demo"
	table := "driver_tx_options"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := New.DB(database, table)
	if err != nil {
		t.Fatalf("Open DB error = %v", err)
	}
	defer db.Close()

	if err = db.Put("k", []byte("v1")); err != nil {
		t.Fatalf("Put(v1) error = %v", err)
	}

	tx, err := db.BeginTxWithOptions(TxOptions{Isolation: TxIsolationReadCommitted})
	if err != nil {
		t.Fatalf("BeginTxWithOptions error = %v", err)
	}

	if err = db.Put("k", []byte("v2")); err != nil {
		t.Fatalf("Put(v2) error = %v", err)
	}

	value, ok, err := tx.Get("k")
	if err != nil {
		t.Fatalf("tx.Get error = %v", err)
	}
	if !ok || string(value) != "v2" {
		t.Fatalf("tx.Get = (%s,%v), want (v2,true)", string(value), ok)
	}

	if err = tx.Rollback(); err != nil {
		t.Fatalf("Rollback error = %v", err)
	}
}

func TestDriver_SetAttrs(t *testing.T) {
	database := "demo"
	table := "driver_attrs"
	_ = os.RemoveAll(filepath.Join(database, table))
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(database, table)) })

	db, err := New.DB(database, table,
		kernal.UUIDVersion(7),
		kernal.UUIDWithHyphen(false),
	)
	if err != nil {
		t.Fatalf("Open DB error = %v", err)
	}
	defer db.Close()

	config := db.GetConfig()
	if config.DefaultUUIDVersion != 7 {
		t.Fatalf("DefaultUUIDVersion = %d, want 7", config.DefaultUUIDVersion)
	}
	if config.DefaultUUIDWithHyphen == nil || *config.DefaultUUIDWithHyphen {
		t.Fatalf("DefaultUUIDWithHyphen = %v, want false", config.DefaultUUIDWithHyphen)
	}
}
