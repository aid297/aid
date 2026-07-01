package rpcServer

import (
	"testing"
)

type (
	Args struct{ A, B int }
	Cal  struct{}
)

func (*Cal) Add(args Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func Test1(t *testing.T) {
	var (
		err        error
		rpcService *Server
	)

	if rpcService, err = new(Server).New(":9999"); err != nil {
		t.Fatalf("连接失败：%s", err.Error())
	}
	defer rpcService.Close()

	{
		rpcService.GetRpc().Register(new(Cal))
		rpcService.Launch()
	}

	t.Log("服务已启动")
}
