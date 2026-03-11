package kernal

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type txWrite struct {
	value   []byte
	deleted bool
}

type TxIsolationLevel string

const (
	TxIsolationSnapshot      TxIsolationLevel = "snapshot"
	TxIsolationReadCommitted TxIsolationLevel = "read_committed"
)

type TxOptions struct {
	ReadOnly  bool
	Isolation TxIsolationLevel
}

type Tx struct {
	db        *SimpleDB
	snapshot  int64
	readOnly  bool
	isolation TxIsolationLevel
	readSet   map[string]int64
	writeSet  map[string]txWrite
	closed    bool
}

func normalizeTxIsolation(level TxIsolationLevel) TxIsolationLevel {
	switch level {
	case TxIsolationReadCommitted:
		return TxIsolationReadCommitted
	case TxIsolationSnapshot, "":
		return TxIsolationSnapshot
	default:
		return TxIsolationSnapshot
	}
}

func (tx *Tx) effectiveSnapshotLocked() int64 {
	if tx.isolation == TxIsolationReadCommitted {
		if tx.db.lastMVCCAt > 0 {
			return tx.db.lastMVCCAt
		}
		return time.Now().UnixNano()
	}
	return tx.snapshot
}

func (tx *Tx) getVisibleValueByKeyLocked(key string) ([]byte, bool, int64, error) {
	if pending, exists := tx.writeSet[key]; exists {
		tx.readSet[key] = tx.snapshot
		if pending.deleted {
			return nil, false, tx.snapshot, nil
		}
		return cloneBytes(pending.value), true, tx.snapshot, nil
	}

	value, exists, versionAt := tx.db.readValueAtSnapshotLocked(key, tx.snapshot)
	if tx.isolation == TxIsolationReadCommitted {
		value, exists, versionAt = tx.db.readValueAtSnapshotLocked(key, tx.effectiveSnapshotLocked())
	}
	tx.readSet[key] = versionAt
	if !exists {
		return nil, false, versionAt, nil
	}
	return value, true, versionAt, nil
}

