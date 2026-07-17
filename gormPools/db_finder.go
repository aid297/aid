package gormPools

import (
	"reflect"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/aid297/aid/v2/operations"
	"github.com/aid297/aid/v2/str"
)

var _ Finder = (*FinderImpl)(nil)

type (
	Finder interface {
		GetDB() *gorm.DB
		Find(ret any, preloads ...string) Finder
		Ex(functions ...func(db *gorm.DB)) Finder
		TryPagination(page, size int) Finder
		TryOrder(orders ...string) Finder
		TryPreload(preloads ...string) Finder
		TryQuery(mode string, fieldName string, values ...any)
		SetTotal(total int64) Finder
		GetTotal() int64
		When(condition bool, query any, args ...any) Finder
		WhenIn(condition bool, query any, args any) Finder
		WhenInPtr(condition bool, query any, args any) Finder
		WhenNotIn(condition bool, query any, args any) Finder
		WhenNotInPtr(condition bool, query any, args any) Finder
		WhenBetween(condition bool, query any, args ...any) Finder
		WhenNotBetween(condition bool, query any, args ...any) Finder
		WhenLike(condition bool, query, arg any) Finder
		WhenLikeLeft(condition bool, query, arg any) Finder
		WhenLikeRight(condition bool, query, arg any) Finder
		WhenFunc(condition bool, fn func(db *gorm.DB)) Finder
		Transaction(functions ...func(db *gorm.DB)) error
		QueryUseMap(queries map[string][]any) Finder
		QueryUseCondition(finderCondition *FinderCondition) Finder
		FindUseMap(queries map[string][]any, preloads []string, orders []string, page, size int, ret any) Finder
		FindUseCondition(finderCondition *FinderCondition, page, size int, ret any) Finder
		FindOnlyCondition(finderCondition *FinderCondition, ret any) Finder
		finderNext() Finder
	}

	// FinderImpl 查询帮助器
	FinderImpl struct {
		db    *gorm.DB
		total int64
	}

	// FinderCondition 查询条件
	FinderCondition struct {
		Table    *string  `json:"table,omitempty"`
		Queries  []Query  `json:"queries,omitempty"`  // 查询条件
		Orders   []string `json:"orders,omitempty"`   // 排序
		Preloads []string `json:"preloads,omitempty"` // 预加载
		Page     int      `json:"page,omitempty"`     // 页码
		PageSize int      `json:"pageSize,omitempty"` // 页容量
	}

	// Condition 查询
	Condition struct {
		Key      string `json:"key"`      // SQL字段名称，如果有别名则需要带有别名
		Operator string `json:"operator"` // 操作符：=、>、<、!=、<=、>=、<>、in、not in、between、not between、like、like%、%like、raw、join
		Values   []any  `json:"values"`   // 查询条件值
	}

	// Query 查询
	Query struct {
		Option     *string     `json:"option,omitempty"`     // 操作：and、or、not
		Conditions []Condition `json:"conditions,omitempty"` // 条件
	}
)

func NewFinder(db *gorm.DB) Finder { return &FinderImpl{db: db, total: -1} }

// // New 实例化：查询帮助器
// func (*FinderImpl) New(db *gorm.DB) *FinderImpl { return &FinderImpl{db: db, total: -1} }

// GetDB 获取 gorm.DB 对象
func (my *FinderImpl) GetDB() *gorm.DB { return my.db }

// Find 查询数据
func (my *FinderImpl) Find(ret any, preloads ...string) Finder {
	my.TryPreload(preloads...)
	my.db.Find(ret)

	return my
}

// Ex 额外操作
func (my *FinderImpl) Ex(functions ...func(db *gorm.DB)) Finder {
	if len(functions) > 0 {
		for _, fn := range functions {
			fn(my.db)
		}
	}

	return my
}

// TryPagination 尝试分页
func (my *FinderImpl) TryPagination(page, size int) Finder {
	if page > 0 && size > 0 {
		if my.total == -1 {
			if my.db.Count(&my.total).Error != nil {
				return my
			}
		}

		my.db.Limit(size).Offset((page - 1) * size)
	}

	return my
}

// TryOrder 尝试排序
func (my *FinderImpl) TryOrder(orders ...string) Finder {
	for _, order := range orders {
		my.db.Order(order)
	}

	return my
}

// TryPreload 尝试深度查询
func (my *FinderImpl) TryPreload(preloads ...string) Finder {
	for _, preload := range preloads {
		my.db.Preload(preload)
	}

	return my
}

