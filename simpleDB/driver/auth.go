package driver

import "github.com/aid297/aid/simpleDB/kernal"

type AuthenticatedUser = kernal.AuthenticatedUser

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
