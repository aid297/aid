package websockets

import (
	"net/http"
	"time"

	"github.com/aid297/aid/dict"

	"github.com/gorilla/websocket"
)

type (
	Client struct {
		err                             error
		requestHeader                   http.Header
		groupName, name, addr           string
		conn                            *websocket.Conn
		status                          WebsocketConnStatus
		closeChan                       chan struct{}
		receiveMessageChan              chan []byte
		asyncReceiveCallbackDict        *dict.AnyDict[string, clientCallbackFn]
		syncMessageTimeout              time.Duration
		heart                           *time.Ticker
		heartCallback                   clientHeartFn
		onConnSuccessCallback           clientStandardSuccessFn
		onConnFailCallback              clientStandardFailFn
		onCloseSuccessCallback          clientStandardSuccessFn
		onCloseFailCallback             clientStandardFailFn
		onReceiveMessageSuccessCallback clientReceiveMessageSuccessFn
		onReceiveMessageFailCallback    clientStandardFailFn
		onSendMessageFailCallback       clientStandardFailFn
	}
)

var ClientApp Client

// New 实例化：链接
func (*Client) New(
	groupName, name, addr string,
	clientCallbackConfig ClientCallbackConfig,
	options ...any,
) (*Client, error) {
	return NewClient(groupName, name, addr, clientCallbackConfig, options...)
}

// NewClient 实例化：链接
//
//go:fix 推荐使用：New方法
func NewClient(
	groupName, name, addr string,
	clientCallbackConfig ClientCallbackConfig,
	options ...any,
) (client *Client, err error) {
	if name == "" || addr == "" {
		return nil, WebsocketConnOptionErr.New("名称或地址不能为空")
	}

	client = &Client{
		addr:                            addr,
		groupName:                       groupName,
		name:                            name,
		conn:                            &websocket.Conn{},
		status:                          Offline,
		closeChan:                       make(chan struct{}, 1),
		receiveMessageChan:              make(chan []byte, 1),
		asyncReceiveCallbackDict:        dict.Make[string, clientCallbackFn](),
		syncMessageTimeout:              5 * time.Second,
		onConnSuccessCallback:           clientCallbackConfig.OnConnSuccessCallback,
		onConnFailCallback:              clientCallbackConfig.OnConnFailCallback,
		onCloseSuccessCallback:          clientCallbackConfig.OnCloseSuccessCallback,
		onCloseFailCallback:             clientCallbackConfig.OnCloseFailCallback,
		onReceiveMessageSuccessCallback: clientCallbackConfig.OnReceiveMessageSuccessCallback,
		onReceiveMessageFailCallback:    clientCallbackConfig.OnReceiveMessageFailCallback,
		onSendMessageFailCallback:       clientCallbackConfig.OnSendMessageFailCallback,
	}

	if len(options) > 0 {
		for i := 0; i < len(options); i++ {
			if v, ok := options[i].(http.Header); ok {
				client.requestHeader = v
			}
		}
	}

	return
}

// GetStatus 获取链接状态
func (my *Client) GetStatus() WebsocketConnStatus { return my.status }

// GetName 获取链接名称
func (my *Client) GetName() string { return my.name }

// GetAddr 获取链接地址
func (my *Client) GetAddr() string { return my.addr }

// GetAddress 获取链接地址
func (my *Client) GetAddress() string { return my.addr }

// GetConn 获取链接本体
func (my *Client) GetConn() *websocket.Conn { return my.conn }

// GetConnection 获取链接本体
func (my *Client) GetConnection() *websocket.Conn { return my.conn }

// GetReqHdr 获取请求头
func (my *Client) GetReqHdr() http.Header { return my.requestHeader }

// GetRequestHeaders 获取请求头
func (my *Client) GetRequestHeaders() http.Header { return my.requestHeader }

// SetReqHdr 设置请求头
func (my *Client) SetReqHdr(hdr http.Header) *Client { return my.SetRequestHeaders(hdr) }

// SetRequestHeaders 设置请求头
func (my *Client) SetRequestHeaders(header http.Header) *Client {
	my.requestHeader = header
	return my
}

// AppendReqHdr 新增请求头
func (my *Client) AppendReqHdr(hdr http.Header) *Client { return my.AppendRequestHeader(hdr) }

// AppendRequestHeader 新增请求头
func (my *Client) AppendRequestHeader(header http.Header) *Client {
	for k, v := range header {
		my.requestHeader[k] = v
	}

	return my
}

// Boot 启动链接，并打开监听
func (my *Client) Boot() *Client {
	var (
		receiveMessage []byte
		messageType    int
	)

	my.conn, _, my.err = websocket.DefaultDialer.Dial(my.addr, my.requestHeader)
	if my.err != nil {
		if my.onConnFailCallback != nil {
			my.onConnFailCallback(my.groupName, my.name, my.conn, my.err)
		}

		return my
	}

	if my.onConnSuccessCallback != nil {
		my.onConnSuccessCallback(my.groupName, my.name, my.conn)
	}

	// 开启监听
	go func(client *Client) {
		for {
			messageType, receiveMessage, client.err = client.conn.ReadMessage()
			if client.err != nil {
				if client.onReceiveMessageFailCallback != nil {
					client.onReceiveMessageFailCallback(client.groupName, client.name, client.conn, client.err)
				}

				return
			}

			switch messageType {
			case websocket.TextMessage, websocket.BinaryMessage:
				// 解析消息
				message := ParseMessage(receiveMessage)
				client.onReceiveMessageSuccessCallback(client.groupName, client.name, message.GetMessage())
				if message.GetAsync() { // 异步消息
					if callback, ok := client.asyncReceiveCallbackDict.Get(message.GetMessageId()); ok {
						callback(client.groupName, client.name, message.GetMessage())       // 执行异步回调
						client.asyncReceiveCallbackDict.RemoveByKey(message.GetMessageId()) // 删除异步回调
					}
				} else { // 同步消息
					client.receiveMessageChan <- message.GetMessage()
				}
			case websocket.CloseMessage:
				client.Close()
			case websocket.PingMessage:
				_ = my.conn.WriteMessage(websocket.TextMessage, []byte{})
			case websocket.PongMessage:
			}
		}
	}(my)

	my.status = Online

	return my
}

