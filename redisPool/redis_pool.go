package redisPool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	rds "github.com/redis/go-redis/v9"

	"github.com/aid297/aid/dict"
	"github.com/aid297/aid/str"
)

type (
	RedisPool struct {
		connections *dict.AnyDict[string, *redisConn]
	}

	redisConn struct {
		prefix string
		conn   *rds.Client
	}
)

var (
	redisPoolIns  *RedisPool
	redisPoolOnce sync.Once
)

func (*RedisPool) Once(redisSetting *RedisSetting) *RedisPool {
	redisPoolOnce.Do(func() {
		redisPoolIns = &RedisPool{}
		redisPoolIns.connections = dict.Make[string, *redisConn]()

		if len(redisSetting.Pool) > 0 {
			for idx := range redisSetting.Pool {
				redisPoolIns.connections.Set(redisSetting.Pool[idx].Key, &redisConn{
					prefix: str.BufferApp.NewString(redisSetting.Prefix).S(":", redisSetting.Pool[idx].Prefix).String(),
					conn: rds.NewClient(&rds.Options{
						Addr:     str.BufferApp.NewAny(redisSetting.Host).Any(":", redisSetting.Port).String(),
						Password: redisSetting.Password,
						DB:       redisSetting.Pool[idx].DbNum,
					}),
				})
			}
		}
	})

	return redisPoolIns
}

// GetClient 获取链接和链接前缀
func (*RedisPool) GetClient(key string) (string, *rds.Client) {
	if client, exist := redisPoolIns.connections.Get(key); exist {
		return client.prefix, client.conn
	}

	return "", nil
}

// Get 获取值
func (*RedisPool) Get(clientName, key string) (string, error) {
	var (
		err         error
		prefix, ret string
		client      *rds.Client
	)

	prefix, client = redisPoolIns.GetClient(clientName)
	if client == nil {
		return "", fmt.Errorf("没有找到redis链接：%s", clientName)
	}

	ret, err = client.Get(context.Background(), fmt.Sprintf("%s:%s", prefix, key)).Result()
	if err != nil {
		if errors.Is(err, rds.Nil) {
			return "", nil
		} else {
			return "", err
		}
	}

	return ret, nil
}

// Set 设置值
func (*RedisPool) Set(clientName, key string, val any, exp time.Duration) (string, error) {
	var (
		prefix string
		client *rds.Client
	)

	prefix, client = redisPoolIns.GetClient(clientName)
	if client == nil {
		return "", fmt.Errorf("没有找到redis链接：%s", clientName)
	}

	return client.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, key), val, exp).Result()
}

// Close 关闭链接
func (my *RedisPool) Close(key string) error {
	if client, exist := redisPoolIns.connections.Get(key); exist {
		return client.conn.Close()
	}

	return nil
}

// Clean 清理链接
func (*RedisPool) Clean() {
	for key, val := range redisPoolIns.connections.ToMap() {
		_ = val.conn.Close()
		redisPoolIns.connections.RemoveByKey(key)
	}
}
