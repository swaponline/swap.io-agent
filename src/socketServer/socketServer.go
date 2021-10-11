package socketServer

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"swap.io-agent/src/auth"

	"swap.io-agent/src/queueEvents"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/subscribersManager"
)

type RawMessage struct{
    Type string `json:"type"`
    Data json.RawMessage `json:"data"`
}
type SocketServer struct {}
type Config struct {
    usersManager redisStore.IUserManager
    queueEvents  *queueEvents.QueueEvents
    subscribersManager *subscribersManager.SubscribesManager
}

var upgrader = websocket.Upgrader{} // use default options

func InitializeServer(config Config) *SocketServer {
    err := config.subscribersManager.LoadAllSubscriptions()
    if err != nil {
        log.Panic(err)
    }
    err = config.queueEvents.ReservQueueForUser("0")
    if err != nil {
        log.Panic(err)
    }

    wsHandle := func(w http.ResponseWriter, r *http.Request) {
        userId, err := auth.AuthenticationRequest(r)
        if err != nil {
            log.Println("ERROR user not connected")
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`failed auth`))
            return
        }

        c, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("upgrade:", err)
            return
        }
        defer c.Close()

        log.Println("connect:", userId)

        sourceTx,
        isOk,
        isStopped,
        stop := config.queueEvents.GetTxEventNotifier(userId)
		go func() {
			for {
				select {
				case tx := <-sourceTx:
					{
						txBytes, _ := json.Marshal(tx)
                        err := c.WriteMessage(websocket.TextMessage, txBytes)
                        if err != nil {
                            log.Println("ERROR",err)
                            break
                        }

						log.Println("send tx for", userId, tx.Hash)
					}
				case <-isStopped:
					return
				}
			}
		}()

        for {
            _, _, err := c.ReadMessage()
            if err != nil {
                log.Println("disconnect:", userId)
                stop()
                break
            }
            select {
            case isOk <- struct{}{}: continue
            case <- isStopped: break
            }
        }
    }

    http.HandleFunc("/ws", wsHandle)

	return &SocketServer{}
}

func (*SocketServer) Start() {}
func (*SocketServer) Status() error {
	return nil
}
func (*SocketServer) Stop() error  {
	return nil
}
