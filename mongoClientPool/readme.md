### MongoDB 连接池

1. 初始化
   ```go
   package main
   
   import (
   	`fmt`
   
   	`github.com/aid297/aid/mongoClientPool`
   )
   
   func main() {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(fmt.Errorf("创建mongo客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(fmt.Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	fmt.Printf("连接池：%v\n", mp)
   	fmt.Printf("客户端：%v\n", mc)
   }
   ```

2. 插入单条数据
   ```go
   package main
   
   import (
   	. `fmt`
   	`log`
   	`testing`
   
   	`github.com/aid297/aid/mongoClientPool`
   	`go.mongodb.org/mongo-driver/mongo`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		err          error
   		insertOneRes *mongo.InsertOneResult
   		mp, mc       = getDB()
   		user         = Student{Name: "张三", Age: 18}
   	)
   	// 清空数据
   	_ = mc.DeleteMany(nil)
   	if mc.Err != nil {
   		panic(Errorf("清空数据失败： %v", err))
   	}
   
   	// 插入单条数据
   	if mc.InsertOne(user, &insertOneRes).Err != nil {
   		panic(Errorf("插入单条数据失败：%v", err))
   	}
   	Printf("插入单条数据成功：%s\n", insertOneRes.InsertedID.(mongoClientPool.OID).String())
   
   	mp.Clean()
   }
   ```

3. 插入多条数据
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/mongoClientPool`
   	`go.mongodb.org/mongo-driver/bson/primitive`
   	`go.mongodb.org/mongo-driver/mongo`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		insertOneRes  *mongo.InsertOneResult
   		insertManyRes *mongo.InsertManyResult
   		mp, mc        = getDB()
   	)
   	// 插入多条数据
   	if mc.SetCollection("classes").InsertOne(Class{Id: primitive.NewObjectID(), Name: "一班"}, &insertOneRes).Err != nil {
   		panic(Errorf("插入班级失败：%v", mc.Err))
   	}
   	Printf("插入班级成功：%s\n", insertOneRes.InsertedID.(mongoClientPool.OID).String())
   
   	if mc.SetCollection("students").InsertMany([]any{
   		Student{Id: primitive.NewObjectID(), Name: "张三", Age: 18, ClassId: insertOneRes.InsertedID.(mongoClientPool.OID)},
   		Student{Id: primitive.NewObjectID(), Name: "李四", Age: 19, ClassId: insertOneRes.InsertedID.(mongoClientPool.OID)},
   	}, &insertManyRes).Err != nil {
   		panic(Errorf("插入多条数据失败：%v", mc.Err))
   	}
   
   	if mc.SetCollection("classes").InsertOne(mongoClientPool.Map{"name": "二班"}, &insertOneRes).Err != nil {
   		panic(Errorf("插入班级失败：%v", mc.Err))
   	}
   
   	if mc.SetCollection("students").InsertMany([]any{
   		Student{Id: primitive.NewObjectID(), Name: "王五", Age: 20, ClassId: insertOneRes.InsertedID.(mongoClientPool.OID)},
   		Student{Id: primitive.NewObjectID(), Name: "赵六", Age: 21, ClassId: insertOneRes.InsertedID.(mongoClientPool.OID)},
   	}, &insertManyRes).Err != nil {
   		panic(Errorf("插入学生失败：%v", mc.Err))
   	}
   
   	Printf("插入多条数据成功：%v\n", insertManyRes.InsertedIDs)
   
   	mp.Clean()
   }
   ```

4. 更新单条数据
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/mongoClientPool`
   	`go.mongodb.org/mongo-driver/mongo`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		updateOneRes *mongo.UpdateResult
   		mp, mc       = getDB()
   	)
   
   	if mc.Where(mongoClientPool.Map{"name": "张三"}).UpdateOne(Student{Name: "张三", Age: 1}, &updateOneRes).Err != nil {
   		panic(Errorf("更新单条数据失败：%v", mc.Err))
   	}
   	Printf("更新成功：%d\n", updateOneRes.ModifiedCount)
   
   	mp.Clean()
   }
   ```

5. 更新多条数据
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/mongoClientPool`
   	`go.mongodb.org/mongo-driver/mongo`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		updateManyRes *mongo.UpdateResult
   		mp, mc        = getDB()
   	)
   
   	if mc.SetCollection("students").Where(mongoClientPool.Map{"name": mongoClientPool.Map{"$ne": "张三"}}).UpdateMany(mongoClientPool.Map{"age": 0}, &updateManyRes).Err != nil {
   		panic(Errorf("更新单条数据失败：%v", mc.Err))
   	}
   	Printf("更新成功：%d\n", updateManyRes.ModifiedCount)
   
   	mp.Clean()
   }
   ```

6. 查询单条数据
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/mongoClientPool`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		student *Student
   		mp, mc  = getDB()
   	)
   
   	if mc.SetCollection("students").Where(mongoClientPool.Map{"name": "张三"}).FindOne(&student, nil).Err != nil {
   		panic(Errorf("查询单条数据失败：%v", mc.Err))
   	}
   	Printf("查询单条数据成功：%#v\n", student)
   
   	mp.Clean()
   }
   ```

