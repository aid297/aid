package websocketPool

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/aid297/aid/str"
)

func ClientDemo() {
	var (
		err error
		wcp *ClientPool
		wci *ClientIns
		wg  sync.WaitGroup
	)

	wcp = OnceClientPool().
		SetOnConnect(func(insName, clientName string) {
			str.NewTerminalLog("[client] 链接成功：%s").Info(clientName)
		}).
		SetOnConnectWrong(func(insName, clientName string, err error) {
			str.NewTerminalLog("[client] 链接失败：%s；%v").Wrong(clientName, err)
		}).
		SetOnReceiveMsgWrong(func(insName, clientName string, bytes []byte, err error) {
			str.NewTerminalLog("[client] 收到消息失败：%s；%v").Wrong(clientName, err)
		}).
		SetOnSendMsgWrong(func(insName, clientName string, err error) {
			str.NewTerminalLog("[client] 发送消息失败：%v").Wrong(err)
		}).
		SetOnCloseClientWrong(func(insName, clientName string, err error) {
			str.NewTerminalLog("[client] 关闭客户端：%s失败；%v").Wrong(clientName, err)
		})

	wci, err = wcp.SetClientIns("timeout")
	if err != nil {
		str.NewTerminalLog("[client] 创建实例失败：%v").Error(err)
	}

	_, err = wci.SetClient("timeout", "127.0.0.1:44444", "", func(insName, clientName string, propertyMessage []byte) ([]byte, error) {
		str.NewTerminalLog("[client timeout] 收到消息：%s；文本消息：%s；原始消息：%v").Info(clientName, string(propertyMessage), propertyMessage)
		return propertyMessage, nil
	}, DefaultHeart(), DefaultMessageTimeout())
	if err != nil {
		str.NewTerminalLog("[client] 创建链接失败：%v").Wrong(err)
	}
	_, err = wcp.SendMsgByName("timeout", "timeout", 1, []byte("abc"))
	if err != nil {
		str.NewTerminalLog("[client] 发送消息失败：%v").Wrong(err)
	}

	_, err = wcp.SetClient(
		"test",
		"01",
		"127.0.0.1:41111",
		"",
		func(insName, clientName string, prototypeMsg []byte) ([]byte, error) {
			str.NewTerminalLog("[client] 收到消息：%s；文本消息：%s；原始消息：%v").Info(clientName, string(prototypeMsg), prototypeMsg)
			return prototypeMsg, nil
		},
		DefaultHeart(),
		DefaultMessageTimeout(),
	)
	if err != nil {
		str.NewTerminalLog("[client] 创建链接失败：%v").Error(err)
	}

	_, err = wcp.SetClient(
		"test",
		"02",
		"127.0.0.1:42222",
		"",
		func(insName, clientName string, prototypeMsg []byte) ([]byte, error) {
			str.NewTerminalLog("[client] 收到消息：%s；文本消息：%s；原始消息：%v").Info(clientName, string(prototypeMsg), prototypeMsg)
			return prototypeMsg, nil
		},
		DefaultHeart(),
		DefaultMessageTimeout(),
	)
	if err != nil {
		str.NewTerminalLog("[client] 创建链接失败：%v").Error(err)
	}

	_, err = wcp.SetClient(
		"test",
		"03",
		"127.0.0.1:43333",
		"",
		func(insName, clientName string, prototypeMsg []byte) ([]byte, error) {
			str.NewTerminalLog("[client] 收到消息：%s；文本消息：%s；原始消息：%v").Info(clientName, string(prototypeMsg), prototypeMsg)
			return prototypeMsg, nil
		},
		DefaultHeart(),
		DefaultMessageTimeout(),
	)
	if err != nil {
		str.NewTerminalLog("[client] 创建链接失败：%v").Error(err)
	}

	for i := 0; i < 3; i++ {
		wg.Add(3)
		for o := 0; o < 3; o++ {
			go func(i, o int) {
				var res []byte

				res, err = wcp.SendMsgByName(
					"test",
					fmt.Sprintf("0%d", o+1),
					1,
					[]byte(fmt.Sprintf("hello world: %d-%d", i+1, o+1)),
				)
				if err != nil {
					str.NewTerminalLog("ERR: %v").Wrong(err)
				}
				str.NewTerminalLog("OK: %s").Info(string(res))
				defer wg.Done()
			}(i, o)
		}
		wg.Wait()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	<-ctx.Done()
}

type ResponseWrt struct{}

func (ResponseWrt) Header() http.Header { return http.Header{} }

func (ResponseWrt) Write([]byte) (int, error) { return 0, nil }

func (ResponseWrt) WriteHeader(statusCode int) {}

func ServerDemo() {
	wsp := APP.ServerPool.
		Once().
		SetOnConnect(func(conn *websocket.Conn) {
			str.NewTerminalLog("[server] 链接成功：%s").Info(conn.RemoteAddr().String())
		}).
		SetOnConnectErr(func(err error) {
			str.NewTerminalLog("[server] 链接失败：%v").Error(err)
		}).
		SetOnSendMsgErr(func(conn *websocket.Conn, err error) {
			str.NewTerminalLog("[server] 发送消息失败：%v").Error(err)
		}).
		SetOnReceiveMsg(func(conn *websocket.Conn, bytes []byte) string {
			return strings.Split(string(bytes), ":")[0]
		})

	wsp.RegisterRouter("ping", func(ws *websocket.Conn) {
		_ = wsp.SendMsgByWsConn(ws, []byte("pong"))
	})

	wsp.Handle(ResponseWrt{}, &http.Request{}, http.Request{}.Header, func() (string, bool) { return "test", true })
}
