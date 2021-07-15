package socketServer

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"log"
	"net/http"
	"swap.io-agent/src/auth"
	"swap.io-agent/src/blockchain"
)

type subscribersStore interface {
	ClearAllUserSubscriptions(userId string) error
}

type Config struct {
	db subscribersStore
	onNotifyUsers chan blockchain.TransactionPipeData
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

	connections := make(map[string]socketio.Conn)
	socketServer.io.OnConnect("/", func(s socketio.Conn) error {
		url := s.URL()
		userId, _ := auth.DecodeAccessToken(
			url.Query().Get("token"),
		)
		connections[userId] = s
		s.SetContext(userId)
		log.Printf("connect: %v", userId)

		return nil
	})
	socketServer.io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		userId := s.Context()
		delete(connections, userId.(string))
		err := config.db.ClearAllUserSubscriptions(userId.(string))
		log.Println(
			err, "then delete all user subscribe",
			"user:", s.Context(),
		)

		log.Printf("disconnect: %v", s.Context())
	})
	go func() {
		for {
			info := <-config.onNotifyUsers
			for _, userId := range info.Subscribers {
				if _, ok := connections[userId]; ok {
					//emit
				}
			}
		}
	}()

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
