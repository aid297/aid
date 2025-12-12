package rbac

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/str"
)

type Role struct {
	ID        uint           `gorm:"primaryKey;" json:"id"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deletedAt"`
	UUID      string         `gorm:"column:uuid;type:char(36);not null;unique;comment:uuid;" json:"uuid"`
	Name      string         `gorm:"type:varchar(255);not null;unique;comment:角色名称" json:"name"`
	Groups    []*Group       `gorm:"foreignKey:group_uuid;references:uuid;" json:"groups"`
	RBACs     []*Edge        `gorm:"foreignKey:group_uuid;references:uuid;" json:"rbacs"`
}

func (Role) New(attrs ...RoleAttributer) Role { return Role{}.SetAttrs(attrs...) }

func (my Role) SetAttrs(attrs ...RoleAttributer) Role {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}

	return my
}

func (Role) TableName() string {
	if edgeIns.tablePrefix != "" {
		return str.APP.Buffer.JoinString(edgeIns.tablePrefix, "_roles")
	} else {
		return "edge_roles"
	}
}

func (Role) DB() *gorm.DB { return edgeIns.db.Model(new(Role)) }

func (my Role) Store(attrs ...RoleAttributer) error {
	var (
		err error
		r   int64
	)

	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	if err = my.SetAttrs(attrs...).DB().Where("name", my.Name).Count(&r).Error; err != nil {
		return fmt.Errorf("%w：%w", ErrCheckRepeat, err)
	}

	if r > 0 {
		return fmt.Errorf("角色%w", ErrRepeatName)
	}

	my.UUID = uuid.Must(uuid.NewV6()).String()
	my.CreatedAt = time.Now()
	my.UpdatedAt = time.Now()
	return my.DB().Create(my).Error
}

func (my Role) Update(attrs ...RoleAttributer) error {
	var (
		err error
		r   int64
	)

	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	if err = my.SetAttrs(attrs...).DB().Where("uuid <> ?", my.UUID).Where("name", my.Name).Count(&r).Error; err != nil {
		return fmt.Errorf("%w：%w", ErrCheckRepeat, err)
	}

	if r > 0 {
		return ErrRepeatName
	}

	my.UpdatedAt = time.Now()
	if err = my.DB().Where("uuid", my.UUID).Save(my).Error; err != nil {
		return fmt.Errorf("%w：%w", ErrUpdate, err)
	}

	return nil
}

func (my Role) Destroy() error {
	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	return my.DB().Where("uuid", my.UUID).Delete(my).Error
}

func (my Role) GetDotsByUUID(edgeUUID uuid.UUID) (permissions []Permission, err error) {
	var groups []*Group
	if err = my.DB().Preload("Groups").Where("uuid", edgeUUID).Find(&groups).Error; err != nil {
		return
	}

	if err = APP.Permission.DB().Where("uuid in (?)", anyArrayV2.Cast(anyArrayV2.NewList(groups).Pluck(func(item *Group) any { return item.PermissionUUID }), func(src any) string { return cast.ToString(src) }).ToSlice()).Find(&permissions).Error; err != nil {
		return
	}

	return
}

func (my Role) GetDotsByUUIDs(edgeUUIDs []uuid.UUID) (permissions []Permission, err error) {
	var groups []*Group
	if err = my.DB().Preload("Groups").Where("uuid in (?)", edgeUUIDs).Find(&groups).Error; err != nil {
		return
	}

	if err = APP.Permission.DB().Where("uuid in (?)", anyArrayV2.Cast(anyArrayV2.NewList(groups).Pluck(func(item *Group) any { return item.PermissionUUID }), func(src any) string { return cast.ToString(src) }).ToSlice()).Find(&permissions).Error; err != nil {
		return
	}

	return
}
