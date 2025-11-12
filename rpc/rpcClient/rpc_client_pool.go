package rpcClient

import (
	"sync"

	"github.com/aid297/aid/dict"
)

type Pool struct {
	err  error
	pool *dict.AnyDict[string, *Client]
	lock sync.RWMutex
}

var (
	poolOnce sync.Once
	poolIns  *Pool
)

func (*Pool) Once() *Pool {
	poolOnce.Do(func() { poolIns = &Pool{pool: dict.Make[string, *Client]()} })
	return poolIns
}
func (*Pool) Error() error { return poolIns.err }

func (*Pool) Set(name string, addr string) (*Client, error) {
	poolIns.lock.RLock()
	defer poolIns.lock.RUnlock()

	var rpcClient *Client
	if rpcClient, poolIns.err = new(Client).New(addr); poolIns.err != nil {
		return nil, poolIns.err
	}
	poolIns.pool.Set(name, rpcClient)

	return rpcClient, nil
}

func (*Pool) Get(name string) *Client {
	poolIns.lock.RLock()
	defer poolIns.lock.RUnlock()

	if rpcClient, ok := poolIns.pool.Get(name); ok {
		return rpcClient
	}
	return nil
}

func (*Pool) Close(key string) *Pool {
	poolIns.lock.Lock()
	defer poolIns.lock.Unlock()

	var rpcClient = poolIns.Get(key)

	if rpcClient != nil {
		if poolIns.err = rpcClient.Close(); poolIns.err != nil {
			return poolIns
		}
		poolIns.pool.RemoveByKey(key)
	}

	return poolIns
}

func (*Pool) Clean() []error {
	poolIns.lock.Lock()
	defer poolIns.lock.Unlock()

	var errs []error
	for key, rpcClient := range poolIns.pool.ToMap() {
		if err := rpcClient.Close(); err != nil {
			errs = append(errs, err)
			continue
		}
		poolIns.pool.RemoveByKey(key)
	}

	return errs
}
