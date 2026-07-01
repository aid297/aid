package rpcServer

import (
	"net"
	"net/rpc"
)

type Server struct {
	rpc      *rpc.Server
	listener net.Listener
}

func (*Server) New(port string) (*Server, error) {
	var (
		err      error
		listener net.Listener
	)

	if listener, err = net.Listen("tcp", port); err != nil {
		return nil, err
	}

	ins := &Server{listener: listener, rpc: rpc.NewServer()}

	return ins, nil
}

func (my *Server) GetRpc() *rpc.Server { return my.rpc }

func (my *Server) Close() error { return my.listener.Close() }

func (my *Server) Launch() { go my.rpc.Accept(my.listener) }
