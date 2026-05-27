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
	ID        int64          `gorm:"column:id;type:bigint unsigned;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:current_timestamp" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;not null" json:"deleted_at"`
	UUID      string         `gorm:"column:uuid;type:char(36);not null;unique" json:"uuid"`
}
