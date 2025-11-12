package mongoClientPool

import (
	"log"
	"testing"

	"github.com/aid297/aid/array"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Student struct {
		Id      OID    `bson:"_id"`
		Name    string `bson:"name"`
		Age     uint64 `bson:"age"`
		ClassId OID    `bson:"class_id"`
		Class   *Class `bson:"-"`
	}

	Class struct {
		Id       OID        `bson:"_id"`
		Name     string     `bson:"name"`
		Students []*Student `bson:"-"`
	}
)

func getDB(t *testing.T) (*MongoClientPool, *MongoClient) {
	var err error
	mp := APP.Pool.Once()

	mc, err := APP.Client.New("mongodb://admin:admin@localhost:27017")
	if err != nil {
		t.Fatalf("创建mongo客户端失败：%v", err)
	}
	if _, err = mp.AppendClient("default", mc); err != nil {
		t.Fatalf("添加mongo客户端失败：%v", err)
	}
	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")

	return mp, mc
}

func Test1InsertOne(t *testing.T) {
	t.Run("操作单条数据", func(t *testing.T) {
		var (
			err          error
			insertOneRes *mongo.InsertOneResult
			mp, mc       = getDB(t)
			user         = Student{Name: "张三", Age: 18}
		)
		// 清空数据
		_ = mc.DeleteMany(nil)
		if mc.Err != nil {
			t.Fatalf("清空数据失败：%v", err)
		}

		// 插入单条数据
		if mc.InsertOne(user, &insertOneRes).Err != nil {
			log.Fatalf("插入单条数据失败：%v", err)
		}
		t.Logf("插入单条数据成功：%s\n", insertOneRes.InsertedID.(OID).String())

		mp.Clean()
	})
}

func Test2InsertMany(t *testing.T) {
	t.Run("操作多条数据", func(t *testing.T) {
		var (
			insertOneRes  *mongo.InsertOneResult
			insertManyRes *mongo.InsertManyResult
			mp, mc        = getDB(t)
		)
		// 插入多条数据
		if mc.SetCollection("classes").InsertOne(Class{Id: primitive.NewObjectID(), Name: "一班"}, &insertOneRes).Err != nil {
			t.Fatalf("插入班级失败：%v", mc.Err)
		}
		t.Logf("插入班级成功：%s\n", insertOneRes.InsertedID.(OID).String())

		if mc.SetCollection("students").InsertMany([]any{
			Student{Id: primitive.NewObjectID(), Name: "张三", Age: 18, ClassId: insertOneRes.InsertedID.(OID)},
			Student{Id: primitive.NewObjectID(), Name: "李四", Age: 19, ClassId: insertOneRes.InsertedID.(OID)},
		}, &insertManyRes).Err != nil {
			t.Fatalf("插入多条数据失败：%v", mc.Err)
		}

		if mc.SetCollection("classes").InsertOne(Map{"name": "二班"}, &insertOneRes).Err != nil {
			t.Fatalf("插入班级失败：%v", mc.Err)
		}

		if mc.SetCollection("students").InsertMany([]any{
			Student{Id: primitive.NewObjectID(), Name: "王五", Age: 20, ClassId: insertOneRes.InsertedID.(OID)},
			Student{Id: primitive.NewObjectID(), Name: "赵六", Age: 21, ClassId: insertOneRes.InsertedID.(OID)},
		}, &insertManyRes).Err != nil {
			t.Fatalf("插入学生失败：%v", mc.Err)
		}

		t.Logf("插入多条数据成功：%v\n", insertManyRes.InsertedIDs)

		mp.Clean()
	})
}

func Test3UpdateOne(t *testing.T) {
	var (
		updateOneRes *mongo.UpdateResult
		mp, mc       = getDB(t)
	)

	if mc.Where(Map{"name": "张三"}).UpdateOne(Student{Name: "张三", Age: 1}, &updateOneRes).Err != nil {
		t.Fatalf("更新单条数据失败：%v", mc.Err)
	}
	t.Logf("更新成功：%d\n", updateOneRes.ModifiedCount)

	mp.Clean()
}

