package websocketPools

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/aid297/aid/v2/anyMaps"
	"github.com/aid297/aid/v2/anySlices"

	"github.com/gorilla/websocket"
)

type (
	// ServerPool websocket 服务端连接池
	ServerPool struct {
		onConnect         func(*websocket.Conn)
		onConnectWrong    func(error)
		onReceiveMsg      func(*websocket.Conn, []byte) string
		onReceiveMsgWrong func(*websocket.Conn, error)
		onRouterWrong     func(*websocket.Conn, error)
		onCloseConnWrong  func(*websocket.Conn, error)
		onSendMsgWrong    func(*websocket.Conn, error)
		onPing            func(*websocket.Conn)
		serverInsList     anyMaps.AnyMapper[string, *ServerIns]
		router            anyMaps.AnyMapper[string, func(ws *websocket.Conn)]
	}

	// ServerIns websocket服务端实例
	ServerIns struct{ Connections anySlices.AnySlicer[*Server] }

	// Server websocket服务端链接
	Server struct {
		done chan struct{}
		Conn *websocket.Conn
	}
)

var (
	serverPoolIns  *ServerPool
	serverPoolOnce sync.Once
	upgrader       = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

// Once 单例化：服务端连接池
func (*ServerPool) Once() *ServerPool {
	serverPoolOnce.Do(func() {
		serverPoolIns = &ServerPool{}
		serverPoolIns.serverInsList = anyMaps.New[string, *ServerIns]()
		serverPoolIns.router = anyMaps.New[string, func(*websocket.Conn)]()
	})

	return serverPoolIns
}

// New 实例化：链接切片
func (*ServerIns) New() *ServerIns {
	return &ServerIns{Connections: anySlices.New[*Server]()}
}

// SetOnConnect 设置回调：链接成功后
func (*ServerPool) SetOnConnect(onConnect func(*websocket.Conn)) *ServerPool {
	serverPoolIns.onConnect = onConnect

	return serverPoolIns
}

// SetOnConnectErr 设置回调：链接失败后
func (*ServerPool) SetOnConnectErr(onConnectWrong func(error)) *ServerPool {
	serverPoolIns.onConnectWrong = onConnectWrong

	return serverPoolIns
}

// SetOnReceiveMsg 设置回调：接收消息
func (*ServerPool) SetOnReceiveMsg(onMessage func(*websocket.Conn, []byte) string) *ServerPool {
	serverPoolIns.onReceiveMsg = onMessage

	return serverPoolIns
}

// SetOnReceiveMsgErr 设置回调：接收消息失败
func (*ServerPool) SetOnReceiveMsgErr(onMessageWrong func(*websocket.Conn, error)) *ServerPool {
	serverPoolIns.onReceiveMsgWrong = onMessageWrong

	return serverPoolIns
}

// SetOnRouterErr 设置回调：路由解析失败
func (*ServerPool) SetOnRouterErr(onRouterWrong func(*websocket.Conn, error)) *ServerPool {
	serverPoolIns.onRouterWrong = onRouterWrong

	return serverPoolIns
}

// SetOnCloseConnErr 设置回调：关闭链接错误
func (*ServerPool) SetOnCloseConnErr(onCloseConnectionWrong func(conn *websocket.Conn, err error)) *ServerPool {
	serverPoolIns.onCloseConnWrong = onCloseConnectionWrong

	return serverPoolIns
}

// SetOnSendMsgErr 设置回调：发送消息失败
func (*ServerPool) SetOnSendMsgErr(onSendMessageWrong func(conn *websocket.Conn, err error)) *ServerPool {
	serverPoolIns.onSendMsgWrong = onSendMessageWrong

	return serverPoolIns
}

// SetOnPing 设置回调：ping
func (*ServerPool) SetOnPing(fn func(*websocket.Conn)) *ServerPool {
	serverPoolIns.onPing = fn

	return serverPoolIns
}

// Handle 消息处理
func (*ServerPool) Handle(
	writer http.ResponseWriter,
	req *http.Request,
	header http.Header,
	condition func() (string, bool),
) {
	var (
		err                  error
		ws                   *websocket.Conn
		message              []byte
		accountOpenId        string
		cond                 bool
		serverInstance, rout any
		wsc                  *ServerIns
		messageType          int
	)

	ws, err = upgrader.Upgrade(writer, req, header)
	if err != nil {
		if serverPoolIns.onConnectWrong != nil {
			serverPoolIns.onConnectWrong(err)
		}
	}

	accountOpenId, cond = condition()
	if cond {
		if serverPoolIns.serverInsList.HasKey(accountOpenId) {
			serverInstance, _ = serverPoolIns.serverInsList.GetValueByKey(accountOpenId)
			serverInstance.(*ServerIns).Connections.Append(&Server{Conn: ws})
		} else {
			wsc = APP.ServerIns.New()
			wsc.Connections.Append(&Server{Conn: ws})
			serverPoolIns.serverInsList.SetDatum(accountOpenId, wsc)
		}

		if serverPoolIns.onConnect != nil {
			serverPoolIns.onConnect(ws)
		}
	}

	for {
		messageType, message, err = ws.ReadMessage()
		if err != nil {
			serverPoolIns.onReceiveMsgWrong(ws, err)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			routerKey := serverPoolIns.onReceiveMsg(ws, message)
			if routerKey != "" {
				if serverPoolIns.router.HasKey(routerKey) {
					rout, _ = serverPoolIns.router.GetValueByKey(routerKey)
					rout.(func(*websocket.Conn))(ws)
				} else {
					if serverPoolIns.onRouterWrong != nil {
						serverPoolIns.onRouterWrong(ws, fmt.Errorf("没有找到路由：%s", routerKey))
					}
				}
			}
		case websocket.BinaryMessage:
			runOnReceiverMsgWrong(ws)
		case websocket.CloseMessage:
			_ = ws.Close()
		case websocket.PingMessage:
			if serverPoolIns.onPing != nil {
				serverPoolIns.onPing(ws)
			}
		default:
			runOnReceiverMsgWrong(ws)
		}
	}
}

func runOnReceiverMsgWrong(ws *websocket.Conn) {
	if serverPoolIns.onReceiveMsgWrong != nil {
		serverPoolIns.onReceiveMsgWrong(ws, errors.New("不支持的消息类型"))
	}
}

// SendMsgByWsConn 通过链接发送消息
func (*ServerPool) SendMsgByWsConn(ws *websocket.Conn, message []byte) error {
	err := ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		if serverPoolIns.onSendMsgWrong != nil {
			serverPoolIns.onSendMsgWrong(ws, fmt.Errorf("发送消息失败：%s ==> %s", err.Error(), ws.RemoteAddr()))
		}
		return fmt.Errorf("发送消息失败：%s ==> %s", err.Error(), ws.RemoteAddr())
	}

	return nil
}

