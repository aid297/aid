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
//   - WhenFunc 回调参数为 *gorm.Statement（泛型 Scopes 签名）
//   - 错误通过底层 *gorm.DB.AddError 传递，finisher 自动返回
type FinderG[T any] struct {
	chain          gorm.ChainInterface[T]
	total          int64
	needPagination bool
}

// ToFinderG 绑定泛型模型并应用属性，返回 *FinderG[T]
func ToFinderG[T Modeler](db *gorm.DB, attrs ...GenericAttr[T]) *FinderG[T] {
	return &FinderG[T]{chain: ToModelG[T](db, attrs...)}
}

// GetChain 获取底层 gorm.ChainInterface[T]，用于需要直接操作链式 API 的场景
func (my *FinderG[T]) GetChain() gorm.ChainInterface[T] { return my.chain }

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
