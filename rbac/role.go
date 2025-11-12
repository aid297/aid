package rbac

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/str"
)

type Role struct {
	ID        uint           `gorm:"primaryKey;" json:"id"`
	CreatedAt time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:新建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deletedAt"`
	UUID      uuid.UUID      `gorm:"type:char(36);not null;unique;comment:uuid;" json:"uuid"`
	Name      string         `gorm:"type:varchar(255);not null;unique;comment:角色名称" json:"name"`
	Groups    []*Group       `gorm:"foreignKey:group_uuid;references:uuid;" json:"groups"`
	RBACs     []*Edge        `gorm:"foreignKey:group_uuid;references:uuid;" json:"rbacs"`
}

func (Role) New(attrs ...RoleAttributer) Role {
	ins := new(Role)
	return ins.Set(attrs...)
}

func (my Role) Set(attrs ...RoleAttributer) Role {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}

	return my
}

func (Role) TableName() string {
	if edgeIns.tablePrefix != "" {
		return str.BufferApp.NewString(edgeIns.tablePrefix).S("_roles").String()
	} else {
		return str.BufferApp.NewString("triangle_roles").String()
	}
}

func (Role) DB() *gorm.DB { return edgeIns.db.Model(new(Role)) }

func (my Role) Store(options ...RoleAttributer) error {
	var (
		err error
		r   int64
	)

	if edgeIns.db == nil {
		return errors.New("数据库链接失败")
	}

	if len(options) > 0 {
		my.Set(options...)
	}

	if err = my.DB().Where("name", my.Name).Count(&r).Error; err != nil {
		return fmt.Errorf("查重失败：%w", err)
	}

	if r > 0 {
		return errors.New("名称重复")
	}

	my.UUID = uuid.Must(uuid.NewV6())
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
		return errors.New("数据库连接失败")
	}

	if len(attrs) > 0 {
		my.Set(attrs...)
	}

	if err = my.DB().Where("uuid <> ?", my.UUID.String()).Where("name", my.Name).Count(&r).Error; err != nil {
		return fmt.Errorf("查重失败：%w", err)
	}

	if r > 0 {
		return errors.New("名称重复")
	}

	my.UpdatedAt = time.Now()
	if err = my.DB().Where("uuid", my.UUID.String()).Save(my).Error; err != nil {
		return fmt.Errorf("编辑失败：%w", err)
	}

	return nil
}

func (my Role) Destroy() error {
	if edgeIns.db == nil {
		return errors.New("数据库连接失败")
	}

	return my.DB().Where("uuid", my.UUID.String()).Delete(my).Error
}

func (my Role) GetDotsByUUID(edgeUUID uuid.UUID) ([]Permission, error) {
	var groups []*Group
	if err := my.DB().Preload("Groups").Where("uuid", edgeUUID).Find(&groups).Error; err != nil {
		return nil, err
	}

	var permissionUUIDs = make([]string, len(groups))
	array.New(groups).Each(func(idx int, item *Group) { permissionUUIDs[idx] = item.PermissionUUID.String() })

	var permissions []Permission
	if err := APP.Permission.DB().Where("uuid in (?)", permissionUUIDs).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (my Role) GetDotsByUUIDs(edgeUUIDs []uuid.UUID) ([]Permission, error) {
	var groups []*Group
	if err := my.DB().Preload("Groups").Where("uuid in (?)", edgeUUIDs).Find(&groups).Error; err != nil {
		return nil, err
	}

	var permissionUUIDs = make([]string, len(groups))
	array.New(groups).Each(func(idx int, item *Group) { permissionUUIDs[idx] = item.PermissionUUID.String() })

	var permissions []Permission
	if err := APP.Permission.DB().Where("uuid in (?)", permissionUUIDs).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}
