package grpcServer

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	Attributer interface {
		Register(server *Server)
	}

	AttrName  struct{ name string }
	AttrAddr  struct{ addr string }
	AttrCreds struct {
		crt, key string
		creds    credentials.TransportCredentials
	}
	AttrUnaryLoggerFn struct {
		fn func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	}
	AttrSteamLoggerFn struct {
		fn func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
	}
	AttrRegisterFn struct {
		funcs []func(server *grpc.Server) error
	}
)

// 简单日志拦截器（Unary）
func defaultUnaryLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("[UNARY] %s took=%s err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

// 双向流日志拦截器（Stream）
func defaultStreamLogger(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()
	err := handler(srv, ss)
	log.Printf("[STREAM] %s took=%s err=%v", info.FullMethod, time.Since(start), err)
	return err
}

func Name(name string) AttrName { return AttrName{name: name} }

func (my AttrName) Register(server *Server) { server.name = my.name }

func Addr(addr string) AttrAddr { return AttrAddr{addr: addr} }

func (my AttrAddr) Register(server *Server) { server.addr = my.addr }

func UnaryLoggerFn(fn func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)) AttrUnaryLoggerFn {
	return AttrUnaryLoggerFn{fn: fn}
}

func (my AttrUnaryLoggerFn) Register(server *Server) { server.unaryLoggerFn = my.fn }

func SteamLoggerFn(fn func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error) AttrSteamLoggerFn {
	return AttrSteamLoggerFn{fn: fn}
}

func (my AttrSteamLoggerFn) Register(server *Server) { server.steamLoggerFn = my.fn }

func Creds(crt, key string) AttrCreds { return AttrCreds{crt: crt, key: key} }

func (my AttrCreds) Register(server *Server) {
	if my.creds == nil && my.crt != "" && my.key != "" {
		my.creds, server.Error = credentials.NewServerTLSFromFile(my.crt, my.key)
	} else {
		my.creds = insecure.NewCredentials()
	}
	server.creds = my.creds
}

func RegisterFuncs(funcs ...func(server *grpc.Server) error) AttrRegisterFn {
	return AttrRegisterFn{funcs: funcs}
}

func (my AttrRegisterFn) Register(server *Server) {
	if len(server.registerFuncs) == 0 {
		server.registerFuncs = make([]func(server *grpc.Server) error, 0)
	} else {
		server.registerFuncs = append(server.registerFuncs, my.funcs...)
	}
}
