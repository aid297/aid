package simpleDBDriver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	metaSchemaKey   = "__meta__:schema"
	metaSequenceKey = "__meta__:sequence"
	rowKeyPrefix    = "__row__:"
	timeLayout      = "15:04:05.999999999Z07:00"
)

type Row map[string]any

type RowUpdate struct {
	PrimaryKey any `json:"primaryKey"`
	Updates    Row `json:"updates"`
}

type ColumnType string

const (
	ColumnTypeAny       ColumnType = "any"
	ColumnTypeString    ColumnType = "string"
	ColumnTypeInt       ColumnType = "int"
	ColumnTypeFloat     ColumnType = "float"
	ColumnTypeBool      ColumnType = "bool"
	ColumnTypeObject    ColumnType = "object"
	ColumnTypeArray     ColumnType = "array"
	ColumnTypeUUID      ColumnType = "uuid"
	ColumnTypeTime      ColumnType = "time"
	ColumnTypeTimestamp ColumnType = "timestamp"
)

const (
	DefaultUUIDVersion       = 6
	DefaultCascadeMaxDepth   = 6
	HardCascadeMaxDepthLimit = 6
)

type DatabaseConfig struct {
	DefaultUUIDVersion     int `json:"defaultUUIDVersion,omitempty"`
	DefaultCascadeMaxDepth int `json:"defaultCascadeMaxDepth,omitempty"`
}

const (
	ColumnExprCurrentTime      = "current_time"
	ColumnExprCurrentTimestamp = "current_timestamp"
)

const (
	ColumnCheckGT     = "gt"
	ColumnCheckGTE    = "gte"
	ColumnCheckLT     = "lt"
	ColumnCheckLTE    = "lte"
	ColumnCheckLenGT  = "len_gt"
	ColumnCheckLenGTE = "len_gte"
	ColumnCheckLenLT  = "len_lt"
	ColumnCheckLenLTE = "len_lte"
	ColumnCheckRegex  = "regex"
)

const (
	QueryOpEQ         = "eq"
	QueryOpNE         = "ne"
	QueryOpGT         = "gt"
	QueryOpGTE        = "gte"
	QueryOpLT         = "lt"
	QueryOpLTE        = "lte"
	QueryOpIn         = "in"
	QueryOpNotIn      = "not_in"
	QueryOpBetween    = "between"
	QueryOpNotBetween = "not_between"
)

type ColumnCheck struct {
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type QueryCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value,omitempty"`
	Values   []any  `json:"values,omitempty"`
	Lower    any    `json:"lower,omitempty"`
	Upper    any    `json:"upper,omitempty"`
}

type ForeignKey struct {
	Name     string `json:"name,omitempty"`
	Field    string `json:"field"`
	RefTable string `json:"refTable"`
	RefField string `json:"refField"`
	Alias    string `json:"alias,omitempty"`
}

type Column struct {
	Name          string        `json:"name"`
	Type          string        `json:"type,omitempty"`
	Default       any           `json:"default,omitempty"`
	DefaultExpr   string        `json:"defaultExpr,omitempty"`
	OnUpdateExpr  string        `json:"onUpdateExpr,omitempty"`
	MinLength     int           `json:"minLength,omitempty"`
	MaxLength     int           `json:"maxLength,omitempty"`
	Enum          []any         `json:"enum,omitempty"`
	Checks        []ColumnCheck `json:"checks,omitempty"`
	Nullable      *bool         `json:"nullable,omitempty"`
	Required      bool          `json:"required,omitempty"`
	PrimaryKey    bool          `json:"primaryKey,omitempty"`
	AutoIncrement bool          `json:"autoIncrement,omitempty"`
	Unique        bool          `json:"unique,omitempty"`
	Indexed       bool          `json:"indexed,omitempty"`
}

type TableSchema struct {
	Columns       []Column     `json:"columns"`
	ForeignKeys   []ForeignKey `json:"foreignKeys,omitempty"`
	PrimaryKey    string       `json:"primaryKey"`
	AutoIncrement bool         `json:"autoIncrement,omitempty"`
}

func (db *SimpleDB) Configure(schema TableSchema) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	normalized, err := normalizeSchema(schema)
	if err != nil {
		return err
	}

	if db.schema != nil {
		if schemasEqual(*db.schema, normalized) {
			return nil
		}
		return ErrSchemaAlreadyExists
	}

	payload, err := json.Marshal(normalized)
	if err != nil {
		return err
	}
	if err = db.putRawLocked(metaSchemaKey, payload); err != nil {
		return err
	}

	if normalized.AutoIncrement && autoIncrementUsesSequence(normalized) {
		if err = db.persistSequenceLocked(0); err != nil {
			return err
		}
	}

	db.schema = &normalized
	db.autoSeq = 0
	db.resetSecondaryIndexesLocked()
	return nil
}

func (db *SimpleDB) GetSchema() (*TableSchema, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}
	if db.schema == nil {
		return nil, ErrSchemaNotConfigured
	}
	cloned := cloneSchema(*db.schema)
	return &cloned, nil
}

func (db *SimpleDB) InsertRow(values Row) (Row, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}
	if db.schema == nil {
		return nil, ErrSchemaNotConfigured
	}
	return db.insertRowLocked(values)
}

func (db *SimpleDB) InsertRows(values []Row) ([]Row, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}
	if db.schema == nil {
		return nil, ErrSchemaNotConfigured
	}
	if len(values) == 0 {
		return nil, ErrBatchEmpty
	}

	results := make([]Row, 0, len(values))
	for _, value := range values {
		row, err := db.insertRowLocked(value)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (db *SimpleDB) insertRowLocked(values Row) (Row, error) {
	row, err := db.prepareInsertRowLocked(values)
	if err != nil {
		return nil, err
	}

	pkValue := row[db.schema.PrimaryKey]
	pkToken, err := valueToken(pkValue)
	if err != nil {
		return nil, err
	}
	rowKey := buildRowKey(pkToken)
	if current, exists := db.index[rowKey]; exists && !current.Deleted {
		return nil, ErrPrimaryKeyConflict
	}

	if err = db.checkUniqueConstraintsLocked(row, pkToken); err != nil {
		return nil, err
	}

	encodedRow, err := encodeRow(row)
	if err != nil {
		return nil, err
	}
	if err = db.putRawLocked(rowKey, encodedRow); err != nil {
		return nil, err
	}

	db.addRowToIndexesLocked(row, pkToken)
	return cloneRow(row), nil
}

func (db *SimpleDB) UpdateRow(primaryKey any, updates Row) (Row, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}
	if db.schema == nil {
		return nil, ErrSchemaNotConfigured
	}
	return db.updateRowLocked(primaryKey, updates)
}

func (db *SimpleDB) UpdateRows(updates []RowUpdate) ([]Row, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}
	if db.schema == nil {
		return nil, ErrSchemaNotConfigured
	}
	if len(updates) == 0 {
		return nil, ErrBatchEmpty
	}

	results := make([]Row, 0, len(updates))
	for _, update := range updates {
		row, err := db.updateRowLocked(update.PrimaryKey, update.Updates)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, nil
}

