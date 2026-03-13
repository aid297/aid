package api

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aid297/aid/simpleDB/driver"
	"github.com/aid297/aid/simpleDB/kernal"
)

var ErrSystemTableAccessDenied = fmt.Errorf("system table access denied: super_admin role required")

const superAdminRoleCode = "super_admin"

type Backend string

const (
	BackendDriver Backend = "driver"
	BackendKernal Backend = "kernal"
)

type Engine struct {
	Database string
	Backend  Backend
	Actor    *driver.AuthenticatedUser
	cache    map[string]tableDB
}

func NewEngine(database string, backend Backend) *Engine {
	if backend == "" {
		backend = BackendDriver
	}
	return &Engine{
		Database: database,
		Backend:  backend,
		cache:    make(map[string]tableDB),
	}
}

func (e *Engine) Close() error {
	for _, db := range e.cache {
		db.Close()
	}
	e.cache = make(map[string]tableDB)
	return nil
}

func (e *Engine) WithActor(user *driver.AuthenticatedUser) *Engine {
	e.Actor = user
	return e
}

func (e *Engine) Parse(sql string) (Statement, error) { return Parse(sql) }

func (e *Engine) Execute(sql string) (ExecResult, error) {
	stmt, err := Parse(sql)
	if err != nil {
		return ExecResult{}, err
	}
	return e.ExecuteStatement(stmt)
}

func (e *Engine) ExecuteStatement(stmt Statement) (ExecResult, error) {
	switch s := stmt.(type) {
	case CreateTableStmt:
		return e.execCreate(s)
	case AlterTableStmt:
		return e.execAlter(s)
	case InsertStmt:
		return e.execInsert(s)
	case UpdateStmt:
		return e.execUpdate(s)
	case DeleteStmt:
		return e.execDelete(s)
	case SelectStmt:
		return e.execSelect(s)
	case DropTableStmt:
		return e.execDrop(s)
	case TruncateTableStmt:
		return e.execTruncate(s)
	default:
		return ExecResult{}, fmt.Errorf("unsupported statement")
	}
}

type tableDB interface {
	Close() error
	CreateTable(schema kernal.TableSchema) error
	AlterTable(plan kernal.AlterTablePlan) error
	InsertRow(values kernal.Row) (kernal.Row, error)
	UpdateRow(primaryKey any, updates kernal.Row) (kernal.Row, error)
	RemoveByCondition(conditions ...kernal.QueryCondition) (int, error)
	Find(conditions ...kernal.QueryCondition) ([]kernal.Row, error)
	DropTable() error
	TruncateTable() error
	SetPersistenceConfig(windowSecs int, windowBytes uint64, threshold uint64)
	GetPath() string
}

func (e *Engine) open(table string) (tableDB, error) {
	return e.openWithScope(table, kernal.TableAccessScopeDML)
}

func (e *Engine) openWithScope(table string, scope kernal.TableAccessScope) (tableDB, error) {
	if err := e.authorizeTableOp(table, scope); err != nil {
		return nil, err
	}

	if db, ok := e.cache[table]; ok {
		return db, nil
	}

	var db tableDB
	var err error
	if e.Backend == BackendKernal {
		db, err = kernal.New.DB(e.Database, table)
	} else {
		db, err = driver.New.DB(e.Database, table)
	}

	if err != nil {
		return nil, err
	}

	e.cache[table] = db
	return db, nil
}

func (e *Engine) authorizeTable(table string) error {
	return e.authorizeTableOp(table, kernal.TableAccessScopeDML)
}

func (e *Engine) authorizeTableOp(table string, scope kernal.TableAccessScope) error {
	table = strings.TrimSpace(table)
	// system tables: only super_admin can access
	if strings.HasPrefix(table, "_sys_") {
		if e.Actor != nil && hasRole(e.Actor.Roles, superAdminRoleCode) {
			return nil
		}
		return ErrSystemTableAccessDenied
	}
	_ = scope
	// user tables: check database-level binding
	return kernal.New.CheckDatabaseBinding(e.Database, e.Actor)
}

func hasRole(roles []string, roleCode string) bool {
	for _, role := range roles {
		if role == roleCode {
			return true
		}
	}
	return false
}

