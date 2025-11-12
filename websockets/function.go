package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	clientCallbackFn              func(groupName, name string, message []byte)
	clientStandardSuccessFn       func(groupName, name string, conn *websocket.Conn)
	clientStandardFailFn          func(groupName, name string, conn *websocket.Conn, err error)
	clientReceiveMessageSuccessFn func(groupName, name string, prototypeMessage []byte)
	clientHeartFn                 func(groupName, name string, client *Client)
	pingFn                        func(conn *websocket.Conn) error
	serverConnectionFailFn        func(err error)
	serverConnectionSuccessFn     func(conn *websocket.Conn) error
	serverConnectionCheckFn       func(header http.Header) (string, error)
	serverReceiveMessageSuccessFn func(server *Server, message Message)
	serverReceiveMessageFailFn    func(conn *websocket.Conn, err error)
	// serverReceivePingFn           func(conn *websocket.Conn)
	serverSendMessageFailFn    func(err error)
	serverSendMessageSuccessFn func(conn *websocket.Conn, message, prototypeMessage []byte)
	serverCloseCallbackFn      func(conn *websocket.Conn)
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
