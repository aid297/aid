package rpcServer

import (
	"sync"

	"github.com/aid297/aid/dict"
)

type Pool struct {
	pool *dict.AnyDict[string, *Server]
	lock sync.RWMutex
	err  error
}

var (
	poolOnce sync.Once
	poolIns  *Pool
)

func (*Pool) Once() *Pool {
	poolOnce.Do(func() { poolIns = &Pool{pool: dict.Make[string, *Server]()} })
	return poolIns
}

func (*Pool) Error() error { return poolIns.err }

func (*Pool) Set(name string, port string) (*Server, error) {
	poolIns.lock.RLock()
	defer poolIns.lock.RUnlock()

	var rpcServer *Server
	if rpcServer, poolIns.err = new(Server).New(port); poolIns.err != nil {
		return nil, poolIns.err
	}
	poolIns.pool.Set(name, rpcServer)

	return rpcServer, nil
}

func (*Pool) Get(key string) *Server {
	poolIns.lock.RLock()
	defer poolIns.lock.RUnlock()

	if val, ok := poolIns.pool.Get(key); ok {
		return val
	}
	return nil
}

func (*Pool) Close(key string) *Pool {
	poolIns.lock.RLock()
	defer poolIns.lock.RUnlock()

	var rpcServer = poolIns.Get(key)

	if rpcServer != nil {
		if poolIns.err = rpcServer.Close(); poolIns.err != nil {
			return poolIns
		}
		poolIns.pool.RemoveByKey(key)
	}

	return poolIns
}

func (*Pool) Clean() []error {
	poolIns.lock.RLock()
	defer poolIns.lock.RUnlock()

	var errs []error

	for key, rpcServer := range poolIns.pool.ToMap() {
		if err := rpcServer.Close(); err != nil {
			errs = append(errs, err)
			continue
		}
		poolIns.pool.RemoveByKey(key)
	}

	return errs
}
