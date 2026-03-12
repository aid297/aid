package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aid297/aid/simpleDB/kernal"
)

var (
	reCreate   = regexp.MustCompile(`(?i)^CREATE\s+TABLE\s+([a-zA-Z_][\w]*)\s*\((.*)\)$`)
	reAlter    = regexp.MustCompile(`(?i)^ALTER\s+TABLE\s+([a-zA-Z_][\w]*)\s+(.+)$`)
	reInsert   = regexp.MustCompile(`(?i)^INSERT\s+INTO\s+([a-zA-Z_][\w]*)\s*\((.*)\)\s*VALUES\s*(.+)$`)
	reUpdate   = regexp.MustCompile(`(?i)^UPDATE\s+([a-zA-Z_][\w]*)\s+SET\s+(.+?)(?:\s+WHERE\s+(.+))?$`)
	reDelete   = regexp.MustCompile(`(?i)^DELETE\s+FROM\s+([a-zA-Z_][\w]*)(?:\s+WHERE\s+(.+))?$`)
	reSelect   = regexp.MustCompile(`(?i)^SELECT\s+(.+?)\s+FROM\s+([a-zA-Z_][\w]*)((?:\s+(?:INNER|LEFT)\s+JOIN\s+[a-zA-Z_][\w]*\s+ON\s+[^\s]+\s*=\s*[^\s]+)*)(?:\s+WHERE\s+(.+?))?(?:\s+GROUP\s+BY\s+(.+?))?(?:\s+ORDER\s+BY\s+([a-zA-Z_][\w.]*)(?:\s+(ASC|DESC))?)?(?:\s+LIMIT\s+(\d+))?(?:\s+OFFSET\s+(\d+))?$`)
	reDrop     = regexp.MustCompile(`(?i)^DROP\s+TABLE\s+([a-zA-Z_][\w]*)$`)
	reTruncate = regexp.MustCompile(`(?i)^TRUNCATE\s+TABLE\s+([a-zA-Z_][\w]*)$`)
)

func Parse(sql string) (Statement, error) {
	sql = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(sql), ";"))
	if sql == "" {
		return nil, fmt.Errorf("empty sql")
	}

	if m := reCreate.FindStringSubmatch(sql); len(m) == 3 {
		return parseCreateTable(m[1], m[2])
	}
	if m := reAlter.FindStringSubmatch(sql); len(m) == 3 {
		return parseAlterTable(m[1], m[2])
	}
	if m := reInsert.FindStringSubmatch(sql); len(m) == 4 {
		return parseInsert(m[1], m[2], m[3])
	}
	if m := reUpdate.FindStringSubmatch(sql); len(m) == 4 {
		return parseUpdate(m[1], m[2], m[3])
	}
	if m := reDelete.FindStringSubmatch(sql); len(m) == 3 {
		return parseDelete(m[1], m[2])
	}
	if m := reSelect.FindStringSubmatch(sql); len(m) == 10 {
		return parseSelect(m[2], m[1], m[3], m[4], m[5], m[6], m[7], m[8], m[9])
	}
	if m := reDrop.FindStringSubmatch(sql); len(m) == 2 {
		return DropTableStmt{Table: m[1]}, nil
	}
	if m := reTruncate.FindStringSubmatch(sql); len(m) == 2 {
		return TruncateTableStmt{Table: m[1]}, nil
	}

	return nil, fmt.Errorf("unsupported sql")
}

func parseCreateTable(table, body string) (Statement, error) {
	parts := splitCSV(body)
	if len(parts) == 0 {
		return nil, fmt.Errorf("create table has no columns")
	}

	schema := kernal.TableSchema{Columns: make([]kernal.Column, 0), ForeignKeys: make([]kernal.ForeignKey, 0)}
	for _, part := range parts {
		p := strings.TrimSpace(part)
		up := strings.ToUpper(p)
		if strings.HasPrefix(up, "FOREIGN KEY") {
			fk, err := parseForeignKeyDef(p)
			if err != nil {
				return nil, err
			}
			schema.ForeignKeys = append(schema.ForeignKeys, fk)
			continue
		}

		col, err := parseColumnDef(p)
		if err != nil {
			return nil, err
		}
		schema.Columns = append(schema.Columns, col)
	}

	return CreateTableStmt{Table: table, Schema: schema}, nil
}

