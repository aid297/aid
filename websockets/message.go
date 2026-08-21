package websockets

import (
	"bytes"
	"strings"

	"github.com/google/uuid"

	"github.com/aid297/aid/v3/operations"
)

type Message struct {
	async            bool
	messageId        string
	message          []byte
	prototypeMessage []byte
}

var MessageApp Message

// New 新建消息
func (*Message) New(async bool, message []byte) Message { return NewMessage(async, message) }

// Parse 解析消息
func (*Message) Parse(prototypeMessage []byte) Message { return ParseMessage(prototypeMessage) }

// NewMessage 新建消息
//
//go:fix 推荐使用：New方法
func NewMessage(async bool, message []byte) Message {
	u := uuid.Must(uuid.NewV6()).String()
	b := bytes.Buffer{}
	b.Write([]byte(u))
	b.WriteByte(':')
	b.Write(message)
	return Message{
		async:            async,
		messageId:        operations.NewTernary(operations.TrueValue(u)).GetByValue(async),
		message:          operations.NewTernary(operations.TrueFn(b.Bytes), operations.FalseValue(message)).GetByValue(async),
		prototypeMessage: message,
	}
}

// ParseMessage 解析消息
//
//go:fix 推荐使用：推荐使用Parse方法
func ParseMessage(prototypeMessage []byte) Message {
	var (
		messages = strings.Split(string(prototypeMessage), ":")
		wm       = Message{}
	)

	if len(messages) == 2 {
		wm.messageId = messages[0]
		wm.message = []byte(messages[1])
		wm.prototypeMessage = prototypeMessage
		wm.async = true
	} else {
		wm.message = prototypeMessage
		wm.prototypeMessage = prototypeMessage
	}

	return wm
}

// GetAsync 获取同步类型
func (my *Message) GetAsync() bool { return my.async }

// GetMessageId 获取消息编号
func (my *Message) GetMessageId() string { return my.messageId }

// GetMessage 获取消息
func (my *Message) GetMessage() []byte {
	return operations.NewTernary(operations.TrueValue(my.message), operations.FalseValue(my.prototypeMessage)).GetByValue(my.async)
}

// GetPrototypeMessage 获取原始消息
func (my *Message) GetPrototypeMessage() []byte { return my.prototypeMessage }
