package grpc

import (
	"github.com/aid297/aid/v2/grpc/grpcClient"
	"github.com/aid297/aid/v2/grpc/grpcServer"
)

var APP struct {
	Client struct {
		grpcClient.Client
		grpcClient.Pool
	}
	Server struct {
		grpcServer.Server
		grpcServer.Pool
	}
}
