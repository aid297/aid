package grpcServer

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	Error         error
	name          string
	addr          string
	unaryLoggerFn func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	steamLoggerFn func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
	listener      net.Listener
	server        *grpc.Server
	online        bool
	creds         credentials.TransportCredentials
	registerFuncs []func(server *grpc.Server) error
}

func (Server) New(attrs ...Attributer) Server {
	defaultAttrs := []Attributer{Name(":50051"), Addr(":50051"), UnaryLoggerFn(defaultUnaryLogger), SteamLoggerFn(defaultStreamLogger)}
	options := make([]Attributer, 0, len(defaultAttrs)+len(attrs))
	options = append(options, defaultAttrs...)
	options = append(options, attrs...)
	ins := Server{registerFuncs: make([]func(server *grpc.Server) error, 0)}.SetAttrs(options...)
	return ins.Connection()
}

func (my Server) SetAttrs(attrs ...Attributer) Server {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}
	return my
}

// Connection 链接
func (my Server) Connection() Server {
	if my.listener, my.Error = net.Listen("tcp", my.addr); my.listener != nil {
		return my
	}

	my.server = grpc.NewServer(
		grpc.Creds(my.creds), // 演示用；生产请换 TLS
		grpc.UnaryInterceptor(my.unaryLoggerFn),
		grpc.StreamInterceptor(my.steamLoggerFn),
	)

	my.online = true

	for idx := range my.registerFuncs {
		if my.registerFuncs[idx] != nil {
			if my.Error = my.registerFuncs[idx](my.server); my.Error != nil {
				return my
			}
		}
	}

	my.Error = my.server.Serve(my.listener)

	return my
}

// Offline 下线
func (my Server) Offline() Server {
	my.server.Stop()
	my.Error = my.listener.Close()
	if my.Error != nil {
		return my
	}

	my.online = false
	return my
}

func (my Server) GetAddr() string { return my.addr }

func (my Server) GetOnline() bool { return my.online }

func (my Server) GetListener() net.Listener { return my.listener }

func (my Server) GetServer() *grpc.Server { return my.server }

func (my Server) GetCreds() credentials.TransportCredentials { return my.creds }
