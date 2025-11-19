package websocketPool

var APP struct {
	ServerPool     ServerPool
	ServerIns      ServerIns
	MessageTimeout MessageTimeout
	Heart          Heart
	Client         Client
	ClientIns      ClientIns
	ClientPool     ClientPool
}
