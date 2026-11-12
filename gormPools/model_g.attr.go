package gormPools

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aid297/aid/v2/str"
)

// GenericAttr 泛型模型属性：作用于 gorm.ChainInterface[T]
//
// 与传统 ModelAttr 的差异：
//   - Joins/Preload/Scopes 因 gorm 泛型 API 签名变化，参数类型与传统版不同
//   - Table 通过 ChainInterface.Table 直接设置（v1.31.0+ ChainInterface 已包含 Table）
type GenericAttr[T any] func(db gorm.ChainInterface[T]) gorm.ChainInterface[T]

// ToModelG 绑定泛型模型并应用属性，返回 gorm.ChainInterface[T]
//
// 返回 ChainInterface[T]：覆盖 Where/Order/Limit/Find/First/Count 等查询能力。
// 如需 Create，请直接使用 gorm.G[T](db).Create(ctx, &obj)。
//
// 示例：
//
//	chain := ToModelG[TestTable1](db, GTable[TestTable1]("t"), GSelect[TestTable1]("id,name"))
//	list, err := chain.Where("id = ?", 1).Find(ctx) // list 类型为 []TestTable1
func ToModelG[T Modeler](db *gorm.DB, attrs ...GenericAttr[T]) gorm.ChainInterface[T] {
	// gorm.G[T](db) 返回 Interface[T]，其 Select/Omit/Table 返回 CreateInterface[T]，
	// 不直接满足 ChainInterface[T]；通过 Scopes() 转换为 ChainInterface[T]（无副作用）。
	chain := gorm.G[T](db).Scopes()
	for _, attr := range attrs {
		chain = attr(chain)
	}
	return chain
}

// GTable 设置表别名，等价于传统 Table 属性
func GTable[T Modeler](alias string) GenericAttr[T] {
	var m T
	tableName := m.TableName()
	expr := str.APP.Buffer.JoinStringLimit(" ", tableName, "as", alias)
	return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
		return db.Table(expr)
	}
}

// GSelect 查询字段，等价于传统 Select 属性
func GSelect[T any](query string, args ...any) GenericAttr[T] {
	return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
		return db.Select(query, args...)
	}
}

// GDistinct 去重，等价于传统 Distinct 属性
func GDistinct[T any](args ...any) GenericAttr[T] {
	return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
		return db.Distinct(args...)
	}
}

// GJoins 连接查询
//
// jt 由 clause.Has(name) / clause.LeftJoin.Association(name) 等构造；
// on 回调通过 JoinBuilder 构造连接条件。
func GJoins[T any](jt clause.JoinTarget, on func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error) GenericAttr[T] {
	return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
		return db.Joins(jt, on)
	}
}

// GPreload 预加载
//
// query 为 nil 时表示无额外条件。
func GPreload[T any](association string, query func(db gorm.PreloadBuilder) error) GenericAttr[T] {
	return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
		return db.Preload(association, query)
	}
}

// GPreloadMap 批量预加载，等价于传统 PreloadMap 属性
func GPreloadMap[T any](preloadMap map[string]func(db gorm.PreloadBuilder) error) GenericAttr[T] {
	return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
		for assoc, fn := range preloadMap {
			db = db.Preload(assoc, fn)
		}
		return db
	}
}

// GScopes 作用域，等价于传统 Scopes 属性
//
// 回调参数为 *gorm.Statement（泛型 API 中 Scopes 签名与传统不同）。
func GScopes[T any](scopes ...func(db *gorm.Statement)) GenericAttr[T] {
	return func(db gorm.ChainInterface[T]) gorm.ChainInterface[T] {
		return db.Scopes(scopes...)
	}
}
