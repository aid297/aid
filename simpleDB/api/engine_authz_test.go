package api

import (
	"errors"
	"testing"

	"github.com/aid297/aid/simpleDB/driver"
)

func TestEngine_SystemTableAccessRequiresSuperAdmin(t *testing.T) {
	db := t.TempDir()

	eng := NewEngine(db, BackendDriver)
	_, err := eng.Execute("SELECT * FROM _sys_users")
	if err == nil {
		t.Fatal("expected system table access denied error")
	}
	if !errors.Is(err, ErrSystemTableAccessDenied) {
		t.Fatalf("unexpected error: %v", err)
	}

	superAdmin := &driver.AuthenticatedUser{Username: "sdb", Roles: []string{"super_admin"}}
	eng.WithActor(superAdmin)
	if _, err = eng.Execute("SELECT * FROM _sys_users"); err != nil {
		t.Fatalf("super_admin should access system table, got error: %v", err)
	}
}