// SendMsgByWsManyConn 通过链接切片发送消息
func (*ServerPool) SendMsgByWsManyConn(servers anySlices.AnySlicer[*Server], message []byte) {
	if servers.Length() > 0 {
		for _, server := range servers.ToSlice() {
			if server != nil {
				err := serverPoolIns.SendMsgByWsConn(server.Conn, message)
				if err != nil {
					if serverPoolIns.onSendMsgWrong != nil {
						serverPoolIns.onSendMsgWrong(server.Conn, err)
					}
				}
			}
		}
	}
}

// SendMsgByAccountOpenId 根据用户openId发送消息
func (*ServerPool) SendMsgByAccountOpenId(accountOpenId string, message []byte) error {
	if serverPoolIns.serverInsList.HasKey(accountOpenId) {
		client, _ := serverPoolIns.serverInsList.GetValueByKey(accountOpenId)
		serverPoolIns.SendMsgByWsManyConn(client.Connections, message)
	}

	return fmt.Errorf("消息接收对象：%s 不存在", accountOpenId)
}

// RegisterRouter 注册路由
func (*ServerPool) RegisterRouter(routerKey string, fn func(ws *websocket.Conn)) *ServerPool {
	if serverPoolIns.router.HasKey(routerKey) {
		serverPoolIns.router.RemoveByKey(routerKey)
	}
	serverPoolIns.router.SetDatum(routerKey, fn)

	return serverPoolIns
}

// Close 关闭连接池
func (*ServerPool) Close() {
	var err error

	serverPoolIns.serverInsList.Each(func(key string, value *ServerIns) {
		value.Connections.Each(func(idx int, item *Server) (isBreak bool) {
			if err = item.Conn.Close(); err != nil {
				if serverPoolIns.onCloseConnWrong != nil {
					serverPoolIns.onCloseConnWrong(item.Conn, err)
				}
				return
			}
			item.done <- struct{}{}

			return
		})
		value.Connections.Clean()
	})
}