func Test4UpdateMany(t *testing.T) {
	var (
		updateManyRes *mongo.UpdateResult
		mp, mc        = getDB(t)
	)

	if mc.SetCollection("students").Where(Map{"name": Map{"$ne": "张三"}}).UpdateMany(Map{"age": 0}, &updateManyRes).Err != nil {
		t.Fatalf("更新单条数据失败：%v", mc.Err)
	}
	t.Logf("更新成功：%d\n", updateManyRes.ModifiedCount)

	mp.Clean()
}

func Test5FindOne(t *testing.T) {
	var (
		student *Student
		mp, mc  = getDB(t)
	)

	t.Run("查询单条数据", func(t *testing.T) {
		if mc.SetCollection("students").Where(Map{"name": "张三"}).FindOne(&student, nil).Err != nil {
			t.Fatalf("查询单条数据失败：%v", mc.Err)
		}
		t.Errorf("查询单条数据成功：%#v\n", student)

		mp.Clean()
	})
}

func Test6FindMany(t *testing.T) {
	var (
		classes  []Map
		classes2 []*Class
		classA   *Class
		students []*Student
		mp, mc   = getDB(t)
	)

	t.Run("查询多条数据1", func(t *testing.T) {
		if mc.SetCollection("classes").
			Where(Map{"$lookup": Map{
				"from":         "students",
				"localField":   "_id",
				"foreignField": "class_id",
				"as":           "students",
			}}, Map{"$match": Map{"name": "一班"}}).
			Aggregate(&classes).Err != nil {
			t.Fatalf("查询多条数据失败：%v", mc.Err)
		}
		t.Logf("查询多条数据成功：%#v\n", classes)
	})

	t.Run("查询多条数据2", func(t *testing.T) {
		if mc.SetCollection("classes").
			Where(Map{"name": "一班"}).
			FindOne(&classA, nil).Err != nil {
			t.Fatalf("查询多条数据失败：%v", mc.Err)
		}

		if mc.SetCollection("students").
			Where(Map{"class_id": classA.Id}).
			FindMany(&students, nil).Err != nil {
			t.Fatalf("查询多条数据失败：%v", mc.Err)
		}

		classA.Students = students
		t.Logf("查询成功：%#v\n", classA)
	})

	t.Run("查询多条数据3", func(t *testing.T) {
		if mc.SetCollection("students").
			FindMany(&students, nil).Err != nil {
			t.Fatalf("查询多条数据失败：%v", mc.Err)
		}

		if mc.SetCollection("classes").
			Where(Map{
				"_id": Map{
					"$in": array.Cast(array.New(students), func(value *Student) OID { return value.ClassId }),
				},
			}).
			FindMany(&classes2, nil).Err != nil {
			t.Fatalf("查询多条数据失败：%v", mc.Err)
		}

		for idx := range students {
			for _, class := range classes2 {
				if students[idx].ClassId == class.Id {
					students[idx].Class = class
				}
			}
		}

		t.Logf("查询成功：%v\n", students)
	})

	mp.Clean()
}

func Test7DeleteOne(t *testing.T) {
	var (
		deleteOneRes *mongo.DeleteResult
		mp, mc       = getDB(t)
	)
	t.Run("删除单条数据", func(t *testing.T) {
		// 删除单条数据
		if mc.Where(Map{"name": "张三"}).DeleteOne(&deleteOneRes).Err != nil {
			t.Fatalf("删除单条数据失败：%v", mc.Err)
		}
		t.Logf("成功删除数据：%d\n", deleteOneRes.DeletedCount)

		mp.Clean()
	})
}

func Test8DeleteMany(t *testing.T) {
	var (
		deleteManyRes *mongo.DeleteResult
		mp, mc        = getDB(t)
	)

	t.Run("删除多条数据", func(t *testing.T) {
		// 删除多条数据
		if mc.Where(Map{"name": Map{"$ne": "张三"}}).DeleteMany(&deleteManyRes).Err != nil {
			t.Fatalf("删除多条数据失败：%v", mc.Err)
		}
		t.Logf("删除数据成功：%d\n", deleteManyRes.DeletedCount)

	})

	mp.Clean()
}
