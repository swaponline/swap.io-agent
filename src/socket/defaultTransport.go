package socket

import (
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"net/http"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

var DefaultTransport = []transport.Transport{
	&polling.Transport{
		CheckOrigin: allowOriginFunc,
	},
	&websocket.Transport{
		CheckOrigin: allowOriginFunc,
	},
}