func (e *Engine) execCreate(s CreateTableStmt) (ExecResult, error) {
	db, err := e.openWithScope(s.Table, kernal.TableAccessScopeDDL)
	if err != nil {
		return ExecResult{}, err
	}
	if err = db.CreateTable(s.Schema); err != nil {
		return ExecResult{}, err
	}
	// register creator as table owner
	if e.Actor != nil && strings.TrimSpace(e.Actor.ID) != "" {
		if err = kernal.New.RegisterTableOwner(e.Database, s.Table, e.Actor.ID); err != nil {
			return ExecResult{}, err
		}
	}
	return ExecResult{Statement: StmtCreateTable, Affected: 1}, nil
}

func (e *Engine) execAlter(s AlterTableStmt) (ExecResult, error) {
	db, err := e.openWithScope(s.Table, kernal.TableAccessScopeDDL)
	if err != nil {
		return ExecResult{}, err
	}
	if err = db.AlterTable(s.Plan); err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtAlterTable, Affected: 1}, nil
}

func (e *Engine) execInsert(s InsertStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	rows := s.Rows
	if len(rows) == 0 && len(s.Row) > 0 {
		rows = []kernal.Row{s.Row}
	}
	if len(rows) == 0 {
		return ExecResult{}, fmt.Errorf("insert values is empty")
	}

	insertedRows := make([]kernal.Row, 0, len(rows))
	for _, values := range rows {
		row, insertErr := db.InsertRow(values)
		if insertErr != nil {
			return ExecResult{}, insertErr
		}
		insertedRows = append(insertedRows, row)
	}

	result := ExecResult{Statement: StmtInsert, Affected: len(insertedRows), InsertedRows: insertedRows}
	if len(insertedRows) == 1 {
		result.Inserted = insertedRows[0]
	}
	return result, nil
}

func (e *Engine) execUpdate(s UpdateStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}

	primaryKeys := append([]any(nil), s.PrimaryKeys...)
	if len(primaryKeys) == 0 {
		primaryKeys = append(primaryKeys, s.PrimaryKey)
	}
	if len(primaryKeys) == 0 {
		return ExecResult{}, fmt.Errorf("update primary key is required")
	}

	updatedRows := make([]kernal.Row, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		row, updateErr := db.UpdateRow(primaryKey, s.Updates)
		if updateErr != nil {
			return ExecResult{}, updateErr
		}
		updatedRows = append(updatedRows, row)
	}

	result := ExecResult{Statement: StmtUpdate, Affected: len(updatedRows), UpdatedRows: updatedRows}
	if len(updatedRows) == 1 {
		result.Updated = updatedRows[0]
	}
	return result, nil
}

func (e *Engine) execDelete(s DeleteStmt) (ExecResult, error) {
	if len(s.Conditions) == 0 {
		return ExecResult{}, fmt.Errorf("delete requires where condition")
	}
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}
	count, err := db.RemoveByCondition(s.Conditions...)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtDelete, Affected: count}, nil
}

func (e *Engine) execSelect(s SelectStmt) (ExecResult, error) {
	db, err := e.open(s.Table)
	if err != nil {
		return ExecResult{}, err
	}

	conditions := s.Conditions // 展开子查询条件，将 IN/NOT IN (SELECT ...) 转为字面量值列表追加到 Conditions
	for _, sc := range s.SubConds {
		subResult, subErr := e.execSelect(sc.SubStmt)
		if subErr != nil {
			return ExecResult{}, fmt.Errorf("subquery error: %v", subErr)
		}
		// 收集子查询第一列的所有值
		values := make([]any, 0, len(subResult.Rows))
		for _, row := range subResult.Rows {
			for _, v := range row {
				values = append(values, v)
				break // 只取第一列
			}
		}
		op := kernal.QueryOpIn
		if sc.NotIn {
			op = kernal.QueryOpNotIn
		}
		conditions = append(conditions, kernal.QueryCondition{Field: sc.Field, Operator: op, Value: values})
	}

	rows, err := db.Find(conditions...)
	if err != nil {
		return ExecResult{}, err
	}

	// 有 JOIN 子句时在内存中做连接
	if len(s.Joins) > 0 {
		rows, err = e.applyJoins(s.Table, rows, s.Joins, conditions)
		if err != nil {
			return ExecResult{}, err
		}
	}

	if len(s.GroupBy) > 0 {
		rows = applyGroupBy(rows, s.Fields, s.GroupBy)
	} else {
		rows = applyProjection(rows, s.Fields)
	}
	if s.OrderBy != "" {
		sort.Slice(rows, func(i, j int) bool {
			cmp := compareRowValues(rows[i][s.OrderBy], rows[j][s.OrderBy])
			if s.OrderDesc {
				return cmp > 0
			}
			return cmp < 0
		})
	}
	if s.Offset > 0 {
		if s.Offset >= len(rows) {
			rows = rows[:0]
		} else {
			rows = rows[s.Offset:]
		}
	}
	if s.Limit > 0 && len(rows) > s.Limit {
		rows = rows[:s.Limit]
	}
	return ExecResult{Statement: StmtSelect, Rows: rows, Affected: len(rows)}, nil
}

