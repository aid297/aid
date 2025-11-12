package rbac

import (
	"errors"
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
	UUID      uuid.UUID      `gorm:"type:char(36);not null;unique;comment:uuid;" json:"uuid"`
	Name      string         `gorm:"type:varchar(50);not null;uniqueIndex;comment:权限名称" json:"name"`
	Identity  string         `gorm:"type:varchar(255);not null;uniqueIndex;comment:权限标识" json:"identity"`
	Faces     []*Group       `gorm:"foreignKey:dot_uuid;references:uuid;" json:"dot"`
}

func (Permission) New(attrs ...PermissionAttributer) Permission {
	ins := new(Permission)
	return ins.Set(attrs...)
}

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
		return str.BufferApp.NewString(edgeIns.tablePrefix).S("_permissions").String()
	} else {
		return str.BufferApp.NewString("triangle_permissions").String()
	}
}

func (Permission) DB() *gorm.DB { return edgeIns.db.Model(new(Permission)) }

func (my Permission) Store(options ...PermissionAttributer) error {
	var (
		err    error
		repeat int64
	)

	if edgeIns.db == nil {
		return errors.New("数据库链接失败")
	}

	if len(options) > 0 {
		my.Set(options...)
	}

	if err = my.DB().Where("name", my.Name).Count(&repeat).Error; err != nil {
		return fmt.Errorf("查重（名称）失败：%w", err)
	}

	if repeat > 0 {
		return errors.New("名称重复")
	}

	if err = my.DB().Where("identity", my.Identity).Count(&repeat).Error; err != nil {
		return fmt.Errorf("查重（标识）失败：%w", err)
	}

	if repeat > 0 {
		return errors.New("标识重复")
	}

	my.UUID = uuid.Must(uuid.NewV6())
	my.CreatedAt = time.Now()
	my.UpdatedAt = time.Now()
	return my.DB().Create(my).Error
}

func (my Permission) Update(options ...PermissionAttributer) error {
	var (
		err    error
		repeat int64
	)

	if edgeIns.db == nil {
		return errors.New("数据库连接失败")
	}

	if len(options) > 0 {
		my.Set(options...)
	}

	if err = my.DB().Where("uuid <> ?", my.UUID.String()).Where("name", my.Name).Count(&repeat).Error; err != nil {
		return fmt.Errorf("查重（名称）失败：%w", err)
	}

	if repeat > 0 {
		return errors.New("名称重复")
	}

	if err = my.DB().Where("uuid <> ?", my.UUID.String()).Where("identity", my.Identity).Count(&repeat).Error; err != nil {
		return fmt.Errorf("查重（标识）失败：%w", err)
	}

	if repeat > 0 {
		return errors.New("标识重复")
	}

	my.UpdatedAt = time.Now()
	if err = my.DB().Where("uuid", my.UUID.String()).Save(my).Error; err != nil {
		return fmt.Errorf("编辑失败：%w", err)
	}

	return nil
}

func (my Permission) Destroy() error {
	if edgeIns.db == nil {
		return errors.New("数据库连接失败")
	}

	return my.DB().Where("uuid", my.UUID.String()).Delete(my).Error
}
