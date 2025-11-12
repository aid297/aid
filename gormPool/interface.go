package gormPool

import "gorm.io/gorm"

type (
	GormPool interface {
		GetConn() *gorm.DB
		getRws() *gorm.DB
		Close() error
	}
)
