package websocketPool

var APP struct {
	ServerPool     ServerPool
	ServerInstance ServerInstance
	MessageTimeout MessageTimeout
	Heart          Heart
	Client         Client
	ClientInstance ClientInstance
	ClientPool     ClientPool
}
