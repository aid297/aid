package websockets

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/aid297/aid/dict"
)

type (
	ServerPool struct {
		connections             *dict.AnyDict[string, *Server]
		addrToAuth              *dict.AnyDict[string, string]
		onConnectionFail        serverConnectionFailFn
		onConnectionSuccess     serverConnectionSuccessFn
		onSendMessageSuccess    serverSendMessageSuccessFn
		onSendMessageFail       serverSendMessageFailFn
		onReceiveMessageFail    serverReceiveMessageFailFn
		onReceiveMessageSuccess serverReceiveMessageSuccessFn
		onCloseCallback         serverCloseCallbackFn
	}
)

var (
	ServerPoolApp  ServerPool
	serverPoolOnce sync.Once
	serverPool     *ServerPool
)

// Once 单例化：websocket服务端
func (*ServerPool) Once(attrs ...ServerCallbackAttributer) *ServerPool {
	serverPoolOnce.Do(func() {
		serverPool = &ServerPool{
			connections: dict.Make[string, *Server](),
			addrToAuth:  dict.Make[string, string](),
		}
	})

	if len(attrs) > 0 {
		for _, attr := range attrs {
			attr.Register()
		}
	}

	return serverPool
}

// OnceServer 单例化：websocket服务端
//
//go:fix 推荐使用：Once方法
func OnceServer(serverCallbackConfig ServerCallbackConfig) *ServerPool {
	serverPoolOnce.Do(func() {
		serverPool = &ServerPool{
			connections:             dict.Make[string, *Server](),
			addrToAuth:              dict.Make[string, string](),
			onConnectionFail:        serverCallbackConfig.OnConnectionFail,
			onConnectionSuccess:     serverCallbackConfig.OnConnectionSuccess,
			onSendMessageSuccess:    serverCallbackConfig.OnSendMessageSuccess,
			onSendMessageFail:       serverCallbackConfig.OnSendMessageFail,
			onReceiveMessageFail:    serverCallbackConfig.OnReceiveMessageFail,
			onReceiveMessageSuccess: serverCallbackConfig.OnReceiveMessageSuccess,
			onCloseCallback:         serverCallbackConfig.OnCloseCallback,
		}
	})

	return serverPool
}

// appendConn 增加连接
func (*ServerPool) appendConn(authId *string, conn *websocket.Conn) (server *Server) {
	server = ServerApp.New(conn)
	serverPool.addrToAuth.Set(conn.RemoteAddr().String(), *authId)
	serverPool.connections.Set(conn.RemoteAddr().String(), server)

	return
}

// removeConn 移除连接
func (*ServerPool) removeConn(addr *string) {
	serverPool.addrToAuth.RemoveByKey(*addr)
	serverPool.connections.RemoveByKey(*addr)
}

// SendMessageByAddr 发送消息：通过地址
func (*ServerPool) SendMessageByAddr(addr *string, prototypeMessage []byte) {
	if server, ok := serverPool.connections.Get(*addr); ok {
		server.AsyncMessage(prototypeMessage, serverPool.onSendMessageSuccess, serverPool.onSendMessageFail)
	} else {
		if serverPool.onSendMessageFail != nil {
			serverPool.onSendMessageFail(fmt.Errorf("没有找到连接：%s", *addr))
		}
	}
}

// SendMessageByAuthId 发送消息：通过认证ID
func (*ServerPool) SendMessageByAuthId(authId *string, prototypeMessage []byte) {
	serverPool.connections.GetValuesByKeys(serverPool.addrToAuth.GetKeysByValues(*authId).ToSlice()...).Each(func(idx int, server *Server) {
		server.AsyncMessage(prototypeMessage, serverPool.onSendMessageSuccess, serverPool.onSendMessageFail)
	})
}

// SetOnConnectionSuccess 设置回调：当连接成功
func (*ServerPool) SetOnConnectionSuccess(onConnectionSuccess serverConnectionSuccessFn) *ServerPool {
	serverPool.onConnectionSuccess = onConnectionSuccess
	return serverPool
}

// SetOnConnectionFail 设置回调：当连接失败
func (*ServerPool) SetOnConnectionFail(onConnectionFail serverConnectionFailFn) *ServerPool {
	serverPool.onConnectionFail = onConnectionFail
	return serverPool
}

// SetOnSendMessageSuccess 设置回调：当发送消息成功
func (*ServerPool) SetOnSendMessageSuccess(onSendMessageSuccess serverSendMessageSuccessFn) *ServerPool {
	serverPool.onSendMessageSuccess = onSendMessageSuccess
	return serverPool
}

// SetOnSendMessageFail 设置回调：当发送消息失败
func (*ServerPool) SetOnSendMessageFail(onSendMessageFail serverSendMessageFailFn) *ServerPool {
	serverPool.onSendMessageFail = onSendMessageFail
	return serverPool
}

// SetOnReceiveMessageSuccess 设置回调：当接收消息成功
func (*ServerPool) SetOnReceiveMessageSuccess(onReceiveMessageSuccess serverReceiveMessageSuccessFn) *ServerPool {
	serverPool.onReceiveMessageSuccess = onReceiveMessageSuccess
	return serverPool
}

// SetOnReceiveMessageFail 设置回调：当接收消息失败
func (*ServerPool) SetOnReceiveMessageFail(onReceiveMessageFail serverReceiveMessageFailFn) *ServerPool {
	serverPool.onReceiveMessageFail = onReceiveMessageFail
	return serverPool
}

// SetOnCloseCallback 设置回调：关闭时回调
func (*ServerPool) SetOnCloseCallback(onCloseCallback serverCloseCallbackFn) *ServerPool {
	serverPool.onCloseCallback = onCloseCallback

	return serverPool
}

// Handle 消息处理
func (*ServerPool) Handle(
	writer http.ResponseWriter,
	req *http.Request,
	header http.Header,
	condition serverConnectionCheckFn,
) *ServerPool {
	var (
		err  error
		conn *websocket.Conn
	)

	if condition == nil {
		serverPool.onConnectionFail(errors.New("验证方法不能为空"))
		return serverPool
	}

	// 升级协议
	conn, err = upgrader.Upgrade(writer, req, header)
	if err != nil {
		if serverPool.onConnectionFail != nil {
			serverPool.onConnectionFail(err)
		}
	}

	// 验证连接
	identity, err := condition(header)
	if err != nil && serverPool.onConnectionFail != nil {
		serverPool.onConnectionFail(err)
		return serverPool
	}

	// 加入连接池
	server := serverPool.appendConn(&identity, conn)

	// 开启接收消息
	if err = server.Boot(
		serverPool.onReceiveMessageSuccess,
		serverPool.onReceiveMessageFail,
		serverPool.onSendMessageFail,
		serverPool.onCloseCallback,
	); err != nil {
		if serverPool.onConnectionFail != nil {
			serverPool.onConnectionFail(err)
		}

		server.Close()
		serverPool.removeConn(&server.addr)
		server = nil
	}

	return serverPool
}
