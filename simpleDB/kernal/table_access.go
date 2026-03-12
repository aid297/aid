package kernal

import (
	"strings"
)

type TableAccessScope string

const (
	TableAccessScopeDML TableAccessScope = "dml"
	TableAccessScopeDDL TableAccessScope = "ddl"
)

type TableAccessGrant struct {
	TableName     string           `json:"tableName"`
	GranteeUserID string           `json:"granteeUserId"`
	Scope         TableAccessScope `json:"scope"`
	OwnerApproved bool             `json:"ownerApproved"`
	AdminApproved bool             `json:"adminApproved"`
}

func (*app) RegisterTableOwner(database, tableName, ownerUserID string) error {
	return registerTableOwner(database, tableName, ownerUserID)
}

func (*app) CheckTableAccess(database, tableName string, actor *AuthenticatedUser, scope TableAccessScope) error {
	return checkTableAccess(database, tableName, actor, scope)
}

func (*app) ApproveTableAccess(database string, approver *AuthenticatedUser, tableName, granteeUsername string, scope TableAccessScope) (*TableAccessGrant, error) {
	return approveTableAccess(database, approver, tableName, granteeUsername, scope)
}

func registerTableOwner(database, tableName, ownerUserID string) error {
	tableName = normalizeTableName(tableName)
	ownerUserID = strings.TrimSpace(ownerUserID)
	if tableName == "" || ownerUserID == "" {
		return ErrInvalidTableAccessGrant
	}
	if err := ensureSystemTables(database); err != nil {
		return err
	}
	ownersDB, err := newSimpleDB(systemDatabaseFor(database), systemTableTableOwners)
	if err != nil {
		return err
	}
	defer ownersDB.Close()

	ownerRow, exists, err := ownersDB.FindOne(QueryCondition{Field: "tableName", Operator: QueryOpEQ, Value: tableName})
	if err != nil {
		return err
	}
	if exists {
		if rowString(ownerRow, "ownerUserId") == ownerUserID {
			return nil
		}
		return ErrTableOwnerAlreadyAssigned
	}

	_, err = ownersDB.InsertRow(Row{
		"tableName":   tableName,
		"ownerUserId": ownerUserID,
	})
	return err
}

func checkTableAccess(database, tableName string, actor *AuthenticatedUser, scope TableAccessScope) error {
	tableName = normalizeTableName(tableName)
	scope = normalizeTableAccessScope(scope)
	if tableName == "" {
		return ErrTableAccessDenied
	}
	if strings.HasPrefix(tableName, "_sys_") {
		if actor != nil && isSystemApprover(actor) {
			return nil
		}
		return ErrTableAccessDenied
	}
	if actor == nil {
		return nil
	}
	if isSystemApprover(actor) {
		return nil
	}
	if !isValidTableAccessScope(scope) {
		return ErrInvalidTableAccessGrant
	}
	if err := ensureSystemTables(database); err != nil {
		return err
	}

	ownerRow, exists, err := findTableOwner(database, tableName)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	if rowString(ownerRow, "ownerUserId") == strings.TrimSpace(actor.ID) {
		return nil
	}

	grantsDB, err := newSimpleDB(systemDatabaseFor(database), systemTableTableAccessGrants)
	if err != nil {
		return err
	}
	defer grantsDB.Close()

	grantRow, approved, err := grantsDB.FindOne(
		QueryCondition{Field: "tableName", Operator: QueryOpEQ, Value: tableName},
		QueryCondition{Field: "granteeUserId", Operator: QueryOpEQ, Value: actor.ID},
		QueryCondition{Field: "scope", Operator: QueryOpEQ, Value: string(scope)},
	)
	if err != nil {
		return err
	}
	if !approved {
		return ErrTableAccessDenied
	}
	if rowBool(grantRow, "ownerApproved") && rowBool(grantRow, "adminApproved") {
		return nil
	}
	return ErrTableAccessDenied
}

