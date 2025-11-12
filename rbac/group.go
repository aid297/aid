package rbac

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aid297/aid/str"
)

type Group struct {
	RoleUUID       uuid.UUID   `gorm:"type:char(36);not null;primaryKey;comment:角色UUID;" json:"roleUUID"`
	Role           *Role       `gorm:"foreignKey:role_uuid;references:uuid;" json:"role"`
	PermissionUUID uuid.UUID   `gorm:"type:char(36);not null;primaryKey;comment:权限UUID;" json:"permissionUUID"`
	Permission     *Permission `gorm:"foreignKey:permission_uuid;references:uuid;" json:"permission"`
}

func (Group) New() Group { return Group{} }

func (Group) TableName() string {
	if edgeIns.tablePrefix != "" {
		return str.BufferApp.NewString(edgeIns.tablePrefix).S("_groups").String()
	} else {
		return str.BufferApp.NewString("triangle_groups").String()
	}
}

func (Group) DB() *gorm.DB { return edgeIns.db.Model(new(Group)) }

func (Group) BindPermissions(edge *Role, dots []*Permission) error {
	if edgeIns.db == nil {
		return errors.New("数据库连接失败")
	}

	return edgeIns.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(new(Group)).Where("group_uuid", edge.UUID.String()).Delete(new(Group)).Error; err != nil {
			return err
		}

		if len(dots) > 0 {
			faces := make([]*Group, 0, len(dots))

			for idx := range dots {
				faces = append(faces, &Group{RoleUUID: edge.UUID, PermissionUUID: dots[idx].UUID})
			}

			if err := tx.Create(&faces).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
