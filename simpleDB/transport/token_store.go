package transport

import (
	"time"

	"github.com/aid297/aid/simpleDB/driver"
)

type TokenStore interface {
	SaveActive(database string, tokenID string, expiresAt int64) error
	GetActiveExpiresAt(database string, tokenID string) (int64, error)
	RemoveActive(database string, tokenID string) error

	SaveRevoked(database string, tokenID string, expiresAt int64) error
	LoadRevoked(database string, now int64) (map[string]int64, error)
	ClearExpired(database string, now int64) error
}

type DBTokenStore struct{}

func NewDBTokenStore() *DBTokenStore { return &DBTokenStore{} }

func (*DBTokenStore) SaveActive(database string, tokenID string, expiresAt int64) error {
	if err := driver.New.EnsureSystemTables(database); err != nil {
		return err
	}
	db, err := driver.New.DB(database, "_sys_active_tokens")
	if err != nil {
		return err
	}
	defer db.Close()
	row, exists, err := db.FindOne(driver.QueryCondition{Field: "tokenId", Operator: driver.QueryOpEQ, Value: tokenID})
	if err != nil {
		return err
	}
	if exists {
		_, err = db.UpdateRow(row["id"], driver.Row{"expiresAt": expiresAt})
		return err
	}
	_, err = db.InsertRow(driver.Row{"tokenId": tokenID, "expiresAt": expiresAt})
	return err
}

func (*DBTokenStore) GetActiveExpiresAt(database string, tokenID string) (int64, error) {
	if err := driver.New.EnsureSystemTables(database); err != nil {
		return 0, err
	}
	db, err := driver.New.DB(database, "_sys_active_tokens")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	row, exists, err := db.FindOne(
		driver.QueryCondition{Field: "tokenId", Operator: driver.QueryOpEQ, Value: tokenID},
	)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, nil
	}
	exp, _ := asInt64(row["expiresAt"])
	return exp, nil
}

func (*DBTokenStore) RemoveActive(database string, tokenID string) error {
	if err := driver.New.EnsureSystemTables(database); err != nil {
		return err
	}
	db, err := driver.New.DB(database, "_sys_active_tokens")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.RemoveByCondition(driver.QueryCondition{Field: "tokenId", Operator: driver.QueryOpEQ, Value: tokenID})
	return err
}

func (*DBTokenStore) SaveRevoked(database string, tokenID string, expiresAt int64) error {
	if err := driver.New.EnsureSystemTables(database); err != nil {
		return err
	}
	db, err := driver.New.DB(database, "_sys_revoked_tokens")
	if err != nil {
		return err
	}
	defer db.Close()
	row, exists, err := db.FindOne(driver.QueryCondition{Field: "tokenId", Operator: driver.QueryOpEQ, Value: tokenID})
	if err != nil {
		return err
	}
	if exists {
		_, err = db.UpdateRow(row["id"], driver.Row{"expiresAt": expiresAt})
		return err
	}
	_, err = db.InsertRow(driver.Row{"tokenId": tokenID, "expiresAt": expiresAt})
	return err
}

func (*DBTokenStore) LoadRevoked(database string, now int64) (map[string]int64, error) {
	if err := driver.New.EnsureSystemTables(database); err != nil {
		return nil, err
	}
	db, err := driver.New.DB(database, "_sys_revoked_tokens")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Find(driver.QueryCondition{Field: "expiresAt", Operator: driver.QueryOpGT, Value: now})
	if err != nil {
		return nil, err
	}
	m := make(map[string]int64, len(rows))
	for _, r := range rows {
		tokenID, _ := r["tokenId"].(string)
		exp, _ := asInt64(r["expiresAt"])
		m[tokenID] = exp
	}
	return m, nil
}

func (*DBTokenStore) ClearExpired(database string, now int64) error {
	if err := driver.New.EnsureSystemTables(database); err != nil {
		return err
	}
	// 1. 清理可用 Token 表中的过期项
	activeDB, err := driver.New.DB(database, "_sys_active_tokens")
	if err != nil {
		return err
	}
	defer activeDB.Close()
	_, err = activeDB.RemoveByCondition(driver.QueryCondition{Field: "expiresAt", Operator: driver.QueryOpLTE, Value: now})
	if err != nil {
		return err
	}

	// 2. 清理已撤销 Token 表中的过期项
	revokedDB, err := driver.New.DB(database, "_sys_revoked_tokens")
	if err != nil {
		return err
	}
	defer revokedDB.Close()
	_, err = revokedDB.RemoveByCondition(driver.QueryCondition{Field: "expiresAt", Operator: driver.QueryOpLTE, Value: now})
	return err
}

func asInt64(v any) (int64, bool) {
	switch t := v.(type) {
	case int:
		return int64(t), true
	case int32:
		return int64(t), true
	case int64:
		return t, true
	case float64:
		return int64(t), true
	case float32:
		return int64(t), true
	case time.Time:
		return t.Unix(), true
	default:
		return 0, false
	}
}
