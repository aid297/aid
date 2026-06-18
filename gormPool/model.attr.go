package gormPool

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/aid297/aid/v2/str"
)

type (
	// Modeler 接口：模型
	Modeler interface{ TableName() string }

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

func SaveOrCreate[model Modeler](tx *gorm.DB, data Modeler, query any, args ...any) (err error) {
	var count int64

	if err = tx.Where(query, args...).Limit(1).Count(&count).Error; err != nil {
		return fmt.Errorf("查询待更新数据存在失败: %w", err)
	}

	if count != 0 {
		if err = tx.Where(query, args...).Save(data).Error; err != nil {
			return fmt.Errorf("更新数据失败: %w", err)
		}
	} else {
		if err = tx.Table(data.TableName()).Create(data).Error; err != nil {
			return fmt.Errorf("创建数据失败: %w", err)
		}
	}

	return
}

func UpdatesOrCreate[model Modeler](tx *gorm.DB, data Modeler, updates any, query any, args ...any) (err error) {
	var count int64

	if err = tx.Where(query, args...).Limit(1).Count(&count).Error; err != nil {
		return fmt.Errorf("查询待更新数据存在失败: %w", err)
	}

	if count != 0 {
		if err = tx.Where(query, args...).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新数据失败: %w", err)
		}
	} else {
		if err = tx.Table(data.TableName()).Create(data).Error; err != nil {
			return fmt.Errorf("创建数据失败: %w", err)
		}
	}

	return
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