func parseAlterTable(table, body string) (Statement, error) {
	parts := splitCSV(body)
	plan := kernal.AlterTablePlan{}
	for _, raw := range parts {
		op := strings.TrimSpace(raw)
		up := strings.ToUpper(op)

		switch {
		case strings.HasPrefix(up, "ADD COLUMN "):
			colDef := strings.TrimSpace(op[len("ADD COLUMN "):])
			col, err := parseColumnDef(colDef)
			if err != nil {
				return nil, err
			}
			plan.AddColumns = append(plan.AddColumns, col)
		case strings.HasPrefix(up, "DROP COLUMN "):
			plan.DropColumns = append(plan.DropColumns, strings.TrimSpace(op[len("DROP COLUMN "):]))
		case strings.HasPrefix(up, "ADD UNIQUE"):
			field, err := parseFieldInParen(op)
			if err != nil {
				return nil, err
			}
			plan.AddUniques = append(plan.AddUniques, field)
		case strings.HasPrefix(up, "DROP UNIQUE"):
			field, err := parseFieldInParen(op)
			if err != nil {
				return nil, err
			}
			plan.DropUniques = append(plan.DropUniques, field)
		case strings.HasPrefix(up, "ADD INDEX"):
			field, err := parseFieldInParen(op)
			if err != nil {
				return nil, err
			}
			plan.AddIndexes = append(plan.AddIndexes, field)
		case strings.HasPrefix(up, "DROP INDEX"):
			field, err := parseFieldInParen(op)
			if err != nil {
				return nil, err
			}
			plan.DropIndexes = append(plan.DropIndexes, field)
		case strings.HasPrefix(up, "ADD FOREIGN KEY"):
			fk, err := parseAlterAddForeignKey(op)
			if err != nil {
				return nil, err
			}
			plan.AddForeignKeys = append(plan.AddForeignKeys, fk)
		case strings.HasPrefix(up, "DROP FOREIGN KEY "):
			name := strings.TrimSpace(op[len("DROP FOREIGN KEY "):])
			plan.DropForeignKeys = append(plan.DropForeignKeys, name)
		default:
			return nil, fmt.Errorf("unsupported alter operation: %s", op)
		}
	}

	return AlterTableStmt{Table: table, Plan: plan}, nil
}

func parseInsert(table, fieldsRaw, valuesRaw string) (Statement, error) {
	fields := splitCSV(fieldsRaw)
	tuples, err := parseValuesTuples(valuesRaw)
	if err != nil {
		return nil, err
	}
	if len(tuples) == 0 {
		return nil, fmt.Errorf("insert values is empty")
	}

	rows := make([]kernal.Row, 0, len(tuples))
	for _, tuple := range tuples {
		values := splitCSV(tuple)
		if len(fields) != len(values) {
			return nil, fmt.Errorf("insert fields/values length mismatch")
		}
		row := kernal.Row{}
		for i := range fields {
			v, parseErr := parseLiteral(values[i])
			if parseErr != nil {
				return nil, parseErr
			}
			row[strings.TrimSpace(fields[i])] = v
		}
		rows = append(rows, row)
	}

	if len(rows) == 1 {
		return InsertStmt{Table: table, Row: rows[0], Rows: rows}, nil
	}
	return InsertStmt{Table: table, Rows: rows}, nil
}