func (db *SimpleDB) updateRowLocked(primaryKey any, updates Row) (Row, error) {
	currentRow, pkToken, err := db.findRowLocked(primaryKey)
	if err != nil {
		return nil, err
	}

	updatedRow := cloneRow(currentRow)
	for field, value := range updates {
		if !db.hasColumn(field) {
			return nil, fmt.Errorf("%w: %s", ErrFieldNotDefined, field)
		}
		if field == db.schema.PrimaryKey {
			if currentToken, tokenErr := valueToken(currentRow[field]); tokenErr != nil {
				return nil, tokenErr
			} else if nextToken, tokenErr := valueToken(value); tokenErr != nil {
				return nil, tokenErr
			} else if currentToken != nextToken {
				return nil, ErrPrimaryKeyImmutable
			}
		}
		updatedRow[field] = value
	}
	db.applyOnUpdateExpressionsLocked(updatedRow)
	updatedRow, err = db.normalizeRowValuesLocked(updatedRow)
	if err != nil {
		return nil, err
	}

	if err = db.checkUniqueConstraintsLocked(updatedRow, pkToken); err != nil {
		return nil, err
	}

	encodedRow, err := encodeRow(updatedRow)
	if err != nil {
		return nil, err
	}
	if err = db.putRawLocked(buildRowKey(pkToken), encodedRow); err != nil {
		return nil, err
	}

	db.removeRowFromIndexesLocked(currentRow, pkToken)
	db.addRowToIndexesLocked(updatedRow, pkToken)
	return cloneRow(updatedRow), nil
}

func (db *SimpleDB) DeleteRow(primaryKey any) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}
	if db.schema == nil {
		return ErrSchemaNotConfigured
	}
	return db.deleteRowLocked(primaryKey)
}

func (db *SimpleDB) DeleteRows(primaryKeys []any) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}
	if db.schema == nil {
		return ErrSchemaNotConfigured
	}
	if len(primaryKeys) == 0 {
		return ErrBatchEmpty
	}

	for _, primaryKey := range primaryKeys {
		if err := db.deleteRowLocked(primaryKey); err != nil {
			return err
		}
	}
	return nil
}

func (db *SimpleDB) deleteRowLocked(primaryKey any) error {
	currentRow, pkToken, err := db.findRowLocked(primaryKey)
	if err != nil {
		return err
	}

	if err = db.deleteRawLocked(buildRowKey(pkToken)); err != nil {
		return err
	}
	db.removeRowFromIndexesLocked(currentRow, pkToken)
	return nil
}

func (db *SimpleDB) FindRow(primaryKey any) (Row, bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, false, err
	}
	if db.schema == nil {
		return nil, false, ErrSchemaNotConfigured
	}

	row, _, err := db.findRowLocked(primaryKey)
	if err == ErrKeyNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return cloneRow(row), true, nil
}

func (db *SimpleDB) FindByUnique(field string, value any) (Row, bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, false, err
	}
	if db.schema == nil {
		return nil, false, ErrSchemaNotConfigured
	}
	return db.findByUniqueLocked(field, value)
}

func (db *SimpleDB) FindByIndex(field string, value any) ([]Row, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}
	if db.schema == nil {
		return nil, ErrSchemaNotConfigured
	}
	if !db.isIndexedField(field) {
		return nil, fmt.Errorf("%w: %s", ErrFieldNotIndexed, field)
	}
	column, _ := db.columnByName(field)
	normalizedValue, err := normalizeColumnValue(column, value)
	if err != nil {
		return nil, err
	}

	token, err := valueToken(normalizedValue)
	if err != nil {
		return nil, err
	}

	if db.isUniqueField(field) {
		row, found, err := db.findByUniqueLocked(field, normalizedValue)
		if err != nil || !found {
			return nil, err
		}
		return []Row{row}, nil
	}

	pkSet := db.indexIdx[field][token]
	if len(pkSet) == 0 {
		return []Row{}, nil
	}

	pkTokens := make([]string, 0, len(pkSet))
	for pkToken := range pkSet {
		pkTokens = append(pkTokens, pkToken)
	}
	sort.Strings(pkTokens)

	rows := make([]Row, 0, len(pkTokens))
	for _, pkToken := range pkTokens {
		row, err := db.getRowByTokenLocked(pkToken)
		if err == ErrKeyNotFound {
			continue
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, cloneRow(row))
	}

	return rows, nil
}

func (db *SimpleDB) FindByConditions(conditions []QueryCondition) ([]Row, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, err
	}
	if db.schema == nil {
		return nil, ErrSchemaNotConfigured
	}

	normalizedConditions, err := db.normalizeQueryConditionsLocked(conditions)
	if err != nil {
		return nil, err
	}

	candidatePKs, err := db.planConditionCandidatePKsLocked(normalizedConditions)
	if err != nil {
		return nil, err
	}

	pkTokens := db.collectRowPKTokensLocked(candidatePKs)

	rows := make([]Row, 0, len(pkTokens))
	for _, pkToken := range pkTokens {
		row, err := db.getRowByTokenLocked(pkToken)
		if err == ErrKeyNotFound {
			continue
		}
		if err != nil {
			return nil, err
		}
		matched, err := db.rowMatchesConditionsLocked(row, normalizedConditions)
		if err != nil {
			return nil, err
		}
		if matched {
			rows = append(rows, cloneRow(row))
		}
	}

	return rows, nil
}

// Find is the unified query entry. It accepts arbitrary conditions and
// automatically chooses index-based candidates when possible, then falls back
// to row scan for non-indexed predicates.
func (db *SimpleDB) Find(conditions ...QueryCondition) ([]Row, error) {
	return db.FindByConditions(conditions)
}

// FindOne returns the first matched row for a given condition set.
// The second returned value indicates whether a row is found.
func (db *SimpleDB) FindOne(conditions ...QueryCondition) (Row, bool, error) {
	rows, err := db.FindByConditions(conditions)
	if err != nil {
		return nil, false, err
	}
	if len(rows) == 0 {
		return nil, false, nil
	}
	return rows[0], true, nil
}

func (db *SimpleDB) planConditionCandidatePKsLocked(conditions []QueryCondition) (map[string]struct{}, error) {
	var candidatePKs map[string]struct{}
	for _, condition := range conditions {
		conditionPKs, usable, err := db.indexCandidatePKsForConditionLocked(condition)
		if err != nil {
			return nil, err
		}
		if !usable {
			continue
		}
		if candidatePKs == nil {
			candidatePKs = conditionPKs
			continue
		}
		candidatePKs = intersectPKSets(candidatePKs, conditionPKs)
		if len(candidatePKs) == 0 {
			return candidatePKs, nil
		}
	}
	return candidatePKs, nil
}

func (db *SimpleDB) indexCandidatePKsForConditionLocked(condition QueryCondition) (map[string]struct{}, bool, error) {
	if !db.isIndexedField(condition.Field) {
		return nil, false, nil
	}

	column, ok := db.columnByName(condition.Field)
	if !ok {
		return nil, false, fmt.Errorf("%w: %s", ErrFieldNotDefined, condition.Field)
	}

	switch condition.Operator {
	case QueryOpEQ:
		return db.indexCandidatePKsForEqualLocked(condition.Field, condition.Value)
	case QueryOpNE:
		matched, _, err := db.indexCandidatePKsForEqualLocked(condition.Field, condition.Value)
		if err != nil {
			return nil, false, err
		}
		return subtractPKSets(db.allRowPKSetLocked(), matched), true, nil
	case QueryOpIn:
		merged := make(map[string]struct{})
		for _, value := range condition.Values {
			pkSet, _, err := db.indexCandidatePKsForEqualLocked(condition.Field, value)
			if err != nil {
				return nil, false, err
			}
			mergePKSets(merged, pkSet)
		}
		return merged, true, nil
	case QueryOpNotIn:
		matched := make(map[string]struct{})
		for _, value := range condition.Values {
			pkSet, _, err := db.indexCandidatePKsForEqualLocked(condition.Field, value)
			if err != nil {
				return nil, false, err
			}
			mergePKSets(matched, pkSet)
		}
		return subtractPKSets(db.allRowPKSetLocked(), matched), true, nil
	case QueryOpGT, QueryOpGTE, QueryOpLT, QueryOpLTE, QueryOpBetween:
		return db.indexCandidatePKsForRangeLocked(column, condition)
	case QueryOpNotBetween:
		matched, _, err := db.indexCandidatePKsForRangeLocked(column, QueryCondition{
			Field:    condition.Field,
			Operator: QueryOpBetween,
			Lower:    condition.Lower,
			Upper:    condition.Upper,
		})
		if err != nil {
			return nil, false, err
		}
		return subtractPKSets(db.allRowPKSetLocked(), matched), true, nil
	default:
		return nil, false, nil
	}
}

