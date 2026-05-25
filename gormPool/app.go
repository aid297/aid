package gormPool

import (
	"time"

	"gorm.io/gorm"
)

var APP struct {
	MySQLPool     MySQLPool
	PGPool        PGPool
	SQLServerPool SQLServerPool
	DBSetting     DBSetting
	Finder        Finder
}

type MySQLBasic struct {
	ID        int64          `gorm:"column:id;type:bigint unsigned;primaryKey"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;not null;default:current_timestamp"`
}