func (tx *Tx) visibleRowKeysLocked() []string {
	keySet := make(map[string]struct{})
	for key := range tx.db.versions {
		if strings.HasPrefix(key, rowKeyPrefix) {
			keySet[key] = struct{}{}
		}
	}
	for key := range tx.writeSet {
		if strings.HasPrefix(key, rowKeyPrefix) {
			keySet[key] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (tx *Tx) ensureSchemaLocked() error {
	if tx.db.schema == nil {
		return ErrSchemaNotConfigured
	}
	return nil
}

func (tx *Tx) normalizePrimaryKeyLocked(primaryKey any) (Column, any, string, error) {
	pkColumn, ok := tx.db.columnByName(tx.db.schema.PrimaryKey)
	if !ok {
		return Column{}, nil, "", ErrPrimaryKeyMissing
	}
	normalizedPK, err := tx.db.normalizeColumnValueLocked(pkColumn, primaryKey)
	if err != nil {
		return Column{}, nil, "", err
	}
	token, err := valueToken(normalizedPK)
	if err != nil {
		return Column{}, nil, "", err
	}
	return pkColumn, normalizedPK, buildRowKey(token), nil
}

func (tx *Tx) visibleRowsLocked() ([]Row, map[string]string, error) {
	keys := tx.visibleRowKeysLocked()
	rows := make([]Row, 0, len(keys))
	rowKeyByPK := make(map[string]string, len(keys))

	for _, rowKey := range keys {
		raw, exists, _, err := tx.getVisibleValueByKeyLocked(rowKey)
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			continue
		}
		row, err := decodeRow(raw)
		if err != nil {
			return nil, nil, err
		}
		row, err = tx.db.normalizeRowValuesLocked(row)
		if err != nil {
			return nil, nil, err
		}
		pk := row[tx.db.schema.PrimaryKey]
		pkToken, err := valueToken(pk)
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, cloneRow(row))
		rowKeyByPK[pkToken] = rowKey
	}

	return rows, rowKeyByPK, nil
}

func (tx *Tx) checkUniqueConstraintsForRowLocked(row Row, selfPKToken string) error {
	rows, _, err := tx.visibleRowsLocked()
	if err != nil {
		return err
	}

	for _, column := range tx.db.schema.Columns {
		if !column.Unique {
			continue
		}

		value, exists := row[column.Name]
		if !exists || value == nil {
			continue
		}
		token, err := valueToken(value)
		if err != nil {
			return err
		}

		for _, candidate := range rows {
			candidatePKToken, err := valueToken(candidate[tx.db.schema.PrimaryKey])
			if err != nil {
				return err
			}
			if candidatePKToken == selfPKToken {
				continue
			}
			candidateValue, exists := candidate[column.Name]
			if !exists || candidateValue == nil {
				continue
			}
			candidateToken, err := valueToken(candidateValue)
			if err != nil {
				return err
			}
			if candidateToken == token {
				return ErrUniqueConflict
			}
		}
	}

	return nil
}

func (db *SimpleDB) appendVersionLocked(key string, value []byte, deleted bool, updatedAt int64) {
	db.versions[key] = append(db.versions[key], versionedEntry{
		Value:     cloneBytes(value),
		Deleted:   deleted,
		UpdatedAt: updatedAt,
	})
	if updatedAt > db.lastMVCCAt {
		db.lastMVCCAt = updatedAt
	}
}

func (db *SimpleDB) latestVersionTimestampLocked(key string) int64 {
	versions := db.versions[key]
	if len(versions) == 0 {
		return 0
	}
	return versions[len(versions)-1].UpdatedAt
}

func (db *SimpleDB) readValueAtSnapshotLocked(key string, snapshot int64) ([]byte, bool, int64) {
	versions := db.versions[key]
	for idx := len(versions) - 1; idx >= 0; idx-- {
		current := versions[idx]
		if current.UpdatedAt > snapshot {
			continue
		}
		if current.Deleted {
			return nil, false, current.UpdatedAt
		}
		return cloneBytes(current.Value), true, current.UpdatedAt
	}
	return nil, false, 0
}

func (db *SimpleDB) BeginTx() (*Tx, error) {
	return db.BeginTxWithOptions(TxOptions{})
}

func (db *SimpleDB) BeginReadOnlyTx() (*Tx, error) {
	return db.BeginTxWithOptions(TxOptions{ReadOnly: true})
}

func (db *SimpleDB) BeginTxWithOptions(options TxOptions) (*Tx, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}

	snapshot := db.lastMVCCAt
	if snapshot == 0 {
		snapshot = time.Now().UnixNano()
	}

	return &Tx{
		db:        db,
		snapshot:  snapshot,
		readOnly:  options.ReadOnly,
		isolation: normalizeTxIsolation(options.Isolation),
		readSet:   make(map[string]int64),
		writeSet:  make(map[string]txWrite),
	}, nil
}

func (db *SimpleDB) WithTx(fn func(tx *Tx) error) error {
	tx, err := db.BeginTx()
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tx *Tx) ensureOpen() error {
	if tx == nil || tx.closed {
		return ErrTxClosed
	}
	return nil
}

func (tx *Tx) Get(key string) ([]byte, bool, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, false, err
	}
	if err := validateKey(key); err != nil {
		return nil, false, err
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return nil, false, err
	}

	value, exists, _, err := tx.getVisibleValueByKeyLocked(key)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	return value, true, nil
}

func (tx *Tx) FindRow(primaryKey any) (Row, bool, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, false, err
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return nil, false, err
	}
	if err := tx.ensureSchemaLocked(); err != nil {
		return nil, false, err
	}

	_, _, rowKey, err := tx.normalizePrimaryKeyLocked(primaryKey)
	if err != nil {
		return nil, false, err
	}

	raw, exists, _, err := tx.getVisibleValueByKeyLocked(rowKey)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}

	row, err := decodeRow(raw)
	if err != nil {
		return nil, false, err
	}
	row, err = tx.db.normalizeRowValuesLocked(row)
	if err != nil {
		return nil, false, err
	}

	return cloneRow(row), true, nil
}

func (tx *Tx) Find(conditions ...QueryCondition) ([]Row, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, err
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return nil, err
	}
	if err := tx.ensureSchemaLocked(); err != nil {
		return nil, err
	}

	normalizedConditions, err := tx.db.normalizeQueryConditionsLocked(conditions)
	if err != nil {
		return nil, err
	}

	rows, _, err := tx.visibleRowsLocked()
	if err != nil {
		return nil, err
	}

	matchedRows := make([]Row, 0, len(rows))
	for _, row := range rows {
		matched, err := tx.db.rowMatchesConditionsLocked(row, normalizedConditions)
		if err != nil {
			return nil, err
		}
		if matched {
			matchedRows = append(matchedRows, cloneRow(row))
		}
	}

	return matchedRows, nil
}

