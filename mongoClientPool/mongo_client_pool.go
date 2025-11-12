package mongoClientPool

import (
	"errors"
	"sync"

	"github.com/aid297/aid/dict"
)

type MongoClientPool struct {
	clients *dict.AnyDict[string, *MongoClient]
}

var (
	mongoClientPool *MongoClientPool
	mongoPoolOnce   sync.Once
)

// Once 单例化：mongodb连接池
func (*MongoClientPool) Once() *MongoClientPool {
	mongoPoolOnce.Do(func() { mongoClientPool = &MongoClientPool{clients: dict.Make[string, *MongoClient]()} })
	return mongoClientPool
}

// AppendClient 增加客户端
func (*MongoClientPool) AppendClient(key string, mongoClient *MongoClient) (*MongoClientPool, error) {
	if mongoClientPool.clients.HasKey(key) {
		return mongoClientPool, errors.New("客户端已存在")
	}

	mongoClientPool.clients.Set(key, mongoClient)

	return mongoClientPool, nil
}

// HasClient 检查客户端是否存在
func (*MongoClientPool) HasClient(key string) bool { return mongoClientPool.clients.HasKey(key) }

// GetClient 获取客户端
func (*MongoClientPool) GetClient(key string) *MongoClient {
	if mongoClient, exist := mongoClientPool.clients.Get(key); exist {
		return mongoClient
	}

	return nil
}

// 清除客户端
func (*MongoClientPool) Remove(key string) (*MongoClientPool, error) {
	if mongoClient, exist := mongoClientPool.clients.Get(key); !exist {
		return mongoClientPool, errors.New("客户端不存在")
	} else {
		if err := mongoClient.Close(); err != nil {
			return mongoClientPool, err
		}

		mongoClientPool.clients.RemoveByKey(key)
	}

	return mongoClientPool, nil
}

// Clean 清理客户端
func (*MongoClientPool) Clean() *MongoClientPool {
	for _, key := range mongoClientPool.clients.GetKeys().ToSlice() {
		_, _ = mongoClientPool.Remove(key)
	}

	return mongoClientPool
}