// TryQuery 尝试查询
func (my *FinderImpl) TryQuery(mode string, fieldName string, values ...any) {
	switch mode {
	case "and", "":
		my.db.Where(fieldName, values...)
	case "or":
		my.db.Or(fieldName, values...)
	case "not":
		my.db.Not(fieldName, values...)
	}
}

// SetTotal 设置总数
func (my *FinderImpl) SetTotal(total int64) Finder {
	my.total = total
	return my
}

// GetTotal 获取总数
func (my *FinderImpl) GetTotal() int64 { return my.total }

// When 当条件满足时执行：where
func (my *FinderImpl) When(condition bool, query any, args ...any) Finder {
	if condition {
		my.db.Where(query, args...)
	}

	return my
}

// WhenIn 当条件满足时执行：where in
func (my *FinderImpl) WhenIn(condition bool, query any, args any) Finder {
	if condition {
		ref := reflect.ValueOf(args)
		if ref.Kind() == reflect.Ptr && !ref.IsNil() {
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "in", "(?)"), ref.Elem().Interface())
			// my.db.Where(fmt.Sprintf("%v in (?)", query), ref.Elem().Interface())
		} else {
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "in", "(?)"), args)
			// my.db.Where(fmt.Sprintf("%v in (?)", query), args)
		}
	}

	return my
}

// WhenInPtr 当条件满足时执行：where in
// args 为指针类型
func (my *FinderImpl) WhenInPtr(condition bool, query any, args any) Finder {
	if condition && args != nil {
		my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "in", "(?)"), reflect.ValueOf(args).Elem().Interface())
		// my.db.Where(fmt.Sprintf("%v in (?)", query), reflect.ValueOf(args).Elem().Interface())
	}

	return my
}

// WhenNotIn 当条件满足时执行：where not in
func (my *FinderImpl) WhenNotIn(condition bool, query any, args any) Finder {
	if condition {
		ref := reflect.ValueOf(args)
		if ref.Kind() == reflect.Ptr && !ref.IsNil() {
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "not", "in", "(?)"), ref.Elem().Interface())
			// my.db.Where(fmt.Sprintf("%v not in (?)", query), ref.Elem().Interface())
		} else {
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "not", "in", "(?)"), args)
			// my.db.Where(fmt.Sprintf("%v not in (?)", query), args)
		}
	}

	return my
}

// WhenNotInPtr 当条件满足时执行：where in
// args 为指针类型
func (my *FinderImpl) WhenNotInPtr(condition bool, query any, args any) Finder {
	if condition && args != nil {
		my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "not", "in", "(?)"), reflect.ValueOf(args).Elem().Interface())
		// my.db.Where(fmt.Sprintf("%v not in (?)", query), reflect.ValueOf(args).Elem().Interface())
	}

	return my
}

// WhenBetween 当条件满足时执行：where between
func (my *FinderImpl) WhenBetween(condition bool, query any, args ...any) Finder {
	if condition {
		my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "between", "?", "and", "?"), args...)
		// my.db.Where(fmt.Sprintf("%v between ? and ?", query), args...)
	}

	return my
}

// WhenNotBetween 当条件满足时执行：where not between
func (my *FinderImpl) WhenNotBetween(condition bool, query any, args ...any) Finder {
	if condition {
		my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "not", "between", "?", "and", "?"), args...)
		// my.db.Where(fmt.Sprintf("%v not between ? and ?", query), args...)
	}

	return my
}

// WhenLike 当条件满足时执行：like %?%
func (my *FinderImpl) WhenLike(condition bool, query, arg any) Finder {
	if condition {
		my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(arg), "%"))
		// my.db.Where(fmt.Sprintf("%v like ?", query), fmt.Sprintf("%%%s%%", arg))
	}

	return my
}

// WhenLikeLeft 当条件满足时执行：like %?
func (my *FinderImpl) WhenLikeLeft(condition bool, query, arg any) Finder {
	if condition {
		my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(arg)))
		// my.db.Where(fmt.Sprintf("%v like ?", query), fmt.Sprintf("%%%s", arg))
	}

	return my
}

// WhenLikeRight 当条件满足时执行：like ?%
func (my *FinderImpl) WhenLikeRight(condition bool, query, arg any) Finder {
	if condition {
		my.db.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "like", "?"), str.APP.Buffer.JoinString(cast.ToString(arg)), "%")
		// my.db.Where(fmt.Sprintf("%v like ?", query), fmt.Sprintf("%s%%", arg))
	}

	return my
}