func (db *SimpleDB) indexCandidatePKsForEqualLocked(field string, value any) (map[string]struct{}, bool, error) {
	if value == nil {
		return map[string]struct{}{}, true, nil
	}
	token, err := valueToken(value)
	if err != nil {
		return nil, false, err
	}
	if db.isUniqueField(field) {
		pkToken, exists := db.uniqueIdx[field][token]
		if !exists {
			return map[string]struct{}{}, true, nil
		}
		return map[string]struct{}{pkToken: {}}, true, nil
	}
	return clonePKSet(db.indexIdx[field][token]), true, nil
}

func (db *SimpleDB) indexCandidatePKsForRangeLocked(column Column, condition QueryCondition) (map[string]struct{}, bool, error) {
	matched := make(map[string]struct{})
	if db.isUniqueField(condition.Field) {
		for token, pkToken := range db.uniqueIdx[condition.Field] {
			bucketValue, err := db.indexBucketValueLocked(column, token, pkToken)
			if err != nil {
				return nil, false, err
			}
			ok, err := evaluateQueryCondition(column, bucketValue, condition)
			if err != nil {
				return nil, false, err
			}
			if ok {
				matched[pkToken] = struct{}{}
			}
		}
		return matched, true, nil
	}

	for token, pkSet := range db.indexIdx[condition.Field] {
		pkToken, exists := firstPKToken(pkSet)
		if !exists {
			continue
		}
		bucketValue, err := db.indexBucketValueLocked(column, token, pkToken)
		if err != nil {
			return nil, false, err
		}
		ok, err := evaluateQueryCondition(column, bucketValue, condition)
		if err != nil {
			return nil, false, err
		}
		if ok {
			mergePKSets(matched, pkSet)
		}
	}
	return matched, true, nil
}

func (db *SimpleDB) indexBucketValueLocked(column Column, token string, pkToken string) (any, error) {
	row, err := db.getRowByTokenLocked(pkToken)
	if err != nil {
		return nil, err
	}
	value, exists := row[column.Name]
	if !exists {
		decoded, decodeErr := decodeValueToken(token)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return normalizeColumnValue(column, decoded)
	}
	return normalizeColumnValue(column, value)
}

func (db *SimpleDB) collectRowPKTokensLocked(candidatePKs map[string]struct{}) []string {
	if candidatePKs != nil {
		pkTokens := make([]string, 0, len(candidatePKs))
		for pkToken := range candidatePKs {
			pkTokens = append(pkTokens, pkToken)
		}
		sort.Strings(pkTokens)
		return pkTokens
	}

	pkTokens := make([]string, 0)
	for key, current := range db.index {
		if current.Deleted || !strings.HasPrefix(key, rowKeyPrefix) {
			continue
		}
		pkTokens = append(pkTokens, strings.TrimPrefix(key, rowKeyPrefix))
	}
	sort.Strings(pkTokens)
	return pkTokens
}

func (db *SimpleDB) allRowPKSetLocked() map[string]struct{} {
	pkSet := make(map[string]struct{})
	for key, current := range db.index {
		if current.Deleted || !strings.HasPrefix(key, rowKeyPrefix) {
			continue
		}
		pkSet[strings.TrimPrefix(key, rowKeyPrefix)] = struct{}{}
	}
	return pkSet
}

func (db *SimpleDB) rebuildStructuredState() error {
	db.schema = nil
	db.autoSeq = 0
	db.uniqueIdx = make(map[string]map[string]string)
	db.indexIdx = make(map[string]map[string]map[string]struct{})

	if schemaEntry, exists := db.index[metaSchemaKey]; exists && !schemaEntry.Deleted {
		var schema TableSchema
		if err := json.Unmarshal(schemaEntry.Value, &schema); err != nil {
			return fmt.Errorf("%w: schema", ErrCorruptedRecord)
		}
		normalized, err := normalizeSchema(schema)
		if err != nil {
			return err
		}
		db.schema = &normalized
		db.resetSecondaryIndexesLocked()
	}

	if seqEntry, exists := db.index[metaSequenceKey]; exists && !seqEntry.Deleted {
		if err := json.Unmarshal(seqEntry.Value, &db.autoSeq); err != nil {
			return fmt.Errorf("%w: sequence", ErrCorruptedRecord)
		}
	}

	if db.schema == nil {
		return nil
	}

	for key, current := range db.index {
		if current.Deleted || !strings.HasPrefix(key, rowKeyPrefix) {
			continue
		}
		row, err := decodeRow(current.Value)
		if err != nil {
			return err
		}
		if row, err = db.normalizeRowValuesLocked(row); err != nil {
			return err
		}
		pkValue, exists := row[db.schema.PrimaryKey]
		if !exists {
			return ErrPrimaryKeyMissing
		}
		pkToken, err := valueToken(pkValue)
		if err != nil {
			return err
		}
		db.addRowToIndexesLocked(row, pkToken)
	}

	return nil
}

func (db *SimpleDB) prepareInsertRowLocked(values Row) (Row, error) {
	row := cloneRow(values)
	supplied := make(map[string]struct{}, len(values))
	for field := range row {
		if !db.hasColumn(field) {
			return nil, fmt.Errorf("%w: %s", ErrFieldNotDefined, field)
		}
		supplied[field] = struct{}{}
	}

	pkField := db.schema.PrimaryKey
	pkColumn, ok := db.columnByName(pkField)
	if !ok {
		return nil, ErrPrimaryKeyMissing
	}
	pkValue, exists := row[pkField]
	if !exists || pkValue == nil {
		if !db.schema.AutoIncrement {
			return nil, ErrPrimaryKeyMissing
		}
		switch ColumnType(normalizeColumnType(pkColumn.Type)) {
		case ColumnTypeUUID:
			version := extractUUIDVersion(pkColumn.Type)
			if version == 0 {
				version = db.getDefaultUUIDVersion()
			}
			generatedUUID, err := generateUUIDByVersion(version)
			if err != nil {
				return nil, fmt.Errorf("failed to generate UUID v%d: %w", version, err)
			}
			row[pkField] = generatedUUID.String()
		default:
			db.autoSeq++
			row[pkField] = db.autoSeq
			if err := db.persistSequenceLocked(db.autoSeq); err != nil {
				return nil, err
			}
		}
	} else if db.schema.AutoIncrement {
		switch ColumnType(normalizeColumnType(pkColumn.Type)) {
		case ColumnTypeUUID:
			uuidValue, err := normalizeUUIDValue(pkColumn, pkValue)
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
			if intValue > db.autoSeq {
				db.autoSeq = intValue
				if err := db.persistSequenceLocked(db.autoSeq); err != nil {
					return nil, err
				}
			}
		}
	}

	row, err := db.applyColumnConstraintsLocked(row, supplied, true, true)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (db *SimpleDB) persistSequenceLocked(sequence int64) error {
	payload, err := json.Marshal(sequence)
	if err != nil {
		return err
	}
	return db.putRawLocked(metaSequenceKey, payload)
}

func (db *SimpleDB) findRowLocked(primaryKey any) (Row, string, error) {
	column, ok := db.columnByName(db.schema.PrimaryKey)
	if !ok {
		return nil, "", ErrPrimaryKeyMissing
	}
	normalizedPrimaryKey, err := normalizeColumnValue(column, primaryKey)
	if err != nil {
		return nil, "", err
	}
	pkToken, err := valueToken(normalizedPrimaryKey)
	if err != nil {
		return nil, "", err
	}
	row, err := db.getRowByTokenLocked(pkToken)
	if err != nil {
		return nil, "", err
	}
	return row, pkToken, nil
}

func (db *SimpleDB) getRowByTokenLocked(pkToken string) (Row, error) {
	current, exists := db.index[buildRowKey(pkToken)]
	if !exists || current.Deleted {
		return nil, ErrKeyNotFound
	}
	row, err := decodeRow(current.Value)
	if err != nil {
		return nil, err
	}
	if db.schema == nil {
		return row, nil
	}
	return db.normalizeRowValuesLocked(row)
}

func (db *SimpleDB) findByUniqueLocked(field string, value any) (Row, bool, error) {
	column, ok := db.columnByName(field)
	if !ok || !column.Unique {
		return nil, false, fmt.Errorf("%w: %s", ErrFieldNotIndexed, field)
	}

	normalizedValue, err := normalizeColumnValue(column, value)
	if err != nil {
		return nil, false, err
	}
	token, err := valueToken(normalizedValue)
	if err != nil {
		return nil, false, err
	}
	pkToken, exists := db.uniqueIdx[field][token]
	if !exists {
		return nil, false, nil
	}

	row, err := db.getRowByTokenLocked(pkToken)
	if err == ErrKeyNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return cloneRow(row), true, nil
}

func (db *SimpleDB) normalizeQueryConditionsLocked(conditions []QueryCondition) ([]QueryCondition, error) {
	normalized := make([]QueryCondition, 0, len(conditions))
	for _, condition := range conditions {
		condition.Field = strings.TrimSpace(condition.Field)
		if condition.Field == "" {
			return nil, fmt.Errorf("%w: field is required", ErrInvalidQueryCondition)
		}
		column, ok := db.columnByName(condition.Field)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrFieldNotDefined, condition.Field)
		}
		normalizedCondition, err := normalizeQueryCondition(column, condition)
		if err != nil {
			return nil, err
		}
		normalized = append(normalized, normalizedCondition)
	}
	return normalized, nil
}