func (tx *Tx) FindByConditions(conditions []QueryCondition) ([]Row, error) {
	return tx.Find(conditions...)
}

func (tx *Tx) FindOne(conditions ...QueryCondition) (Row, bool, error) {
	rows, err := tx.Find(conditions...)
	if err != nil {
		return nil, false, err
	}
	if len(rows) == 0 {
		return nil, false, nil
	}
	return rows[0], true, nil
}

func (tx *Tx) InsertRow(values Row) (Row, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, err
	}
	if tx.readOnly {
		return nil, ErrTxReadOnly
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return nil, err
	}
	if err := tx.ensureSchemaLocked(); err != nil {
		return nil, err
	}

	row := cloneRow(values)
	supplied := make(map[string]struct{}, len(values))
	for field := range row {
		if !tx.db.hasColumn(field) {
			return nil, fmt.Errorf("%w: %s", ErrFieldNotDefined, field)
		}
		supplied[field] = struct{}{}
	}

	pkField := tx.db.schema.PrimaryKey
	pkColumn, ok := tx.db.columnByName(pkField)
	if !ok {
		return nil, ErrPrimaryKeyMissing
	}

	pkValue, exists := row[pkField]
	if !exists || pkValue == nil {
		if !tx.db.schema.AutoIncrement {
			return nil, ErrPrimaryKeyMissing
		}
		switch ColumnType(normalizeColumnType(pkColumn.Type)) {
		case ColumnTypeUUID:
			version := extractUUIDVersion(pkColumn.Type)
			if version == 0 {
				version = tx.db.getDefaultUUIDVersion()
			}
			generatedUUID, err := generateUUIDByVersion(version)
			if err != nil {
				return nil, fmt.Errorf("failed to generate UUID v%d: %w", version, err)
			}
			row[pkField] = tx.db.formatUUIDLocked(generatedUUID)
		default:
			rows, _, err := tx.visibleRowsLocked()
			if err != nil {
				return nil, err
			}
			var maxSeq int64
			for _, current := range rows {
				if value, ok := asInt64(current[pkField]); ok && value > maxSeq {
					maxSeq = value
				}
			}
			if tx.db.autoSeq > maxSeq {
				maxSeq = tx.db.autoSeq
			}
			row[pkField] = maxSeq + 1
		}
	} else if tx.db.schema.AutoIncrement {
		switch ColumnType(normalizeColumnType(pkColumn.Type)) {
		case ColumnTypeUUID:
			uuidValue, err := tx.db.normalizeColumnValueLocked(pkColumn, pkValue)
			if err != nil {
				return nil, err
			}
			row[pkField] = uuidValue
		default:
			intValue, ok := asInt64(pkValue)
			if !ok {
				return nil, fmt.Errorf("%w: auto increment primary key must be integer", ErrFieldTypeMismatch)
			}
			row[pkField] = intValue
		}
	}

	normalizedRow, err := tx.db.applyColumnConstraintsLocked(row, supplied, true, true)
	if err != nil {
		return nil, err
	}

	pkToken, err := valueToken(normalizedRow[pkField])
	if err != nil {
		return nil, err
	}
	if err = tx.checkUniqueConstraintsForRowLocked(normalizedRow, pkToken); err != nil {
		return nil, err
	}

	encoded, err := encodeRow(normalizedRow)
	if err != nil {
		return nil, err
	}
	tx.writeSet[buildRowKey(pkToken)] = txWrite{value: encoded}

	return cloneRow(normalizedRow), nil
}

