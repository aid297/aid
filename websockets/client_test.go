package websockets

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func onLine() (*Client, error) {
	client, err := NewClient(
		"",
		"client-test1",
		"ws://127.0.0.1:8080",
		ClientCallbackConfig{
			OnConnSuccessCallback: func(groupName, name string, conn *websocket.Conn) {
				log.Printf("[%s:%s] 链接：成功\n", groupName, name)
			},
			OnConnFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				log.Fatalf("[%s:%s] 链接失败：%v", groupName, name, err)
			},
			OnCloseSuccessCallback: func(groupName, name string, conn *websocket.Conn) {
				log.Printf("[%s:%s] 关闭链接：成功\n", groupName, name)
			},
			OnCloseFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				log.Printf("[%s:%s] 关闭链接失败：%v\n", groupName, name, err)
			},
			OnReceiveMessageSuccessCallback: func(groupName, name string, prototypeMessage []byte) {
				log.Printf("[%s:%s] 接收消息：成功 -> %s\n", groupName, name, prototypeMessage)
			},
			OnReceiveMessageFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				log.Printf("[%s:%s] 接收消息失败：%v", groupName, name, err)
			},
			OnSendMessageFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				log.Printf("[%s:%s] 发送消息失败：%v", groupName, name, err)
			},
		},
		func(groupName, name string, conn *websocket.Conn) {
			log.Printf("[%s:%s] 链接成功\n", groupName, name)
		},
		func(groupName, name string, conn *websocket.Conn, err error) {
			log.Fatalf("[%s:%s] 链接失败：%v", groupName, name, err)
		},
	)
	if err != nil {
		return nil, err

	}
	if client == nil {
		return nil, errors.New("创建链接失败")
	}

	return client, nil
}

func offLine(client *Client) error {
	if err := client.Close().Error(); err != nil {
		return fmt.Errorf("关闭链接失败：%v", err)
	}
	return nil
}

func Test(t *testing.T) {
	client, err := NewClient(
		"groupA",
		"aa",
		"ws://127.0.0.1:8080/ws",
		ClientCallbackConfig{
			OnReceiveMessageSuccessCallback: func(groupName, name string, prototypeMessage []byte) {
				t.Logf("发送消息成功：[%s:%s] %s", groupName, name, prototypeMessage)
			},
			OnReceiveMessageFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				t.Logf("接收消息失败：[%s:%s] %s -> %v", groupName, name, conn.RemoteAddr().String(), err)
			},
			OnConnSuccessCallback: func(groupName, name string, conn *websocket.Conn) {
				t.Logf("连接成功：[%s:%s] %s", groupName, name, conn.RemoteAddr().String())
			},
			OnConnFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				t.Logf("连接失败：[%s:%s] %s -> %v", groupName, name, conn.RemoteAddr().String(), err)
			},
			OnSendMessageFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				t.Logf("发送消息失败：[%s:%s] %s -> %v", groupName, name, conn.RemoteAddr().String(), err)
			},
			OnCloseSuccessCallback: func(groupName, name string, conn *websocket.Conn) {
				t.Logf("关闭成功：[%s:%s] %s", groupName, name, conn.RemoteAddr().String())
			},
			OnCloseFailCallback: func(groupName, name string, conn *websocket.Conn, err error) {
				t.Logf("关闭失败：[%s:%s] %s -> %v", groupName, name, conn.RemoteAddr().String(), err)
			},
		})
	if err != nil {
		t.Fatalf("创建websocket客户端失败：%v", err)
	}

	client.AppendRequestHeader(http.Header{"Identity": []string{"abc"}}).Boot()

	ci := NewClientInstance("groupA")
	cip := OnceClientInstancePool()
	if err = ci.Append(client); err != nil {
		t.Fatalf("添加客户端到实例失败：%v", err)
	}

	if err = cip.Append(ci); err != nil {
		t.Fatalf("添加客户端实例到客户端实例池失败：%v", err)
	}

	client.AsyncMessage([]byte("这里是异步消息"), func(groupName, name string, message []byte) {
		t.Logf("异步消息回调：%s", message)
	}, time.Second)
}

// Test1Conn 测试：创建和关闭链接
func Test1Conn(t *testing.T) {
	t.Run("websocket客户端测试", func(t *testing.T) {
		client, err := onLine()
		if err != nil {
			t.Error(err)
		}

		client.AppendRequestHeader(http.Header{"Identity": []string{"abc"}}).Boot()

		if err = offLine(client); err != nil {
			t.Error(err)
		}
	})
}

// Test2Sync 测试：同步消息
func Test2Sync(t *testing.T) {
	client, err := onLine()
	if err != nil {
		t.Error(err)
	}

	_, err = client.Boot().SyncMessage([]byte("hello"), time.Second) // 1秒超时
	if err != nil {
		if !errors.Is(err, SyncMessageTimeoutErr.New("")) {
			t.Errorf("发送消息失败：%v", err)
		}
	}

	if err = offLine(client); err != nil {
		t.Errorf("关闭错误：%v", err)
	}
}

// Test3Heart 测试：心跳
func Test3Heart(t *testing.T) {
	t.Run("心跳", func(t *testing.T) {
		client, err := onLine()
		if err != nil {
			t.Errorf("获取链接失败：%v", err)
		}

		client.Boot().Heart(time.Second, func(groupName, name string, client *Client) {
			err = client.Ping(func(conn *websocket.Conn) error {
				return conn.WriteMessage(websocket.TextMessage, []byte(time.Now().GoString()))
			}).Error()
			if err != nil {
				t.Errorf("[%s:%s] 心跳失败：%v", groupName, name, err)
			} else {
				log.Printf("[%s:%s] 心跳成功\n", groupName, name)
			}
		})

		timer := time.After(5 * time.Second)

		<-timer
		log.Printf("测试成功\n")
		if err = offLine(client); err != nil {
			t.Errorf("关闭错误：%v", err)
		}
	})
}

func Test3Async(t *testing.T) {
	t.Run("心跳", func(t *testing.T) {
		client, err := onLine()
		if err != nil {
			t.Errorf("获取链接失败：%v", err)
		}

		closeSign := make(chan struct{}, 1)

		if err = client.Boot().AsyncMessage([]byte("123"), func(groupName, name string, message []byte) {
			log.Printf("[%s:%s] 回调成功 -> %s", groupName, name, message)
			closeSign <- struct{}{}
		}, 60*time.Second).Error(); err != nil {
			t.Errorf("异步消息错误：%v", err)
		}

		<-closeSign
		if err = offLine(client); err != nil {
			t.Errorf("关闭错误：%v", err)
		}
	})
}