func (db *SimpleDB) rowMatchesConditionsLocked(row Row, conditions []QueryCondition) (bool, error) {
	for _, condition := range conditions {
		column, ok := db.columnByName(condition.Field)
		if !ok {
			return false, fmt.Errorf("%w: %s", ErrFieldNotDefined, condition.Field)
		}
		value, exists := row[condition.Field]
		if !exists {
			value = nil
		}
		matched, err := evaluateQueryCondition(column, value, condition)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

func (db *SimpleDB) checkUniqueConstraintsLocked(row Row, selfPKToken string) error {
	for _, column := range db.schema.Columns {
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
		if otherPK, exists := db.uniqueIdx[column.Name][token]; exists && otherPK != selfPKToken {
			return fmt.Errorf("%w: %s", ErrUniqueConflict, column.Name)
		}
	}
	return nil
}

func (db *SimpleDB) addRowToIndexesLocked(row Row, pkToken string) {
	for _, column := range db.schema.Columns {
		value, exists := row[column.Name]
		if !exists || value == nil {
			continue
		}
		token, err := valueToken(value)
		if err != nil {
			continue
		}
		if column.Unique {
			if db.uniqueIdx[column.Name] == nil {
				db.uniqueIdx[column.Name] = make(map[string]string)
			}
			db.uniqueIdx[column.Name][token] = pkToken
		}
		if column.Indexed {
			if db.indexIdx[column.Name] == nil {
				db.indexIdx[column.Name] = make(map[string]map[string]struct{})
			}
			if db.indexIdx[column.Name][token] == nil {
				db.indexIdx[column.Name][token] = make(map[string]struct{})
			}
			db.indexIdx[column.Name][token][pkToken] = struct{}{}
		}
	}
}

func (db *SimpleDB) removeRowFromIndexesLocked(row Row, pkToken string) {
	for _, column := range db.schema.Columns {
		value, exists := row[column.Name]
		if !exists || value == nil {
			continue
		}
		token, err := valueToken(value)
		if err != nil {
			continue
		}
		if column.Unique {
			delete(db.uniqueIdx[column.Name], token)
		}
		if column.Indexed {
			pkSet := db.indexIdx[column.Name][token]
			delete(pkSet, pkToken)
			if len(pkSet) == 0 {
				delete(db.indexIdx[column.Name], token)
			}
		}
	}
}

func (db *SimpleDB) resetSecondaryIndexesLocked() {
	db.uniqueIdx = make(map[string]map[string]string)
	db.indexIdx = make(map[string]map[string]map[string]struct{})
	if db.schema == nil {
		return
	}
	for _, column := range db.schema.Columns {
		if column.Unique {
			db.uniqueIdx[column.Name] = make(map[string]string)
		}
		if column.Indexed {
			db.indexIdx[column.Name] = make(map[string]map[string]struct{})
		}
	}
}

func (db *SimpleDB) hasColumn(field string) bool {
	_, ok := db.columnByName(field)
	return ok
}

func (db *SimpleDB) isUniqueField(field string) bool {
	column, ok := db.columnByName(field)
	return ok && column.Unique
}

func (db *SimpleDB) isIndexedField(field string) bool {
	column, ok := db.columnByName(field)
	return ok && (column.Indexed || column.Unique)
}

func (db *SimpleDB) columnByName(field string) (Column, bool) {
	if db.schema == nil {
		return Column{}, false
	}
	for _, column := range db.schema.Columns {
		if column.Name == field {
			return column, true
		}
	}
	return Column{}, false
}

func normalizeSchema(schema TableSchema) (TableSchema, error) {
	if len(schema.Columns) == 0 {
		return TableSchema{}, ErrInvalidSchema
	}

	columnNames := make(map[string]struct{}, len(schema.Columns))
	primaryKeyCount := 0
	normalized := TableSchema{Columns: make([]Column, 0, len(schema.Columns))}
	for _, column := range schema.Columns {
		column.Name = strings.TrimSpace(column.Name)
		column.Type = normalizeColumnType(column.Type)
		column.DefaultExpr = normalizeColumnExpression(column.DefaultExpr)
		column.OnUpdateExpr = normalizeColumnExpression(column.OnUpdateExpr)
		if column.Name == "" {
			return TableSchema{}, fmt.Errorf("%w: empty column name", ErrInvalidSchema)
		}
		if _, exists := columnNames[column.Name]; exists {
			return TableSchema{}, fmt.Errorf("%w: duplicate column %s", ErrInvalidSchema, column.Name)
		}
		columnNames[column.Name] = struct{}{}
		if column.PrimaryKey {
			primaryKeyCount++
			normalized.PrimaryKey = column.Name
			column.Unique = true
			column.Nullable = boolPtr(false)
			if !column.AutoIncrement {
				column.Required = true
			}
		}
		if column.Required {
			column.Nullable = boolPtr(false)
		}
		if column.AutoIncrement {
			if !column.PrimaryKey {
				return TableSchema{}, fmt.Errorf("%w: auto increment field %s must be primary key", ErrInvalidSchema, column.Name)
			}
			if column.Type != string(ColumnTypeAny) && column.Type != string(ColumnTypeInt) && column.Type != string(ColumnTypeUUID) {
				return TableSchema{}, fmt.Errorf("%w: auto increment field %s must be int or uuid", ErrInvalidSchema, column.Name)
			}
			if column.Type == string(ColumnTypeAny) {
				column.Type = string(ColumnTypeInt)
			}
			normalized.AutoIncrement = true
		}
		if !isSupportedColumnType(ColumnType(column.Type)) {
			return TableSchema{}, fmt.Errorf("%w: %s=%s", ErrUnsupportedFieldType, column.Name, column.Type)
		}
		if err := validateColumnConstraintCompatibility(column); err != nil {
			return TableSchema{}, err
		}
		if hasColumnDefault(column) && column.DefaultExpr != "" {
			return TableSchema{}, fmt.Errorf("%w: %s cannot set both default and defaultExpr", ErrInvalidSchema, column.Name)
		}
		if hasColumnDefault(column) {
			normalizedDefault, err := normalizeColumnValue(column, column.Default)
			if err != nil {
				return TableSchema{}, fmt.Errorf("%w: invalid default for %s", ErrInvalidSchema, column.Name)
			}
			if err = validateNormalizedColumnValue(column, normalizedDefault); err != nil {
				return TableSchema{}, fmt.Errorf("%w: invalid default for %s", ErrInvalidSchema, column.Name)
			}
			column.Default = normalizedDefault
		}
		if len(column.Enum) > 0 {
			normalizedEnum, err := normalizeEnumValues(column)
			if err != nil {
				return TableSchema{}, err
			}
			column.Enum = normalizedEnum
		}
		if len(column.Checks) > 0 {
			normalizedChecks, err := normalizeColumnChecks(column)
			if err != nil {
				return TableSchema{}, err
			}
			column.Checks = normalizedChecks
		}
		if column.DefaultExpr != "" {
			if !isSupportedColumnExpression(column.DefaultExpr) {
				return TableSchema{}, fmt.Errorf("%w: invalid defaultExpr for %s", ErrInvalidSchema, column.Name)
			}
			if !expressionMatchesColumnType(column, column.DefaultExpr) {
				return TableSchema{}, fmt.Errorf("%w: defaultExpr type mismatch for %s", ErrInvalidSchema, column.Name)
			}
		}
		if column.OnUpdateExpr != "" {
			if !isSupportedColumnExpression(column.OnUpdateExpr) {
				return TableSchema{}, fmt.Errorf("%w: invalid onUpdateExpr for %s", ErrInvalidSchema, column.Name)
			}
			if !expressionMatchesColumnType(column, column.OnUpdateExpr) {
				return TableSchema{}, fmt.Errorf("%w: onUpdateExpr type mismatch for %s", ErrInvalidSchema, column.Name)
			}
		}
		normalized.Columns = append(normalized.Columns, column)
	}

	if primaryKeyCount != 1 {
		return TableSchema{}, fmt.Errorf("%w: exactly one primary key is required", ErrInvalidSchema)
	}

	normalized.ForeignKeys = make([]ForeignKey, 0, len(schema.ForeignKeys))
	foreignKeyNames := make(map[string]struct{}, len(schema.ForeignKeys))
	for _, foreignKey := range schema.ForeignKeys {
		foreignKey.Name = strings.TrimSpace(foreignKey.Name)
		foreignKey.Field = strings.TrimSpace(foreignKey.Field)
		foreignKey.RefTable = strings.TrimSpace(foreignKey.RefTable)
		foreignKey.RefField = strings.TrimSpace(foreignKey.RefField)
		foreignKey.Alias = strings.TrimSpace(foreignKey.Alias)
		if foreignKey.Field == "" || foreignKey.RefTable == "" || foreignKey.RefField == "" {
			return TableSchema{}, ErrInvalidForeignKey
		}
		if _, ok := columnNames[foreignKey.Field]; !ok {
			return TableSchema{}, fmt.Errorf("%w: field %s not defined", ErrInvalidForeignKey, foreignKey.Field)
		}
		if foreignKey.Name != "" {
			if _, exists := foreignKeyNames[foreignKey.Name]; exists {
				return TableSchema{}, fmt.Errorf("%w: duplicate foreign key %s", ErrInvalidForeignKey, foreignKey.Name)
			}
			foreignKeyNames[foreignKey.Name] = struct{}{}
		}
		for index, column := range normalized.Columns {
			if column.Name == foreignKey.Field {
				column.Indexed = true
				normalized.Columns[index] = column
				break
			}
		}
		normalized.ForeignKeys = append(normalized.ForeignKeys, foreignKey)
	}

	if schema.PrimaryKey != "" && schema.PrimaryKey != normalized.PrimaryKey {
		return TableSchema{}, fmt.Errorf("%w: primary key mismatch", ErrInvalidSchema)
	}

	return normalized, nil
}

func schemasEqual(left, right TableSchema) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	if leftErr != nil || rightErr != nil {
		return false
	}
	return string(leftJSON) == string(rightJSON)
}

func cloneSchema(schema TableSchema) TableSchema {
	cloned := TableSchema{
		PrimaryKey:    schema.PrimaryKey,
		AutoIncrement: schema.AutoIncrement,
		Columns:       make([]Column, len(schema.Columns)),
		ForeignKeys:   make([]ForeignKey, len(schema.ForeignKeys)),
	}
	copy(cloned.Columns, schema.Columns)
	copy(cloned.ForeignKeys, schema.ForeignKeys)
	return cloned
}

func autoIncrementUsesSequence(schema TableSchema) bool {
	if !schema.AutoIncrement {
		return false
	}
	for _, column := range schema.Columns {
		if column.Name == schema.PrimaryKey {
			return ColumnType(normalizeColumnType(column.Type)) != ColumnTypeUUID
		}
	}
	return true
}

func cloneRow(row Row) Row {
	cloned := make(Row, len(row))
	for key, value := range row {
		cloned[key] = value
	}
	return cloned
}

func encodeRow(row Row) ([]byte, error) {
	return json.Marshal(row)
}

func decodeRow(raw []byte) (Row, error) {
	var row Row
	if err := json.Unmarshal(raw, &row); err != nil {
		return nil, fmt.Errorf("%w: row", ErrCorruptedRecord)
	}
	return row, nil
}

func (db *SimpleDB) normalizeRowValuesLocked(row Row) (Row, error) {
	return db.applyColumnConstraintsLocked(row, nil, false, false)
}

func (db *SimpleDB) applyColumnConstraintsLocked(row Row, supplied map[string]struct{}, enforceRequired bool, applyDefaults bool) (Row, error) {
	normalized := cloneRow(row)
	for _, column := range db.schema.Columns {
		_, wasSupplied := supplied[column.Name]
		value, exists := normalized[column.Name]
		if !exists && applyDefaults && column.DefaultExpr != "" {
			evaluated, err := evaluateColumnExpression(column, column.DefaultExpr, time.Now().UTC())
			if err != nil {
				return nil, err
			}
			normalized[column.Name] = evaluated
			value = evaluated
			exists = true
		}
		if !exists && applyDefaults && hasColumnDefault(column) {
			clonedDefault, err := cloneAny(column.Default)
			if err != nil {
				return nil, err
			}
			normalized[column.Name] = clonedDefault
			value = clonedDefault
			exists = true
		}
		if enforceRequired && column.Required && !wasSupplied {
			return nil, fmt.Errorf("%w: %s", ErrFieldRequired, column.Name)
		}
		if !exists {
			if !columnAllowsNull(column) {
				if column.PrimaryKey {
					return nil, ErrPrimaryKeyMissing
				}
				return nil, fmt.Errorf("%w: %s", ErrFieldRequired, column.Name)
			}
			continue
		}
		if value == nil {
			if !columnAllowsNull(column) {
				return nil, fmt.Errorf("%w: %s", ErrFieldNotNullable, column.Name)
			}
			continue
		}
		nextValue, err := normalizeColumnValue(column, value)
		if err != nil {
			return nil, err
		}
		if err = validateNormalizedColumnValue(column, nextValue); err != nil {
			return nil, err
		}
		normalized[column.Name] = nextValue
	}
	return normalized, nil
}

func normalizeColumnType(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return string(ColumnTypeAny)
	}
	if strings.HasPrefix(raw, "uuid:") || strings.HasPrefix(raw, "uuid/") {
		return string(ColumnTypeUUID)
	}
	switch raw {
	case "integer":
		return string(ColumnTypeInt)
	case "number", "double":
		return string(ColumnTypeFloat)
	case "boolean":
		return string(ColumnTypeBool)
	case "datetime":
		return string(ColumnTypeTimestamp)
	case "guid":
		return string(ColumnTypeUUID)
	default:
		return raw
	}
}

