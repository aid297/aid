package kernal

import "strings"

func (*app) BindUserDatabase(database string, approver *AuthenticatedUser, username string) error {
	return bindUserDatabase(database, approver, username)
}

func (*app) CheckDatabaseBinding(database string, actor *AuthenticatedUser) error {
	return checkDatabaseBinding(database, actor)
}

func bindUserDatabase(database string, approver *AuthenticatedUser, username string) error {
	database = strings.TrimSpace(database)
	username = strings.TrimSpace(username)
	if database == "" || username == "" || approver == nil {
		return ErrInvalidTableAccessGrant
	}
	if !isSystemApprover(approver) {
		return ErrTableAccessDenied
	}
	if err := ensureSystemTables(database); err != nil {
		return err
	}

	userRow, exists, err := findUserByUsername(database, username)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}
	userID := rowString(userRow, "id")
	if userID == "" {
		return ErrUserNotFound
	}

	bindingsDB, err := newSimpleDB(systemDatabaseFor(database), systemTableUserDBBindings)
	if err != nil {
		return err
	}
	defer bindingsDB.Close()

	row, linked, err := bindingsDB.FindOne(
		QueryCondition{Field: "userId", Operator: QueryOpEQ, Value: userID},
		QueryCondition{Field: "databaseName", Operator: QueryOpEQ, Value: database},
	)
	if err != nil {
		return err
	}
	if !linked {
		_, err = bindingsDB.InsertRow(Row{"userId": userID, "databaseName": database, "enabled": true})
		return err
	}
	if rowBool(row, "enabled") {
		return nil
	}
	_, err = bindingsDB.UpdateRow(row["id"], Row{"enabled": true})
	return err
}

func checkDatabaseBinding(database string, actor *AuthenticatedUser) error {
	database = strings.TrimSpace(database)
	if database == "" {
		return ErrDBPathEmpty
	}
	if actor == nil {
		return nil
	}
	if isSystemApprover(actor) {
		return nil
	}
	if err := ensureSystemTables(database); err != nil {
		return err
	}

	bindingsDB, err := newSimpleDB(systemDatabaseFor(database), systemTableUserDBBindings)
	if err != nil {
		return err
	}
	defer bindingsDB.Close()

	row, linked, err := bindingsDB.FindOne(
		QueryCondition{Field: "userId", Operator: QueryOpEQ, Value: strings.TrimSpace(actor.ID)},
		QueryCondition{Field: "databaseName", Operator: QueryOpEQ, Value: database},
	)
	if err != nil {
		return err
	}
	if !linked || !rowBool(row, "enabled") {
		return ErrTableAccessDenied
	}
	return nil
}
