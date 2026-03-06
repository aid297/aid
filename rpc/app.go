package rpc

import (
	"github.com/aid297/aid/rpc/rpcClient"
	"github.com/aid297/aid/rpc/rpcServer"
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