func extractUUIDVersion(columnType string) int {
	columnType = strings.TrimSpace(strings.ToLower(columnType))
	if strings.HasPrefix(columnType, "uuid:") {
		versionStr := strings.TrimPrefix(columnType, "uuid:")
		versionStr = strings.TrimPrefix(versionStr, "v")
		if version := parseInt(versionStr); version >= 1 && version <= 8 {
			return int(version)
		}
	}
	if strings.HasPrefix(columnType, "uuid/") {
		versionStr := strings.TrimPrefix(columnType, "uuid/")
		versionStr = strings.TrimPrefix(versionStr, "v")
		if version := parseInt(versionStr); version >= 1 && version <= 8 {
			return int(version)
		}
	}
	return 0
}

func generateUUIDByVersion(version int) (uuid.UUID, error) {
	switch version {
	case 1:
		return uuid.NewUUID()
	case 4:
		return uuid.NewRandom()
	case 6:
		return uuid.NewV6()
	case 7:
		return uuid.NewV7()
	default:
		return uuid.Nil, fmt.Errorf("unsupported UUID version: %d (supported: 1, 4, 6, 7)", version)
	}
}

func parseInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var result int64
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0
		}
		result = result*10 + int64(ch-'0')
	}
	return result
}

func isSupportedColumnType(kind ColumnType) bool {
	switch kind {
	case ColumnTypeAny, ColumnTypeString, ColumnTypeInt, ColumnTypeFloat, ColumnTypeBool, ColumnTypeObject, ColumnTypeArray, ColumnTypeUUID, ColumnTypeTime, ColumnTypeTimestamp:
		return true
	default:
		return false
	}
}

