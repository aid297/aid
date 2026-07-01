package rpcServer

import "testing"

func Test2(t *testing.T) {
	var (
		err           error
		rpcServer     *Server
		rpcServerPool *Pool
	)

	rpcServerPool = new(Pool).Once()

	if rpcServer, err = rpcServerPool.Set("server1", ":9999"); err != nil {
		t.Fatalf("连接失败：%s", err.Error())
	}

	rpcServer.GetRpc().Register(new(Cal))

	rpcServer.Launch()

	t.Log("服务已启动")
}