// applyJoins 逐条执行 JOIN，支持 INNER JOIN 和 LEFT JOIN
// 列名冲突时以 "table.field" 形式存入合并行
func (e *Engine) applyJoins(mainTable string, mainRows []kernal.Row, joins []JoinClause, _ []kernal.QueryCondition) ([]kernal.Row, error) {
	current := mainRows
	for _, j := range joins {
		jdb, err := e.open(j.Table)
		if err != nil {
			return nil, err
		}
		rightRows, err := jdb.Find()
		jdb.Close()
		if err != nil {
			return nil, err
		}

		// 构建右表索引：ON 右侧字段值 → []row
		rightField := resolveField(j.RightAlias, j.Table)
		rightIndex := make(map[string][]kernal.Row, len(rightRows))
		for _, rr := range rightRows {
			key := fmt.Sprintf("%v", rr[rightField])
			rightIndex[key] = append(rightIndex[key], rr)
		}

		leftField := resolveField(j.LeftAlias, mainTable)
		merged := make([]kernal.Row, 0, len(current))
		for _, lr := range current {
			key := fmt.Sprintf("%v", lr[leftField])
			matches := rightIndex[key]
			if len(matches) == 0 {
				if j.Type == JoinLeft {
					// LEFT JOIN：保留主表行，右侧字段填 nil
					newRow := mergeRows(mainTable, lr, j.Table, nil)
					merged = append(merged, newRow)
				}
				// INNER JOIN：丢弃无匹配行
				continue
			}
			for _, rr := range matches {
				newRow := mergeRows(mainTable, lr, j.Table, rr)
				merged = append(merged, newRow)
			}
		}
		current = merged
	}
	return current, nil
}

// resolveField 从 "table.field" 或 "field" 中提取字段名
func resolveField(alias, _ string) string {
	if idx := strings.LastIndex(alias, "."); idx >= 0 {
		return alias[idx+1:]
	}
	return alias
}

// mergeRows 将左右两行合并为一行，字段名冲突时加 "table." 前缀
func mergeRows(leftTable string, left kernal.Row, rightTable string, right kernal.Row) kernal.Row {
	row := make(kernal.Row, len(left)+len(right))
	// 先收集所有键
	rightKeys := make(map[string]struct{}, len(right))
	for k := range right {
		rightKeys[k] = struct{}{}
	}
	for k, v := range left {
		if _, conflict := rightKeys[k]; conflict {
			row[leftTable+"."+k] = v
		} else {
			row[k] = v
		}
	}
	for k, v := range right {
		if _, conflict := left[k]; conflict {
			row[rightTable+"."+k] = v
		} else {
			row[k] = v
		}
	}
	return row
}

func applyProjection(rows []kernal.Row, fields []SelectField) []kernal.Row {
	if len(fields) == 0 {
		return rows
	}
	hasStar, hasAgg := false, false
	for _, f := range fields {
		if f.Star {
			hasStar = true
		}
		if f.Agg != AggNone {
			hasAgg = true
		}
	}
	if hasStar && !hasAgg {
		return rows
	}
	// 无 GROUP BY 的纯聚合：返回单行
	if hasAgg {
		resultRow := kernal.Row{}
		for _, f := range fields {
			if f.Agg != AggNone {
				key := f.Alias
				if key == "" {
					key = string(f.Agg) + "(" + f.AggField + ")"
				}
				resultRow[key] = computeAggregate(f.Agg, f.AggField, rows)
			} else if f.Star {
				if len(rows) > 0 {
					for k, v := range rows[0] {
						resultRow[k] = v
					}
				}
			} else if f.Field != "" {
				key := f.Alias
				if key == "" {
					key = f.Field
				}
				if len(rows) > 0 {
					resultRow[key] = rows[0][f.Field]
				}
			}
		}
		return []kernal.Row{resultRow}
	}
	// 纯字段投影
	result := make([]kernal.Row, 0, len(rows))
	for _, row := range rows {
		newRow := kernal.Row{}
		for _, f := range fields {
			if f.Star {
				for k, v := range row {
					newRow[k] = v
				}
			} else if f.Field != "" {
				key := f.Alias
				if key == "" {
					key = f.Field
				}
				newRow[key] = row[f.Field]
			}
		}
		result = append(result, newRow)
	}
	return result
}

