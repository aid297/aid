package websockets

var APP struct {
	ClientInstancePool ClientInstancePool
	ClientInstance     ClientInstance
	Client             Client
	Message            Message
	ServerPool         ServerPool
	Server             Server
}
