package kernal

import (
	"fmt"
	"slices"
	"strings"

	json "github.com/json-iterator/go"
)

type CascadeQuery struct {
	Conditions []QueryCondition `json:"conditions,omitempty"`
	Includes   []CascadeInclude `json:"includes,omitempty"`
	MaxDepth   int              `json:"maxDepth,omitempty"`
}

type CascadeInclude struct {
	Table      string           `json:"table"`
	Alias      string           `json:"alias,omitempty"`
	ForeignKey string           `json:"foreignKey,omitempty"`
	Conditions []QueryCondition `json:"conditions,omitempty"`
	Includes   []CascadeInclude `json:"includes,omitempty"`
}

type relationDirection string

const (
	relationDirectionParent relationDirection = "parent"
	relationDirectionChild  relationDirection = "child"
)

type resolvedRelation struct {
	direction  relationDirection
	foreignKey ForeignKey
	alias      string
	targetDB   *SimpleDB
}

func (db *SimpleDB) FindByConditionsJSON(conditions []QueryCondition) ([]byte, error) {
	var (
		err  error
		rows []Row
	)
	if rows, err = db.FindByConditions(conditions); err != nil {
		return nil, err
	}

	return json.Marshal(rows)
}

func (db *SimpleDB) QueryCascadeJSON(query CascadeQuery) ([]byte, error) {
	var (
		err      error
		maxDepth int
		rows     []map[string]any
	)

	if maxDepth, err = db.normalizeCascadeMaxDepth(query.MaxDepth); err != nil {
		return nil, err
	}

	if rows, err = db.queryCascadeObjects(query, maxDepth, 0, []string{db.table}); err != nil {
		return nil, err
	}

	return json.Marshal(rows)
}

// QueryCascade 返回级联查询的对象结构（非 JSON）。
// 如需序列化输出，可使用 QueryCascadeJSON。
func (db *SimpleDB) QueryCascade(query CascadeQuery) ([]map[string]any, error) {
	var (
		err      error
		maxDepth int
		rows     []map[string]any
	)

	if maxDepth, err = db.normalizeCascadeMaxDepth(query.MaxDepth); err != nil {
		return nil, err
	}

	if rows, err = db.queryCascadeObjects(query, maxDepth, 0, []string{db.table}); err != nil {
		return nil, err
	}

	return rows, nil
}

func (db *SimpleDB) queryCascadeObjects(query CascadeQuery, maxDepth int, depth int, path []string) ([]map[string]any, error) {
	var (
		err      error
		rows     []Row
		objects  []map[string]any
		expanded Row
	)
	if rows, err = db.FindByConditions(query.Conditions); err != nil {
		return nil, err
	}

	objects = make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		if expanded, err = db.expandCascadeRow(row, query.Includes, maxDepth, depth, path); err != nil {
			return nil, err
		}
		objects = append(objects, expanded)
	}

	return objects, nil
}

func (db *SimpleDB) expandCascadeRow(row Row, includes []CascadeInclude, maxDepth int, depth int, path []string) (map[string]any, error) {
	var (
		err         error
		result      = make(map[string]any, len(row)+len(includes))
		targetTable string
		relation    *resolvedRelation
		closeTarget func() error
		nestedValue any
	)

	for key, value := range row {
		result[key] = value
	}

	for _, include := range includes {
		if targetTable = strings.TrimSpace(include.Table); targetTable != "" && tableInPath(path, targetTable) {
			return nil, fmt.Errorf("%w: %s", ErrCascadeCycleNotAllow, targetTable)
		}

		if relation, closeTarget, err = db.resolveCascadeRelation(include); err != nil {
			return nil, err
		}
		if closeTarget != nil {
			defer closeTarget()
		}
		if depth+1 > maxDepth {
			return nil, fmt.Errorf("%w: maxDepth=%d", ErrCascadeDepthExceeded, maxDepth)
		}
		if tableInPath(path, relation.targetDB.table) {
			return nil, fmt.Errorf("%w: %s", ErrCascadeCycleNotAllow, relation.targetDB.table)
		}

		if nestedValue, err = db.fetchCascadeRelationValue(
			row,
			include,
			relation,
			maxDepth,
			depth+1,
			append(path, relation.targetDB.table),
		); err != nil {
			return nil, err
		}
		result[relation.alias] = nestedValue
	}

	return result, nil
}