func normalizeColumnValue(column Column, value any) (any, error) {
	columnType := ColumnType(normalizeColumnType(column.Type))
	switch columnType {
	case ColumnTypeAny:
		return value, nil
	case ColumnTypeString:
		stringValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s requires string", ErrFieldTypeMismatch, column.Name)
		}
		return stringValue, nil
	case ColumnTypeInt:
		intValue, ok := asInt64(value)
		if !ok {
			return nil, fmt.Errorf("%w: %s requires int", ErrFieldTypeMismatch, column.Name)
		}
		return intValue, nil
	case ColumnTypeFloat:
		floatValue, ok := asFloat64(value)
		if !ok {
			return nil, fmt.Errorf("%w: %s requires float", ErrFieldTypeMismatch, column.Name)
		}
		return floatValue, nil
	case ColumnTypeBool:
		boolValue, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("%w: %s requires bool", ErrFieldTypeMismatch, column.Name)
		}
		return boolValue, nil
	case ColumnTypeObject:
		objectValue, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%w: %s requires object", ErrFieldTypeMismatch, column.Name)
		}
		return objectValue, nil
	case ColumnTypeArray:
		arrayValue, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("%w: %s requires array", ErrFieldTypeMismatch, column.Name)
		}
		return arrayValue, nil
	case ColumnTypeUUID:
		return normalizeUUIDValue(column, value)
	case ColumnTypeTime:
		return normalizeTimeValue(column, value)
	case ColumnTypeTimestamp:
		return normalizeTimestampValue(column, value)
	default:
		return nil, fmt.Errorf("%w: %s=%s", ErrUnsupportedFieldType, column.Name, column.Type)
	}
}

func normalizeUUIDValue(column Column, value any) (any, error) {
	switch typed := value.(type) {
	case uuid.UUID:
		return typed.String(), nil
	case string:
		parsed, err := uuid.Parse(strings.TrimSpace(typed))
		if err != nil {
			return nil, fmt.Errorf("%w: %s requires uuid", ErrFieldTypeMismatch, column.Name)
		}
		return parsed.String(), nil
	default:
		return nil, fmt.Errorf("%w: %s requires uuid", ErrFieldTypeMismatch, column.Name)
	}
}

func normalizeTimeValue(column Column, value any) (any, error) {
	switch typed := value.(type) {
	case time.Time:
		return typed.UTC().Format(timeLayout), nil
	case string:
		parsed, err := time.Parse(timeLayout, typed)
		if err != nil {
			return nil, fmt.Errorf("%w: %s requires time", ErrFieldTypeMismatch, column.Name)
		}
		return parsed.UTC().Format(timeLayout), nil
	default:
		return nil, fmt.Errorf("%w: %s requires time", ErrFieldTypeMismatch, column.Name)
	}
}

func normalizeTimestampValue(column Column, value any) (any, error) {
	switch typed := value.(type) {
	case time.Time:
		return typed.UTC().Format(time.RFC3339Nano), nil
	case string:
		parsed, err := time.Parse(time.RFC3339Nano, typed)
		if err != nil {
			return nil, fmt.Errorf("%w: %s requires timestamp", ErrFieldTypeMismatch, column.Name)
		}
		return parsed.UTC().Format(time.RFC3339Nano), nil
	default:
		return nil, fmt.Errorf("%w: %s requires timestamp", ErrFieldTypeMismatch, column.Name)
	}
}

