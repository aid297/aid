package rpcs

import (
	"github.com/aid297/aid/v2/rpcs/rpcClients"
	"github.com/aid297/aid/v2/rpcs/rpcServers"
)

var APP struct {
	RPCClient struct {
		Client rpcClients.Client
		Pool   rpcClients.Pool
	}
	RPCServer struct {
		Server rpcServers.Server
		Pool   rpcServers.Pool
	}
}
