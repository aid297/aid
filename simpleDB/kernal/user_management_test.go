package kernal

import (
	"errors"
	"testing"
)

func TestAssignRoles_SuperAdminReservedForSDB(t *testing.T) {
	dir := t.TempDir()

	if _, err := New.DB(dir, "biz_users"); err != nil {
		t.Fatalf("init db: %v", err)
	}
	if _, err := New.RegisterUser(dir, "alice", "123456", "Alice"); err != nil {
		t.Fatalf("register alice: %v", err)
	}

	_, err := New.AssignRoles(dir, "alice", []string{defaultSystemRoleCode})
	if err == nil {
		t.Fatal("expected ErrSuperAdminRoleReserved")
	}
	if !errors.Is(err, ErrSuperAdminRoleReserved) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitSDBPassword(t *testing.T) {
	dir := t.TempDir()

	if _, err := New.DB(dir, "biz_users"); err != nil {
		t.Fatalf("init db: %v", err)
	}

	if err := New.InitSDBPassword(dir); err != nil {
		t.Fatalf("init sdb password: %v", err)
	}

	if _, err := New.Authenticate(dir, defaultSystemAdminUsername, defaultSystemAdminPassword); err != nil {
		t.Fatalf("authenticate sdb after init password: %v", err)
	}
}