// WhenFunc 当条件满足时执行：通过回调执行
func (my *FinderImpl) WhenFunc(condition bool, fn func(db *gorm.DB)) Finder {
	if condition {
		fn(my.db)
	}

	return my
}

// Transaction 执行一组数据库事务操作
// 参数 funcs 为需要在事务中执行的函数切片,每个函数接收一个 *gorm.DB 参数
// 如果任一函数执行出错,将回滚整个事务并返回错误
// 所有函数执行成功后提交事务
// 返回 error,nil 表示事务执行成功,非 nil 表示事务执行失败
func (my *FinderImpl) Transaction(functions ...func(db *gorm.DB)) error {
	my.db.Begin()

	for _, fn := range functions {
		fn(my.db)
		if my.db.Error != nil {
			my.db.Rollback()
			return my.db.Error
		}
	}

	my.db.Commit()

	return nil
}

// QueryUseMap 从map中解析参数并查询
func (my *FinderImpl) QueryUseMap(queries map[string][]any) Finder {
	for key, value := range queries {
		var (
			ok       = false
			operator string
		)

		if operator, ok = value[0].(string); !ok {
			continue
		}

		switch operator {
		case "alias":
			// 表别名：{"tableName": ["alias", "aliasName"]}
			// tableAlias := fmt.Sprintf("%s as %s", key, value[1])
			my.db.Table(str.APP.Buffer.JoinStringLimit(" ", key, "as", cast.ToString(value[1])))
		case "=", ">", "<", "!=", "<=", ">=", "<>":
			// 常规比较操作：{"fieldName": ["=", "value"]}
			// my.db.Where(fmt.Sprintf("%s %s ?", key, operator), value[1])
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", key, operator, "?"), value[1])
		case "in", "not in":
			// 包含或不包含操作：{"fieldName": ["in", ["value1", "value2"]]}
			// my.db.Where(fmt.Sprintf("%s %s (?)", key, operator), value[1])
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", key, operator, "(?)"), value[1])
		case "between", "not between":
			// 范围查询：{"fieldName": ["between", ["value1", "value2"]]}
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", key, operator, "?", "and", "?"), value[1], value[2])
			// my.db.Where(fmt.Sprintf("%s %s ? and ?", key, operator), value[1], value[2])
		case "like":
			// 模糊查询：{"fieldName": ["like", "value"]}
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", key, "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(value[1]), "%"))
			// my.db.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%%%s%%", value[1]))
		case "like%":
			// 模糊查询：{"fieldName": ["like%", "value"]}
			// my.db.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%s%%", value[1]))
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", key, "like", "?"), str.APP.Buffer.JoinString(cast.ToString(value[1]), "%"))
		case "%like":
			// 模糊查询：{"fieldName": ["%like", "value"]}
			// my.db.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%%%s", value[1]))
			my.db.Where(str.APP.Buffer.JoinStringLimit(" ", key, "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(value[1])))
		case "join":
			// 连接查询：{"otherTableName": ["join", "joinSql", "onCondition"]}
			my.db.Joins(key, value[1:]...)
		case "raw":
			// 原生查询：{"fieldName": ["raw", "> ?", 100]}
			my.db.Where(key, value[1:]...)
		}
	}

	return my
}

