package simpleDBDriver

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSimpleDB_CRUDAndReload(t *testing.T) {
	dataDir := "./demo"

	db, err := newSimpleDB(dataDir)
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

	reloaded, err := newSimpleDB(dataDir)
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
	dataDir := t.TempDir()

	db, err := newSimpleDB(dataDir)
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

	before, err := os.Stat(filepath.Join(dataDir, defaultDataFile))
	if err != nil {
		t.Fatalf("Stat before compact error = %v", err)
	}

	if err = db.Compact(); err != nil {
		t.Fatalf("Compact() error = %v", err)
	}

	after, err := os.Stat(filepath.Join(dataDir, defaultDataFile))
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
	dataDir := t.TempDir()

	db, err := newSimpleDB(dataDir)
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
