package kernal

import (
	"sort"
	"strings"

	"github.com/aid297/aid/digest"
)

type AuthenticatedUser struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName,omitempty"`
	Status      string   `json:"status,omitempty"`
	IsAdmin     bool     `json:"isAdmin"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func (*app) Authenticate(database, username, password string) (*AuthenticatedUser, error) {
	return authenticate(database, username, password)
}

func authenticate(database, username, password string) (*AuthenticatedUser, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	if err := ensureSystemTables(database); err != nil {
		return nil, err
	}

	usersDB, err := newSimpleDB(database, systemTableUsers)
	if err != nil {
		return nil, err
	}
	defer usersDB.Close()

	row, ok, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: username})
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}

	passwordHash, _ := row["passwordHash"].(string)
	if passwordHash == "" || !digest.BcryptCheck(password, passwordHash) {
		return nil, ErrInvalidCredentials
	}

	user := &AuthenticatedUser{
		ID:          rowString(row, "id"),
		Username:    rowString(row, "username"),
		DisplayName: rowString(row, "displayName"),
		Status:      rowString(row, "status"),
		IsAdmin:     rowBool(row, "isAdmin"),
	}
	if user.Status == "" {
		user.Status = defaultSystemStatus
	}
	if user.Status != defaultSystemStatus {
		return nil, ErrUserInactive
	}

	roles, permissions, err := collectUserAccess(database, row["id"])
	if err != nil {
		return nil, err
	}
	user.Roles = roles
	user.Permissions = permissions

	return user, nil
}

func collectUserAccess(database string, userID any) ([]string, []string, error) {
	userRolesDB, err := newSimpleDB(database, systemTableUserRoles)
	if err != nil {
		return nil, nil, err
	}
	defer userRolesDB.Close()

	userRoleRows, err := userRolesDB.Find(QueryCondition{Field: "userId", Operator: QueryOpEQ, Value: userID})
	if err != nil {
		return nil, nil, err
	}
	if len(userRoleRows) == 0 {
		return nil, nil, nil
	}

	rolesDB, err := newSimpleDB(database, systemTableRoles)
	if err != nil {
		return nil, nil, err
	}
	defer rolesDB.Close()

	rolePermissionsDB, err := newSimpleDB(database, systemTableRolePermissions)
	if err != nil {
		return nil, nil, err
	}
	defer rolePermissionsDB.Close()

	permissionsDB, err := newSimpleDB(database, systemTablePermissions)
	if err != nil {
		return nil, nil, err
	}
	defer permissionsDB.Close()

	roleSet := make(map[string]struct{})
	permissionSet := make(map[string]struct{})

	for _, userRoleRow := range userRoleRows {
		roleID, ok := userRoleRow["roleId"]
		if !ok {
			continue
		}

		roleRow, found, err := rolesDB.FindRow(roleID)
		if err != nil {
			return nil, nil, err
		}
		if found {
			if code := rowString(roleRow, "code"); code != "" {
				roleSet[code] = struct{}{}
			}
		}

		permissionLinks, err := rolePermissionsDB.Find(QueryCondition{Field: "roleId", Operator: QueryOpEQ, Value: roleID})
		if err != nil {
			return nil, nil, err
		}
		for _, permissionLink := range permissionLinks {
			permissionID, exists := permissionLink["permissionId"]
			if !exists {
				continue
			}
			permissionRow, found, err := permissionsDB.FindRow(permissionID)
			if err != nil {
				return nil, nil, err
			}
			if found {
				if code := rowString(permissionRow, "code"); code != "" {
					permissionSet[code] = struct{}{}
				}
			}
		}
	}

	roles := make([]string, 0, len(roleSet))
	for code := range roleSet {
		roles = append(roles, code)
	}
	sort.Strings(roles)

	permissions := make([]string, 0, len(permissionSet))
	for code := range permissionSet {
		permissions = append(permissions, code)
	}
	sort.Strings(permissions)

	return roles, permissions, nil
}

func rowString(row Row, field string) string {
	value, ok := row[field]
	if !ok || value == nil {
		return ""
	}
	text, _ := value.(string)
	return text
}

func rowBool(row Row, field string) bool {
	value, ok := row[field]
	if !ok || value == nil {
		return false
	}
	b, _ := value.(bool)
	return b
}
