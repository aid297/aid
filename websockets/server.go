package websockets

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type (
	Server struct {
		addr               string
		conn               *websocket.Conn
		closeChan          chan struct{}
		receiveMessageChan chan []byte
		status             WebsocketConnStatus
	}

	ServerReceiveMessage struct {
		Target string `json:"target"`
	}
)

var ServerApp Server

func (*Server) New(conn *websocket.Conn) *Server { return NewServer(conn) }

// NewServer 实例话：websocket服务端
//
//go:fix 推荐使用：推荐使用New方法
func NewServer(conn *websocket.Conn) *Server {
	return &Server{
		addr:               conn.RemoteAddr().String(),
		conn:               conn,
		closeChan:          make(chan struct{}, 1),
		receiveMessageChan: make(chan []byte, 1),
		status:             Offline,
	}
}

// IsOnline 是否在线
func (my *Server) IsOnline() bool { return my.status == Online }

// IsOffline 是否离线
func (my *Server) IsOffline() bool { return my.status == Offline }

// Conn 获取链接
func (my *Server) Conn() *websocket.Conn { return my.conn }

// SyncMessage 发送消息：同步
func (my *Server) SyncMessage(prototypeMessage []byte, onSuccess serverSendMessageSuccessFn, onFail serverSendMessageFailFn) {
	if my.IsOffline() && onFail != nil {
		onFail(fmt.Errorf("发送失败：连接离线：%s -> %s", my.addr, prototypeMessage))
		return
	}

	message := NewMessage(false, prototypeMessage)
	if err := my.conn.WriteMessage(websocket.TextMessage, message.GetMessage()); err != nil && onFail != nil {
		onFail(fmt.Errorf("发送失败：%s [%s -> %s] %s", err.Error(), my.addr, message.GetMessage(), message.GetPrototypeMessage()))
		return
	}
	if onSuccess != nil {
		onSuccess(my.conn, message.GetMessage(), message.GetPrototypeMessage())
	}
}

// AsyncMessage 发送消息：异步
func (my *Server) AsyncMessage(prototypeMessage []byte, onSuccess serverSendMessageSuccessFn, onFail serverSendMessageFailFn) {
	if my.IsOffline() && onFail != nil {
		onFail(fmt.Errorf("发送失败：连接离线：%s -> %s", my.addr, prototypeMessage))
		return
	}

	message := MessageApp.New(true, prototypeMessage)
	if err := my.conn.WriteMessage(websocket.TextMessage, message.GetMessage()); err != nil && onFail != nil {
		onFail(fmt.Errorf("发送失败：%s [%s -> %s] %s", err.Error(), my.addr, message.GetMessage(), message.GetPrototypeMessage()))
		return
	}
	if onSuccess != nil {
		onSuccess(my.conn, message.GetMessage(), message.GetPrototypeMessage())
	}
}

// Close 关闭
func (my *Server) Close() *Server {
	my.closeChan <- struct{}{}
	return my
}

// Boot 启动
func (my *Server) Boot(
	onReceiveMessageSuccess serverReceiveMessageSuccessFn,
	onReceiveMessageFail serverReceiveMessageFailFn,
	onSendMessageFail serverSendMessageFailFn,
	onCloseCallback serverCloseCallbackFn,
) error {
	if onReceiveMessageSuccess == nil {
		return errors.New("解析消息函数不能为空：onReceiveMessageSuccess")
	}

	go func(
		onReceiveMessageSuccess serverReceiveMessageSuccessFn,
		onReceiveMessageFail serverReceiveMessageFailFn,
		onSendMessageFail serverSendMessageFailFn,
		onCloseCallback serverCloseCallbackFn,
	) {
		defer my.conn.Close() // 确保 goroutine 结束时关闭连接
		my.status = Online

		for {
			select {
			case <-my.closeChan:
				close(my.closeChan)
				my.status = Offline
				if onCloseCallback != nil {
					onCloseCallback(my.conn)
				}
				return
			default:
				messageType, prototypeMessage, err := my.conn.ReadMessage()
				if err != nil {
					if err.Error() == "websocket: close 1006 (abnormal closure): unexpected EOF" {
						close(my.closeChan)
						my.status = Offline
						if onCloseCallback != nil {
							onCloseCallback(my.conn)
						}
						return
					}
					if onReceiveMessageFail != nil {
						onReceiveMessageFail(my.conn, err)
					}
					break
				}

				switch messageType {
				case websocket.TextMessage:
					message := ParseMessage(prototypeMessage)
					go onReceiveMessageSuccess(my, message)
				case websocket.BinaryMessage:
				case websocket.CloseMessage:
					return
				case websocket.PingMessage:
					if err = my.conn.WriteMessage(websocket.TextMessage, []byte{}); err != nil {
						if onSendMessageFail != nil {
							onSendMessageFail(fmt.Errorf("发送消息失败(pong)：%s", my.conn.RemoteAddr().String()))
						}
					}
				case websocket.PongMessage:
				default:
					if onReceiveMessageFail != nil {
						onReceiveMessageFail(my.conn, fmt.Errorf("不支持的消息类型：%d", messageType))
					}
				}
			}
		}
	}(
		onReceiveMessageSuccess,
		onReceiveMessageFail,
		onSendMessageFail,
		onCloseCallback,
	)

	return nil
}
