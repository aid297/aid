package kernal

import (
	"errors"
	"strings"

	"github.com/aid297/aid/digest"
)

const defaultInactiveStatus = "inactive"

func (*app) RegisterUser(database, username, password, displayName string) (*AuthenticatedUser, error) {
	return registerUser(database, username, password, displayName)
}

func (*app) ActivateUser(database, username string) (*AuthenticatedUser, error) {
	return activateUser(database, username)
}

func (*app) DeactivateUser(database, username string) (*AuthenticatedUser, error) {
	return deactivateUser(database, username)
}

func (*app) AssignRoles(database, username string, roleCodes []string) (*AuthenticatedUser, error) {
	return assignRoles(database, username, roleCodes)
}

func (*app) InitSDBPassword(database string) error {
	return initSDBPassword(database)
}

func registerUser(database, username, password, displayName string) (*AuthenticatedUser, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	displayName = strings.TrimSpace(displayName)
	if username == "" || password == "" {
		return nil, ErrInvalidRegistration
	}

	if err := ensureSystemTables(database); err != nil {
		return nil, err
	}

	usersDB, err := newSimpleDB(database, systemTableUsers)
	if err != nil {
		return nil, err
	}
	defer usersDB.Close()

	_, exists, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: username})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	inserted, err := usersDB.InsertRow(Row{
		"username":     username,
		"passwordHash": digest.BcryptHash(password),
		"displayName":  displayName,
		"isAdmin":      false,
		"status":       defaultInactiveStatus,
	})
	if err != nil {
		if errors.Is(err, ErrUniqueConflict) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return &AuthenticatedUser{
		ID:          rowString(inserted, "id"),
		Username:    rowString(inserted, "username"),
		DisplayName: rowString(inserted, "displayName"),
		Status:      rowString(inserted, "status"),
		IsAdmin:     rowBool(inserted, "isAdmin"),
	}, nil
}

func activateUser(database, username string) (*AuthenticatedUser, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrUserNotFound
	}

	if err := ensureSystemTables(database); err != nil {
		return nil, err
	}

	usersDB, err := newSimpleDB(database, systemTableUsers)
	if err != nil {
		return nil, err
	}
	defer usersDB.Close()

	row, exists, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: username})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	userID, ok := row["id"]
	if !ok {
		return nil, ErrUserNotFound
	}

	updated, err := usersDB.UpdateRow(userID, Row{"status": defaultSystemStatus})
	if err != nil {
		return nil, err
	}

	roles, permissions, err := collectUserAccess(database, userID)
	if err != nil {
		return nil, err
	}

	return &AuthenticatedUser{
		ID:          rowString(updated, "id"),
		Username:    rowString(updated, "username"),
		DisplayName: rowString(updated, "displayName"),
		Status:      rowString(updated, "status"),
		IsAdmin:     rowBool(updated, "isAdmin"),
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func deactivateUser(database, username string) (*AuthenticatedUser, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrUserNotFound
	}

	if err := ensureSystemTables(database); err != nil {
		return nil, err
	}

	usersDB, err := newSimpleDB(database, systemTableUsers)
	if err != nil {
		return nil, err
	}
	defer usersDB.Close()

	row, exists, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: username})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	userID, ok := row["id"]
	if !ok {
		return nil, ErrUserNotFound
	}

	updated, err := usersDB.UpdateRow(userID, Row{"status": defaultInactiveStatus})
	if err != nil {
		return nil, err
	}

	return &AuthenticatedUser{
		ID:          rowString(updated, "id"),
		Username:    rowString(updated, "username"),
		DisplayName: rowString(updated, "displayName"),
		Status:      rowString(updated, "status"),
		IsAdmin:     rowBool(updated, "isAdmin"),
	}, nil
}

func assignRoles(database, username string, roleCodes []string) (*AuthenticatedUser, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, ErrUserNotFound
	}

	cleanRoleCodes := normalizeRoleCodes(roleCodes)
	if len(cleanRoleCodes) == 0 {
		return nil, ErrRoleNotFound
	}

	if err := ensureSystemTables(database); err != nil {
		return nil, err
	}

	usersDB, err := newSimpleDB(database, systemTableUsers)
	if err != nil {
		return nil, err
	}
	defer usersDB.Close()

	userRow, exists, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: username})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}
	userID, userIDOK := userRow["id"]
	if !userIDOK {
		return nil, ErrUserNotFound
	}

	rolesDB, err := newSimpleDB(database, systemTableRoles)
	if err != nil {
		return nil, err
	}
	defer rolesDB.Close()

	userRolesDB, err := newSimpleDB(database, systemTableUserRoles)
	if err != nil {
		return nil, err
	}
	defer userRolesDB.Close()

	for _, roleCode := range cleanRoleCodes {
		if roleCode == defaultSystemRoleCode && username != defaultSystemAdminUsername {
			return nil, ErrSuperAdminRoleReserved
		}

		roleRow, found, roleErr := rolesDB.FindOne(QueryCondition{Field: "code", Operator: QueryOpEQ, Value: roleCode})
		if roleErr != nil {
			return nil, roleErr
		}
		if !found {
			return nil, ErrRoleNotFound
		}

		roleID, roleIDOK := roleRow["id"]
		if !roleIDOK {
			return nil, ErrRoleNotFound
		}

		_, linked, linkErr := userRolesDB.FindOne(
			QueryCondition{Field: "userId", Operator: QueryOpEQ, Value: userID},
			QueryCondition{Field: "roleId", Operator: QueryOpEQ, Value: roleID},
		)
		if linkErr != nil {
			return nil, linkErr
		}
		if linked {
			continue
		}

		if _, linkErr = userRolesDB.InsertRow(Row{"userId": userID, "roleId": roleID}); linkErr != nil {
			return nil, linkErr
		}
	}

	_ = rolesDB.Close()
	_ = userRolesDB.Close()

	roles, permissions, err := collectUserAccess(database, userID)
	if err != nil {
		return nil, err
	}

	status := rowString(userRow, "status")
	if status == "" {
		status = defaultSystemStatus
	}

	return &AuthenticatedUser{
		ID:          rowString(userRow, "id"),
		Username:    rowString(userRow, "username"),
		DisplayName: rowString(userRow, "displayName"),
		Status:      status,
		IsAdmin:     rowBool(userRow, "isAdmin"),
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func initSDBPassword(database string) error {
	if err := ensureSystemTables(database); err != nil {
		return err
	}

	usersDB, err := newSimpleDB(database, systemTableUsers)
	if err != nil {
		return err
	}
	defer usersDB.Close()

	row, exists, err := usersDB.FindOne(QueryCondition{Field: "username", Operator: QueryOpEQ, Value: defaultSystemAdminUsername})
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}

	userID, ok := row["id"]
	if !ok {
		return ErrUserNotFound
	}

	_, err = usersDB.UpdateRow(userID, Row{
		"passwordHash": digest.BcryptHash(defaultSystemAdminPassword),
		"status":       defaultSystemStatus,
		"isAdmin":      true,
	})
	if err != nil {
		return err
	}

	return nil
}

func normalizeRoleCodes(roleCodes []string) []string {
	result := make([]string, 0, len(roleCodes))
	seen := make(map[string]struct{}, len(roleCodes))
	for _, code := range roleCodes {
		trimmed := strings.TrimSpace(code)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}