func (tx *Tx) InsertRows(values []Row) ([]Row, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, err
	}
	if tx.readOnly {
		return nil, ErrTxReadOnly
	}
	if len(values) == 0 {
		return nil, ErrBatchEmpty
	}

	rows := make([]Row, 0, len(values))
	for _, value := range values {
		row, err := tx.InsertRow(value)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (tx *Tx) UpdateRow(primaryKey any, updates Row) (Row, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, err
	}
	if tx.readOnly {
		return nil, ErrTxReadOnly
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return nil, err
	}
	if err := tx.ensureSchemaLocked(); err != nil {
		return nil, err
	}

	currentRow, exists, err := tx.FindRow(primaryKey)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrKeyNotFound
	}

	updatedRow := cloneRow(currentRow)
	for field, value := range updates {
		if !tx.db.hasColumn(field) {
			return nil, fmt.Errorf("%w: %s", ErrFieldNotDefined, field)
		}
		if field == tx.db.schema.PrimaryKey {
			currentToken, tokenErr := valueToken(currentRow[field])
			if tokenErr != nil {
				return nil, tokenErr
			}
			nextToken, tokenErr := valueToken(value)
			if tokenErr != nil {
				return nil, tokenErr
			}
			if currentToken != nextToken {
				return nil, ErrPrimaryKeyImmutable
			}
		}
		updatedRow[field] = value
	}

	tx.db.applyOnUpdateExpressionsLocked(updatedRow)
	normalizedRow, err := tx.db.normalizeRowValuesLocked(updatedRow)
	if err != nil {
		return nil, err
	}

	pkToken, err := valueToken(normalizedRow[tx.db.schema.PrimaryKey])
	if err != nil {
		return nil, err
	}
	if err = tx.checkUniqueConstraintsForRowLocked(normalizedRow, pkToken); err != nil {
		return nil, err
	}

	encoded, err := encodeRow(normalizedRow)
	if err != nil {
		return nil, err
	}
	tx.writeSet[buildRowKey(pkToken)] = txWrite{value: encoded}

	return cloneRow(normalizedRow), nil
}

func (tx *Tx) UpdateRows(updates []RowUpdate) ([]Row, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, err
	}
	if tx.readOnly {
		return nil, ErrTxReadOnly
	}
	if len(updates) == 0 {
		return nil, ErrBatchEmpty
	}

	rows := make([]Row, 0, len(updates))
	for _, update := range updates {
		row, err := tx.UpdateRow(update.PrimaryKey, update.Updates)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func (tx *Tx) DeleteRow(primaryKey any) error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if tx.readOnly {
		return ErrTxReadOnly
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return err
	}
	if err := tx.ensureSchemaLocked(); err != nil {
		return err
	}

	_, _, rowKey, err := tx.normalizePrimaryKeyLocked(primaryKey)
	if err != nil {
		return err
	}

	if _, exists, _, err := tx.getVisibleValueByKeyLocked(rowKey); err != nil {
		return err
	} else if !exists {
		return ErrKeyNotFound
	}

	tx.writeSet[rowKey] = txWrite{deleted: true}
	return nil
}

func (tx *Tx) DeleteRows(primaryKeys []any) error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if tx.readOnly {
		return ErrTxReadOnly
	}
	if len(primaryKeys) == 0 {
		return ErrBatchEmpty
	}

	for _, primaryKey := range primaryKeys {
		if err := tx.DeleteRow(primaryKey); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) RemoveByCondition(conditions ...QueryCondition) (int, error) {
	if err := tx.ensureOpen(); err != nil {
		return 0, err
	}
	if tx.readOnly {
		return 0, ErrTxReadOnly
	}

	rows, err := tx.Find(conditions...)
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}

	deletedCount := 0
	pkField := tx.db.schema.PrimaryKey
	for _, row := range rows {
		pkValue := row[pkField]
		if err = tx.DeleteRow(pkValue); err != nil {
			return 0, err
		}
		deletedCount++
	}

	return deletedCount, nil
}

func (tx *Tx) RemoveOneByCondition(conditions ...QueryCondition) (bool, error) {
	rows, err := tx.Find(conditions...)
	if err != nil {
		return false, err
	}
	if len(rows) == 0 {
		return false, nil
	}

	pkField := tx.db.schema.PrimaryKey
	if err = tx.DeleteRow(rows[0][pkField]); err != nil {
		return false, err
	}

	return true, nil
}

func (tx *Tx) Put(key string, value []byte) error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if tx.readOnly {
		return ErrTxReadOnly
	}
	if err := validateKey(key); err != nil {
		return err
	}

	tx.writeSet[key] = txWrite{value: cloneBytes(value)}
	return nil
}

