package httpClient

import (
	"sync"
)

type Multiple struct {
	clients []*HttpClient
}

var MultipleApp Multiple

func (*Multiple) New() *Multiple { return NewMultiple() }

// NewMultiple 实例化：批量请求对象
//
//go:fix 推荐使用New方法
func NewMultiple() *Multiple { return &MultipleApp }

// Append 添加httpClient对象
func (my *Multiple) Append(hc *HttpClient) *Multiple {
	my.clients = append(my.clients, hc)

	return my
}

// SetClients 设置httpClient对象
func (my *Multiple) SetClients(clients []*HttpClient) *Multiple {
	my.clients = clients

	return my
}

// Send 批量发送
func (my *Multiple) Send() *Multiple {
	if len(my.clients) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(my.clients))

		for _, client := range my.clients {
			go func(client *HttpClient) {
				defer wg.Done()

				client.Send()
			}(client)
		}

		wg.Wait()
	}

	return my
}

// GetClients 获取链接池
func (my *Multiple) GetClients() []*HttpClient { return my.clients }
