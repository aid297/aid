package sqlapi

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
	reInsert   = regexp.MustCompile(`(?i)^INSERT\s+INTO\s+([a-zA-Z_][\w]*)\s*\((.*)\)\s*VALUES\s*\((.*)\)$`)
	reUpdate   = regexp.MustCompile(`(?i)^UPDATE\s+([a-zA-Z_][\w]*)\s+SET\s+(.+?)(?:\s+WHERE\s+(.+))?$`)
	reDelete   = regexp.MustCompile(`(?i)^DELETE\s+FROM\s+([a-zA-Z_][\w]*)(?:\s+WHERE\s+(.+))?$`)
	reSelect   = regexp.MustCompile(`(?i)^SELECT\s+\*\s+FROM\s+([a-zA-Z_][\w]*)(?:\s+WHERE\s+(.+))?$`)
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
	if m := reSelect.FindStringSubmatch(sql); len(m) == 3 {
		return parseSelect(m[1], m[2])
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
	values := splitCSV(valuesRaw)
	if len(fields) != len(values) {
		return nil, fmt.Errorf("insert fields/values length mismatch")
	}
	row := kernal.Row{}
	for i := range fields {
		v, err := parseLiteral(values[i])
		if err != nil {
			return nil, err
		}
		row[strings.TrimSpace(fields[i])] = v
	}
	return InsertStmt{Table: table, Row: row}, nil
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
	conditions, err := parseConditions(whereRaw)
	if err != nil {
		return nil, err
	}
	if len(conditions) != 1 || conditions[0].Operator != kernal.QueryOpEQ {
		return nil, fmt.Errorf("update requires exactly one '=' where condition on primary key")
	}

	return UpdateStmt{Table: table, PrimaryKey: conditions[0].Value, Updates: updates}, nil
}

func parseDelete(table, whereRaw string) (Statement, error) {
	conditions, err := parseConditions(whereRaw)
	if err != nil {
		return nil, err
	}
	return DeleteStmt{Table: table, Conditions: conditions}, nil
}

func parseSelect(table, whereRaw string) (Statement, error) {
	conditions, err := parseConditions(whereRaw)
	if err != nil {
		return nil, err
	}
	return SelectStmt{Table: table, Conditions: conditions}, nil
}

func parseConditions(whereRaw string) ([]kernal.QueryCondition, error) {
	whereRaw = strings.TrimSpace(whereRaw)
	if whereRaw == "" {
		return nil, nil
	}
	parts := splitByAND(whereRaw)
	conditions := make([]kernal.QueryCondition, 0, len(parts))
	for _, part := range parts {
		field, op, value, err := parseSingleCondition(part)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, kernal.QueryCondition{Field: field, Operator: op, Value: value})
	}
	return conditions, nil
}

func parseSingleCondition(part string) (string, string, any, error) {
	part = strings.TrimSpace(part)
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
