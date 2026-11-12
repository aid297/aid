package gormPools

import (
	"context"
	"reflect"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/aid297/aid/v2/str"
)

// FinderG 泛型查询帮助器，包装 gorm.ChainInterface[T]
//
// 与传统 Finder 的差异：
//   - finisher 方法（Find/First/Count）直接返回 (结果, error)，而非通过 Error() 获取
//   - WhenFunc / Ex 回调参数为 *gorm.Statement（泛型 Scopes 签名，Statement 嵌入 *DB，可直接调用 DB 方法）
//   - 校验类错误（如 TryOrder）通过底层 *gorm.DB.AddError 传递，finisher 自动返回
type FinderG[T any] struct {
	db             *gorm.DB
	chain          gorm.ChainInterface[T]
	total          int64
	err            error
	needPagination bool
}

// ToFinderG 绑定泛型模型并应用属性，返回 *FinderG[T]
func ToFinderG[T Modeler](db *gorm.DB, attrs ...GenericAttr[T]) *FinderG[T] {
	return &FinderG[T]{db: db, chain: ToModelG[T](db, attrs...)}
}

// GetChain 获取底层 gorm.ChainInterface[T]，用于需要直接操作链式 API 的场景
func (my *FinderG[T]) GetChain() gorm.ChainInterface[T] { return my.chain }

// Ex 额外操作，回调参数为 *gorm.Statement（嵌入 *gorm.DB，可直接调用 DB 方法）
//
// 与传统 Finder.Ex(func(db *gorm.DB)) 对应；泛型 API 中通过 Scopes 注入，
// 回调内可用 stmt.DB 访问底层 *gorm.DB，或直接使用 stmt 的 Statement 级方法。
func (my *FinderG[T]) Ex(functions ...func(db *gorm.Statement)) *FinderG[T] {
	if len(functions) > 0 {
		my.chain = my.chain.Scopes(functions...)
	}
	return my
}

// TryQuery 尝试查询
//
// mode 取值："and"/""（默认 where）、"or"、"not"
func (my *FinderG[T]) TryQuery(mode string, fieldName string, values ...any) *FinderG[T] {
	switch mode {
	case "and", "":
		my.chain = my.chain.Where(fieldName, values...)
	case "or":
		my.chain = my.chain.Or(fieldName, values...)
	case "not":
		my.chain = my.chain.Not(fieldName, values...)
	}
	return my
}

// Error 返回客户端校验类错误（如 TryOrder 校验失败）
//
// 注意：SQL 执行类错误由 Find/First/Count 等 finisher 方法直接返回。
func (my *FinderG[T]) Error() error { return my.err }

// Where 添加查询条件
func (my *FinderG[T]) Where(query any, args ...any) *FinderG[T] {
	my.chain = my.chain.Where(query, args...)
	return my
}

// Not 添加 not 条件
func (my *FinderG[T]) Not(query any, args ...any) *FinderG[T] {
	my.chain = my.chain.Not(query, args...)
	return my
}

// Or 添加 or 条件
func (my *FinderG[T]) Or(query any, args ...any) *FinderG[T] {
	my.chain = my.chain.Or(query, args...)
	return my
}

// When 当条件满足时执行 where
func (my *FinderG[T]) When(condition bool, query any, args ...any) *FinderG[T] {
	if condition {
		my.chain = my.chain.Where(query, args...)
	}
	return my
}

// WhenIn 当条件满足时执行 where in
//
// args 支持切片或指向切片的指针（与传统 Finder.WhenIn 行为一致）
func (my *FinderG[T]) WhenIn(condition bool, query any, args any) *FinderG[T] {
	if condition {
		ref := reflect.ValueOf(args)
		if ref.Kind() == reflect.Pointer && !ref.IsNil() {
			my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "in", "(?)"), ref.Elem().Interface())
		} else {
			my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "in", "(?)"), args)
		}
	}
	return my
}

// WhenNotIn 当条件满足时执行 where not in
func (my *FinderG[T]) WhenNotIn(condition bool, query any, args any) *FinderG[T] {
	if condition {
		ref := reflect.ValueOf(args)
		if ref.Kind() == reflect.Pointer && !ref.IsNil() {
			my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "not", "in", "(?)"), ref.Elem().Interface())
		} else {
			my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "not", "in", "(?)"), args)
		}
	}
	return my
}