// AsyncMsg 发送消息：异步
func (my *Client) AsyncMsg(msg []byte, fn clientCallbackFn, to time.Duration) *Client {
	return my.AsyncMessage(msg, fn, to)
}

// AsyncMessage 发送消息：异步
func (my *Client) AsyncMessage(message []byte, fn clientCallbackFn, timeout time.Duration) *Client {
	msg := NewMessage(true, message)

	if fn == nil {
		my.err = AsyncMessageCallbackEmptyErr.New("")

		return my
	}

	if timeout <= 0 {
		my.err = AsyncMessageCallbackEmptyErr.New("")

		return my
	}

	my.err = my.conn.WriteMessage(websocket.TextMessage, msg.GetMessage()) // 发送消息
	if my.err != nil {
		if my.onSendMessageFailCallback != nil {
			my.onSendMessageFailCallback(my.groupName, my.name, my.conn, my.err) // 执行发送失败回调

			return my
		}
	}

	my.asyncReceiveCallbackDict.Set(msg.GetMessageId(), fn) // 配置异步回调
	timer := time.After(timeout)                            // 设置超时

	go func(messageId string) {
		<-timer // 超时删除异步回调方法
		my.asyncReceiveCallbackDict.RemoveByKey(messageId)
		my.onSendMessageFailCallback(my.groupName, my.name, my.conn, AsyncMessageTimeoutErr.New("")) // 执行发送消息回调
	}(msg.GetMessageId())

	return my
}

// SyncMsg 发送消息：同步
func (my *Client) SyncMsg(msg []byte, options ...any) ([]byte, error) {
	return my.SyncMessage(msg, options...)
}

// SyncMessage 发送消息：同步
func (my *Client) SyncMessage(message []byte, options ...any) ([]byte, error) {
	var (
		err     error
		timeout = my.syncMessageTimeout
		msg     = NewMessage(false, message)
	)

	if my.conn == nil || my.status == Offline {
		if my.onSendMessageFailCallback != nil {
			my.onSendMessageFailCallback(my.groupName, my.name, my.conn, WebsocketOfflineErr.New(""))
		}

		return nil, WebsocketOfflineErr.New("")
	}

	err = my.conn.WriteMessage(websocket.TextMessage, msg.GetMessage()) // 发送消息
	if err != nil {
		if my.onSendMessageFailCallback != nil {
			my.onSendMessageFailCallback(my.groupName, my.name, my.conn, err)
		}

		return nil, err
	}

	for i := range options {
		if v, ok := options[i].(time.Duration); ok && v > 0 {
			timeout = v
		}
	}

	timeoutTimer := time.After(timeout)

	select {
	case receiveMessage := <-my.receiveMessageChan:
		return receiveMessage, nil
	case <-timeoutTimer:
		if my.onSendMessageFailCallback != nil {
			my.onSendMessageFailCallback(my.groupName, my.name, my.conn, SyncMessageTimeoutErr.New(""))
		}

		return nil, SyncMessageTimeoutErr.New("")
	}
}

// Cls 关闭链接
func (my *Client) Cls() *Client { return my.Close() }

// Close 关闭链接
func (my *Client) Close() *Client {
	if my.conn != nil && my.status == Online {
		my.err = my.conn.Close()
		if my.err != nil {
			if my.onCloseFailCallback != nil {
				my.onCloseFailCallback(my.groupName, my.name, my.conn, my.err)
			}
			my.status = Online
		} else {
			my.conn = nil
			my.status = Offline
			close(my.receiveMessageChan) // 关闭同步消息通道
		}
	} else {
		my.conn = nil
		my.status = Offline
		close(my.receiveMessageChan)
	}

	if my.onCloseSuccessCallback != nil {
		my.onCloseSuccessCallback(my.groupName, my.name, my.conn)
	}

	my.closeChan <- struct{}{}

	return my
}

// Err 获取错误
func (my *Client) Err() error { return my.Error() }

// Error 获取错误
func (my *Client) Error() error {
	err := my.err
	my.err = nil

	return err
}

// Ping 发送ping
func (my *Client) Ping(fn pingFn) *Client {
	if fn != nil {
		my.err = fn(my.conn)
	} else {
		my.err = my.conn.WriteMessage(websocket.TextMessage, []byte(time.Now().String()))
	}

	return my
}

// Heart 设置或重置心跳
func (my *Client) Heart(interval time.Duration, fn clientHeartFn) *Client {
	if interval > 0 {
		if my.heart != nil {
			my.heart.Reset(interval)
		} else {
			my.heart = time.NewTicker(interval)
		}
	}

	if fn != nil {
		my.heartCallback = fn
	}

	if my.heart != nil && my.heartCallback != nil {
		go func(client *Client) {
			for {
				select {
				case <-my.closeChan: // 接收链接关闭信号，停止心跳
					my.heart.Stop()
					return
				case <-client.heart.C:
					my.heartCallback(client.groupName, client.name, client)
				}
			}
		}(my)
	}

	return my
}

// ReBoot 重连
func (my *Client) ReBoot() error {
	{
		my.Close()
		if my.err != nil {
			return my.Error()
		}
	}

	{
		my.Boot()
		if my.err != nil {
			return my.Error()
		}
	}

	return nil
}
