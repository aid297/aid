package websockets

import (
	"reflect"

	"github.com/aid297/aid/v2/anySlice"

	"github.com/aid297/aid/v2/myError"
)

type (
	WebsocketConnOption                                 struct{ myError.MyError }
	SyncMessageTimeout                                  struct{ myError.MyError }
	WebsocketOffline                                    struct{ myError.MyError }
	AsyncMessageCallbackEmpty                           struct{ myError.MyError }
	AsyncMessageTimeout                                 struct{ myError.MyError }
	WebsocketClientExist                                struct{ myError.MyError }
	WebsocketClientNotExist                             struct{ myError.MyError }
	WebsocketServerConnConditionFuncEmpty               struct{ myError.MyError }
	WebsocketServerConnTagEmpty                         struct{ myError.MyError }
	WebsocketServerConnTagExist                         struct{ myError.MyError }
	WebsocketServerOnReceiveMessageSuccessCallbackEmpty struct{ myError.MyError }
)

var (
	WebsocketConnOptionErr                                 WebsocketConnOption
	SyncMessageTimeoutErr                                  SyncMessageTimeout
	WebsocketOfflineErr                                    WebsocketOffline
	AsyncMessageCallbackEmptyErr                           AsyncMessageCallbackEmpty
	AsyncMessageTimeoutErr                                 AsyncMessageTimeout
	WebsocketClientExistErr                                WebsocketClientExist
	WebsocketClientNotExistErr                             WebsocketClientNotExist
	WebsocketServerConnConditionFuncEmptyErr               WebsocketServerConnConditionFuncEmpty
	WebsocketServerConnTagEmptyErr                         WebsocketServerConnTagEmpty
	WebsocketServerConnTagExistErr                         WebsocketServerConnTagExist
	WebsocketServerOnReceiveMessageSuccessCallbackEmptyErr WebsocketServerOnReceiveMessageSuccessCallbackEmpty
)

func (*WebsocketConnOption) New(msg string) myError.IMyError {
	return &WebsocketConnOption{myError.MyError{Msg: anySlice.NewItems("websocket连接参数错误", msg).JoinNotEmpty("：")}}
}

