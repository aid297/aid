package kernal

import "testing"

func TestAssignRolePermissions(t *testing.T) {
	dir := t.TempDir()

	if _, err := New.DB(dir, "biz_orders"); err != nil {
		t.Fatalf("init db: %v", err)
	}

	if err := New.AssignRolePermissions(dir, defaultSystemRoleCode, []string{"table.orders.select", "table.orders.update"}); err != nil {
		t.Fatalf("assign role permissions: %v", err)
	}

	permissionsDB, err := newSimpleDB(dir, systemTablePermissions)
	if err != nil {
		t.Fatalf("open permissions db: %v", err)
	}
	defer permissionsDB.Close()

	if _, ok, err := permissionsDB.FindOne(QueryCondition{Field: "code", Operator: QueryOpEQ, Value: "table.orders.select"}); err != nil {
		t.Fatalf("find permission table.orders.select: %v", err)
	} else if !ok {
		t.Fatal("permission table.orders.select should exist")
	}

	rolesDB, err := newSimpleDB(dir, systemTableRoles)
	if err != nil {
		t.Fatalf("open roles db: %v", err)
	}
	defer rolesDB.Close()
	roleRow, ok, err := rolesDB.FindOne(QueryCondition{Field: "code", Operator: QueryOpEQ, Value: defaultSystemRoleCode})
	if err != nil {
		t.Fatalf("find super_admin role: %v", err)
	}
	if !ok {
		t.Fatal("super_admin role should exist")
	}

	rolePermissionsDB, err := newSimpleDB(dir, systemTableRolePermissions)
	if err != nil {
		t.Fatalf("open role_permissions db: %v", err)
	}
	defer rolePermissionsDB.Close()

	links, err := rolePermissionsDB.Find(QueryCondition{Field: "roleId", Operator: QueryOpEQ, Value: roleRow["id"]})
	if err != nil {
		t.Fatalf("find role permission links: %v", err)
	}
	if len(links) == 0 {
		t.Fatal("expected role permission links for super_admin")
	}
}
