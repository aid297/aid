package gormPool

import (
	"gorm.io/gorm"

	"github.com/aid297/aid/str"
)

type (
	// Modeler 接口：模型
	Modeler interface{ TableName() string }

	ModelAttributer interface {
		Register(model Modeler, db *gorm.DB) *gorm.DB
	}

	AttrTable struct{ table string }
	AttrJoins struct {
		join string
		args []any
	}
	AttrPreload struct{ preloads []string }
	AttrSelect  struct {
		query any
		args  []any
	}
	AttrDistinct struct{ args []any }
)

func DefaultModel[model Modeler](db *gorm.DB, attrs ...ModelAttributer) *gorm.DB {
	ins := new(model)
	db = db.Model(ins)

	for _, attr := range attrs {
		db = attr.Register(*ins, db)
	}

	return db
}

func DefaultFinder[model Modeler](db *gorm.DB, attrs ...ModelAttributer) *Finder {
	return FinderApp.New(DefaultModel[model](db, attrs...))
}

func Table(table string) *AttrTable { return &AttrTable{table: table} }

func (my *AttrTable) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Table(str.APP.Buffer.JoinStringLimit(" ", model.TableName(), "as", my.table))
}

func Joins(join string, args ...any) *AttrJoins { return &AttrJoins{join: join, args: args} }

func (my *AttrJoins) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Joins(my.join, my.args...)
}

func Preload(preloads ...string) *AttrPreload { return &AttrPreload{preloads: preloads} }

func (my *AttrPreload) Register(model Modeler, db *gorm.DB) *gorm.DB {
	for _, preload := range my.preloads {
		db = db.Preload(preload)
	}
	return db
}

func Select(query any, args ...any) *AttrSelect { return &AttrSelect{query: query, args: args} }

func (my *AttrSelect) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Select(my.query, my.args...)
}

func Distinct(args ...any) *AttrDistinct { return &AttrDistinct{args: args} }

func (my *AttrDistinct) Register(model Modeler, db *gorm.DB) *gorm.DB {
	return db.Distinct(my.args...)
}