func (*WebsocketConnOption) Wrap(err error) myError.IMyError {
	return &WebsocketConnOption{myError.MyError{Msg: anySlice.NewItems("websocket链接参数错误", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketConnOption) Panic() myError.IMyError {
	return &WebsocketConnOption{myError.MyError{Msg: "websocket连接参数错误"}}
}

func (my *WebsocketConnOption) Error() string { return my.Msg }

func (*WebsocketConnOption) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketConnOptionErr)
}

func (*SyncMessageTimeout) New(msg string) myError.IMyError {
	return &SyncMessageTimeout{myError.MyError{Msg: anySlice.NewItems("消息同步超时", msg).JoinNotEmpty("：")}}
}

func (*SyncMessageTimeout) Wrap(err error) myError.IMyError {
	return &SyncMessageTimeout{myError.MyError{Msg: anySlice.NewItems("同步消息", err.Error()).JoinNotEmpty("：")}}
}

func (*SyncMessageTimeout) Panic() myError.IMyError {
	return &SyncMessageTimeout{myError.MyError{Msg: "消息同步超时"}}
}

func (my *SyncMessageTimeout) Error() string { return my.Msg }

func (*SyncMessageTimeout) Is(target error) bool {
	return reflect.DeepEqual(target, &SyncMessageTimeoutErr)
}

func (*WebsocketOffline) New(msg string) myError.IMyError {
	return &WebsocketOffline{myError.MyError{Msg: anySlice.NewItems("连接不在线", msg).JoinNotEmpty("：")}}
}

func (*WebsocketOffline) Wrap(err error) myError.IMyError {
	return &WebsocketOffline{myError.MyError{Msg: anySlice.NewItems("连接不在线", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketOffline) Panic() myError.IMyError {
	return &WebsocketOffline{myError.MyError{Msg: "连接不在线"}}
}

func (my *WebsocketOffline) Error() string { return my.Msg }

func (*WebsocketOffline) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketOfflineErr)
}

func (*AsyncMessageCallbackEmpty) New(msg string) myError.IMyError {
	return &AsyncMessageCallbackEmpty{myError.MyError{Msg: anySlice.NewItems("异步消息回调不能为空", msg).JoinNotEmpty("：")}}
}

func (*AsyncMessageCallbackEmpty) Wrap(err error) myError.IMyError {
	return &AsyncMessageCallbackEmpty{myError.MyError{Msg: anySlice.NewItems("异步消息回调不能为空", err.Error()).JoinNotEmpty("：")}}
}

func (*AsyncMessageCallbackEmpty) Panic() myError.IMyError {
	return &AsyncMessageCallbackEmpty{myError.MyError{Msg: "异步消息回调不能为空"}}
}

func (my *AsyncMessageCallbackEmpty) Error() string { return my.Msg }

func (*AsyncMessageCallbackEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &AsyncMessageCallbackEmptyErr)
}

func (*AsyncMessageTimeout) New(msg string) myError.IMyError {
	return &AsyncMessageTimeout{myError.MyError{Msg: anySlice.NewItems("异步消息回调超时必须大于0", msg).JoinNotEmpty("：")}}
}

func (*AsyncMessageTimeout) Wrap(err error) myError.IMyError {
	return &AsyncMessageTimeout{myError.MyError{Msg: anySlice.NewItems("异步消息回调超时必须大于0", err.Error()).JoinNotEmpty("：")}}
}

func (*AsyncMessageTimeout) Panic() myError.IMyError {
	return &AsyncMessageTimeout{myError.MyError{Msg: "异步消息回调超时必须大于0"}}
}

func (my *AsyncMessageTimeout) Error() string { return my.Msg }

func (*AsyncMessageTimeout) Is(target error) bool {
	return reflect.DeepEqual(target, &AsyncMessageTimeoutErr)
}

func (*WebsocketClientExist) New(msg string) myError.IMyError {
	return &WebsocketClientExist{myError.MyError{Msg: anySlice.NewItems("websocket客户端已存在", msg).JoinNotEmpty("：")}}
}

func (*WebsocketClientExist) Wrap(err error) myError.IMyError {
	return &WebsocketClientExist{myError.MyError{Msg: anySlice.NewItems("websocket客户端已经存在", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketClientExist) Panic() myError.IMyError {
	return &WebsocketClientExist{myError.MyError{Msg: "websocket客户端已存在"}}
}

func (my *WebsocketClientExist) Error() string { return my.Msg }

func (*WebsocketClientExist) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketClientExistErr)
}

func (*WebsocketClientNotExist) New(msg string) myError.IMyError {
	return &WebsocketClientNotExist{myError.MyError{Msg: anySlice.NewItems("websocket客户端不存在", msg).JoinNotEmpty("：")}}
}

func (*WebsocketClientNotExist) Wrap(err error) myError.IMyError {
	return &WebsocketClientNotExist{myError.MyError{Msg: anySlice.NewItems("websocket客户端不存在", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketClientNotExist) Panic() myError.IMyError {
	return &WebsocketClientNotExist{myError.MyError{Msg: "websocket客户端不存在"}}
}

func (my *WebsocketClientNotExist) Error() string { return my.Msg }

func (*WebsocketClientNotExist) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketClientExistErr)
}

func (*WebsocketServerConnConditionFuncEmpty) New(msg string) myError.IMyError {
	return &WebsocketServerConnConditionFuncEmpty{myError.MyError{Msg: anySlice.NewItems("websocket服务端连接函数不能为空", msg).JoinNotEmpty("：")}}
}

func (*WebsocketServerConnConditionFuncEmpty) Wrap(err error) myError.IMyError {
	return &WebsocketServerConnConditionFuncEmpty{myError.MyError{Msg: anySlice.NewItems("websocket服务端函数不能为空", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketServerConnConditionFuncEmpty) Panic() myError.IMyError {
	return &WebsocketServerConnConditionFuncEmpty{myError.MyError{Msg: "websocket服务端连接函数不能为空"}}
}

func (my *WebsocketServerConnConditionFuncEmpty) Error() string { return my.Msg }

func (*WebsocketServerConnConditionFuncEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerConnConditionFuncEmptyErr)
}

func (*WebsocketServerConnTagEmpty) New(msg string) myError.IMyError {
	return &WebsocketServerConnTagEmpty{myError.MyError{Msg: anySlice.NewItems("websocket服务端连接标识不能为空", msg).JoinNotEmpty("：")}}
}

func (*WebsocketServerConnTagEmpty) Wrap(err error) myError.IMyError {
	return &WebsocketServerConnTagEmpty{myError.MyError{Msg: anySlice.NewItems("websocket服务端连接标识不能为空", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketServerConnTagEmpty) Panic() myError.IMyError {
	return &WebsocketServerConnTagEmpty{myError.MyError{Msg: "websocket服务端连接标识不能为空"}}
}

func (my *WebsocketServerConnTagEmpty) Error() string { return my.Msg }

func (*WebsocketServerConnTagEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerConnTagEmptyErr)
}

func (*WebsocketServerConnTagExist) New(msg string) myError.IMyError {
	return &WebsocketServerConnTagExist{myError.MyError{Msg: anySlice.NewItems("websocket服务端连接标识重复", msg).JoinNotEmpty("：")}}
}

func (*WebsocketServerConnTagExist) Wrap(err error) myError.IMyError {
	return &WebsocketServerConnTagExist{myError.MyError{Msg: anySlice.NewItems("websocket服务端连接标识重复", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketServerConnTagExist) Panic() myError.IMyError {
	return &WebsocketServerConnTagExist{myError.MyError{Msg: "websocket服务端连接标识重复"}}
}

func (my *WebsocketServerConnTagExist) Error() string { return my.Msg }

func (*WebsocketServerConnTagExist) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerConnTagExistErr)
}

func (*WebsocketServerOnReceiveMessageSuccessCallbackEmpty) New(msg string) myError.IMyError {
	return &WebsocketServerOnReceiveMessageSuccessCallbackEmpty{myError.MyError{Msg: anySlice.NewItems("websocket服务端接收消息成功回调不能为空", msg).JoinNotEmpty("：")}}
}

func (*WebsocketServerOnReceiveMessageSuccessCallbackEmpty) Wrap(err error) myError.IMyError {
	return &WebsocketServerOnReceiveMessageSuccessCallbackEmpty{myError.MyError{Msg: anySlice.NewItems("websocket服务端接收消息成功回调不能为空", err.Error()).JoinNotEmpty("：")}}
}

func (*WebsocketServerOnReceiveMessageSuccessCallbackEmpty) Panic() myError.IMyError {
	return &WebsocketServerOnReceiveMessageSuccessCallbackEmpty{myError.MyError{Msg: "websocket服务端接收消息成功回调不能为空"}}
}

func (my *WebsocketServerOnReceiveMessageSuccessCallbackEmpty) Error() string { return my.Msg }

func (*WebsocketServerOnReceiveMessageSuccessCallbackEmpty) Is(target error) bool {
	return reflect.DeepEqual(target, &WebsocketServerOnReceiveMessageSuccessCallbackEmptyErr)
}
