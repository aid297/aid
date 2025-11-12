package websockets

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestServer(t *testing.T) {
	serverPool := ServerPoolApp.Once(
		OnConnectionSuccess(func(conn *websocket.Conn) error {
			t.Logf("连接成功：%s\n", conn.RemoteAddr().String())
			return nil
		}),
		OnConnectionFail(func(err error) { t.Logf("连接失败：%v\n", err) }),
		OnReceiveMessageSuccess(func(server *Server, message Message) { t.Logf("接收消息成功：%s\n", message.GetMessage()) }),
		OnReceiveMessageFail(func(conn *websocket.Conn, err error) {
			t.Logf("接收消息失败：%s -> %v\n", conn.RemoteAddr().String(), err)
		}),
		OnSendMessageSuccess(func(conn *websocket.Conn, message, prototypeMessage []byte) {
			t.Logf("发送消息成功：%s -> %s\n", conn.RemoteAddr().String(), string(message))
		}),
		OnSendMessageFail(func(err error) { t.Logf("发送消息失败：%v\n", err) }),
		OnCloseCallback(func(conn *websocket.Conn) { t.Logf("关闭连接：%s\n", conn.RemoteAddr().String()) }),
	)

	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		serverPool.Handle(
			c.Writer,
			c.Request,
			c.Request.Header,
			func(header http.Header) (string, error) {
				t.Logf("header: %#v\n", header)
				authId := header.Get("Identity")
				return authId, nil
			},
		)
	})
	if err := r.Run(); err != nil {
		t.Fatalf("启动失败：%v\n", err)
	}

	timer := time.NewTimer(30 * time.Second)
	<-timer.C
	timer.Stop()
	t.Log("测试结束")
}
