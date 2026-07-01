package gormPools

import (
	"gorm.io/gorm"
)

type (
	GORMPool interface {
		GetConn() *gorm.DB
		getRws() *gorm.DB
		Close() error
	}
)
