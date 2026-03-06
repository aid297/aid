package websocketPool

import (
	"errors"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type (
	// Client websocket 客户端链接
	Client struct {
		url           url.URL
		insName, Name string
		Conn          *websocket.Conn
		mu            sync.Mutex    // 同步锁
		closeChan     chan struct{} // 关闭信号
		syncChan      chan []byte   // 同步消息
		onReceiveMsg  func(insName, clientName string, prototypeMsg []byte) ([]byte, error)
		heart         *Heart
		timeout       *MessageTimeout
	}

	// PendingRequest 待处理请求
	PendingRequest struct {
		Uuid uuid.UUID
		Chan chan []byte
		Done chan struct{}
		Err  error
	}
)

func (*Client) New(
	insName, name, host, path string,
	receiveMessageFunc func(insName, clientName string, prototypeMsg []byte) ([]byte, error),
) (*Client, error) {
	return NewClient(insName, name, host, path, receiveMessageFunc)
}

// NewClient 实例化：websocket 客户端链接
//
//go:fix 推荐使用：推荐使用New方法
func NewClient(
	insName, name, host, path string,
	receiveMessageFunc func(insName, clientName string, prototypeMsg []byte) ([]byte, error),
) (*Client, error) {
	var (
		err    error
		conn   *websocket.Conn
		client *Client = &Client{
			insName:      insName,
			Name:         name,
			url:          url.URL{},
			Conn:         nil,
			onReceiveMsg: receiveMessageFunc,
		}
	)
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:38000", Path: "/cbit_db/ws"}
	log.Printf("真实连接：%s", u.String())

	// 建立连接
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	client.Conn = conn
	client.syncChan = make(chan []byte, 1)
	client.closeChan = make(chan struct{}, 1)

	return client, nil
}

// SendMsg 发送消息：通过链接
func (my *Client) SendMsg(msgType int, msg []byte) ([]byte, error) {
	var (
		err error
		res []byte
	)

	if my.timeout == nil || my.timeout.interval == 0 {
		clientPoolIns.Error = errors.New("同步消息，需要设置超时时间")
		return nil, errors.New("同步消息，需要设置超时时间")
	}

	my.mu.Lock()
	defer my.mu.Unlock()

	err = my.Conn.WriteMessage(msgType, msg)
	if err != nil {
		if clientPoolIns.onSendMsgWrong != nil {
			clientPoolIns.onSendMsgWrong(my.insName, my.Name, err)
		}
		clientPoolIns.Error = err
		return nil, err
	}

	timer := time.After(my.timeout.interval)
	select {
	case <-timer:
		clientPoolIns.Error = errors.New("请求超时")
		return nil, errors.New("请求超时")
	case res = <-my.syncChan:
		if my.onReceiveMsg != nil {
			return my.onReceiveMsg(my.insName, my.Name, res)
		}
		return res, nil
	}
}

// Close 关闭链接
func (my *Client) Close() error {
	var err error

	// 发送关闭消息
	err = my.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	if err = my.Conn.Close(); err != nil {
		if clientPoolIns.onCloseWrong != nil {
			clientPoolIns.onCloseWrong(my.insName, my.Name, err)
		}
		my.closeChan <- struct{}{}
		return err
	}

	return nil
}