func normalizeColumnExpression(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func normalizeCheckOperator(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func normalizeQueryOperator(raw string) string {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	normalized = strings.ReplaceAll(normalized, "-", "_")
	return normalized
}

func isSupportedColumnExpression(expr string) bool {
	switch normalizeColumnExpression(expr) {
	case "", ColumnExprCurrentTime, ColumnExprCurrentTimestamp:
		return true
	default:
		return false
	}
}

func expressionMatchesColumnType(column Column, expr string) bool {
	switch normalizeColumnExpression(expr) {
	case ColumnExprCurrentTime:
		return ColumnType(normalizeColumnType(column.Type)) == ColumnTypeTime
	case ColumnExprCurrentTimestamp:
		return ColumnType(normalizeColumnType(column.Type)) == ColumnTypeTimestamp
	default:
		return false
	}
}

func evaluateColumnExpression(column Column, expr string, now time.Time) (any, error) {
	switch normalizeColumnExpression(expr) {
	case ColumnExprCurrentTime:
		return normalizeTimeValue(column, now.UTC())
	case ColumnExprCurrentTimestamp:
		return normalizeTimestampValue(column, now.UTC())
	default:
		return nil, fmt.Errorf("%w: unsupported expression for %s", ErrInvalidSchema, column.Name)
	}
}

func (db *SimpleDB) applyOnUpdateExpressionsLocked(row Row) {
	now := time.Now().UTC()
	for _, column := range db.schema.Columns {
		if column.OnUpdateExpr == "" {
			continue
		}
		if value, err := evaluateColumnExpression(column, column.OnUpdateExpr, now); err == nil {
			row[column.Name] = value
		}
	}
}

func columnAllowsNull(column Column) bool {
	if column.PrimaryKey {
		return false
	}
	if column.Nullable == nil {
		return true
	}
	return *column.Nullable
}

func hasColumnDefault(column Column) bool {
	return column.Default != nil
}

func validateColumnConstraintCompatibility(column Column) error {
	columnType := ColumnType(normalizeColumnType(column.Type))
	if column.MinLength < 0 || column.MaxLength < 0 {
		return fmt.Errorf("%w: %s length must be >= 0", ErrInvalidSchema, column.Name)
	}
	if column.MaxLength > 0 && column.MinLength > column.MaxLength {
		return fmt.Errorf("%w: %s minLength > maxLength", ErrInvalidSchema, column.Name)
	}
	if (column.MinLength > 0 || column.MaxLength > 0) && !supportsLengthConstraint(columnType) {
		return fmt.Errorf("%w: %s length only supports string/array", ErrInvalidSchema, column.Name)
	}
	return nil
}

func supportsRangeComparison(columnType ColumnType) bool {
	switch columnType {
	case ColumnTypeInt, ColumnTypeFloat, ColumnTypeString, ColumnTypeTime, ColumnTypeTimestamp:
		return true
	default:
		return false
	}
}

func normalizeQueryCondition(column Column, condition QueryCondition) (QueryCondition, error) {
	condition.Operator = normalizeQueryOperator(condition.Operator)
	if condition.Operator == "" {
		condition.Operator = QueryOpEQ
	}

	columnType := ColumnType(normalizeColumnType(column.Type))
	switch condition.Operator {
	case QueryOpEQ, QueryOpNE:
		if condition.Value != nil {
			normalizedValue, err := normalizeColumnValue(column, condition.Value)
			if err != nil {
				return QueryCondition{}, err
			}
			condition.Value = normalizedValue
		}
		return condition, nil
	case QueryOpGT, QueryOpGTE, QueryOpLT, QueryOpLTE:
		if !supportsRangeComparison(columnType) {
			return QueryCondition{}, fmt.Errorf("%w: %s does not support %s", ErrInvalidQueryCondition, column.Name, condition.Operator)
		}
		if condition.Value == nil {
			return QueryCondition{}, fmt.Errorf("%w: %s requires value", ErrInvalidQueryCondition, condition.Operator)
		}
		normalizedValue, err := normalizeColumnValue(column, condition.Value)
		if err != nil {
			return QueryCondition{}, err
		}
		condition.Value = normalizedValue
		return condition, nil
	case QueryOpIn, QueryOpNotIn:
		if len(condition.Values) == 0 {
			return QueryCondition{}, fmt.Errorf("%w: %s requires values", ErrInvalidQueryCondition, condition.Operator)
		}
		normalizedValues := make([]any, 0, len(condition.Values))
		for _, item := range condition.Values {
			if item == nil {
				normalizedValues = append(normalizedValues, nil)
				continue
			}
			normalizedValue, err := normalizeColumnValue(column, item)
			if err != nil {
				return QueryCondition{}, err
			}
			normalizedValues = append(normalizedValues, normalizedValue)
		}
		condition.Values = normalizedValues
		return condition, nil
	case QueryOpBetween, QueryOpNotBetween:
		if !supportsRangeComparison(columnType) {
			return QueryCondition{}, fmt.Errorf("%w: %s does not support %s", ErrInvalidQueryCondition, column.Name, condition.Operator)
		}
		if condition.Lower == nil || condition.Upper == nil {
			return QueryCondition{}, fmt.Errorf("%w: %s requires lower and upper", ErrInvalidQueryCondition, condition.Operator)
		}
		lower, err := normalizeColumnValue(column, condition.Lower)
		if err != nil {
			return QueryCondition{}, err
		}
		upper, err := normalizeColumnValue(column, condition.Upper)
		if err != nil {
			return QueryCondition{}, err
		}
		cmp, err := compareScalarValues(lower, upper)
		if err != nil {
			return QueryCondition{}, fmt.Errorf("%w: %s lower/upper incomparable", ErrInvalidQueryCondition, column.Name)
		}
		if cmp > 0 {
			return QueryCondition{}, fmt.Errorf("%w: %s lower > upper", ErrInvalidQueryCondition, column.Name)
		}
		condition.Lower = lower
		condition.Upper = upper
		return condition, nil
	default:
		return QueryCondition{}, fmt.Errorf("%w: unsupported operator %s", ErrInvalidQueryCondition, condition.Operator)
	}
}

func supportsLengthConstraint(columnType ColumnType) bool {
	switch columnType {
	case ColumnTypeString, ColumnTypeArray:
		return true
	default:
		return false
	}
}

func normalizeEnumValues(column Column) ([]any, error) {
	normalized := make([]any, 0, len(column.Enum))
	seen := make(map[string]struct{}, len(column.Enum))
	for _, item := range column.Enum {
		nextValue, err := normalizeColumnValue(column, item)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid enum for %s", ErrInvalidSchema, column.Name)
		}
		if err = validateNormalizedColumnValue(Column{Type: column.Type, Name: column.Name, MinLength: column.MinLength, MaxLength: column.MaxLength}, nextValue); err != nil {
			return nil, fmt.Errorf("%w: invalid enum for %s", ErrInvalidSchema, column.Name)
		}
		token, err := valueToken(nextValue)
		if err != nil {
			return nil, err
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		normalized = append(normalized, nextValue)
	}
	return normalized, nil
}

func normalizeColumnChecks(column Column) ([]ColumnCheck, error) {
	normalized := make([]ColumnCheck, 0, len(column.Checks))
	for _, check := range column.Checks {
		check.Operator = normalizeCheckOperator(check.Operator)
		if check.Operator == "" {
			return nil, fmt.Errorf("%w: empty check operator for %s", ErrInvalidSchema, column.Name)
		}
		nextValue, err := normalizeCheckValue(column, check)
		if err != nil {
			return nil, err
		}
		check.Value = nextValue
		normalized = append(normalized, check)
	}
	return normalized, nil
}

func normalizeCheckValue(column Column, check ColumnCheck) (any, error) {
	columnType := ColumnType(normalizeColumnType(column.Type))
	switch check.Operator {
	case ColumnCheckGT, ColumnCheckGTE, ColumnCheckLT, ColumnCheckLTE:
		switch columnType {
		case ColumnTypeInt, ColumnTypeFloat, ColumnTypeString, ColumnTypeTime, ColumnTypeTimestamp:
			return normalizeColumnValue(column, check.Value)
		default:
			return nil, fmt.Errorf("%w: unsupported check %s for %s", ErrInvalidSchema, check.Operator, column.Name)
		}
	case ColumnCheckLenGT, ColumnCheckLenGTE, ColumnCheckLenLT, ColumnCheckLenLTE:
		if !supportsLengthConstraint(columnType) {
			return nil, fmt.Errorf("%w: unsupported length check for %s", ErrInvalidSchema, column.Name)
		}
		length, ok := asInt64(check.Value)
		if !ok || length < 0 {
			return nil, fmt.Errorf("%w: invalid length check for %s", ErrInvalidSchema, column.Name)
		}
		return length, nil
	case ColumnCheckRegex:
		if columnType != ColumnTypeString {
			return nil, fmt.Errorf("%w: regex only supports string for %s", ErrInvalidSchema, column.Name)
		}
		pattern, ok := check.Value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: regex must be string for %s", ErrInvalidSchema, column.Name)
		}
		if _, err := regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("%w: invalid regex for %s", ErrInvalidSchema, column.Name)
		}
		return pattern, nil
	default:
		return nil, fmt.Errorf("%w: unsupported check operator %s", ErrInvalidSchema, check.Operator)
	}
}

func validateNormalizedColumnValue(column Column, value any) error {
	if value == nil {
		return nil
	}
	if err := validateLengthConstraint(column, value); err != nil {
		return err
	}
	if err := validateEnumConstraint(column, value); err != nil {
		return err
	}
	if err := validateCheckConstraints(column, value); err != nil {
		return err
	}
	return nil
}

func validateLengthConstraint(column Column, value any) error {
	if column.MinLength == 0 && column.MaxLength == 0 {
		return nil
	}
	length, ok := valueLength(value)
	if !ok {
		return fmt.Errorf("%w: %s length unsupported", ErrFieldLengthViolation, column.Name)
	}
	if column.MinLength > 0 && length < column.MinLength {
		return fmt.Errorf("%w: %s length < %d", ErrFieldLengthViolation, column.Name, column.MinLength)
	}
	if column.MaxLength > 0 && length > column.MaxLength {
		return fmt.Errorf("%w: %s length > %d", ErrFieldLengthViolation, column.Name, column.MaxLength)
	}
	return nil
}

func validateEnumConstraint(column Column, value any) error {
	if len(column.Enum) == 0 {
		return nil
	}
	targetToken, err := valueToken(value)
	if err != nil {
		return err
	}
	for _, item := range column.Enum {
		token, err := valueToken(item)
		if err != nil {
			return err
		}
		if token == targetToken {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrFieldEnumViolation, column.Name)
}

func validateCheckConstraints(column Column, value any) error {
	for _, check := range column.Checks {
		ok, err := evaluateCheck(value, check)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%w: %s %s", ErrFieldCheckViolation, column.Name, check.Operator)
		}
	}
	return nil
}

