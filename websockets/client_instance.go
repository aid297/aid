package websockets

import (
	"github.com/aid297/aid/dict"
)

type ClientInstance struct {
	name        string
	connections *dict.AnyDict[string, *Client]
}

var ClientInstanceApp ClientInstance

// New 实例化：websocket客户端实例
func (*ClientInstance) New(name string) *ClientInstance { return NewClientInstance(name) }

// NewClientInstance 实例化：websocket客户端实例
//
//go:fix 推荐使用：New方法
func NewClientInstance(name string) *ClientInstance {
	return &ClientInstance{name: name, connections: dict.Make[string, *Client]()}
}

// Append 增加客户端
func (my *ClientInstance) Append(client *Client) error {
	if my.connections.HasKey(client.name) {
		return WebsocketClientExistErr.New(client.name)
	}

	my.connections.Set(client.name, client)

	return nil
}

// Remove 删除客户端
func (my *ClientInstance) Remove(name string) error {
	if !my.connections.HasKey(name) {
		return WebsocketClientNotExistErr.New(name)
	}

	my.connections.RemoveByKey(name)

	return nil
}

// Get 获取客户端
func (my *ClientInstance) Get(name string) (*Client, error) {
	if !my.connections.HasKey(name) {
		return nil, WebsocketClientNotExistErr.New(name)
	}

	client, _ := my.connections.Get(name)

	return client, nil
}

// Has 检查客户端是否存在
func (my *ClientInstance) Has(name string) bool {
	return my.connections.HasKey(name)
}

// Close 关闭客户端
func (my *ClientInstance) Close(name string) error {
	if client, err := my.Get(name); err != nil {
		return err
	} else {
		err = client.Close().Error()
		my.connections.RemoveByKey(client.name)

		return err
	}
}

// Clean 清空客户端
func (my *ClientInstance) Clean() []error {
	var errorList []error
	my.connections.Each(func(key string, client *Client) {
		if err := client.Close().Error(); err != nil {
			errorList = append(errorList, err)
		} else {
			my.connections.RemoveByKey(client.name)
		}
	})

	return errorList
}