func approveTableAccess(database string, approver *AuthenticatedUser, tableName, granteeUsername string, scope TableAccessScope) (*TableAccessGrant, error) {
	tableName = normalizeTableName(tableName)
	scope = normalizeTableAccessScope(scope)
	granteeUsername = strings.TrimSpace(granteeUsername)
	if approver == nil || tableName == "" || granteeUsername == "" || !isValidTableAccessScope(scope) {
		return nil, ErrInvalidTableAccessGrant
	}
	if strings.HasPrefix(tableName, "_sys_") {
		return nil, ErrTableAccessDenied
	}
	if err := ensureSystemTables(database); err != nil {
		return nil, err
	}

	ownerRow, exists, err := findTableOwner(database, tableName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTableOwnerNotFound
	}
	ownerUserID := rowString(ownerRow, "ownerUserId")
	isOwnerApprover := ownerUserID != "" && ownerUserID == strings.TrimSpace(approver.ID)
	isAdminApprover := isSystemApprover(approver)
	if !isOwnerApprover && !isAdminApprover {
		return nil, ErrTableAccessDenied
	}

	granteeRow, exists, err := findUserByUsername(database, granteeUsername)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}
	granteeUserID := rowString(granteeRow, "id")
	if granteeUserID == "" {
		return nil, ErrUserNotFound
	}
	if granteeUserID == ownerUserID {
		return &TableAccessGrant{TableName: tableName, GranteeUserID: granteeUserID, Scope: scope, OwnerApproved: true, AdminApproved: true}, nil
	}

	grantsDB, err := newSimpleDB(systemDatabaseFor(database), systemTableTableAccessGrants)
	if err != nil {
		return nil, err
	}
	defer grantsDB.Close()

	grantRow, exists, err := grantsDB.FindOne(
		QueryCondition{Field: "tableName", Operator: QueryOpEQ, Value: tableName},
		QueryCondition{Field: "granteeUserId", Operator: QueryOpEQ, Value: granteeUserID},
		QueryCondition{Field: "scope", Operator: QueryOpEQ, Value: string(scope)},
	)
	if err != nil {
		return nil, err
	}
	if !exists {
		inserted, insertErr := grantsDB.InsertRow(Row{
			"tableName":     tableName,
			"granteeUserId": granteeUserID,
			"scope":         string(scope),
			"ownerApproved": isOwnerApprover,
			"adminApproved": isAdminApprover,
		})
		if insertErr != nil {
			return nil, insertErr
		}
		return grantRowToModel(inserted, scope), nil
	}

	updates := Row{}
	if isOwnerApprover {
		updates["ownerApproved"] = true
	}
	if isAdminApprover {
		updates["adminApproved"] = true
	}
	updated, err := grantsDB.UpdateRow(grantRow["id"], updates)
	if err != nil {
		return nil, err
	}
	return grantRowToModel(updated, scope), nil
}

func findTableOwner(database, tableName string) (Row, bool, error) {
	ownersDB, err := newSimpleDB(systemDatabaseFor(database), systemTableTableOwners)
	if err != nil {
		return nil, false, err
	}
	defer ownersDB.Close()
	return ownersDB.FindOne(QueryCondition{Field: "tableName", Operator: QueryOpEQ, Value: normalizeTableName(tableName)})
}

func findUserByUsername(database, username string) (Row, bool, error) {
	usersDB, err := newSimpleDB(systemDatabaseFor(database), systemTableUsers)
	if err != nil {
		return nil, false, err
	}
	defer usersDB.Close()
	return usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: strings.TrimSpace(username)})
}

func grantRowToModel(row Row, scope TableAccessScope) *TableAccessGrant {
	return &TableAccessGrant{
		TableName:     rowString(row, "tableName"),
		GranteeUserID: rowString(row, "granteeUserId"),
		Scope:         scope,
		OwnerApproved: rowBool(row, "ownerApproved"),
		AdminApproved: rowBool(row, "adminApproved"),
	}
}

func normalizeTableName(tableName string) string {
	return strings.TrimSpace(tableName)
}

func isValidTableAccessScope(scope TableAccessScope) bool {
	switch normalizeTableAccessScope(scope) {
	case TableAccessScopeDML, TableAccessScopeDDL:
		return true
	default:
		return false
	}
}

func normalizeTableAccessScope(scope TableAccessScope) TableAccessScope {
	return TableAccessScope(strings.ToLower(strings.TrimSpace(string(scope))))
}

func isSystemApprover(actor *AuthenticatedUser) bool {
	if actor == nil {
		return false
	}
	if strings.TrimSpace(actor.Username) == defaultSystemAdminUsername {
		return true
	}
	return hasRoleCode(actor.Roles, defaultSystemRoleCode)
}

func hasRoleCode(roles []string, roleCode string) bool {
	for _, role := range roles {
		if strings.TrimSpace(role) == roleCode {
			return true
		}
	}
	return false
}
