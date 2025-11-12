package websocketPool

import (
	"log"
	"testing"

	"github.com/gorilla/websocket"
)

func online() {
	client, err := NewClient(
		"client-test1",
		"client-test1",
		"127.0.0.1:8080",
		"ws",
		func(instanceName, clientName string, prototypeMsg []byte) ([]byte, error) {
			log.Printf("收到消息[%s:%s]：%s", instanceName, clientName, prototypeMsg)
			return prototypeMsg, nil
		},
	)
	if err != nil {
		log.Fatalf("链接失败：%v", err)
	}

	_, _ = client.SendMsg(websocket.TextMessage, []byte("hello world"))
}

func Test1(t *testing.T) {
	online()
}
