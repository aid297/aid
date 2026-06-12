package gormPool

import (
	"time"

	"gorm.io/gorm"
)

var APP struct {
	DBSetting DBSetting
}

type MySQLBasic struct {
	ID        int64          `gorm:"column:id;type:bigint unsigned;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:current_timestamp;autoUpdateTime:true" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;not null" json:"deleted_at"`
	UUID      string         `gorm:"column:uuid;type:char(36);not null;unique" json:"uuid"`
}