func (tx *Tx) Delete(key string) error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	if tx.readOnly {
		return ErrTxReadOnly
	}
	if err := validateKey(key); err != nil {
		return err
	}

	if pending, exists := tx.writeSet[key]; exists {
		if pending.deleted {
			return nil
		}
		tx.writeSet[key] = txWrite{deleted: true}
		return nil
	}

	if _, exists, err := tx.Get(key); err != nil {
		return err
	} else if !exists {
		return ErrKeyNotFound
	}

	tx.writeSet[key] = txWrite{deleted: true}
	return nil
}

func (tx *Tx) Query(prefix string) (map[string][]byte, error) {
	if err := tx.ensureOpen(); err != nil {
		return nil, err
	}

	tx.db.mu.RLock()
	defer tx.db.mu.RUnlock()

	if err := tx.db.ensureOpen(); err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for key := range tx.db.versions {
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}

		if pending, exists := tx.writeSet[key]; exists {
			tx.readSet[key] = tx.snapshot
			if pending.deleted {
				continue
			}
			result[key] = cloneBytes(pending.value)
			continue
		}

		value, exists, versionAt := tx.db.readValueAtSnapshotLocked(key, tx.snapshot)
		tx.readSet[key] = versionAt
		if !exists {
			continue
		}
		result[key] = value
	}

	for key, pending := range tx.writeSet {
		if _, exists := result[key]; exists {
			continue
		}
		if pending.deleted {
			continue
		}
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}
		result[key] = cloneBytes(pending.value)
	}

	return result, nil
}

func (tx *Tx) Keys() ([]string, error) {
	rows, err := tx.Query("")
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(rows))
	for key := range rows {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

func (tx *Tx) Commit() error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}

	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()

	if err := tx.db.ensureOpen(); err != nil {
		return err
	}

	if tx.isolation == TxIsolationSnapshot {
		for key := range tx.readSet {
			if tx.db.latestVersionTimestampLocked(key) > tx.snapshot {
				return ErrTxConflict
			}
		}
	}
	for key := range tx.writeSet {
		if tx.db.latestVersionTimestampLocked(key) > tx.snapshot {
			return ErrTxConflict
		}
	}

	if len(tx.writeSet) == 0 {
		tx.closed = true
		return nil
	}

	if tx.db.schema != nil {
		localUnique := make(map[string]map[string]string)
		for key, write := range tx.writeSet {
			if write.deleted || !strings.HasPrefix(key, rowKeyPrefix) {
				continue
			}

			row, err := decodeRow(write.value)
			if err != nil {
				return err
			}
			row, err = tx.db.normalizeRowValuesLocked(row)
			if err != nil {
				return err
			}

			pkValue, exists := row[tx.db.schema.PrimaryKey]
			if !exists {
				return ErrPrimaryKeyMissing
			}
			pkToken, err := valueToken(pkValue)
			if err != nil {
				return err
			}

			if err = tx.db.checkUniqueConstraintsLocked(row, pkToken); err != nil {
				return ErrTxConflict
			}

			for _, column := range tx.db.schema.Columns {
				if !column.Unique {
					continue
				}
				value, exists := row[column.Name]
				if !exists || value == nil {
					continue
				}
				token, err := valueToken(value)
				if err != nil {
					return err
				}
				if localUnique[column.Name] == nil {
					localUnique[column.Name] = make(map[string]string)
				}
				if otherPK, exists := localUnique[column.Name][token]; exists && otherPK != pkToken {
					return ErrTxConflict
				}
				localUnique[column.Name][token] = pkToken
			}
		}
	}

	commitAt := time.Now().UnixNano()
	if commitAt <= tx.db.lastMVCCAt {
		commitAt = tx.db.lastMVCCAt + 1
	}

	for key, write := range tx.writeSet {
		var err error
		if write.deleted {
			err = tx.db.deleteRawAtLocked(key, commitAt)
		} else {
			err = tx.db.putRawAtLocked(key, write.value, commitAt)
		}
		if err != nil {
			return err
		}
	}

	if err := tx.db.rebuildStructuredState(); err != nil {
		return err
	}

	tx.closed = true
	return nil
}

func (tx *Tx) Rollback() error {
	if err := tx.ensureOpen(); err != nil {
		return err
	}
	tx.closed = true
	tx.writeSet = nil
	tx.readSet = nil
	return nil
}