func (db *SimpleDB) resolveCascadeRelation(include CascadeInclude) (*resolvedRelation, func() error, error) {
	var (
		err           error
		targetTable   string
		currentSchema *TableSchema
		targetSchema  *TableSchema
		foreignKey    ForeignKey
		ok            bool
		targetDB      *SimpleDB
		closeTarget   func() error
	)

	if targetTable = strings.TrimSpace(include.Table); targetTable == "" {
		return nil, nil, fmt.Errorf("%w: include table is required", ErrRelationNotFound)
	}

	if currentSchema, err = db.GetSchema(); err != nil {
		return nil, nil, err
	}

	if foreignKey, ok, err = selectForeignKey(currentSchema.ForeignKeys, func(foreignKey ForeignKey) bool {
		return foreignKey.RefTable == targetTable
	}, include.ForeignKey); err != nil {
		return nil, nil, err
	} else if ok {
		if targetDB, closeTarget, err = db.openRelatedTable(targetTable); err != nil {
			return nil, nil, err
		}
		return &resolvedRelation{
			direction:  relationDirectionParent,
			foreignKey: foreignKey,
			alias:      relationAlias(include, foreignKey, targetTable),
			targetDB:   targetDB,
		}, closeTarget, nil
	}

	if targetDB, closeTarget, err = db.openRelatedTable(targetTable); err != nil {
		return nil, nil, err
	}

	if targetSchema, err = targetDB.GetSchema(); err != nil {
		if closeTarget != nil {
			_ = closeTarget()
		}
		return nil, nil, err
	}

	if foreignKey, ok, err = selectForeignKey(targetSchema.ForeignKeys, func(foreignKey ForeignKey) bool {
		return foreignKey.RefTable == db.table
	}, include.ForeignKey); err != nil {
		if closeTarget != nil {
			_ = closeTarget()
		}
		return nil, nil, err
	}
	if !ok {
		if closeTarget != nil {
			_ = closeTarget()
		}
		return nil, nil, fmt.Errorf("%w: %s -> %s", ErrRelationNotFound, db.table, targetTable)
	}

	return &resolvedRelation{
		direction:  relationDirectionChild,
		foreignKey: foreignKey,
		alias:      relationAlias(include, foreignKey, targetTable),
		targetDB:   targetDB,
	}, closeTarget, nil
}

func (db *SimpleDB) fetchCascadeRelationValue(
	row Row,
	include CascadeInclude,
	relation *resolvedRelation,
	maxDepth int,
	depth int,
	path []string,
) (any, error) {
	var (
		err        error
		conditions = make([]QueryCondition, 0, len(include.Conditions)+1)
		results    []map[string]any
	)

	switch relation.direction {
	case relationDirectionParent:
		value, exists := row[relation.foreignKey.Field]
		if !exists || value == nil {
			return nil, nil
		}
		conditions = append(conditions, QueryCondition{Field: relation.foreignKey.RefField, Operator: QueryOpEQ, Value: value})
		conditions = append(conditions, include.Conditions...)
		if results, err = relation.
			targetDB.
			queryCascadeObjects(
				CascadeQuery{
					Conditions: conditions,
					Includes:   include.Includes,
				},
				maxDepth,
				depth,
				path,
			); err != nil {
			return nil, err
		}
		switch len(results) {
		case 0:
			return nil, nil
		case 1:
			return results[0], nil
		}

		return results, nil
	case relationDirectionChild:
		value, exists := row[relation.foreignKey.RefField]
		if !exists || value == nil {
			return []map[string]any{}, nil
		}
		conditions = append(conditions, QueryCondition{Field: relation.foreignKey.Field, Operator: QueryOpEQ, Value: value})
		conditions = append(conditions, include.Conditions...)
		return relation.
			targetDB.
			queryCascadeObjects(
				CascadeQuery{
					Conditions: conditions,
					Includes:   include.Includes,
				},
				maxDepth,
				depth,
				path,
			)
	default:
		return nil, fmt.Errorf("%w: unsupported relation direction", ErrRelationNotFound)
	}
}

func (db *SimpleDB) normalizeCascadeMaxDepth(value int) (int, error) {
	if value == 0 {
		return db.getDefaultCascadeMaxDepth(), nil
	}
	if value < 0 {
		return 0, fmt.Errorf("%w: maxDepth must be >= 0", ErrInvalidQueryCondition)
	}
	if value > HardCascadeMaxDepthLimit {
		return 0, fmt.Errorf("%w: maxDepth > %d", ErrCascadeDepthExceeded, HardCascadeMaxDepthLimit)
	}
	return value, nil
}

func tableInPath(path []string, table string) bool { return slices.Contains(path, table) }

func (db *SimpleDB) openRelatedTable(table string) (*SimpleDB, func() error, error) {
	var (
		err     error
		related *SimpleDB
	)

	table = strings.TrimSpace(table)
	switch table {
	case "":
		return nil, nil, fmt.Errorf("%w: empty target table", ErrRelationNotFound)
	case db.table:
		return db, nil, nil
	default:
		if related, err = newSimpleDB(db.database, table); err != nil {
			return nil, nil, err
		}
		return related, related.Close, nil
	}
}

func selectForeignKey(
	foreignKeys []ForeignKey,
	predicate func(ForeignKey) bool,
	requested string,
) (ForeignKey, bool, error) {
	var matches = make([]ForeignKey, 0, len(foreignKeys))

	requested = strings.TrimSpace(requested)

	for _, foreignKey := range foreignKeys {
		if !predicate(foreignKey) || (requested != "" && requested != foreignKey.Name && requested != foreignKey.Field && requested != foreignKey.Alias) {
			continue
		}
		matches = append(matches, foreignKey)
	}
	if len(matches) == 0 {
		return ForeignKey{}, false, nil
	}
	if len(matches) > 1 && requested == "" {
		return ForeignKey{}, false, fmt.Errorf("%w: ambiguous relation", ErrRelationNotFound)
	}
	return matches[0], true, nil
}

func relationAlias(include CascadeInclude, foreignKey ForeignKey, fallback string) string {
	var alias string

	if alias = strings.TrimSpace(include.Alias); alias != "" {
		return alias
	}

	if alias = strings.TrimSpace(foreignKey.Alias); alias != "" {
		return alias
	}

	return fallback
}
