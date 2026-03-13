package kernal

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aid297/aid/digest"
	json "github.com/json-iterator/go"
)

const (
	systemTableUsers             = "_sys_users"
	systemTableRevokedTokens     = "_sys_revoked_tokens"
	systemTableActiveTokens      = "_sys_active_tokens"
	systemTableRoles             = "_sys_roles"
	systemTablePermissions       = "_sys_permissions"
	systemTableUserRoles         = "_sys_user_roles"
	systemTableRolePermissions   = "_sys_role_permissions"
	systemTableUserDBBindings    = "_sys_user_db_bindings"
	systemTableDatabaseOwners    = "_sys_db_owners"
	systemTableTableOwners       = "_sys_table_owners"
	systemTableTableAccessGrants = "_sys_table_access_grants"

	systemDatabaseName = "__simpledb_sys"

	defaultSystemAdminUsername = "sdb"
	defaultSystemAdminPassword = "simpleDB"
	defaultSystemAdminName     = "simpleDB administrator"
	defaultSystemRoleCode      = "super_admin"
	defaultSystemRoleName      = "Super Administrator"
	defaultSystemStatus        = "active"
)

var systemBootstrapMu sync.Mutex

type systemTableDefinition struct {
	name   string
	schema TableSchema
}

func ensureSystemTables(database string) error {
	systemDatabase := systemDatabaseFor(database)
	if strings.TrimSpace(systemDatabase) == "" {
		return ErrDBPathEmpty
	}

	systemBootstrapMu.Lock()
	defer systemBootstrapMu.Unlock()

	for _, definition := range systemTableDefinitions() {
		if err := ensureSystemTable(systemDatabase, definition); err != nil {
			return err
		}
	}

	if err := ensureDefaultAdmin(systemDatabase); err != nil {
		return err
	}

	return nil
}

func systemDatabaseFor(database string) string {
	trimmed := strings.TrimSpace(database)
	if trimmed == "" {
		return ""
	}
	cleaned := filepath.Clean(trimmed)
	if filepath.Base(cleaned) == systemDatabaseName {
		return trimmed
	}
	parent := filepath.Dir(cleaned)
	if parent == "." || strings.TrimSpace(parent) == "" {
		return systemDatabaseName
	}
	return filepath.Join(parent, systemDatabaseName)
}

func systemTableDefinitions() []systemTableDefinition {
	return []systemTableDefinition{
		{name: systemTableUsers, schema: systemUsersSchema()},
		{name: systemTableRevokedTokens, schema: systemRevokedTokensSchema()},
		{name: systemTableActiveTokens, schema: systemActiveTokensSchema()},
		{name: systemTableRoles, schema: systemRolesSchema()},
		{name: systemTablePermissions, schema: systemPermissionsSchema()},
		{name: systemTableUserRoles, schema: systemUserRolesSchema()},
		{name: systemTableRolePermissions, schema: systemRolePermissionsSchema()},
		{name: systemTableDatabaseOwners, schema: systemDatabaseOwnersSchema()},
		{name: systemTableUserDBBindings, schema: systemUserDBBindingsSchema()},
		{name: systemTableTableOwners, schema: systemTableOwnersSchema()},
		{name: systemTableTableAccessGrants, schema: systemTableAccessGrantsSchema()},
	}
}

func systemRevokedTokensSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "tokenId", Type: "string", Required: true, Unique: true},
			{Name: "expiresAt", Type: "int", Required: true},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemActiveTokensSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "tokenId", Type: "string", Required: true, Unique: true},
			{Name: "expiresAt", Type: "int", Required: true},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemDatabaseOwnersSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "databaseName", Type: "string", Required: true, Unique: true},
			{Name: "ownerUserId", Type: "uuid", Required: true},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
			{Name: "updatedAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp, OnUpdateExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemUserDBBindingsSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "userId", Type: "uuid", Required: true},
			{Name: "databaseName", Type: "string", Required: true},
			{Name: "enabled", Type: "bool", Default: true},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
			{Name: "updatedAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp, OnUpdateExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemTableOwnersSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "tableName", Type: "string", Required: true, Unique: true},
			{Name: "ownerUserId", Type: "uuid", Required: true},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemTableAccessGrantsSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "tableName", Type: "string", Required: true},
			{Name: "granteeUserId", Type: "uuid", Required: true},
			{Name: "scope", Type: "string", Required: true},
			{Name: "ownerApproved", Type: "bool", Default: false},
			{Name: "adminApproved", Type: "bool", Default: false},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
			{Name: "updatedAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp, OnUpdateExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemUsersSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "username", Type: "string", Required: true, Unique: true},
			{Name: "passwordHash", Type: "string", Required: true},
			{Name: "displayName", Type: "string", Default: ""},
			{Name: "isAdmin", Type: "bool", Default: false},
			{Name: "status", Type: "string", Default: defaultSystemStatus},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
			{Name: "updatedAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp, OnUpdateExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemRolesSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "code", Type: "string", Required: true, Unique: true},
			{Name: "name", Type: "string", Required: true},
			{Name: "description", Type: "string", Default: ""},
			{Name: "isSystem", Type: "bool", Default: false},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
			{Name: "updatedAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp, OnUpdateExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemPermissionsSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "code", Type: "string", Required: true, Unique: true},
			{Name: "name", Type: "string", Required: true},
			{Name: "description", Type: "string", Default: ""},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
			{Name: "updatedAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp, OnUpdateExpr: ColumnExprCurrentTimestamp},
		}}
}

func systemUserRolesSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "userId", Type: "uuid", Required: true},
			{Name: "roleId", Type: "uuid", Required: true},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
		},
		ForeignKeys: []ForeignKey{
			{Name: "fk_sys_user_roles_user", Field: "userId", RefTable: systemTableUsers, RefField: "id", Alias: "user"},
			{Name: "fk_sys_user_roles_role", Field: "roleId", RefTable: systemTableRoles, RefField: "id", Alias: "role"},
		},
	}
}

