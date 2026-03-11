package kernal

import (
	"errors"
	"os"
	"path/filepath"
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