func parseUpdate(table, setRaw, whereRaw string) (Statement, error) {
	updates := kernal.Row{}
	for _, item := range splitCSV(setRaw) {
		pair := strings.SplitN(item, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid set clause: %s", item)
		}
		value, err := parseLiteral(pair[1])
		if err != nil {
			return nil, err
		}
		updates[strings.TrimSpace(pair[0])] = value
	}

	if strings.TrimSpace(whereRaw) == "" {
		return nil, fmt.Errorf("update requires where primary key condition")
	}
	conditions, _, err := parseConditions(whereRaw)
	if err != nil {
		return nil, err
	}
	if len(conditions) != 1 {
		return nil, fmt.Errorf("update requires exactly one where condition on primary key")
	}

	condition := conditions[0]
	if condition.Operator == kernal.QueryOpEQ {
		return UpdateStmt{Table: table, PrimaryKey: condition.Value, Updates: updates}, nil
	}
	if condition.Operator == kernal.QueryOpIn {
		values, ok := condition.Value.([]any)
		if !ok || len(values) == 0 {
			return nil, fmt.Errorf("update IN condition requires non-empty values")
		}
		return UpdateStmt{Table: table, PrimaryKeys: values, Updates: updates}, nil
	}

	return nil, fmt.Errorf("update requires where condition operator '=' or 'IN' on primary key")
}

func parseValuesTuples(valuesRaw string) ([]string, error) {
	valuesRaw = strings.TrimSpace(valuesRaw)
	if valuesRaw == "" {
		return nil, nil
	}

	result := make([]string, 0)
	start := -1
	depth := 0
	inQuote := false

	for i := 0; i < len(valuesRaw); i++ {
		ch := valuesRaw[i]
		if ch == '\'' {
			inQuote = !inQuote
		}
		if inQuote {
			continue
		}

		switch ch {
		case '(':
			if depth == 0 {
				start = i + 1
			}
			depth++
		case ')':
			if depth == 0 {
				return nil, fmt.Errorf("invalid insert values syntax")
			}
			depth--
			if depth == 0 {
				result = append(result, strings.TrimSpace(valuesRaw[start:i]))
				start = -1
			}
		}
	}
	if depth != 0 {
		return nil, fmt.Errorf("invalid insert values syntax")
	}
	return result, nil
}

func parseDelete(table, whereRaw string) (Statement, error) {
	conditions, _, err := parseConditions(whereRaw)
	if err != nil {
		return nil, err
	}
	if len(conditions) == 0 {
		return nil, fmt.Errorf("delete requires where condition")
	}
	return DeleteStmt{Table: table, Conditions: conditions}, nil
}

func parseSelect(table, fieldsRaw, joinRaw, whereRaw, groupByRaw, orderField, orderDir, limitRaw, offsetRaw string) (Statement, error) {
	fields, err := parseSelectFields(fieldsRaw)
	if err != nil {
		return nil, err
	}
	joins, err := parseJoinClauses(joinRaw)
	if err != nil {
		return nil, err
	}
	conditions, subConds, err := parseConditions(whereRaw)
	if err != nil {
		return nil, err
	}
	var groupBy []string
	if s := strings.TrimSpace(groupByRaw); s != "" {
		for _, g := range splitCSV(s) {
			groupBy = append(groupBy, strings.TrimSpace(g))
		}
	}
	desc := strings.EqualFold(strings.TrimSpace(orderDir), "desc")
	var limit, offset int
	if s := strings.TrimSpace(limitRaw); s != "" {
		n, e := strconv.Atoi(s)
		if e != nil {
			return nil, fmt.Errorf("invalid LIMIT value: %s", s)
		}
		limit = n
	}
	if s := strings.TrimSpace(offsetRaw); s != "" {
		n, e := strconv.Atoi(s)
		if e != nil {
			return nil, fmt.Errorf("invalid OFFSET value: %s", s)
		}
		offset = n
	}
	return SelectStmt{
		Table:      table,
		Fields:     fields,
		Joins:      joins,
		Conditions: conditions,
		SubConds:   subConds,
		GroupBy:    groupBy,
		OrderBy:    strings.TrimSpace(orderField),
		OrderDesc:  desc,
		Limit:      limit,
		Offset:     offset,
	}, nil
}

