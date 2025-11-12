package websocketPool

import (
	"errors"
	"fmt"

	"github.com/aid297/aid/dict"
)

type (
	// ClientPool websocket 客户端连接池
	ClientPool struct {
		onConnect       func(instanceName string, clientName string)
		onConnectErr    func(instanceName, clientName string, err error)
		onSendMsgErr    func(instanceName, clientName string, err error)
		onCloseErr      func(instanceName, clientName string, err error)
		onReceiveMsgErr func(instanceName, clientName string, prototypeMsg []byte, err error)
		clientInstances *dict.AnyDict[string, *ClientInstance]
		Error           error
	}
)

var ClientPoolApp ClientPool

func (*ClientPool) Once() *ClientPool { return OnceClientPool() }

// OnceClientPool 单例化：websocket 客户端连接池
//
//go:fix 推荐使用：Once方法
func OnceClientPool() *ClientPool {
	clientPoolOnce.Do(func() {
		clientPoolIns = &ClientPool{}
		clientPoolIns.clientInstances = dict.Make[string, *ClientInstance]()
	})

	return clientPoolIns
}

// SetOnConnect 设置回调：成功创建链接
func (*ClientPool) SetOnConnect(fn func(instanceName, clientName string)) *ClientPool {
	clientPoolIns.onConnect = fn

	return clientPoolIns
}

// SetOnConnectErr 设置回调：链接错误
func (*ClientPool) SetOnConnectErr(fn func(instanceName, clientName string, err error)) *ClientPool {
	clientPoolIns.onConnectErr = fn

	return clientPoolIns
}

// SetOnCloseClientErr 设置回调：关闭客户端链接错
func (*ClientPool) SetOnCloseClientErr(fn func(instanceName, clientName string, err error)) *ClientPool {
	clientPoolIns.onCloseErr = fn

	return clientPoolIns
}

// SetOnSendMsgErr 设置回调：发送消息错误
func (*ClientPool) SetOnSendMsgErr(fn func(instanceName, clientName string, err error)) *ClientPool {
	clientPoolIns.onSendMsgErr = fn

	return clientPoolIns
}

// SetOnReceiveMsgErr 设置回调：接收消息错误
func (*ClientPool) SetOnReceiveMsgErr(fn func(instanceName, clientName string, propertyMessage []byte, err error)) *ClientPool {
	clientPoolIns.onReceiveMsgErr = fn

	return clientPoolIns
}

// GetClientInstance 获取链接实例
func (*ClientPool) GetClientInstance(instanceName string) (*ClientInstance, bool) {
	return clientPoolIns.clientInstances.Get(instanceName)
}

// SetClientInstance 设置实例链接
func (*ClientPool) SetClientInstance(instanceName string) (*ClientInstance, error) {
	var (
		clientInstance *ClientInstance
		exist          bool
	)

	_, exist = clientPoolIns.clientInstances.Get(instanceName)
	if exist {
		return nil, fmt.Errorf("创建实例失败：%s已经存在", instanceName)
	}

	clientInstance = NewClientInstance(instanceName)
	clientPoolIns.clientInstances.Set(instanceName, clientInstance)

	return clientInstance, nil
}

// GetClient 获取客户端链接
func (*ClientPool) GetClient(instanceName, clientName string) *Client {
	var (
		clientInstance *ClientInstance
		client         *Client
		exist          bool
	)

	clientInstance, exist = clientPoolIns.clientInstances.Get(instanceName)
	if !exist {
		clientPoolIns.Error = fmt.Errorf("实例不存在：%s", instanceName)
		return nil
	}

	client, exist = clientInstance.GetClient(clientName)
	if !exist {
		clientPoolIns.Error = fmt.Errorf("链接不存在：%s", clientName)
		return nil
	}

	return client
}

// SetClient 设置websocket客户端链接
func (*ClientPool) SetClient(
	instanceName,
	clientName,
	host,
	path string,
	receiveMessageFn func(instanceName, clientName string, prototypeMsg []byte) ([]byte, error),
	heart *Heart,
	timeout *MessageTimeout,
) (*Client, error) {
	var (
		exist          bool
		clientInstance *ClientInstance
	)

	clientInstance, exist = clientPoolIns.clientInstances.Get(instanceName)
	if !exist {
		clientInstance = NewClientInstance(instanceName)
		clientPoolIns.clientInstances.Set(instanceName, clientInstance)
	}

	return clientInstance.SetClient(clientName, host, path, receiveMessageFn, heart, timeout)
}

// SendMsgByName 发送消息：通过名称
func (*ClientPool) SendMsgByName(instanceName, clientName string, msgType int, msg []byte) ([]byte, error) {
	var (
		exist          bool
		clientInstance *ClientInstance
	)
	clientInstance, exist = clientPoolIns.clientInstances.Get(instanceName)
	if !exist {
		if clientPoolIns.onSendMsgErr != nil {
			clientPoolIns.onSendMsgErr(instanceName, clientName, errors.New("没有找到客户端实例"))
		}
	}

	return clientInstance.SendMsgByName(clientName, msgType, msg)
}

// Close 关闭客户端实例池
func (*ClientPool) Close() {
	clientPoolIns.clientInstances.Each(func(key string, clientInstance *ClientInstance) {
		clientInstance.Close()
	})

	clientPoolIns.clientInstances.Clean()
}

// CloseClient 关闭链接
func (*ClientPool) CloseClient(instanceName, clientName string) error {
	var (
		exist          bool
		clientInstance *ClientInstance
		client         *Client
	)
	clientInstance, exist = clientPoolIns.clientInstances.Get(instanceName)
	if !exist {
		clientPoolIns.onCloseErr(instanceName, clientName, errors.New("没有找到客户端实例"))
		return errors.New("没有找到客户端实例")
	}

	client, exist = clientInstance.Clients.Get(clientName)
	if !exist {
		clientPoolIns.onCloseErr(instanceName, clientName, errors.New("没有找到客户端链接"))
		return errors.New("没有找到客户端链接")
	}

	return client.Close()
}
