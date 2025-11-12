package rpcClient

import (
	"log"
	"testing"
)

type Args struct{ A, B int }

func Test1(t *testing.T) {
	var (
		err           error
		rpcClientPool = new(Pool).Once()
		rc            *Client
	)
	if rc, err = rpcClientPool.Set("client1", "localhost:9999"); err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	args := Args{A: 7, B: 8}
	var reply int

	if err = rc.Call("Cal.Add", args, &reply); err != nil {
		log.Fatalf("调用失败: %v", err)
	}

	log.Printf("计算结果：%d + %d = %d", args.A, args.B, reply)
}