func applyGroupBy(rows []kernal.Row, fields []SelectField, groupBy []string) []kernal.Row {
	type entry struct{ rows []kernal.Row }
	orderMap := make(map[string]*entry)
	order := make([]string, 0)
	for _, row := range rows {
		parts := make([]string, 0, len(groupBy))
		for _, g := range groupBy {
			parts = append(parts, fmt.Sprintf("%v", row[g]))
		}
		key := strings.Join(parts, "\x00")
		if e, ok := orderMap[key]; ok {
			e.rows = append(e.rows, row)
		} else {
			orderMap[key] = &entry{rows: []kernal.Row{row}}
			order = append(order, key)
		}
	}
	result := make([]kernal.Row, 0, len(order))
	for _, key := range order {
		groupRows := orderMap[key].rows
		newRow := kernal.Row{}
		for _, f := range fields {
			if f.Star {
				for k, v := range groupRows[0] {
					newRow[k] = v
				}
			} else if f.Agg != AggNone {
				colKey := f.Alias
				if colKey == "" {
					colKey = string(f.Agg) + "(" + f.AggField + ")"
				}
				newRow[colKey] = computeAggregate(f.Agg, f.AggField, groupRows)
			} else if f.Field != "" {
				colKey := f.Alias
				if colKey == "" {
					colKey = f.Field
				}
				newRow[colKey] = groupRows[0][f.Field]
			}
		}
		result = append(result, newRow)
	}
	return result
}

func computeAggregate(agg AggFunc, field string, rows []kernal.Row) any {
	switch agg {
	case AggCount:
		if field == "*" {
			return int64(len(rows))
		}
		count := int64(0)
		for _, row := range rows {
			if row[field] != nil {
				count++
			}
		}
		return count
	case AggSum:
		sum := float64(0)
		for _, row := range rows {
			sum += toFloat64(row[field])
		}
		return sum
	case AggAvg:
		if len(rows) == 0 {
			return nil
		}
		sum := float64(0)
		for _, row := range rows {
			sum += toFloat64(row[field])
		}
		return sum / float64(len(rows))
	case AggMin:
		if len(rows) == 0 {
			return nil
		}
		min := rows[0][field]
		for _, row := range rows[1:] {
			if compareRowValues(row[field], min) < 0 {
				min = row[field]
			}
		}
		return min
	case AggMax:
		if len(rows) == 0 {
			return nil
		}
		max := rows[0][field]
		for _, row := range rows[1:] {
			if compareRowValues(row[field], max) > 0 {
				max = row[field]
			}
		}
		return max
	}
	return nil
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case int64:
		return float64(val)
	case float64:
		return val
	case int:
		return float64(val)
	}
	return 0
}

func compareRowValues(a, b any) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	switch av := a.(type) {
	case int64:
		switch bv := b.(type) {
		case int64:
			if av < bv {
				return -1
			}
			if av > bv {
				return 1
			}
			return 0
		case float64:
			fav := float64(av)
			if fav < bv {
				return -1
			}
			if fav > bv {
				return 1
			}
			return 0
		}
	case float64:
		var bv float64
		switch bval := b.(type) {
		case float64:
			bv = bval
		case int64:
			bv = float64(bval)
		default:
			return strings.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
		}
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
		return 0
	case string:
		if bv, ok := b.(string); ok {
			return strings.Compare(av, bv)
		}
	}
	return strings.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
}

func (e *Engine) execDrop(s DropTableStmt) (ExecResult, error) {
	db, err := e.openWithScope(s.Table, kernal.TableAccessScopeDDL)
	if err != nil {
		return ExecResult{}, err
	}
	if err = db.DropTable(); err != nil {
		return ExecResult{}, err
	}
	delete(e.cache, s.Table)
	return ExecResult{Statement: StmtDropTable, Affected: 1}, nil
}

func (e *Engine) execTruncate(s TruncateTableStmt) (ExecResult, error) {
	db, err := e.openWithScope(s.Table, kernal.TableAccessScopeDDL)
	if err != nil {
		return ExecResult{}, err
	}
	if err = db.TruncateTable(); err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Statement: StmtTruncate, Affected: 1}, nil
}