// parseJoinClauses 解析一或多个 "[INNER|LEFT] JOIN table ON left = right" 子句
func parseJoinClauses(raw string) ([]JoinClause, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	// 按 JOIN 关键字切割（保留类型词）
	reJoin := regexp.MustCompile(`(?i)\b(INNER\s+JOIN|LEFT\s+JOIN)\s+([a-zA-Z_][\w]*)\s+ON\s+([^\s=]+)\s*=\s*([^\s]+)`)
	matches := reJoin.FindAllStringSubmatch(raw, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid join clause: %s", raw)
	}
	joins := make([]JoinClause, 0, len(matches))
	for _, m := range matches {
		joinTypeRaw := strings.ToUpper(strings.Fields(m[1])[0]) // INNER or LEFT
		jt := JoinInner
		if joinTypeRaw == "LEFT" {
			jt = JoinLeft
		}
		joins = append(joins, JoinClause{
			Type:       jt,
			Table:      m[2],
			LeftAlias:  strings.TrimSpace(m[3]),
			RightAlias: strings.TrimSpace(m[4]),
		})
	}
	return joins, nil
}

func parseSelectFields(raw string) ([]SelectField, error) {
	raw = strings.TrimSpace(raw)
	if raw == "*" {
		return []SelectField{{Star: true}}, nil
	}
	parts := splitCSV(raw)
	reAgg := regexp.MustCompile(`(?i)^(COUNT|SUM|AVG|MIN|MAX)\(([^)]+)\)(?:\s+AS\s+([a-zA-Z_][\w]*))?$`)
	reField := regexp.MustCompile(`(?i)^([a-zA-Z_][\w]*)(?:\s+AS\s+([a-zA-Z_][\w]*))?$`)
	fields := make([]SelectField, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "*" {
			fields = append(fields, SelectField{Star: true})
			continue
		}
		if m := reAgg.FindStringSubmatch(part); len(m) == 4 {
			fields = append(fields, SelectField{
				Agg:      AggFunc(strings.ToUpper(m[1])),
				AggField: strings.TrimSpace(m[2]),
				Alias:    m[3],
			})
			continue
		}
		if m := reField.FindStringSubmatch(part); len(m) == 3 {
			fields = append(fields, SelectField{Field: m[1], Alias: m[2]})
			continue
		}
		return nil, fmt.Errorf("invalid select field: %s", part)
	}
	return fields, nil
}

func parseConditions(whereRaw string) ([]kernal.QueryCondition, []SubqueryCondition, error) {
	whereRaw = strings.TrimSpace(whereRaw)
	if whereRaw == "" {
		return nil, nil, nil
	}
	parts := splitByAND(whereRaw)
	conditions := make([]kernal.QueryCondition, 0, len(parts))
	subConds := make([]SubqueryCondition, 0)
	for _, part := range parts {
		field, op, value, err := parseSingleCondition(part)
		if err != nil {
			return nil, nil, err
		}
		if op == "__subquery_in__" || op == "__subquery_not_in__" {
			subSQL, _ := value.(string)
			subStmt, err := parseSelect("", "", "", subSQL, "", "", "", "", "")
			if err != nil {
				return nil, nil, fmt.Errorf("invalid subquery: %v", err)
			}
			sel, ok := subStmt.(SelectStmt)
			if !ok {
				return nil, nil, fmt.Errorf("subquery must be a SELECT statement")
			}
			subConds = append(subConds, SubqueryCondition{
				Field:   field,
				NotIn:   op == "__subquery_not_in__",
				SubStmt: sel,
			})
			continue
		}
		conditions = append(conditions, kernal.QueryCondition{Field: field, Operator: op, Value: value})
	}
	return conditions, subConds, nil
}

