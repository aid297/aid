package mongoClientPool

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongoClient struct {
		url               string
		client            *mongo.Client
		CurrentDatabase   *mongo.Database
		CurrentCollection *mongo.Collection
		conditions        []Map
		Err               error
	}

	Data   = primitive.D
	Entity = primitive.E
	Map    = primitive.M
	OID    = primitive.ObjectID
)

// New 实例化：mongo客户端
func (*MongoClient) New(url string) (*MongoClient, error) {
	var (
		err           error
		mc            = &MongoClient{url: url, conditions: []Map{}}
		clientOptions = options.Client().ApplyURI(mc.url)
	)

	// 连接到 MongoDB
	if mc.client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		return nil, err
	}

	// 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mc.client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return mc, nil
}

// Close 关闭客户端
func (my *MongoClient) Close() error { return my.client.Disconnect(context.Background()) }

// GetClient 获取客户端链接
func (my *MongoClient) GetClient() *mongo.Client { return my.client }

// Ping 测试链接
func (my *MongoClient) Ping() error {
	// 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	my.Err = my.client.Ping(ctx, nil)

	return my.Err
}

// SetDatabase 设置数据库
func (my *MongoClient) SetDatabase(database string, opts ...*options.DatabaseOptions) *MongoClient {
	my.CurrentDatabase = my.client.Database(database, opts...)

	return my
}

// SetCollection 设置文档
func (my *MongoClient) SetCollection(collection string, opts ...*options.CollectionOptions) *MongoClient {
	my.CurrentCollection = my.CurrentDatabase.Collection(collection, opts...)

	return my
}

// InsertOne 插入一条数据
func (my *MongoClient) InsertOne(data any, res **mongo.InsertOneResult) *MongoClient {
	*res, my.Err = my.CurrentCollection.InsertOne(context.TODO(), data)

	return my
}

// InsertMany 插入多条数据
func (my *MongoClient) InsertMany(data []any, res **mongo.InsertManyResult) *MongoClient {
	*res, my.Err = my.CurrentCollection.InsertMany(context.TODO(), data)

	return my
}

// UpdateOne 修改一条数据
func (my *MongoClient) UpdateOne(data any, res **mongo.UpdateResult, opts ...*options.UpdateOptions) *MongoClient {
	*res, my.Err = my.CurrentCollection.UpdateOne(context.TODO(), my.GetFirstCondition(), Map{"$set": data}, opts...)

	return my
}

// UpdateMany 修改多条数据
func (my *MongoClient) UpdateMany(data any, res **mongo.UpdateResult, opts ...*options.UpdateOptions) *MongoClient {
	*res, my.Err = my.CurrentCollection.UpdateMany(context.TODO(), my.GetFirstCondition(), Map{"$set": data}, opts...)

	return my
}

// Where 设置查询条件
func (my *MongoClient) Where(condition ...Map) *MongoClient {
	my.CleanConditions()
	my.conditions = append(my.conditions, condition...)

	return my
}

// CleanConditions 清理查询条件
func (my *MongoClient) CleanConditions() { my.conditions = []Map{} }

// GetFirstCondition 获取第一个查询条件（非聚合条件）
func (my *MongoClient) GetFirstCondition() Map {
	if len(my.conditions) > 0 {
		return my.conditions[0]
	} else {
		return nil
	}
}

// GetConditions 获取全部查询条件（聚合条件）
func (my *MongoClient) GetConditions() []Map { return my.conditions }

// FindOne 查询一条数据
func (my *MongoClient) FindOne(result any, findOneOptionFn func(opt *options.FindOneOptions) *options.FindOneOptions) *MongoClient {
	var findOneOption *options.FindOneOptions

	defer my.CleanConditions()

	if findOneOptionFn != nil {
		findOneOption = findOneOptionFn(options.FindOne())
	}

	my.Err = my.CurrentCollection.FindOne(context.TODO(), my.GetFirstCondition(), findOneOption).Decode(result)

	return my
}

// FindMany 查询多条数据
func (my *MongoClient) FindMany(results any, findOptionFn func(opt *options.FindOptions) *options.FindOptions) *MongoClient {
	var (
		findOption *options.FindOptions
		cursor     *mongo.Cursor
	)

	defer my.CleanConditions()

	if findOptionFn != nil {
		findOption = findOptionFn(options.Find())
	}

	cursor, my.Err = my.CurrentCollection.Find(context.TODO(), my.GetFirstCondition(), findOption)

	if my.Err != nil {
		return my
	}

	my.Err = cursor.All(context.TODO(), results)

	return my
}

// Aggregate 聚合查询
func (my *MongoClient) Aggregate(results any) *MongoClient {
	var cursor *mongo.Cursor

	defer my.CleanConditions()

	cursor, my.Err = my.CurrentCollection.Aggregate(context.TODO(), my.GetConditions())
	if my.Err != nil {
		return my
	}

	my.Err = cursor.All(context.TODO(), results)

	return my
}

// DeleteOne 删除单条数据
func (my *MongoClient) DeleteOne(res **mongo.DeleteResult) *MongoClient {
	defer my.CleanConditions()

	if res == nil {
		_, my.Err = my.CurrentCollection.DeleteOne(context.TODO(), my.GetFirstCondition())
	} else {
		*res, my.Err = my.CurrentCollection.DeleteOne(context.TODO(), my.GetFirstCondition())
	}

	return my
}

// DeleteMany 删除多条数据
func (my *MongoClient) DeleteMany(res **mongo.DeleteResult) *MongoClient {
	defer my.CleanConditions()

	if res == nil {
		_, my.Err = my.CurrentCollection.DeleteMany(context.TODO(), my.GetFirstCondition())
	} else {
		*res, my.Err = my.CurrentCollection.DeleteMany(context.TODO(), my.GetFirstCondition())
	}

	return my
}
