package grpcClient

import (
	"sync"

	"github.com/aid297/aid/array/anyArrayV2"
)

type (
	Pool struct {
		pool anyArrayV2.AnyArray[Client]
		mu   *sync.RWMutex
	}
)

var (
	poolOnce sync.Once
	poolIns  Pool
)

func (Pool) Once(client ...Client) Pool {
	poolOnce.Do(func() {
		poolIns = Pool{mu: &sync.RWMutex{}, pool: anyArrayV2.New(anyArrayV2.List(client))}
	})
	return poolIns
}

func (Pool) SetClients(clients ...Client) Pool {
	poolIns.mu.Lock()
	defer poolIns.mu.Unlock()

	poolIns.pool.Append(clients...)
	return poolIns
}

func (my Pool) GetClientByName(name string) Client {
	poolIns.mu.RLock()
	defer poolIns.mu.RUnlock()

	for _, client := range poolIns.pool.ToSlice() {
		if client.name == name {
			return client
		}
	}

	return Client{}
}

func (my Pool) AppendClients(clients ...Client) Pool {
	poolIns.mu.Lock()
	defer poolIns.mu.Unlock()

	poolIns.pool.Append(clients...)
	return poolIns
}
