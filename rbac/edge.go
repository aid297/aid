package rbac

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aid297/aid/dict"
	"github.com/aid297/aid/str"
)

type Edge struct {
	tablePrefix   string
	db            *gorm.DB
	RoleUUID      uuid.UUID `gorm:"column:role_uuid;type:char(36);not null;primaryKey;comment:组UUID;" json:"roleUUID"`
	Role          *Role     `gorm:"foreignKey:role_uuid;references:uuid;" json:"role"`
	Intersection1 string    `gorm:"column:intersection1;type:varchar(255);not null;primaryKey;comment:交点1;" json:"intersection1"`
	Intersection2 string    `gorm:"column:intersection2;type:varchar(255);not null;primaryKey;comment:交点2;" json:"intersection2"`
}

var (
	edgeOnce sync.Once
	edgeIns  *Edge
)

func (*Edge) Once(options ...EdgeAttributer) *Edge {
	edgeOnce.Do(func() { edgeIns = new(Edge) })
	return edgeIns.Set(options...)
}

func (*Edge) Set(options ...EdgeAttributer) *Edge {
	if len(options) > 0 {
		for idx := range options {
			options[idx].Register(edgeIns)
		}
	}
	return edgeIns
}

func (*Edge) TableName() string {
	if edgeIns.tablePrefix != "" {
		return str.BufferApp.NewString(edgeIns.tablePrefix).S("_edges").String()
	} else {
		return "rbac_edges"
	}
}

func (*Edge) DB() *gorm.DB { return edgeIns.db.Model(new(Edge)) }

func (*Edge) AutoMigrate() error {
	if edgeIns.db == nil {
		return errors.New("数据库连接失败")
	}

	return edgeIns.db.AutoMigrate(new(Edge), new(Group), new(Role), new(Permission))
}

func (*Edge) Bind(roleUUID uuid.UUID, intersections map[string][]string) error {
	if edgeIns.db == nil {
		return errors.New("数据库连接失败")
	}

	return edgeIns.db.Transaction(func(tx *gorm.DB) error {
		intersectionsDict := dict.New(intersections)
		if err := tx.Model(new(Edge)).
			Where("group_uuid = ?", roleUUID).
			Where("intersection1 IN ?", intersectionsDict.GetKeys().ToSlice()).
			Delete(new(Edge)).Error; err != nil {
			return err
		}

		count := 0
		intersectionsDict.Each(func(_ string, value []string) { count += len(value) })

		var triangles = make([]*Edge, 0, count)
		intersectionsDict.Each(func(key string, value []string) {
			for idx := range value {
				triangles = append(triangles, &Edge{
					RoleUUID:      roleUUID,
					Intersection1: key,
					Intersection2: value[idx],
				})
			}
		})

		return tx.Model(new(Edge)).Create(&triangles).Error
	})
}

func (*Edge) GetByIntersection1(intersection1 string) ([]*Edge, error) {
	var (
		err       error
		triangles []*Edge
	)
	if edgeIns.db == nil {
		return nil, errors.New("数据库连接失败")
	}

	if err = edgeIns.DB().
		Preload("RBACs").
		Preload("Groups").
		Preload("Groups.Permission").
		Preload("Groups.Role").
		Where("edges.intersection1", intersection1).
		Find(&triangles).Error; err != nil {
		return nil, err
	}

	return triangles, nil
}

func (*Edge) getGroupDB(intersection1, intersection2, identity string) *gorm.DB {
	return APP.Group.New().
		DB().
		Table(str.BufferApp.JoinString(APP.Edge.TableName(), " as e")).
		Joins(str.BufferApp.JoinString("join ", APP.Group.TableName(), " as g on g.role_uuid = rbacs.role_uuid")).
		Joins(str.BufferApp.JoinString("join ", APP.Role.TableName(), " as r on r.uuid = g.role_uuid")).
		Joins(str.BufferApp.JoinString("join ", APP.Permission.TableName(), " as p on p.uuid = g.permission_uuid")).
		Where("p.identity", identity).
		Where("e.intersection1", intersection1).
		Where("e.intersection2", intersection2)
}

func (*Edge) CheckPermission(intersection1, intersection2, identity string) (bool, error) {
	var (
		err   error
		count int64
	)
	if edgeIns.db == nil {
		return false, errors.New("数据库连接失败")
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (*Edge) CheckRolePermission(intersection1, intersection2, identity string, roleUUID uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)
	if edgeIns.db == nil {
		return false, errors.New("数据库连接失败")
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid", roleUUID.String()).Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (*Edge) CheckInRolesPermission(intersection1, intersection2, identity string, roleUUIDs ...uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)
	if edgeIns.db == nil {
		return false, errors.New("数据库连接失败")
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid in (?)", roleUUIDs).Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (*Edge) CheckNotInRolesPermission(intersection1, intersection2, identity string, roleUUIDs ...uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)
	if edgeIns.db == nil {
		return false, errors.New("数据库连接失败")
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid not in (?)", roleUUIDs).Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (*Edge) CheckAllRolesPermission(intersection1, intersection2, identity string, roleUUIDs ...uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)
	if edgeIns.db == nil {
		return false, errors.New("数据库连接失败")
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid in (?)", roleUUIDs).Count(&count).Error; err != nil {
		return false, err
	}

	if count == int64(len(roleUUIDs)) {
		return true, nil
	}

	return false, nil
}
