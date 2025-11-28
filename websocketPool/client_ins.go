package websocketPool

import (
	"errors"
	"io"

	"github.com/gorilla/websocket"

	"github.com/aid297/aid/dict"
)

// ClientIns websocket 客户端链接实例
type ClientIns struct {
	Name    string
	Clients *dict.AnyDict[string, *Client]
}

// New 实例化：websocket客户端实例
func (*ClientIns) New(insName string) *ClientIns {
	return &ClientIns{Name: insName, Clients: dict.Make[string, *Client]()}
}

// GetClient 获取websocket客户端链接
func (my *ClientIns) GetClient(clientName string) (*Client, bool) {
	websocketClient, exist := my.Clients.Get(clientName)
	if !exist {
		return nil, exist
	}

	return websocketClient, true
}

// SetClient 创建新链接
func (my *ClientIns) SetClient(
	clientName, host, path string,
	receiveMessageFn func(insName, clientName string, propertyMessage []byte) ([]byte, error),
	heart *Heart,
	timeout *MessageTimeout,
) (*Client, error) {
	var (
		err          error
		exist        bool
		client       *Client
		prototypeMsg []byte
	)

	client, exist = my.Clients.Get(clientName)
	if exist {
		if err = client.Conn.Close(); err != nil {
			return nil, err
		}
		my.Clients.RemoveByKey(clientName)
	}

	if client, err = NewClient(my.Name, clientName, host, path, receiveMessageFn); err != nil {
		return nil, err
	}
	my.Clients.Set(clientName, client)

	if clientPoolIns.onConnect != nil {
		clientPoolIns.onConnect(my.Name, clientName)
	}

	if heart == nil {
		heart = DefaultHeart()
	}
	client.heart = heart
	if timeout != nil {
		client.timeout = timeout
	}

	// 开启协程：接收消息
	go func() {
		for {
			select {
			case <-client.closeChan:
				// 关闭链接
				client.heart.ticker.Stop()
				my.Clients.RemoveByKey(clientName)
				return
			case <-client.heart.ticker.C:
				// 执行心跳
				if client.heart.fn != nil {
					client.heart.fn(client)
				}
			default:
				var reader io.Reader
				if _, reader, err = client.Conn.NextReader(); err != nil {
					if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) && clientPoolIns.onReceiveMsgWrong != nil {
						clientPoolIns.onReceiveMsgWrong(my.Name, clientName, prototypeMsg, err)
					}

					continue
				}
				if prototypeMsg, err = io.ReadAll(reader); err != nil {
					if clientPoolIns.onReceiveMsgWrong != nil {
						clientPoolIns.onReceiveMsgWrong(my.Name, clientName, prototypeMsg, err)
					}

					continue
				}

				client.syncChan <- prototypeMsg

				// if _, prototypeMsg, err = client.Conn.ReadMessage(); err != nil {
				// 	if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) && clientPoolIns.onReceiveMsgWrong != nil {
				// 		clientPoolIns.onReceiveMsgWrong(my.Name, clientName, prototypeMsg, err)
				// 	}
				// 	// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 	// 	client.closeChan <- struct{}{} // 链接被意外关闭
				// 	// } else {
				// 	// 	if clientPoolIns.onReceiveMsgWrong != nil {
				// 	// 		clientPoolIns.onReceiveMsgWrong(my.Name, clientName, prototypeMsg, err)
				// 	// 	}
				// 	// }
				// 	continue
				// } else {
				// 	client.syncChan <- prototypeMsg
				// }
			}
		}
	}()

	return client, nil
}

// SendMsgByName 发送消息：通过名称
func (my *ClientIns) SendMsgByName(clientName string, msgType int, msg []byte) ([]byte, error) {
	var (
		exist  bool
		client *Client
	)

	client, exist = my.Clients.Get(clientName)
	if !exist {
		if clientPoolIns.onSendMsgWrong != nil {
			clientPoolIns.onSendMsgWrong(my.Name, clientName, errors.New("没有找到客户端链接"))
		}
	}

	return client.SendMsg(msgType, msg)
}

// Close 关闭客户端实例
func (my *ClientIns) Close() {
	my.Clients.Each(func(key string, value *Client) {
		_ = value.Close()
	})

	my.Clients.Clean()
	clientPoolIns.clientInsList.RemoveByKey(my.Name)
}
