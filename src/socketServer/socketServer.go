package socketServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"swap.io-agent/src/auth"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/subscribersManager"
)

type Config struct {
	subscribeManager *subscribersManager.SubscribesManager
	onNotifyUsers    chan *blockchain.TransactionPipeData
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

	connections := sync.Map{}
	socketServer.io.OnConnect("/", func(s socketio.Conn) error {
		url := s.URL()
		userId, _ := auth.DecodeAccessToken(
			url.Query().Get("token"),
		)
		connections.Store(userId, s)
		s.SetContext(userId)
		err := config.subscribeManager.LoadSubscriptions(
			userId,
		)
		if err != nil {
			return err
		}
		log.Printf("connect: %v", userId)

		return nil
	})
	socketServer.io.OnEvent("/", "subscribe", func(
		s socketio.Conn,
		payload SubscribeEventPayload,
	) string {
		userId := s.Context().(string)
		err := config.subscribeManager.SubscribeUserToAddress(
			userId,
			payload.Address,
			true,
		)
		if err != nil {
			log.Println(
				"err",
				err, "then subscribe user",
				"user:", s.Context(),
			)
			return "error"
		}
		log.Printf(`%#v`, payload)
		return ""
	})
	socketServer.io.OnEvent("/", "unsubscribe", func(
		s socketio.Conn,
		payload SubscribeEventPayload,
	) string {
		userId := s.Context().(string)
		err := config.subscribeManager.UnsubscribeUserToAddress(
			userId,
			payload.Address,
		)
		if err != nil {
			log.Println(
				"err",
				err, "then unsubscribe user",
				"user:", s.Context(),
			)
			return "error"
		}
		log.Printf(`%#v`, payload)
		return "{}"
	})
	socketServer.io.OnEvent("/", "subscriptionsSize", func(
		s socketio.Conn,
	) string {
		userId := s.Context().(string)
		size, err := config.subscribeManager.GetSubscriptionsSize(
			userId,
		)
		if err != nil {
			log.Println(
				"err",
				err, "then unsubscribe user",
				"user:", s.Context(),
			)
			return "error"
		}

		return fmt.Sprintf(`{"size": %v}`, size)
	})
	socketServer.io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		userId := s.Context()
		connections.Delete(userId)
		err := config.subscribeManager.ClearAllUserSubscriptions(userId.(string))
		if err != nil {
			log.Println(
				"err",
				err, "then delete all user subscribe",
				"user:", s.Context(),
			)
		}

		log.Printf("disconnect: %v", s.Context())
	})
	go func() {
		for {
			transactionInfo := <-config.onNotifyUsers
			transactionsJson, err := json.Marshal(transactionInfo.Transaction)
			if err != nil {
				log.Println("err", err)
				continue
			}
			transactionsStr := string(transactionsJson)
			for _, userId := range transactionInfo.Subscribers {
				if connection, ok := connections.Load(userId); ok && connection != nil {
					connection.(socketio.Conn).Emit(
						"newTransaction",
						transactionsStr,
					)
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

func (*SocketServer) Start() {}
func (socketServer *SocketServer) Stop() error {
	return socketServer.io.Close()
}
func (*SocketServer) Status() error {
	return nil
}