// WhenBetween 当条件满足时执行 where between
func (my *FinderG[T]) WhenBetween(condition bool, query any, args ...any) *FinderG[T] {
	if condition {
		my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "between", "?", "and", "?"), args...)
	}
	return my
}

// WhenNotBetween 当条件满足时执行 where not between
func (my *FinderG[T]) WhenNotBetween(condition bool, query any, args ...any) *FinderG[T] {
	if condition {
		my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "not", "between", "?", "and", "?"), args...)
	}
	return my
}

// WhenLike 当条件满足时执行 like %?%
func (my *FinderG[T]) WhenLike(condition bool, query, arg any) *FinderG[T] {
	if condition {
		my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(arg), "%"))
	}
	return my
}

// WhenLikeLeft 当条件满足时执行 like %?
func (my *FinderG[T]) WhenLikeLeft(condition bool, query, arg any) *FinderG[T] {
	if condition {
		my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "like", "?"), str.APP.Buffer.JoinString("%", cast.ToString(arg)))
	}
	return my
}

// WhenLikeRight 当条件满足时执行 like ?%
func (my *FinderG[T]) WhenLikeRight(condition bool, query, arg any) *FinderG[T] {
	if condition {
		my.chain = my.chain.Where(str.APP.Buffer.JoinStringLimit(" ", cast.ToString(query), "like", "?"), str.APP.Buffer.JoinString(cast.ToString(arg), "%"))
	}
	return my
}

// WhenFunc 当条件满足时通过 Scopes 执行回调
//
// 回调参数为 *gorm.Statement（泛型 Scopes 签名），与传统 Finder.WhenFunc 的 *gorm.DB 不同。
func (my *FinderG[T]) WhenFunc(condition bool, fn func(db *gorm.Statement)) *FinderG[T] {
	if condition {
		my.chain = my.chain.Scopes(fn)
	}
	return my
}

// TryPagination 尝试分页
//
// 当 page > 0 且 size > 0 时生效；若此前已设置 needPagination，会先执行 Count 获取总数。
func (my *FinderG[T]) TryPagination(page, size int) *FinderG[T] {
	if page > 0 && size > 0 {
		if my.needPagination {
			total, err := my.chain.Count(context.Background(), "*")
			if err != nil {
				return my
			}
			my.total = total
		}
		my.chain = my.chain.Limit(size).Offset((page - 1) * size)
	}
	return my
}

// TryOrder 尝试排序
//
// 校验排序格式（复用 validateOrder），非法时通过 Statement.AddError 注入错误，
// 后续 finisher 方法会返回该错误。
func (my *FinderG[T]) TryOrder(orders ...string) *FinderG[T] {
	for _, order := range orders {
		if err := validateOrder(order, nil); err != nil {
			my.err = err
			my.chain = my.chain.Scopes(func(stmt *gorm.Statement) {
				stmt.AddError(err)
			})
			return my
		}
		my.chain = my.chain.Order(order)
	}
	return my
}

// TryPreload 尝试预加载（无额外条件）
func (my *FinderG[T]) TryPreload(preloads ...string) *FinderG[T] {
	for _, preload := range preloads {
		my.chain = my.chain.Preload(preload, nil)
	}
	return my
}

// SetTotal 设置总数
func (my *FinderG[T]) SetTotal(total int64) *FinderG[T] {
	my.total = total
	return my
}

// GetTotal 获取总数
func (my *FinderG[T]) GetTotal() int64 { return my.total }

// Transaction 执行一组数据库事务操作
//
// 与传统 Finder.Transaction 对应，使用 gorm.DB.Transaction 自动管理 begin/commit/rollback；
// 任一函数执行出错则回滚并返回错误。注意：回调内须使用传入的 tx 执行操作方可纳入事务。
func (my *FinderG[T]) Transaction(functions ...func(db *gorm.DB)) error {
	return my.db.Transaction(func(tx *gorm.DB) error {
		for _, fn := range functions {
			fn(tx)
			if tx.Error != nil {
				return tx.Error
			}
		}
		return nil
	})
}

// Find 查询数据，返回 []T
func (my *FinderG[T]) Find(ctx context.Context) ([]T, error) {
	return my.chain.Find(ctx)
}

// First 查询首条，返回 T
func (my *FinderG[T]) First(ctx context.Context) (T, error) {
	return my.chain.First(ctx)
}

// Count 统计数量
func (my *FinderG[T]) Count(ctx context.Context) (int64, error) {
	return my.chain.Count(ctx, "*")
}
