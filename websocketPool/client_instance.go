package websocketPool

import (
	"errors"

	"github.com/aid297/aid/dict"
)

// ClientInstance websocket 客户端链接实例
type ClientInstance struct {
	Name    string
	Clients *dict.AnyDict[string, *Client]
}

var ClientInstanceApp ClientInstance

// New 实例化：websocket客户端实例
func (*ClientInstance) New(instanceName string) *ClientInstance {
	return NewClientInstance(instanceName)
}

// NewClientInstance 实例化：websocket 客户端实例
//
//go:fix 推荐使用：推荐使用New方法
func NewClientInstance(instanceName string) *ClientInstance {
	return &ClientInstance{Name: instanceName, Clients: dict.Make[string, *Client]()}
}

// GetClient 获取websocket客户端链接
func (my *ClientInstance) GetClient(clientName string) (*Client, bool) {
	websocketClient, exist := my.Clients.Get(clientName)
	if !exist {
		return nil, exist
	}

	return websocketClient, true
}

// SetClient 创建新链接
func (my *ClientInstance) SetClient(
	clientName, host, path string,
	receiveMessageFn func(instanceName, clientName string, propertyMessage []byte) ([]byte, error),
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
				if _, prototypeMsg, err = client.Conn.ReadMessage(); err != nil {
					if clientPoolIns.onReceiveMsgErr != nil {
						clientPoolIns.onReceiveMsgErr(my.Name, clientName, prototypeMsg, err)
					}
				} else {
					client.syncChan <- prototypeMsg
				}
			}
		}
	}()

	return client, nil
}

// SendMsgByName 发送消息：通过名称
func (my *ClientInstance) SendMsgByName(clientName string, msgType int, msg []byte) ([]byte, error) {
	var (
		exist  bool
		client *Client
	)

	client, exist = my.Clients.Get(clientName)
	if !exist {
		if clientPoolIns.onSendMsgErr != nil {
			clientPoolIns.onSendMsgErr(my.Name, clientName, errors.New("没有找到客户端链接"))
		}
	}

	return client.SendMsg(msgType, msg)
}

// Close 关闭客户端实例
func (my *ClientInstance) Close() {
	my.Clients.Each(func(key string, value *Client) {
		_ = value.Close()
	})

	my.Clients.Clean()
	clientPoolIns.clientInstances.RemoveByKey(my.Name)
}
