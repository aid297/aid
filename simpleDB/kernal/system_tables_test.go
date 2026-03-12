package kernal

import (
	"errors"
	"testing"

	"github.com/aid297/aid/digest"
)

func TestSystemTables_AutoBootstrapAndDefaultAdmin(t *testing.T) {
	dir := t.TempDir()

	db, err := New.DB(dir, "biz_users")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	for _, definition := range systemTableDefinitions() {
		tableDB, err := newSimpleDB(systemDatabaseFor(dir), definition.name)
		if err != nil {
			t.Fatalf("open %s: %v", definition.name, err)
		}

		plan, exists, err := tableDB.SchemaDiff(definition.schema)
		_ = tableDB.Close()
		if err != nil {
			t.Fatalf("schema diff %s: %v", definition.name, err)
		}
		if !exists {
			t.Fatalf("system table %s not created", definition.name)
		}
		if plan != nil {
			t.Fatalf("system table %s schema mismatch: %+v", definition.name, plan)
		}
	}

	usersDB, err := newSimpleDB(systemDatabaseFor(dir), systemTableUsers)
	if err != nil {
		t.Fatalf("open users db: %v", err)
	}
	defer usersDB.Close()

	row, ok, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: defaultSystemAdminUsername})
	if err != nil {
		t.Fatalf("find default admin: %v", err)
	}
	if !ok {
		t.Fatal("default admin not found")
	}
	passwordHash, _ := row["passwordHash"].(string)
	if passwordHash == "" {
		t.Fatal("default admin password hash is empty")
	}
	if !digest.BcryptCheck(defaultSystemAdminPassword, passwordHash) {
		t.Fatal("default admin password hash mismatch")
	}
	isAdmin, _ := row["isAdmin"].(bool)
	if !isAdmin {
		t.Fatal("default admin should be admin")
	}

	roleLinksDB, err := newSimpleDB(systemDatabaseFor(dir), systemTableUserRoles)
	if err != nil {
		t.Fatalf("open user role db: %v", err)
	}
	defer roleLinksDB.Close()

	_, ok, err = roleLinksDB.FindOne(QueryCondition{Field: "userId", Operator: QueryOpEQ, Value: row["id"]})
	if err != nil {
		t.Fatalf("find admin role link: %v", err)
	}
	if !ok {
		t.Fatal("default admin role link not found")
	}
}

func TestSystemTables_SchemaMismatchReturnsError(t *testing.T) {
	dir := t.TempDir()

	usersDB, err := newSimpleDB(systemDatabaseFor(dir), systemTableUsers)
	if err != nil {
		t.Fatalf("open malformed users db: %v", err)
	}
	if err = usersDB.CreateTable(TableSchema{Columns: []Column{
		{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
		{Name: "username", Type: "string", Required: true, Unique: true},
	}}); err != nil {
		_ = usersDB.Close()
		t.Fatalf("create malformed users schema: %v", err)
	}
	_ = usersDB.Close()

	_, err = New.DB(dir, "biz_orders")
	if err == nil {
		t.Fatal("expected schema mismatch error")
	}
	if !errors.Is(err, ErrSystemTableSchema) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSystemTables_RecreateMissingDefaultAdmin(t *testing.T) {
	dir := t.TempDir()

	db, err := New.DB(dir, "biz_accounts")
	if err != nil {
		t.Fatalf("initial open: %v", err)
	}
	_ = db.Close()

	usersDB, err := newSimpleDB(systemDatabaseFor(dir), systemTableUsers)
	if err != nil {
		t.Fatalf("open users db: %v", err)
	}
	if _, err = usersDB.RemoveByCondition(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: defaultSystemAdminUsername}); err != nil {
		_ = usersDB.Close()
		t.Fatalf("remove default admin: %v", err)
	}
	_ = usersDB.Close()

	db, err = New.DB(dir, "biz_logs")
	if err != nil {
		t.Fatalf("re-open db: %v", err)
	}
	_ = db.Close()

	usersDB, err = newSimpleDB(systemDatabaseFor(dir), systemTableUsers)
	if err != nil {
		t.Fatalf("re-open users db: %v", err)
	}
	defer usersDB.Close()

	row, ok, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: defaultSystemAdminUsername})
	if err != nil {
		t.Fatalf("find recreated default admin: %v", err)
	}
	if !ok {
		t.Fatal("default admin should be recreated")
	}
	passwordHash, _ := row["passwordHash"].(string)
	if !digest.BcryptCheck(defaultSystemAdminPassword, passwordHash) {
		t.Fatal("recreated default admin password hash mismatch")
	}
}
