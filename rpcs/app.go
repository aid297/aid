package rpcs

import (
	"github.com/aid297/aid/v2/rpc/rpcClient"
	"github.com/aid297/aid/v2/rpc/rpcServer"
)

var APP struct {
	RPCClient struct {
		Client rpcClient.Client
		Pool   rpcClient.Pool
	}
	RPCServer struct {
		Server rpcServer.Server
		Pool   rpcServer.Pool
	}
}
