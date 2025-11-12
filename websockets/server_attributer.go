package websockets

type (
	ServerCallbackAttributer interface{ Register() }

	AttrOnConnectionFail        struct{ fn serverConnectionFailFn }
	AttrOnConnectionSuccess     struct{ fn serverConnectionSuccessFn }
	AttrOnSendMessageSuccess    struct{ fn serverSendMessageSuccessFn }
	AttrSendMessageFail         struct{ fn serverSendMessageFailFn }
	AttrOnReceiveMessageFail    struct{ fn serverReceiveMessageFailFn }
	AttrOnReceiveMessageSuccess struct{ fn serverReceiveMessageSuccessFn }
	AttrOnCloseCallback         struct{ fn serverCloseCallbackFn }
)

func OnConnectionFail(fn serverConnectionFailFn) AttrOnConnectionFail {
	return AttrOnConnectionFail{fn: fn}
}
func (my AttrOnConnectionFail) Register() { serverPool.onConnectionFail = my.fn }

func OnConnectionSuccess(fn serverConnectionSuccessFn) AttrOnConnectionSuccess {
	return AttrOnConnectionSuccess{fn: fn}
}
func (my AttrOnConnectionSuccess) Register() { serverPool.onConnectionSuccess = my.fn }

func OnSendMessageSuccess(fn serverSendMessageSuccessFn) AttrOnSendMessageSuccess {
	return AttrOnSendMessageSuccess{fn: fn}
}
func (my AttrOnSendMessageSuccess) Register() { serverPool.onSendMessageSuccess = my.fn }

func OnSendMessageFail(fn serverSendMessageFailFn) AttrSendMessageFail {
	return AttrSendMessageFail{fn: fn}
}
func (my AttrSendMessageFail) Register() { serverPool.onSendMessageFail = my.fn }

func OnReceiveMessageFail(fn serverReceiveMessageFailFn) AttrOnReceiveMessageFail {
	return AttrOnReceiveMessageFail{fn: fn}
}
func (my AttrOnReceiveMessageFail) Register() { serverPool.onReceiveMessageFail = my.fn }

func OnReceiveMessageSuccess(fn serverReceiveMessageSuccessFn) AttrOnReceiveMessageSuccess {
	return AttrOnReceiveMessageSuccess{fn: fn}
}
func (my AttrOnReceiveMessageSuccess) Register() { serverPool.onReceiveMessageSuccess = my.fn }

func OnCloseCallback(fn serverCloseCallbackFn) AttrOnCloseCallback {
	return AttrOnCloseCallback{fn: fn}
}
func (my AttrOnCloseCallback) Register() { serverPool.onCloseCallback = my.fn }
