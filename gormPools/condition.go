package gormPools

import (
	_fmt "fmt"
	_sync "sync"

	_gormSchema "gorm.io/gorm/schema"
)

var conditionSchemaCache = _sync.Map{}

// WithAllowFields 依据模型列名填充 FinderCondition.AllowFields 白名单
// condition 为 nil 时返回 nil；同时放行 "<表名>.<列名>" 形式以兼容带别名的查询条件
func WithAllowFields(condition *FinderCondition, model any, exFields ...string) *FinderCondition {
	return WithAllowFieldsAlias(condition, model, "", exFields...)
}

// WithAllowFieldsAlias 同 WithAllowFields，额外按指定别名放行 "<别名>.<列名>" 形式，
// 用于查询使用 Table("table as alias") 的场景
func WithAllowFieldsAlias(condition *FinderCondition, model any, alias string, exFields ...string) *FinderCondition {
	if condition == nil {
		return nil
	}

	var fields []string
	if parsed, err := parseConditionSchema(model); err == nil {
		table := parsed.Table
		for _, column := range parsed.DBNames {
			fields = append(fields, column, table+"."+column)
			if alias != "" {
				fields = append(fields, alias+"."+column)
			}
		}
	}

	condition.AllowFields = append(fields, exFields...)

	return condition
}

func parseConditionSchema(model any) (*_gormSchema.Schema, error) {
	key := _fmt.Sprintf("%T", model)
	if cached, ok := conditionSchemaCache.Load(key); ok {
		return cached.(*_gormSchema.Schema), nil
	}

	parsed, err := _gormSchema.Parse(model, &conditionSchemaCache, _gormSchema.NamingStrategy{})
	if err != nil {
		return nil, err
	}

	conditionSchemaCache.Store(key, parsed)

	return parsed, nil
}
