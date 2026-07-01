package redisPools

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	rds "github.com/redis/go-redis/v9"

	"github.com/aid297/aid/v2/anySlices"
)

type (
	RedisPool struct {
		redisClients sync.Map
		addr         string
		username     string
		password     string
		prefix       string
		pools        anySlices.AnySlicer[redisPoolSetting]
	}

	redisClient struct {
		prefix string
		conn   *rds.Client
	}

	redisPoolSetting struct {
		ClientName string
		DBNum      int
	}
)

// New 实例化：redis连接池
func New(addr, appName string, attrs ...RedisPoolAttr) *RedisPool {
	var ins = &RedisPool{
		addr:         addr,
		prefix:       appName,
		redisClients: sync.Map{},
		pools:        anySlices.New[redisPoolSetting](),
	}

	if len(attrs) > 0 {
		ins.SetAttrs(attrs...)
	}

	if ins.pools.NotEmpty() {
		ins.pools.Each(func(_ int, item redisPoolSetting) (isBreak bool) {
			ins.redisClients.Store(item.ClientName, &redisClient{
				prefix: fmt.Sprintf("%s:%s", ins.prefix, item.ClientName),
				conn:   rds.NewClient(&rds.Options{Addr: addr, Username: ins.username, Password: ins.password, DB: item.DBNum}),
			})
			return
		})
	}

	return ins
}

// SetAttrs 设置属性
func (my *RedisPool) SetAttrs(attrs ...RedisPoolAttr) {
	for _, attr := range attrs {
		attr(my)
	}
}

func (my *RedisPool) SetAddr(addr string)         { my.addr = addr }
func (my *RedisPool) SetUsername(username string) { my.username = username }
func (my *RedisPool) SetPassword(password string) { my.password = password }
func (my *RedisPool) SetPrefix(prefix string)     { my.prefix = prefix }
func (my *RedisPool) SetPool(clientName string, dbNum int) {
	my.pools.Append(redisPoolSetting{ClientName: clientName, DBNum: dbNum})
}

// GetClient 获取链接和链接前缀
func (my *RedisPool) GetClient(clientName string) (string, *rds.Client) {
	if value, exist := my.redisClients.Load(clientName); exist {
		return value.(*redisClient).prefix, value.(*redisClient).conn
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

	pipeliner = client.TxPipeline()

	return
}

// Pipe 执行管道操作
func (my *RedisPool) Pipe(clientName string, fn func(prefix string, pipe rds.Pipeliner) (err error)) error {
	prefix, pipeliner, err := my.GetPipe(clientName)
	if err != nil {
		return err
	}

	return fn(prefix, pipeliner)
}

// Close 关闭链接
func (my *RedisPool) Close(key string) (err error) {
	if client, exist := my.redisClients.Load(key); exist {
		return client.(*redisClient).conn.Close()
	}

	return
}

// Clean 清理链接
func (my *RedisPool) Clean() {
	my.redisClients.Range(func(key, value any) bool {
		_ = value.(*redisClient).conn.Close()
		my.redisClients.Delete(key)
		return true
	})
}