// QueryUseCondition 从请求体中获取查询条件
func (my *FinderImpl) QueryUseCondition(finderCondition *FinderCondition) Finder {
	if finderCondition == nil {
		return my
	}

	// 设置表名
	if finderCondition.Table != nil && *finderCondition.Table != "" {
		my.db.Table(*finderCondition.Table)
	}

	// 设置查询条件
	if len(finderCondition.Queries) > 0 {
		for idx := range finderCondition.Queries {
			for _, condition := range finderCondition.Queries[idx].Conditions {
				if condition.Key != "" && len(condition.Values) > 0 {
					switch condition.Operator {
					case "=", ">", "<", "!=", "<=", ">=", "<>":
						// {key:"fieldName", operator:"=", values:["value"]}
						// my.TryQuery(*query.Option, fmt.Sprintf("%s %s ?", condition.Key, condition.Operator), condition.Values[0])
						my.TryQuery(*finderCondition.Queries[idx].Option, str.APP.Buffer.JoinStringLimit(" ", condition.Key, condition.Operator, "?"), condition.Values[0])
					case "in", "not in":
						// {key:"fieldName", operator:"in", values:["value1", "value2"]}
						// my.TryQuery(*query.Option, fmt.Sprintf("%s %s (?)", condition.Key, condition.Operator), condition.Values[0])
						my.TryQuery(*finderCondition.Queries[idx].Option, str.APP.Buffer.JoinStringLimit(" ", condition.Key, condition.Operator, "(?)"), condition.Values[0])
					case "between", "not between":
						// {key:"fieldName", operator:"between", values:["value1", "value2"]}
						// my.TryQuery(*query.Option, fmt.Sprintf("%s %s ? and ?", condition.Key, condition.Operator), condition.Values...)
						my.TryQuery(*finderCondition.Queries[idx].Option, str.APP.Buffer.JoinStringLimit(" ", condition.Key, condition.Operator, "?", "and", "?"), condition.Values...)
					case "like":
						// {key:"fieldName", operator:"like", values:["value"]}
						// my.TryQuery(*query.Option, fmt.Sprintf("%s like ?", condition.Key), fmt.Sprintf("%%%s%%", condition.Values[0]))
						my.TryQuery(*finderCondition.Queries[idx].Option, str.APP.Buffer.JoinStringLimit(" ", condition.Key, "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(condition.Values[0]), "%"))
					case "like%":
						// {key:"fieldName", operator:"like%", values:["value"]}
						// my.TryQuery(*query.Option, fmt.Sprintf("%s like ?", condition.Key), fmt.Sprintf("%s%%", condition.Values[0]))
						my.TryQuery(*finderCondition.Queries[idx].Option, str.APP.Buffer.JoinStringLimit(" ", condition.Key, "like", "?"), str.APP.Buffer.JoinString(cast.ToString(condition.Values[0]), "%"))
					case "%like":
						// {key:"fieldName", operator:"%like", values:["value"]}
						// my.TryQuery(*query.Option, fmt.Sprintf("%s like ?", condition.Key), fmt.Sprintf("%%%s", condition.Values[0]))
						my.TryQuery(*finderCondition.Queries[idx].Option, str.APP.Buffer.JoinStringLimit(" ", condition.Key, "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(condition.Values[0])))
					case "join", "left join", "right join", "inner join", "outer join":
						// {key:"otherName o", operator:"join", values["o.xxx = tableName.xxx where xxx = ? and yyy = ?","xxx-value","yyy-value"]}
						// my.TryQuery(*query.Option, fmt.Sprintf("%s %s", condition.Operator, condition.Key), condition.Values...)
						my.TryQuery(*finderCondition.Queries[idx].Option, str.APP.Buffer.JoinStringLimit(" ", condition.Operator, condition.Key), condition.Values...)
					case "raw":
						// {key:"fieldName", operator:"raw", values:["> ?", 100]}
						my.TryQuery(*finderCondition.Queries[idx].Option, condition.Key, condition.Values...)
					}
				}
			}
		}
	}

	// 设置排序
	if len(finderCondition.Orders) > 0 {
		my.TryOrder(finderCondition.Orders...)
	}

	// 设置预加载
	if len(finderCondition.Preloads) > 0 {
		my.TryPreload(finderCondition.Preloads...)
	}

	return my
}

// FindUseMap 自动填充查询条件并查询：使用map[string][]any
func (my *FinderImpl) FindUseMap(queries map[string][]any, preloads []string, orders []string, page, size int, ret any) Finder {
	return my.QueryUseMap(queries).TryPagination(page, size).TryOrder(orders...).Find(ret, preloads...).finderNext()
}

// FindUseCondition 自动填充查询条件并查询：使用FinderCondition
func (my *FinderImpl) FindUseCondition(finderCondition *FinderCondition, page, size int, ret any) Finder {
	if finderCondition != nil && finderCondition.Page > 0 && finderCondition.PageSize > 0 {
		page = finderCondition.Page
		size = finderCondition.PageSize
	}
	return my.QueryUseCondition(finderCondition).TryPagination(page, size).Find(ret).finderNext()
}

// FindOnlyCondition  自动填充查询条件并查询：使用FinderCondition
func (my *FinderImpl) FindOnlyCondition(finderCondition *FinderCondition, ret any) Finder {
	if finderCondition == nil {
		my.db.Find(ret)
		return my.finderNext()
	}

	return my.QueryUseCondition(finderCondition).TryPagination(finderCondition.Page, finderCondition.PageSize).Find(ret).finderNext()
}

func (my *FinderImpl) finderNext() Finder {
	my.total = operations.NewTernary(operations.TrueValue[int64](0), operations.FalseValue(my.total)).GetByValue(my.total == -1)
	return my
}
