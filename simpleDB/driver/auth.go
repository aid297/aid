package driver

import "github.com/aid297/aid/simpleDB/kernal"

type AuthenticatedUser = kernal.AuthenticatedUser
type TableAccessScope = kernal.TableAccessScope
type TableAccessGrant = kernal.TableAccessGrant

const (
	TableAccessScopeDML = kernal.TableAccessScopeDML
	TableAccessScopeDDL = kernal.TableAccessScopeDDL
)

func (*app) ApproveTableAccess(database string, approver *AuthenticatedUser, tableName, granteeUsername string, scope TableAccessScope) (*TableAccessGrant, error) {
	grant, err := kernal.New.ApproveTableAccess(database, approver, tableName, granteeUsername, scope)
	if err != nil {
		return nil, wrapError(err)
	}
	return grant, nil
}

func (*app) BindUserDatabase(database string, approver *AuthenticatedUser, username string) error {
	err := kernal.New.BindUserDatabase(database, approver, username)
	if err != nil {
		return wrapError(err)
	}
	return nil
}

func (*app) RevokeUserDatabase(database string, approver *AuthenticatedUser, username string) error {
	err := kernal.New.RevokeUserDatabase(database, approver, username)
	if err != nil {
		return wrapError(err)
	}
	return nil
}

func (*app) Authenticate(database, username, password string) (*AuthenticatedUser, error) {
	user, err := kernal.New.Authenticate(database, username, password)
	if err != nil {
		return nil, wrapError(err)
	}
	return user, nil
}

func (*app) RegisterUser(database, username, password, displayName string) (*AuthenticatedUser, error) {
	user, err := kernal.New.RegisterUser(database, username, password, displayName)
	if err != nil {
		return nil, wrapError(err)
	}
	return user, nil
}

func (*app) ActivateUser(database, username string) (*AuthenticatedUser, error) {
	user, err := kernal.New.ActivateUser(database, username)
	if err != nil {
		return nil, wrapError(err)
	}
	return user, nil
}

func (*app) DeactivateUser(database, username string) (*AuthenticatedUser, error) {
	user, err := kernal.New.DeactivateUser(database, username)
	if err != nil {
		return nil, wrapError(err)
	}
	return user, nil
}

func (*app) AssignRoles(database, username string, roleCodes []string) (*AuthenticatedUser, error) {
	user, err := kernal.New.AssignRoles(database, username, roleCodes)
	if err != nil {
		return nil, wrapError(err)
	}
	return user, nil
}

func (*app) AssignRolePermissions(database, roleCode string, permissionCodes []string) error {
	err := kernal.New.AssignRolePermissions(database, roleCode, permissionCodes)
	if err != nil {
		return wrapError(err)
	}
	return nil
}

func (*app) InitSDBPassword(database string) error {
	err := kernal.New.InitSDBPassword(database)
	if err != nil {
		return wrapError(err)
	}
	return nil
}