func parseSingleCondition(part string) (string, string, any, error) {
	part = strings.TrimSpace(part)

	// IN/NOT IN 子查询：field [NOT] IN (SELECT ...)
	reSubIn := regexp.MustCompile(`(?i)^([a-zA-Z_][\w]*)\s+(NOT\s+IN|IN)\s*\(\s*(SELECT\s+.+)\)$`)
	if m := reSubIn.FindStringSubmatch(part); len(m) == 4 {
		// 用特殊 op 占位，engine 层识别后展开子查询
		notIn := strings.Contains(strings.ToUpper(m[2]), "NOT")
		opMark := "__subquery_in__"
		if notIn {
			opMark = "__subquery_not_in__"
		}
		return strings.TrimSpace(m[1]), opMark, strings.TrimSpace(m[3]), nil
	}

	// NOT IN 字面量：field NOT IN (v1, v2, ...)
	reNotIn := regexp.MustCompile(`(?i)^([a-zA-Z_][\w]*)\s+NOT\s+IN\s*\((.*)\)$`)
	if m := reNotIn.FindStringSubmatch(part); len(m) == 3 {
		items := splitCSV(m[2])
		values := make([]any, 0, len(items))
		for _, item := range items {
			value, err := parseLiteral(item)
			if err != nil {
				return "", "", nil, err
			}
			values = append(values, value)
		}
		if len(values) == 0 {
			return "", "", nil, fmt.Errorf("invalid condition: %s", part)
		}
		return strings.TrimSpace(m[1]), kernal.QueryOpNotIn, values, nil
	}

	// IN 字面量：field IN (v1, v2, ...)
	reIn := regexp.MustCompile(`(?i)^([a-zA-Z_][\w]*)\s+IN\s*\((.*)\)$`)
	if m := reIn.FindStringSubmatch(part); len(m) == 3 {
		items := splitCSV(m[2])
		values := make([]any, 0, len(items))
		for _, item := range items {
			value, err := parseLiteral(item)
			if err != nil {
				return "", "", nil, err
			}
			values = append(values, value)
		}
		if len(values) == 0 {
			return "", "", nil, fmt.Errorf("invalid condition: %s", part)
		}
		return strings.TrimSpace(m[1]), kernal.QueryOpIn, values, nil
	}

	ops := []struct {
		token string
		op    string
	}{{">=", kernal.QueryOpGTE}, {"<=", kernal.QueryOpLTE}, {"!=", kernal.QueryOpNE}, {">", kernal.QueryOpGT}, {"<", kernal.QueryOpLT}, {"=", kernal.QueryOpEQ}}
	for _, item := range ops {
		idx := strings.Index(part, item.token)
		if idx <= 0 {
			continue
		}
		field := strings.TrimSpace(part[:idx])
		valueRaw := strings.TrimSpace(part[idx+len(item.token):])
		value, err := parseLiteral(valueRaw)
		if err != nil {
			return "", "", nil, err
		}
		return field, item.op, value, nil
	}
	return "", "", nil, fmt.Errorf("invalid condition: %s", part)
}

func parseColumnDef(def string) (kernal.Column, error) {
	parts := splitSpaceAware(def)
	if len(parts) < 2 {
		return kernal.Column{}, fmt.Errorf("invalid column def: %s", def)
	}
	col := kernal.Column{Name: parts[0], Type: parts[1]}
	for i := 2; i < len(parts); i++ {
		tok := strings.ToUpper(parts[i])
		switch tok {
		case "PRIMARY":
			if i+1 < len(parts) && strings.ToUpper(parts[i+1]) == "KEY" {
				col.PrimaryKey = true
				i++
			}
		case "AUTO_INCREMENT":
			col.AutoIncrement = true
		case "UNIQUE":
			col.Unique = true
		case "INDEXED", "INDEX":
			col.Indexed = true
		case "REQUIRED", "NOT":
			col.Required = true
		case "NULLABLE":
			v := true
			col.Nullable = &v
		case "DEFAULT":
			if i+1 >= len(parts) {
				return kernal.Column{}, fmt.Errorf("default value missing in %s", def)
			}
			v, err := parseLiteral(parts[i+1])
			if err != nil {
				return kernal.Column{}, err
			}
			col.Default = v
			i++
		}
	}
	return col, nil
}

