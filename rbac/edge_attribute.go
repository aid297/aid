package rbac

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	EdgeAttributer interface {
		Register(triangle *Edge)
	}

	AttrTablePrefix   struct{ tablePrefix string }
	AttrDB            struct{ db *gorm.DB }
	AttrGroupUUID     struct{ faceUUID uuid.UUID }
	AttrIntersection1 struct{ intersection1 string }
	AttrIntersection2 struct{ intersection2 string }
)

func TablePrefix(tablePrefix string) EdgeAttributer {
	return AttrTablePrefix{tablePrefix: tablePrefix}
}

func (my AttrTablePrefix) Register(triangle *Edge) {
	triangle.tablePrefix = my.tablePrefix
}

func DB(db *gorm.DB) EdgeAttributer { return AttrDB{db: db} }

func (my AttrDB) Register(triangle *Edge) { triangle.db = my.db }

func GroupUUID(face uuid.UUID) EdgeAttributer { return AttrGroupUUID{faceUUID: face} }

func (my AttrGroupUUID) Register(triangle *Edge) { triangle.RoleUUID = my.faceUUID }

func Intersection1(intersection1 string) EdgeAttributer {
	return AttrIntersection1{intersection1: intersection1}
}

func (my AttrIntersection1) Register(triangle *Edge) {
	triangle.Intersection1 = my.intersection1
}

func Intersection2(intersection2 string) EdgeAttributer {
	return AttrIntersection2{intersection2: intersection2}
}

func (my AttrIntersection2) Register(triangle *Edge) {
	triangle.Intersection2 = my.intersection2
}