7. 查询多条数据
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/array`
   	`github.com/aid297/aid/mongoClientPool`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		classes  []mongoClientPool.Map
   		classes2 []*Class
   		classA   *Class
   		students []*Student
   		mp, mc   = getDB()
   	)
   
   	// 查询多条数据1
   	func() {
   		if mc.SetCollection("classes").
   			Where(mongoClientPool.Map{"$lookup": mongoClientPool.Map{
   				"from":         "students",
   				"localField":   "_id",
   				"foreignField": "class_id",
   
   				"as": "students",
   			}}, mongoClientPool.Map{"$match": mongoClientPool.Map{"name": "一班"}}).
   			Aggregate(&classes).Err != nil {
   			panic(Errorf("查询多条数据失败：%v", mc.Err))
   		}
   		Printf("查询多条数据成功：%#v\n", classes)
   	}()
   
   	// 查询多条数据2
   	func() {
   		if mc.SetCollection("classes").
   			Where(mongoClientPool.Map{"name": "一班"}).
   			FindOne(&classA, nil).Err != nil {
   			panic(Errorf("查询多条数据失败：%v", mc.Err))
   		}
   
   		if mc.SetCollection("students").
   			Where(mongoClientPool.Map{"class_id": classA.Id}).
   			FindMany(&students, nil).Err != nil {
   			panic(Errorf("查询多条数据失败：%v", mc.Err))
   		}
   
   		classA.Students = students
   		Printf("查询成功：%#v\n", classA)
   	}()
   
   	// 查询多条数据3
   	func() {
   		if mc.SetCollection("students").
   			FindMany(&students, nil).Err != nil {
   			panic(Errorf("查询多条数据失败：%v", mc.Err))
   		}
   
   		if mc.SetCollection("classes").
   			Where(mongoClientPool.Map{
   				"_id": mongoClientPool.Map{
   					"$in": array.Cast(array.New(students), func(value *Student) mongoClientPool.OID { return value.ClassId }),
   				},
   			}).
   			FindMany(&classes2, nil).Err != nil {
   			panic(Errorf("查询多条数据失败：%v", mc.Err))
   		}
   
   		for idx := range students {
   			for _, class := range classes2 {
   				if students[idx].ClassId == class.Id {
   					students[idx].Class = class
   				}
   			}
   		}
   
   		Printf("查询成功：%v\n", students)
   	}()
   
   	mp.Clean()
   }
   ```

8. 删除单条数据
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/mongoClientPool`
   	`go.mongodb.org/mongo-driver/mongo`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		deleteOneRes *mongo.DeleteResult
   		mp, mc       = getDB()
   	)
   
   	// 删除单条数据
   	if mc.Where(mongoClientPool.Map{"name": "张三"}).DeleteOne(&deleteOneRes).Err != nil {
   		panic(Errorf("删除单条数据失败：%v", mc.Err))
   	}
   	Printf("成功删除数据：%d\n", deleteOneRes.DeletedCount)
   
   	mp.Clean()
   }
   ```

9. 删除多条数据
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/mongoClientPool`
   	`go.mongodb.org/mongo-driver/mongo`
   )
   
   type (
   	Student struct {
   		Id      mongoClientPool.OID `bson:"_id"`
   		Name    string              `bson:"name"`
   		Age     uint64              `bson:"age"`
   		ClassId mongoClientPool.OID `bson:"class_id"`
   		Class   *Class              `bson:"-"`
   	}
   
   	Class struct {
   		Id       mongoClientPool.OID `bson:"_id"`
   		Name     string              `bson:"name"`
   		Students []*Student          `bson:"-"`
   	}
   )
   
   func getDB() (*mongoClientPool.MongoClientPool, *mongoClientPool.MongoClient) {
   	var err error
   	mp := mongoClientPool.APP.Pool.Once()
   
   	mc, err := mongoClientPool.APP.Client.New("mongodb://admin:admin@localhost:27017")
   	if err != nil {
   		panic(Errorf("创建 mongo 客户端失败：%v", err))
   	}
   	if _, err = mp.AppendClient("default", mc); err != nil {
   		panic(Errorf("添加mongo客户端失败：%v", err))
   	}
   	mc = mp.GetClient("default").SetDatabase("test_db").SetCollection("test_collection")
   
   	return mp, mc
   }
   
   func main() {
   	var (
   		deleteManyRes *mongo.DeleteResult
   		mp, mc        = getDB()
   	)
   
   	// 删除多条数据
   	if mc.Where(mongoClientPool.Map{"name": mongoClientPool.Map{"$ne": "张三"}}).DeleteMany(&deleteManyRes).Err != nil {
   		panic(Errorf("删除多条数据失败：%v", mc.Err))
   	}
   	Printf("删除数据成功：%d\n", deleteManyRes.DeletedCount)
   
   	mp.Clean()
   }
   ```