func parseForeignKeyDef(def string) (kernal.ForeignKey, error) {
	re := regexp.MustCompile(`(?i)^FOREIGN\s+KEY\s*\((\w+)\)\s*REFERENCES\s+(\w+)\s*\((\w+)\)(?:\s+AS\s+(\w+))?(?:\s+NAME\s+(\w+))?$`)
	m := re.FindStringSubmatch(strings.TrimSpace(def))
	if len(m) == 0 {
		return kernal.ForeignKey{}, fmt.Errorf("invalid foreign key def: %s", def)
	}
	fk := kernal.ForeignKey{Field: m[1], RefTable: m[2], RefField: m[3]}
	if len(m) > 4 {
		fk.Alias = strings.TrimSpace(m[4])
	}
	if len(m) > 5 {
		fk.Name = strings.TrimSpace(m[5])
	}
	return fk, nil
}

func parseAlterAddForeignKey(op string) (kernal.ForeignKey, error) {
	tail := strings.TrimSpace(op[len("ADD "):])
	return parseForeignKeyDef(tail)
}

func parseFieldInParen(op string) (string, error) {
	start := strings.Index(op, "(")
	end := strings.LastIndex(op, ")")
	if start < 0 || end <= start {
		return "", fmt.Errorf("invalid operation syntax: %s", op)
	}
	return strings.TrimSpace(op[start+1 : end]), nil
}

func parseLiteral(raw string) (any, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if strings.EqualFold(raw, "null") {
		return nil, nil
	}
	if strings.EqualFold(raw, "true") {
		return true, nil
	}
	if strings.EqualFold(raw, "false") {
		return false, nil
	}
	if (strings.HasPrefix(raw, "'") && strings.HasSuffix(raw, "'")) || (strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`)) {
		return raw[1 : len(raw)-1], nil
	}
	if strings.Contains(raw, ".") {
		if f, err := strconv.ParseFloat(raw, 64); err == nil {
			return f, nil
		}
	}
	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return i, nil
	}
	return raw, nil
}

func splitByAND(input string) []string {
	items := []string{}
	curr := strings.Builder{}
	inQuote := false
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if ch == '\'' {
			inQuote = !inQuote
			curr.WriteByte(ch)
			continue
		}
		if !inQuote && i+3 <= len(input) && strings.EqualFold(input[i:i+3], "AND") {
			prevSpace := i == 0 || input[i-1] == ' '
			nextSpace := i+3 == len(input) || input[i+3] == ' '
			if prevSpace && nextSpace {
				items = append(items, strings.TrimSpace(curr.String()))
				curr.Reset()
				i += 2
				continue
			}
		}
		curr.WriteByte(ch)
	}
	if strings.TrimSpace(curr.String()) != "" {
		items = append(items, strings.TrimSpace(curr.String()))
	}
	return items
}

func splitCSV(input string) []string {
	parts := []string{}
	curr := strings.Builder{}
	level := 0
	inQuote := false
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if ch == '\'' {
			inQuote = !inQuote
			curr.WriteByte(ch)
			continue
		}
		if !inQuote {
			if ch == '(' {
				level++
			} else if ch == ')' {
				if level > 0 {
					level--
				}
			} else if ch == ',' && level == 0 {
				parts = append(parts, strings.TrimSpace(curr.String()))
				curr.Reset()
				continue
			}
		}
		curr.WriteByte(ch)
	}
	if strings.TrimSpace(curr.String()) != "" {
		parts = append(parts, strings.TrimSpace(curr.String()))
	}
	return parts
}

func splitSpaceAware(input string) []string {
	res := []string{}
	curr := strings.Builder{}
	inQuote := false
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if ch == '\'' {
			inQuote = !inQuote
			curr.WriteByte(ch)
			continue
		}
		if ch == ' ' && !inQuote {
			if curr.Len() > 0 {
				res = append(res, curr.String())
				curr.Reset()
			}
			continue
		}
		curr.WriteByte(ch)
	}
	if curr.Len() > 0 {
		res = append(res, curr.String())
	}
	return res
}
