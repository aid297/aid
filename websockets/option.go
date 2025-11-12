package websockets

type (
	// ClientCallbackConfig 客户端回调
	ClientCallbackConfig struct {
		OnConnSuccessCallback           clientStandardSuccessFn
		OnConnFailCallback              clientStandardFailFn
		OnCloseSuccessCallback          clientStandardSuccessFn
		OnCloseFailCallback             clientStandardFailFn
		OnReceiveMessageSuccessCallback clientReceiveMessageSuccessFn
		OnReceiveMessageFailCallback    clientStandardFailFn
		OnSendMessageFailCallback       clientStandardFailFn
	}

	ServerCallbackConfig struct {
		OnConnectionFail        serverConnectionFailFn
		OnConnectionSuccess     serverConnectionSuccessFn
		OnSendMessageSuccess    serverSendMessageSuccessFn
		OnSendMessageFail       serverSendMessageFailFn
		OnReceiveMessageFail    serverReceiveMessageFailFn
		OnReceiveMessageSuccess serverReceiveMessageSuccessFn
		OnCloseCallback         serverCloseCallbackFn
	}
)
