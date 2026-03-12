package driver

import (
	"errors"
	"testing"

	"github.com/aid297/aid/simpleDB/kernal"
)

func TestDriver_SystemTableSchemaErrorWrapped(t *testing.T) {
	dir := t.TempDir()

	err := kernal.New.EnsureSystemTables(dir)
	if err != nil {
		t.Fatalf("ensure system tables: %v", err)
	}

	broken, err := kernal.New.DB(dir, "_broken_seed")
	if err != nil {
		t.Fatalf("seed open db: %v", err)
	}
	_ = broken.Close()

	raw, err := kernal.New.DB(dir, "_sys_users")
	if err != nil {
		t.Fatalf("open system users: %v", err)
	}
	if err = raw.DropTable(); err != nil {
		_ = raw.Close()
		t.Fatalf("drop users table: %v", err)
	}
	if err = raw.CreateTable(kernal.TableSchema{Columns: []kernal.Column{
		{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Required: true, Unique: true},
	}}); err != nil {
		_ = raw.Close()
		t.Fatalf("create malformed users table: %v", err)
	}
	_ = raw.Close()

	_, err = New.DB(dir, "business")
	if err == nil {
		t.Fatal("expected wrapped system table schema error")
	}
	var driverErr *DriverError
	if !errors.As(err, &driverErr) {
		t.Fatalf("unexpected error type: %T", err)
	}
	if driverErr.Code != ErrorCodeDDL {
		t.Fatalf("unexpected driver error code: %s", driverErr.Code)
	}
	if !errors.Is(err, kernal.ErrSystemTableSchema) {
		t.Fatalf("unexpected wrapped error: %v", err)
	}
}
