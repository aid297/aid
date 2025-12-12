package rbac

import (
	"gorm.io/gorm"
)

type (
	EdgeAttributer interface{ Register(edge *Edge) }

	AttrTablePrefix   struct{ tablePrefix string }
	AttrDB            struct{ db *gorm.DB }
	AttrGroupUUID     struct{ faceUUID string }
	AttrIntersection1 struct{ intersection1 string }
	AttrIntersection2 struct{ intersection2 string }
)

func TablePrefix(tablePrefix string) EdgeAttributer { return AttrTablePrefix{tablePrefix} }
func (my AttrTablePrefix) Register(edge *Edge)      { edge.tablePrefix = my.tablePrefix }

func DB(db *gorm.DB) EdgeAttributer   { return AttrDB{db} }
func (my AttrDB) Register(edge *Edge) { edge.db = my.db }

func GroupUUID(face string) EdgeAttributer   { return AttrGroupUUID{face} }
func (my AttrGroupUUID) Register(edge *Edge) { edge.RoleUUID = my.faceUUID }

func Intersection1(intersection1 string) EdgeAttributer { return AttrIntersection1{intersection1} }
func (my AttrIntersection1) Register(edge *Edge)        { edge.Intersection1 = my.intersection1 }

func Intersection2(intersection2 string) EdgeAttributer { return AttrIntersection2{intersection2} }
func (my AttrIntersection2) Register(edge *Edge)        { edge.Intersection2 = my.intersection2 }
