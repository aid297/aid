package rbac

import (
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aid297/aid/dict"
	"github.com/aid297/aid/str"
)

type Edge struct {
	tablePrefix   string
	db            *gorm.DB
	RoleUUID      string `gorm:"column:role_uuid;type:char(36);not null;primaryKey;comment:组UUID;" json:"roleUUID"`
	Role          *Role  `gorm:"foreignKey:role_uuid;references:uuid;" json:"role"`
	Intersection1 string `gorm:"column:intersection1;type:varchar(255);not null;primaryKey;comment:交点1;" json:"intersection1"`
	Intersection2 string `gorm:"column:intersection2;type:varchar(255);not null;primaryKey;comment:交点2;" json:"intersection2"`
}

var (
	edgeOnce sync.Once
	edgeIns  *Edge
)

func (*Edge) Once(options ...EdgeAttributer) *Edge {
	edgeOnce.Do(func() { edgeIns = new(Edge) })
	return edgeIns.SetAttrs(options...)
}

func (*Edge) SetAttrs(attrs ...EdgeAttributer) *Edge {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(edgeIns)
		}
	}
	return edgeIns
}

func (*Edge) TableName() string {
	if edgeIns.tablePrefix != "" {
		return str.APP.Buffer.JoinString(edgeIns.tablePrefix, "_edges")
	} else {
		return "rbac_edges"
	}
}

func (*Edge) DB() *gorm.DB { return edgeIns.db.Model(new(Edge)) }

func (*Edge) AutoMigrate() error {
	if edgeIns.db == nil {
		return ErrDBConnFailed
	}

	return edgeIns.db.AutoMigrate(new(Edge), new(Group), new(Role), new(Permission))
}

func (*Edge) Bind(roleUUID string, intersections map[string][]string) error {
	if edgeIns.db == nil {
		return ErrDBConnFailed
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

		var edges = make([]*Edge, 0, count)
		intersectionsDict.Each(func(key string, value []string) {
			for idx := range value {
				edges = append(edges, &Edge{
					RoleUUID:      roleUUID,
					Intersection1: key,
					Intersection2: value[idx],
				})
			}
		})

		return tx.Model(new(Edge)).Create(&edges).Error
	})
}

func (*Edge) GetByIntersection1(intersection1 string) ([]*Edge, error) {
	var (
		err   error
		edges []*Edge
	)

	if edgeIns.db == nil {
		return nil, ErrDBConnFailed
	}

	if err = edgeIns.DB().
		Preload("RBACs").
		Preload("Groups").
		Preload("Groups.Permission").
		Preload("Groups.Role").
		Where(str.APP.Buffer.JoinString(new(Edge).TableName(), ".intersection1"), intersection1).
		Find(&edges).Error; err != nil {
		return nil, err
	}

	return edges, nil
}

func (*Edge) getGroupDB(intersection1, intersection2, identity string) *gorm.DB {
	return APP.Group.New().DB().
		Table(str.BufferApp.JoinStringLimit(APP.Edge.TableName(), "as", "e")).
		Joins(str.BufferApp.JoinStringLimit("join", APP.Group.TableName(), "as", "g", "on", "g.role_uuid = rbacs.role_uuid")).
		Joins(str.BufferApp.JoinStringLimit("join", APP.Role.TableName(), "as", "r", "on", "r.uuid = g.role_uuid")).
		Joins(str.BufferApp.JoinStringLimit("join ", APP.Permission.TableName(), "as", "p", "on", "p.uuid = g.permission_uuid")).
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
		return false, ErrDBConnFailed
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (*Edge) CheckRolePermission(intersection1, intersection2, identity string, roleUUID uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)

	if edgeIns.db == nil {
		return false, ErrDBConnFailed
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid", roleUUID.String()).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (*Edge) CheckInRolesPermission(intersection1, intersection2, identity string, roleUUIDs ...uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)

	if edgeIns.db == nil {
		return false, ErrDBConnFailed
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid in (?)", roleUUIDs).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (*Edge) CheckNotInRolesPermission(intersection1, intersection2, identity string, roleUUIDs ...uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)

	if edgeIns.db == nil {
		return false, ErrDBConnFailed
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid not in (?)", roleUUIDs).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (*Edge) CheckAllRolesPermission(intersection1, intersection2, identity string, roleUUIDs ...uuid.UUID) (bool, error) {
	var (
		err   error
		count int64
	)

	if edgeIns.db == nil {
		return false, ErrDBConnFailed
	}

	if err = APP.Edge.getGroupDB(intersection1, intersection2, identity).Where("r.uuid in (?)", roleUUIDs).Count(&count).Error; err != nil {
		return false, err
	}

	return count == int64(len(roleUUIDs)), nil
}
