package grpcServer

import (
	"sync"

	"github.com/aid297/aid/array/anyArrayV2"
)

type Pool struct {
	pool anyArrayV2.AnyArray[Server]
	mu   *sync.RWMutex
}

var (
	poolOnce sync.Once
	poolIns  Pool
)

func (Pool) Once(servers ...Server) Pool {
	poolOnce.Do(func() {
		poolIns = Pool{mu: &sync.RWMutex{}, pool: anyArrayV2.New(anyArrayV2.List(servers))}
	})
	return poolIns
}

func (my Pool) SetServers(servers ...Server) Pool {
	poolIns.mu.Lock()
	defer poolIns.mu.Unlock()

	if len(servers) > 0 {
		poolIns.pool = anyArrayV2.New(anyArrayV2.List(servers))
	}
	return poolIns
}

func (my Pool) GetServerByName(name string) Server {
	poolIns.mu.RLock()
	defer poolIns.mu.RUnlock()

	for _, server := range poolIns.pool.ToSlice() {
		if server.name == name {
			return server
		}
	}

	return Server{}
}

func (my Pool) AppendSevers(servers ...Server) Pool {
	poolIns.mu.Lock()
	defer poolIns.mu.Unlock()

	poolIns.pool.Append(servers...)
	return poolIns
}
