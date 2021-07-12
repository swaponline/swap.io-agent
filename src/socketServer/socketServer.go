package socketServer

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"log"
	"net/http"
	"swap.io-agent/src/auth"
)

type Config struct {
	db socketServerDb
}
type SocketServer struct {
	io *socketio.Server
}
func InitializeServer(config Config) *SocketServer {
	socketServer := SocketServer{
		io: socketio.NewServer(&engineio.Options{
			Transports:     DefaultTransport,
			RequestChecker: AuthenticationSocketConnect,
		}),
	}

	socketServer.io.OnConnect("/", func(s socketio.Conn) error {
		url := s.URL()
		id, _ := auth.DecodeAccessToken(
			url.Query().Get("token"),
		)
		log.Printf("connect: %v", id)

		return nil
	})
	socketServer.io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("disconnect: %v", s.Context())
	})

	go func() {
		if err := socketServer.io.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	http.Handle("/socket.io/", socketServer.io)

	return &socketServer
}

func (_ *SocketServer) Start() {}
func (socketServer *SocketServer) Stop() error {
	return socketServer.io.Close()
}
func (_ *SocketServer) Status() error {
	return nil
}
