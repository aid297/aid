package kernal

import "strings"

func (*app) BindUserDatabase(database string, approver *AuthenticatedUser, username string) error {
	return bindUserDatabase(database, approver, username)
}

func (*app) RevokeUserDatabase(database string, approver *AuthenticatedUser, username string) error {
	return revokeUserDatabase(database, approver, username)
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
	if err := ensureSystemTables(database); err != nil {
		return err
	}
	if err := checkDatabaseOperator(database, approver); err != nil {
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

	ownerRow, ownerExists, err := findDatabaseOwner(database)
	if err != nil {
		return err
	}
	if !ownerExists {
		if !isSystemApprover(approver) {
			return ErrTableAccessDenied
		}
		if err = assignDatabaseOwner(database, userID); err != nil {
			return err
		}
	} else {
		ownerUserID := rowString(ownerRow, "ownerUserId")
		if ownerUserID != "" && ownerUserID == userID {
			return nil
		}
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

func revokeUserDatabase(database string, approver *AuthenticatedUser, username string) error {
	database = strings.TrimSpace(database)
	username = strings.TrimSpace(username)
	if database == "" || username == "" || approver == nil {
		return ErrInvalidTableAccessGrant
	}
	if err := ensureSystemTables(database); err != nil {
		return err
	}
	if err := checkDatabaseOperator(database, approver); err != nil {
		return err
	}

	userRow, exists, err := findUserByUsername(database, username)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}
	targetUserID := rowString(userRow, "id")
	if targetUserID == "" {
		return ErrUserNotFound
	}

	ownerRow, ownerExists, err := findDatabaseOwner(database)
	if err != nil {
		return err
	}
	if ownerExists && rowString(ownerRow, "ownerUserId") == targetUserID {
		return ErrTableAccessDenied
	}

	bindingsDB, err := newSimpleDB(systemDatabaseFor(database), systemTableUserDBBindings)
	if err != nil {
		return err
	}
	defer bindingsDB.Close()

	row, linked, err := bindingsDB.FindOne(
		QueryCondition{Field: "userId", Operator: QueryOpEQ, Value: targetUserID},
		QueryCondition{Field: "databaseName", Operator: QueryOpEQ, Value: database},
	)
	if err != nil {
		return err
	}
	if !linked {
		return nil
	}
	if !rowBool(row, "enabled") {
		return nil
	}
	_, err = bindingsDB.UpdateRow(row["id"], Row{"enabled": false})
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

func checkDatabaseOperator(database string, approver *AuthenticatedUser) error {
	if approver == nil {
		return ErrTableAccessDenied
	}
	if isSystemApprover(approver) {
		return nil
	}
	ownerRow, exists, err := findDatabaseOwner(database)
	if err != nil {
		return err
	}
	if !exists {
		return ErrTableAccessDenied
	}
	if rowString(ownerRow, "ownerUserId") == strings.TrimSpace(approver.ID) {
		return nil
	}
	return ErrTableAccessDenied
}

func findDatabaseOwner(database string) (Row, bool, error) {
	ownersDB, err := newSimpleDB(systemDatabaseFor(database), systemTableDatabaseOwners)
	if err != nil {
		return nil, false, err
	}
	defer ownersDB.Close()
	return ownersDB.FindOne(QueryCondition{Field: "databaseName", Operator: QueryOpEQ, Value: strings.TrimSpace(database)})
}

func assignDatabaseOwner(database, ownerUserID string) error {
	ownersDB, err := newSimpleDB(systemDatabaseFor(database), systemTableDatabaseOwners)
	if err != nil {
		return err
	}
	defer ownersDB.Close()

	row, exists, err := ownersDB.FindOne(QueryCondition{Field: "databaseName", Operator: QueryOpEQ, Value: strings.TrimSpace(database)})
	if err != nil {
		return err
	}
	if !exists {
		_, err = ownersDB.InsertRow(Row{"databaseName": strings.TrimSpace(database), "ownerUserId": strings.TrimSpace(ownerUserID)})
		return err
	}
	if rowString(row, "ownerUserId") == strings.TrimSpace(ownerUserID) {
		return nil
	}
	return ErrTableOwnerAlreadyAssigned
}
