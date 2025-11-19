package websocketPool

import (
	"errors"
	"fmt"

	"github.com/aid297/aid/dict"
)

type (
	// ClientPool websocket 客户端连接池
	ClientPool struct {
		onConnect         func(insName string, clientName string)
		onConnectWrong    func(insName, clientName string, err error)
		onSendMsgWrong    func(insName, clientName string, err error)
		onCloseWrong      func(insName, clientName string, err error)
		onReceiveMsgWrong func(insName, clientName string, prototypeMsg []byte, err error)
		clientInsList     *dict.AnyDict[string, *ClientIns]
		Error             error
	}
)

func (*ClientPool) Once() *ClientPool { return OnceClientPool() }

// OnceClientPool 单例化：websocket 客户端连接池
//
//go:fix 推荐使用：Once方法
func OnceClientPool() *ClientPool {
	clientPoolOnce.Do(func() {
		clientPoolIns = &ClientPool{}
		clientPoolIns.clientInsList = dict.Make[string, *ClientIns]()
	})

	return clientPoolIns
}

// SetOnConnect 设置回调：成功创建链接
func (*ClientPool) SetOnConnect(fn func(insName, clientName string)) *ClientPool {
	clientPoolIns.onConnect = fn

	return clientPoolIns
}

// SetOnConnectErr 设置回调：链接错误
func (*ClientPool) SetOnConnectErr(fn func(insName, clientName string, err error)) *ClientPool {
	clientPoolIns.onConnectWrong = fn

	return clientPoolIns
}

// SetOnCloseClientErr 设置回调：关闭客户端链接错
func (*ClientPool) SetOnCloseClientErr(fn func(insName, clientName string, err error)) *ClientPool {
	clientPoolIns.onCloseWrong = fn

	return clientPoolIns
}

// SetOnSendMsgWrong 设置回调：发送消息错误
func (*ClientPool) SetOnSendMsgWrong(fn func(insName, clientName string, err error)) *ClientPool {
	clientPoolIns.onSendMsgWrong = fn

	return clientPoolIns
}

// SetOnReceiveMsgWrong 设置回调：接收消息错误
func (*ClientPool) SetOnReceiveMsgWrong(fn func(insName, clientName string, propertyMessage []byte, err error)) *ClientPool {
	clientPoolIns.onReceiveMsgWrong = fn

	return clientPoolIns
}

// GetClientIns 获取链接实例
func (*ClientPool) GetClientIns(insName string) (*ClientIns, bool) {
	return clientPoolIns.clientInsList.Get(insName)
}

// SetClientIns 设置实例链接
func (*ClientPool) SetClientIns(insName string) (*ClientIns, error) {
	var (
		clientIns *ClientIns
		exist     bool
	)

	_, exist = clientPoolIns.clientInsList.Get(insName)
	if exist {
		return nil, fmt.Errorf("创建实例失败：%s已经存在", insName)
	}

	clientIns = APP.ClientIns.New(insName)
	clientPoolIns.clientInsList.Set(insName, clientIns)

	return clientIns, nil
}

// GetClient 获取客户端链接
func (*ClientPool) GetClient(insName, clientName string) *Client {
	var (
		clientIns *ClientIns
		client    *Client
		exist     bool
	)

	clientIns, exist = clientPoolIns.clientInsList.Get(insName)
	if !exist {
		clientPoolIns.Error = fmt.Errorf("实例不存在：%s", insName)
		return nil
	}

	client, exist = clientIns.GetClient(clientName)
	if !exist {
		clientPoolIns.Error = fmt.Errorf("链接不存在：%s", clientName)
		return nil
	}

	return client
}

// SetClient 设置websocket客户端链接
func (*ClientPool) SetClient(
	insName,
	clientName,
	host,
	path string,
	receiveMessageFn func(insName, clientName string, prototypeMsg []byte) ([]byte, error),
	heart *Heart,
	timeout *MessageTimeout,
) (*Client, error) {
	var (
		exist     bool
		clientIns *ClientIns
	)

	clientIns, exist = clientPoolIns.clientInsList.Get(insName)
	if !exist {
		clientIns = APP.ClientIns.New(insName)
		clientPoolIns.clientInsList.Set(insName, clientIns)
	}

	return clientIns.SetClient(clientName, host, path, receiveMessageFn, heart, timeout)
}

// SendMsgByName 发送消息：通过名称
func (*ClientPool) SendMsgByName(insName, clientName string, msgType int, msg []byte) ([]byte, error) {
	var (
		exist     bool
		clientIns *ClientIns
	)
	clientIns, exist = clientPoolIns.clientInsList.Get(insName)
	if !exist {
		if clientPoolIns.onSendMsgWrong != nil {
			clientPoolIns.onSendMsgWrong(insName, clientName, errors.New("没有找到客户端实例"))
		}
	}

	return clientIns.SendMsgByName(clientName, msgType, msg)
}

// Close 关闭客户端实例池
func (*ClientPool) Close() {
	clientPoolIns.clientInsList.Each(func(key string, clientInstance *ClientIns) {
		clientInstance.Close()
	})

	clientPoolIns.clientInsList.Clean()
}

// CloseClient 关闭链接
func (*ClientPool) CloseClient(insName, clientName string) error {
	var (
		exist     bool
		clientIns *ClientIns
		client    *Client
	)
	clientIns, exist = clientPoolIns.clientInsList.Get(insName)
	if !exist {
		clientPoolIns.onCloseWrong(insName, clientName, errors.New("没有找到客户端实例"))
		return errors.New("没有找到客户端实例")
	}

	client, exist = clientIns.Clients.Get(clientName)
	if !exist {
		clientPoolIns.onCloseWrong(insName, clientName, errors.New("没有找到客户端链接"))
		return errors.New("没有找到客户端链接")
	}

	return client.Close()
}
