package redisPool

import (
	"context"
	"errors"
	"fmt"
	"time"

	rds "github.com/redis/go-redis/v9"

	"github.com/aid297/aid/v2/anyMap"
	"github.com/aid297/aid/v2/anySlice"
	"github.com/aid297/aid/v2/str"
)

type (
	RedisPool struct {
		redisClients anyMap.AnyMapper[string, *redisClient]
		addr         string
		password     string
		prefix       string
		pools        anySlice.AnySlicer[redisPoolSetting]
	}

	redisClient struct {
		prefix string
		conn   *rds.Client
	}

	redisPoolSetting struct {
		ClientName string
		Prefix     string
		DBNum      int
	}
)

// NewRedisPool 实例化：redis连接池
func NewRedisPool(attrs ...RedisPoolAttr) *RedisPool {
	var redisPool = &RedisPool{
		redisClients: anyMap.New[string, *redisClient](),
		pools:        anySlice.New[redisPoolSetting](),
	}

	if len(attrs) > 0 {
		redisPool.SetAttrs(attrs...)
	}

	if redisPool.pools.NotEmpty() {
		for _, redisPoolSetting := range redisPool.pools.ToSlice() {
			redisPool.redisClients.SetDatum(redisPoolSetting.ClientName, &redisClient{
				prefix: str.APP.Buffer.JoinString(redisPool.prefix, ":", redisPoolSetting.Prefix),
				conn: rds.NewClient(&rds.Options{
					Addr:     redisPool.addr,
					Password: redisPool.password,
					DB:       redisPoolSetting.DBNum,
				}),
			})
		}
	}

	return redisPool
}

// SetAttrs 设置属性
func (my *RedisPool) SetAttrs(attrs ...RedisPoolAttr) {
	for _, attr := range attrs {
		attr(my)
	}
}

func (my *RedisPool) SetAddr(addr string)         { my.addr = addr }
func (my *RedisPool) SetPassword(password string) { my.password = password }
func (my *RedisPool) SetPrefix(prefix string)     { my.prefix = prefix }
func (my *RedisPool) SetPool(clientName, prefix string, dbNum int) {
	my.pools.Append(redisPoolSetting{ClientName: clientName, Prefix: prefix, DBNum: dbNum})
}

// GetClient 获取链接和链接前缀
func (my *RedisPool) GetClient(clientName string) (string, *rds.Client) {
	if client, exist := my.redisClients.GetValueByKey(clientName); exist {
		return client.prefix, client.conn
	}

	return "", nil
}

// Get 获取值
func (my *RedisPool) Get(ctx context.Context, clientName, key string) (string, error) {
	var (
		err         error
		prefix, ret string
		client      *rds.Client
	)

	prefix, client = my.GetClient(clientName)
	if client == nil {
		return "", fmt.Errorf("没有找到redis链接：%s", clientName)
	}

	if ret, err = client.Get(ctx, fmt.Sprintf("%s:%s", prefix, key)).Result(); err != nil {
		if errors.Is(err, rds.Nil) {
			return "", nil
		} else {
			return "", err
		}
	}

	return ret, nil
}

// Set 设置值
func (my *RedisPool) Set(ctx context.Context, clientName, key string, val any, exp time.Duration) (string, error) {
	var (
		prefix string
		client *rds.Client
	)

	prefix, client = my.GetClient(clientName)
	if client == nil {
		return "", fmt.Errorf("没有找到redis链接：%s", clientName)
	}

	return client.Set(ctx, fmt.Sprintf("%s:%s", prefix, key), val, exp).Result()
}

// GetPipe 获取管道
func (my *RedisPool) GetPipe(clientName string) (prefix string, pipeliner rds.Pipeliner, err error) {
	var client *rds.Client

	prefix, client = my.GetClient(clientName)
	if client == nil {
		return "", nil, fmt.Errorf("没有找到redis链接：%s", clientName)
	}

	return
}

// Pipe 执行管道操作
func (my *RedisPool) Pipe(clientName string, fn func(prefix string, pipe rds.Pipeliner) error) error {
	prefix, pipeliner, err := my.GetPipe(clientName)
	if err != nil {
		return err
	}

	return fn(prefix, pipeliner)
}

// Close 关闭链接
func (my *RedisPool) Close(key string) (err error) {
	if client, exist := my.redisClients.GetValueByKey(key); exist {
		return client.conn.Close()
	}

	return
}

// Clean 清理链接
func (my *RedisPool) Clean() {
	for key, val := range my.redisClients.ToMap() {
		_ = val.conn.Close()
		my.redisClients.RemoveByKey(key)
	}
}
