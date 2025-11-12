package grpcClient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	Attributer interface {
		Register(client *Client)
	}

	AttrName  struct{ name string }
	AttrAddr  struct{ addr string }
	AttrCreds struct {
		crt, key string
		creds    credentials.TransportCredentials
	}
	AttrRegisterFuncs struct {
		funcs []func(conn *grpc.ClientConn) error
	}
)

func Name(name string) AttrName { return AttrName{name: name} }

func (my AttrName) Register(client *Client) { client.name = my.name }

func Addr(addr string) AttrAddr { return AttrAddr{addr: addr} }

func (my AttrAddr) Register(client *Client) { client.addr = my.addr }

func RegisterFuncs(funcs ...func(conn *grpc.ClientConn) error) AttrRegisterFuncs {
	return AttrRegisterFuncs{funcs: funcs}
}

func (my AttrRegisterFuncs) Register(client *Client) {
	if len(client.registerFuncs) == 0 {
		client.registerFuncs = make([]func(conn *grpc.ClientConn) error, 0)
	} else {
		client.registerFuncs = append(client.registerFuncs, my.funcs...)
	}
}
