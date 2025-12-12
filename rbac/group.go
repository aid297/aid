package rbac

import (
	"gorm.io/gorm"

	"github.com/aid297/aid/str"
)

type Group struct {
	RoleUUID       string      `gorm:"column:role_uuid;type:char(36);not null;primaryKey;comment:角色UUID;" json:"roleUUID"`
	Role           *Role       `gorm:"foreignKey:role_uuid;references:uuid;" json:"role"`
	PermissionUUID string      `gorm:"column:permission_uuid;type:char(36);not null;primaryKey;comment:权限UUID;" json:"permissionUUID"`
	Permission     *Permission `gorm:"foreignKey:PermissionUUID;references:UUID;" json:"permission"`
}

func (Group) New() Group { return Group{} }

func (Group) TableName() string {
	if edgeIns.tablePrefix != "" {
		return str.APP.Buffer.JoinString(edgeIns.tablePrefix, "_groups")
	} else {
		return "edge_groups"
	}
}

func (Group) DB() *gorm.DB { return edgeIns.db.Model(new(Group)) }

func (Group) BindPermissions(edge *Role, dots []*Permission) error {
	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	return edgeIns.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(new(Group)).Where("group_uuid", edge.UUID).Delete(new(Group)).Error; err != nil {
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
