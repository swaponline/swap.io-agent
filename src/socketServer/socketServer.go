package socketServer

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"log"
	"net/http"
	"swap.io-agent/src/auth"
)

type subscribersStore interface {
	ClearAllUserSubscriptions(userId string) error
}

type Config struct {
	db subscribersStore
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
		userId, _ := auth.DecodeAccessToken(
			url.Query().Get("token"),
		)
		s.SetContext(userId)
		log.Printf("connect: %v", userId)

		return nil
	})
	socketServer.io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		userId := s.Context()
		err := config.db.ClearAllUserSubscriptions(userId.(string))
		log.Println(
			err, "then delete all user subscribe",
			"user:", s.Context(),
		)
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
