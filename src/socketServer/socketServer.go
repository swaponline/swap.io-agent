package socketServer

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"swap.io-agent/src/auth"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/queueEvents"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/subscribersManager"
)

type Config struct {
	usersManager     redisStore.IUserManager
	queueEvents      *queueEvents.QueueEvents
	subscribeManager *subscribersManager.SubscribesManager
	onNotifyUsers    chan *blockchain.TransactionPipeData
}

type SocketServer struct {
	io *socketio.Server
}
type SocketServerUser struct {
	id              string
	conn            socketio.Conn
	eventIsReceived chan<- struct{}
	isStopped       <-chan struct{}
	stopEvents      context.CancelFunc
}

func InitializeServer(config Config) *SocketServer {
	err := config.subscribeManager.LoadAllSubscriptions()
	if err != nil {
		log.Panicln(err)
	}
	config.queueEvents.ReseiveQueueForUser("0") // leger

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

		sourceTx, isOk, isStopped, stop := config.queueEvents.GetTxEventNotifier(userId)
		go func() {
			for {
				select {
				case tx := <-sourceTx:
					{
						txBytes, _ := json.Marshal(tx)
						s.Emit("newTransaction", string(txBytes))
					}
				case <-isStopped:
					return
				}
			}
		}()

		connections.Store(userId, SocketServerUser{
			id:              userId,
			conn:            s,
			eventIsReceived: isOk,
			isStopped:       isStopped,
			stopEvents:      stop,
		})
		s.SetContext(userId)

		config.usersManager.ActiveUser(userId)

		log.Printf("connect: %v", userId)

		return nil
	})
	socketServer.io.OnEvent("/", "receivedTx", func(s socketio.Conn) {
		if userData, ok := connections.Load(s.Context().(string)); ok {
			select {
			case userData.(SocketServerUser).eventIsReceived <- struct{}{}:
			case <-userData.(SocketServerUser).isStopped:
				{
					return
				}
			}
		}
	})
	socketServer.io.OnEvent("/", "subscribe", func(
		s socketio.Conn,
		payload SubscriptionEventPayload,
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
		return feedBackOk(``)
	})
	socketServer.io.OnEvent("/", "unsubscribe", func(
		s socketio.Conn,
		payload SubscriptionEventPayload,
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
			return feedBackErr(``)
		}
		log.Printf(`%#v`, payload)
		return feedBackOk(``)
	})
	socketServer.io.OnDisconnect("/", func(s socketio.Conn, reason string) {
		userId := s.Context()
		if userData, ok := connections.Load(userId); ok {
			userData.(SocketServerUser).stopEvents()
		}
		connections.Delete(userId)

		err := config.usersManager.DeactiveUser(userId.(string))
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