func systemRolePermissionsSchema() TableSchema {
	return TableSchema{
		Engine: EngineDisk,
		Columns: []Column{
			{Name: "id", Type: "uuid", PrimaryKey: true, AutoIncrement: true},
			{Name: "roleId", Type: "uuid", Required: true},
			{Name: "permissionId", Type: "uuid", Required: true},
			{Name: "createdAt", Type: "timestamp", DefaultExpr: ColumnExprCurrentTimestamp},
		},
		ForeignKeys: []ForeignKey{
			{Name: "fk_sys_role_permissions_role", Field: "roleId", RefTable: systemTableRoles, RefField: "id", Alias: "role"},
			{Name: "fk_sys_role_permissions_permission", Field: "permissionId", RefTable: systemTablePermissions, RefField: "id", Alias: "permission"},
		},
	}
}

func ensureSystemTable(database string, definition systemTableDefinition) error {
	db, err := newSimpleDB(database, definition.name)
	if err != nil {
		return fmt.Errorf("%w: 打开 %s 失败: %v", ErrSystemBootstrap, definition.name, err)
	}
	defer db.Close()

	plan, exists, err := db.SchemaDiff(definition.schema)
	if err != nil {
		return fmt.Errorf("%w: 校验 %s 失败: %v", ErrSystemBootstrap, definition.name, err)
	}
	if !exists {
		if err = db.CreateTable(definition.schema); err != nil {
			return fmt.Errorf("%w: 创建 %s 失败: %v", ErrSystemBootstrap, definition.name, err)
		}
		return nil
	}
	if plan != nil {
		payload, marshalErr := json.Marshal(plan)
		if marshalErr != nil {
			return fmt.Errorf("%w: %s", ErrSystemTableSchema, definition.name)
		}
		return fmt.Errorf("%w: %s diff=%s", ErrSystemTableSchema, definition.name, string(payload))
	}
	return nil
}

func ensureDefaultAdmin(database string) error {
	usersDB, err := newSimpleDB(database, systemTableUsers)
	if err != nil {
		return fmt.Errorf("%w: 打开 %s 失败: %v", ErrSystemBootstrap, systemTableUsers, err)
	}
	defer usersDB.Close()

	adminRow, ok, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: defaultSystemAdminUsername})
	if err != nil {
		return fmt.Errorf("%w: 查询默认管理员失败: %v", ErrSystemBootstrap, err)
	}
	if !ok {
		if _, err = usersDB.InsertRow(Row{
			"username":     defaultSystemAdminUsername,
			"passwordHash": digest.BcryptHash(defaultSystemAdminPassword),
			"displayName":  defaultSystemAdminName,
			"isAdmin":      true,
			"status":       defaultSystemStatus,
		}); err != nil {
			return fmt.Errorf("%w: 创建默认管理员失败: %v", ErrSystemBootstrap, err)
		}
		adminRow = nil
	}

	if err = ensureDefaultAdminRole(database, usersDB, adminRow, ok); err != nil {
		return err
	}

	return nil
}

func ensureDefaultAdminRole(database string, usersDB *SimpleDB, adminRow Row, adminExists bool) error {
	rolesDB, err := newSimpleDB(database, systemTableRoles)
	if err != nil {
		return fmt.Errorf("%w: 打开 %s 失败: %v", ErrSystemBootstrap, systemTableRoles, err)
	}
	defer rolesDB.Close()

	roleRow, ok, err := rolesDB.FindOne(QueryCondition{Field: "code", Operator: QueryOpEQ, Value: defaultSystemRoleCode})
	if err != nil {
		return fmt.Errorf("%w: 查询默认角色失败: %v", ErrSystemBootstrap, err)
	}
	if !ok {
		roleRow, err = rolesDB.InsertRow(Row{
			"code":        defaultSystemRoleCode,
			"name":        defaultSystemRoleName,
			"description": "built-in simpleDB administrator role",
			"isSystem":    true,
		})
		if err != nil {
			return fmt.Errorf("%w: 创建默认角色失败: %v", ErrSystemBootstrap, err)
		}
	}

	if !adminExists {
		adminRow, _, err = usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: defaultSystemAdminUsername})
		if err != nil {
			return fmt.Errorf("%w: 重新查询默认管理员失败: %v", ErrSystemBootstrap, err)
		}
	}

	userID, userOK := adminRow["id"]
	roleID, roleOK := roleRow["id"]
	if !userOK || !roleOK {
		return fmt.Errorf("%w: 默认管理员或角色缺少主键", ErrSystemBootstrap)
	}

	userRolesDB, err := newSimpleDB(database, systemTableUserRoles)
	if err != nil {
		return fmt.Errorf("%w: 打开 %s 失败: %v", ErrSystemBootstrap, systemTableUserRoles, err)
	}
	defer userRolesDB.Close()

	_, ok, err = userRolesDB.FindOne(
		QueryCondition{Field: "userId", Operator: QueryOpEQ, Value: userID},
		QueryCondition{Field: "roleId", Operator: QueryOpEQ, Value: roleID},
	)
	if err != nil {
		return fmt.Errorf("%w: 查询默认管理员角色关联失败: %v", ErrSystemBootstrap, err)
	}
	if ok {
		return nil
	}

	if _, err = userRolesDB.InsertRow(Row{"userId": userID, "roleId": roleID}); err != nil {
		return fmt.Errorf("%w: 创建默认管理员角色关联失败: %v", ErrSystemBootstrap, err)
	}

	return nil
}
