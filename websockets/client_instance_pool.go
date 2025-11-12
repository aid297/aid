package websockets

import (
	"sync"

	"github.com/aid297/aid/dict"
)

type ClientInstancePool struct {
	pool *dict.AnyDict[string, *ClientInstance]
}

var (
	ClientInstancePoolApp  ClientInstancePool
	clientInstancePoolOnce sync.Once
	clientInstancePool     *ClientInstancePool
)

func (*Client) Once() *ClientInstancePool { return OnceClientInstancePool() }

// OnceClientInstancePool 单例化：websocket客户端实例池
//
//go:fix 推荐使用：Once方法
func OnceClientInstancePool() *ClientInstancePool {
	clientInstancePoolOnce.Do(func() { clientInstancePool = &ClientInstancePool{pool: dict.Make[string, *ClientInstance]()} })

	return clientInstancePool
}

// Append 增加客户端
func (*ClientInstancePool) Append(clientInstance *ClientInstance) error {
	if clientInstance.connections.HasKey(clientInstance.name) {
		return WebsocketClientExistErr.New(clientInstance.name)
	}

	clientInstancePool.pool.Set(clientInstance.name, clientInstance)

	return nil
}

// Remove 删除客户端
func (*ClientInstancePool) Remove(name string) error {
	if !clientInstancePool.pool.HasKey(name) {
		return WebsocketClientNotExistErr.New(name)
	}

	clientInstancePool.pool.RemoveByKey(name)

	return nil
}

// Get 获取客户端
func (*ClientInstancePool) Get(name string) (*ClientInstance, error) {
	if clientInstance, exists := clientInstancePool.pool.Get(name); !exists {
		return nil, WebsocketClientNotExistErr.New(name)
	} else {
		return clientInstance, nil
	}
}

// Has 检查客户端是否存在
func (*ClientInstancePool) Has(name string) bool { return clientInstancePool.pool.HasKey(name) }

// Close 关闭客户端
func (*ClientInstancePool) Close(name string) error {
	if clientInstance, err := clientInstancePool.Get(name); err != nil {
		return err
	} else {
		err = clientInstance.Close(name)
		clientInstancePool.pool.RemoveByKey(clientInstance.name)

		return err
	}
}

// Clean 清空客户端实例
func (*ClientInstancePool) Clean() []error {
	var errorList []error
	clientInstancePool.pool.Each(func(key string, clientInstance *ClientInstance) {
		err := clientInstance.Clean()
		if len(err) > 0 {
			errorList = append(errorList, err...)
		} else {
			clientInstance.connections.RemoveByKey(clientInstance.name)
		}
	})

	return errorList
}
