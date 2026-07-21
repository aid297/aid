package gormPools

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/aid297/aid/v2/str"
)

type (
	// Modeler 接口：模型
	Modeler interface {
		TableName() string
	}

	ModelAttr interface {
		Register(model Modeler, db *gorm.DB) *gorm.DB
	}

	AttrTable struct{ table string }
	AttrJoins struct {
		join string
		args []any
	}
	AttrPreload struct {
		preload string
		args    []any
	}
	AttrPreloadMap struct{ preloadMap map[string][]any }
	AttrSelect     struct {
		query any
		args  []any
	}
	AttrDistinct struct{ args []any }
	AttrScopes   struct{ scopes []func(db *gorm.DB) *gorm.DB }
)

func ToModel[model Modeler](db *gorm.DB, attrs ...ModelAttr) *gorm.DB {
	ins := new(model)
	db = db.Model(ins)

	for _, attr := range attrs {
		db = attr.Register(*ins, db)
	}

	return db
}

func ToFinder[model Modeler](db *gorm.DB, attrs ...ModelAttr) Finder {
	return NewFinder(ToModel[model](db, attrs...))
}

func Exist(tx *gorm.DB, query any, args ...any) (exist bool, err error) {
	err = tx.Where(query, args...).Limit(1).Find(nil).Error
	return tx.RowsAffected != 0, fmt.Errorf("查询数据失败：%v", err)
}

func SaveOrCreate[T Modeler](tx *gorm.DB, create, save T, query any, args ...any) (T, error) {
	var ret T

	exist, err := Exist(tx, query, args...)
	if err != nil {
		return ret, err
	}

	if exist {
		if err = tx.Where(query, args...).Save(&save).Error; err != nil {
			return ret, fmt.Errorf("保存数据失败: %w", err)
		}

		return save, nil
	} else {
		if err = tx.Create(&create).Error; err != nil {
			return ret, fmt.Errorf("新建数据失败: %w", err)
		}

		return create, nil
	}
}

func UpdatesOrCreate[T Modeler](tx *gorm.DB, create T, update, query any, args ...any) (T, error) {
	var ret T

	exist, err := Exist(tx, query, args...)
	if err != nil {
		return ret, err
	}

	if exist {
		if err = tx.Where(query, args...).Updates(update).Error; err != nil {
			return ret, fmt.Errorf("更新数据失败: %w", err)
		}
	} else {
		if err = tx.Create(&create).Error; err != nil {
			return ret, fmt.Errorf("新建数据失败: %w", err)
		}
	}

	if err = tx.Where(query, args...).First(&ret).Error; err != nil {
		return ret, fmt.Errorf("同步数据失败：%v", err)
	}

	return ret, nil
}

func UpdateColumnsOrCreate[T Modeler](tx *gorm.DB, create T, update, query any, args ...any) (T, error) {
	var ret T

	exist, err := Exist(tx, query, args...)
	if err != nil {
		return ret, err
	}

	if exist {
		if err = tx.Where(query, args...).UpdateColumns(update).Error; err != nil {
			return ret, fmt.Errorf("更新数据失败: %w", err)
		}
	} else {
		if err = tx.Create(&create).Error; err != nil {
			return ret, fmt.Errorf("新建数据失败: %w", err)
		}
	}

	if err = tx.Where(query, args...).First(&ret).Error; err != nil {
		return ret, fmt.Errorf("同步数据失败：%v", err)
	}

	return ret, nil
}

func Table(table string) AttrTable { return AttrTable{table: table} }

func (my AttrTable) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Table(str.APP.Buffer.JoinStringLimit(" ", model.TableName(), "as", my.table))
}

func Joins(join string, args ...any) AttrJoins { return AttrJoins{join: join, args: args} }

func (my AttrJoins) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Joins(my.join, my.args...)
}

func Preload(preload string, args ...any) AttrPreload {
	return AttrPreload{preload: preload, args: args}
}

func (my AttrPreload) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Preload(my.preload, my.args...)
}

func PreloadMap(preloadMap map[string][]any) AttrPreloadMap {
	return AttrPreloadMap{preloadMap: preloadMap}
}

func (my AttrPreloadMap) Register(model Modeler, db *gorm.DB) *gorm.DB {
	for key := range my.preloadMap {
		db = db.Preload(key, my.preloadMap[key]...)
	}
	return db
}

func Select(query any, args ...any) AttrSelect { return AttrSelect{query: query, args: args} }

func (my AttrSelect) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Select(my.query, my.args...)
}

func Distinct(args ...any) AttrDistinct { return AttrDistinct{args: args} }

func (my AttrDistinct) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Distinct(my.args...)
}

func Scopes(scopes ...func(db *gorm.DB) *gorm.DB) AttrScopes { return AttrScopes{scopes: scopes} }

func (my AttrScopes) Register(model Modeler, db *gorm.DB) *gorm.DB {
	for _, scope := range my.scopes {
		db = scope(db)
	}
	return db
}