func evaluateQueryCondition(column Column, rowValue any, condition QueryCondition) (bool, error) {
	switch condition.Operator {
	case QueryOpEQ:
		return equalValues(rowValue, condition.Value)
	case QueryOpNE:
		matched, err := equalValues(rowValue, condition.Value)
		return !matched, err
	case QueryOpGT, QueryOpGTE, QueryOpLT, QueryOpLTE:
		if rowValue == nil {
			return false, nil
		}
		normalizedRowValue, err := normalizeColumnValue(column, rowValue)
		if err != nil {
			return false, err
		}
		cmp, err := compareScalarValues(normalizedRowValue, condition.Value)
		if err != nil {
			return false, err
		}
		switch condition.Operator {
		case QueryOpGT:
			return cmp > 0, nil
		case QueryOpGTE:
			return cmp >= 0, nil
		case QueryOpLT:
			return cmp < 0, nil
		case QueryOpLTE:
			return cmp <= 0, nil
		}
	case QueryOpIn, QueryOpNotIn:
		matched := false
		for _, candidate := range condition.Values {
			eq, err := equalValues(rowValue, candidate)
			if err != nil {
				return false, err
			}
			if eq {
				matched = true
				break
			}
		}
		if condition.Operator == QueryOpNotIn {
			return !matched, nil
		}
		return matched, nil
	case QueryOpBetween, QueryOpNotBetween:
		if rowValue == nil {
			return false, nil
		}
		normalizedRowValue, err := normalizeColumnValue(column, rowValue)
		if err != nil {
			return false, err
		}
		lowerCmp, err := compareScalarValues(normalizedRowValue, condition.Lower)
		if err != nil {
			return false, err
		}
		upperCmp, err := compareScalarValues(normalizedRowValue, condition.Upper)
		if err != nil {
			return false, err
		}
		matched := lowerCmp >= 0 && upperCmp <= 0
		if condition.Operator == QueryOpNotBetween {
			return !matched, nil
		}
		return matched, nil
	}
	return false, fmt.Errorf("%w: unsupported operator %s", ErrInvalidQueryCondition, condition.Operator)
}

func equalValues(left any, right any) (bool, error) {
	leftToken, err := valueToken(left)
	if err != nil {
		return false, err
	}
	rightToken, err := valueToken(right)
	if err != nil {
		return false, err
	}
	return leftToken == rightToken, nil
}

func clonePKSet(source map[string]struct{}) map[string]struct{} {
	cloned := make(map[string]struct{}, len(source))
	for key := range source {
		cloned[key] = struct{}{}
	}
	return cloned
}

func mergePKSets(target map[string]struct{}, source map[string]struct{}) {
	for key := range source {
		target[key] = struct{}{}
	}
}

func intersectPKSets(left map[string]struct{}, right map[string]struct{}) map[string]struct{} {
	if len(left) > len(right) {
		left, right = right, left
	}
	intersection := make(map[string]struct{})
	for key := range left {
		if _, exists := right[key]; exists {
			intersection[key] = struct{}{}
		}
	}
	return intersection
}

func subtractPKSets(source map[string]struct{}, removed map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{}, len(source))
	for key := range source {
		if _, exists := removed[key]; exists {
			continue
		}
		result[key] = struct{}{}
	}
	return result
}

func firstPKToken(pkSet map[string]struct{}) (string, bool) {
	for pkToken := range pkSet {
		return pkToken, true
	}
	return "", false
}

func decodeValueToken(token string) (any, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	var value any
	if err = json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func evaluateCheck(value any, check ColumnCheck) (bool, error) {
	switch check.Operator {
	case ColumnCheckGT, ColumnCheckGTE, ColumnCheckLT, ColumnCheckLTE:
		cmp, err := compareScalarValues(value, check.Value)
		if err != nil {
			return false, err
		}
		switch check.Operator {
		case ColumnCheckGT:
			return cmp > 0, nil
		case ColumnCheckGTE:
			return cmp >= 0, nil
		case ColumnCheckLT:
			return cmp < 0, nil
		case ColumnCheckLTE:
			return cmp <= 0, nil
		}
	case ColumnCheckLenGT, ColumnCheckLenGTE, ColumnCheckLenLT, ColumnCheckLenLTE:
		length, ok := valueLength(value)
		if !ok {
			return false, fmt.Errorf("%w: unsupported length check", ErrFieldCheckViolation)
		}
		expected := check.Value.(int64)
		switch check.Operator {
		case ColumnCheckLenGT:
			return int64(length) > expected, nil
		case ColumnCheckLenGTE:
			return int64(length) >= expected, nil
		case ColumnCheckLenLT:
			return int64(length) < expected, nil
		case ColumnCheckLenLTE:
			return int64(length) <= expected, nil
		}
	case ColumnCheckRegex:
		text, ok := value.(string)
		if !ok {
			return false, fmt.Errorf("%w: regex requires string", ErrFieldCheckViolation)
		}
		return regexp.MatchString(check.Value.(string), text)
	}
	return false, fmt.Errorf("%w: unsupported operator %s", ErrFieldCheckViolation, check.Operator)
}

func compareScalarValues(left any, right any) (int, error) {
	switch leftValue := left.(type) {
	case int64:
		rightValue, ok := right.(int64)
		if !ok {
			return 0, fmt.Errorf("%w: incomparable values", ErrFieldCheckViolation)
		}
		switch {
		case leftValue < rightValue:
			return -1, nil
		case leftValue > rightValue:
			return 1, nil
		default:
			return 0, nil
		}
	case float64:
		rightValue, ok := right.(float64)
		if !ok {
			return 0, fmt.Errorf("%w: incomparable values", ErrFieldCheckViolation)
		}
		switch {
		case leftValue < rightValue:
			return -1, nil
		case leftValue > rightValue:
			return 1, nil
		default:
			return 0, nil
		}
	case string:
		rightValue, ok := right.(string)
		if !ok {
			return 0, fmt.Errorf("%w: incomparable values", ErrFieldCheckViolation)
		}
		switch {
		case leftValue < rightValue:
			return -1, nil
		case leftValue > rightValue:
			return 1, nil
		default:
			return 0, nil
		}
	default:
		return 0, fmt.Errorf("%w: unsupported comparable type", ErrFieldCheckViolation)
	}
}

func valueLength(value any) (int, bool) {
	switch typed := value.(type) {
	case string:
		return len([]rune(typed)), true
	case []any:
		return len(typed), true
	default:
		return 0, false
	}
}

func boolPtr(value bool) *bool {
	ptr := new(bool)
	*ptr = value
	return ptr
}

func cloneAny(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var cloned any
	if err = json.Unmarshal(raw, &cloned); err != nil {
		return nil, err
	}
	return cloned, nil
}

func buildRowKey(pkToken string) string {
	return rowKeyPrefix + pkToken
}

func valueToken(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func asInt64(value any) (int64, bool) {
	switch typed := value.(type) {
	case int:
		return int64(typed), true
	case int8:
		return int64(typed), true
	case int16:
		return int64(typed), true
	case int32:
		return int64(typed), true
	case int64:
		return typed, true
	case uint:
		return int64(typed), true
	case uint8:
		return int64(typed), true
	case uint16:
		return int64(typed), true
	case uint32:
		return int64(typed), true
	case uint64:
		return int64(typed), typed <= uint64(^uint64(0)>>1)
	case float64:
		asInt := int64(typed)
		return asInt, float64(asInt) == typed
	case float32:
		asInt := int64(typed)
		return asInt, float32(asInt) == typed
	default:
		return 0, false
	}
}

func asFloat64(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float32:
		return float64(typed), true
	case float64:
		return typed, true
	default:
		return 0, false
	}
}
