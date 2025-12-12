package rbac

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aid297/aid/str"
)

type Permission struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deletedAt"`
	UUID      string         `gorm:"column:uuid;type:char(36);not null;unique;comment:uuid;" json:"uuid"`
	Name      string         `gorm:"type:varchar(50);not null;uniqueIndex;comment:权限名称" json:"name"`
	Identity  string         `gorm:"type:varchar(255);not null;uniqueIndex;comment:权限标识" json:"identity"`
	Faces     []*Group       `gorm:"foreignKey:dot_uuid;references:uuid;" json:"dot"`
}

func (Permission) New(attrs ...PermissionAttributer) Permission { return Permission{}.Set(attrs...) }

func (my Permission) Set(attrs ...PermissionAttributer) Permission {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}

	return my
}

func (Permission) TableName() string {
	if edgeIns.tablePrefix != "" {
		return str.APP.Buffer.JoinString(edgeIns.tablePrefix, "_permissions")
	} else {
		return "edge_permissions"
	}
}

func (Permission) DB() *gorm.DB { return edgeIns.db.Model(new(Permission)) }

func (my Permission) Store(attrs ...PermissionAttributer) error {
	var (
		err    error
		repeat int64
	)

	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	if err = my.Set(attrs...).DB().Where("name", my.Name).Count(&repeat).Error; err != nil {
		return fmt.Errorf("%w（名称）：%w", ErrCheckRepeat, err)
	}

	if repeat > 0 {
		return fmt.Errorf("权限%w", ErrRepeatName)
	}

	if err = my.DB().Where("identity", my.Identity).Count(&repeat).Error; err != nil {
		return fmt.Errorf("%w（标识）：%w", ErrCheckRepeat, err)
	}

	if repeat > 0 {
		return ErrRepeatIdentity
	}

	my.UUID = uuid.Must(uuid.NewV6()).String()
	my.CreatedAt = time.Now()
	my.UpdatedAt = time.Now()
	return my.DB().Create(my).Error
}

func (my Permission) Update(attrs ...PermissionAttributer) error {
	var (
		err    error
		repeat int64
	)

	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	if err = my.Set(attrs...).DB().Where("uuid <> ?", my.UUID).Where("name", my.Name).Count(&repeat).Error; err != nil {
		return fmt.Errorf("%w（名称）：%w", ErrCheckRepeat, err)
	}

	if repeat > 0 {
		return ErrRepeatName
	}

	if err = my.DB().Where("uuid <> ?", my.UUID).Where("identity", my.Identity).Count(&repeat).Error; err != nil {
		return fmt.Errorf("%w（标识）：%w", ErrCheckRepeat, err)
	}

	if repeat > 0 {
		return ErrRepeatIdentity
	}

	my.UpdatedAt = time.Now()
	if err = my.DB().Where("uuid", my.UUID).Save(my).Error; err != nil {
		return fmt.Errorf("%w：%w", ErrUpdate, err)
	}

	return nil
}

func (my Permission) Destroy() error {
	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	return my.DB().Where("uuid", my.UUID).Delete(my).Error
}
