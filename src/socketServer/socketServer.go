package socketServer

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"swap.io-agent/src/auth"
    "time"

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

const writePeriod = time.Second * 1
const readPeriod  = time.Second * 2

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

        ticker := time.NewTicker(writePeriod)
        defer ticker.Stop()

        sourceTx,
        isOk,
        isStopped,
        stopTxSource := config.queueEvents.GetTxEventNotifier(userId)
		go func() {
            for {
				select {
				case <-isStopped.Done():
                    {
                        return
                    }
                case tx, ok := <-sourceTx:
					{
                        if !ok {
                            return
                        }
						txBytes, _ := json.Marshal(tx)
                        err := c.WriteMessage(websocket.TextMessage, txBytes)
                        if err != nil {
                            log.Println("ERROR",err)
                            return
                        }
						log.Println("send tx for", userId, tx.Hash)
                    }
                case _, ok := <-ticker.C:
                    {
                        if !ok {
                            return
                        }
                        c.SetWriteDeadline(time.Now().Add(writePeriod))
                        if err := c.WriteMessage(
                            websocket.PingMessage, nil,
                        ); err != nil {
                            log.Println("ERROR", err)
                            stopTxSource()
                            return
                        }
                        log.Println("ping")
                    }
				}
			}
		}()

        c.SetReadDeadline(time.Now().Add(readPeriod));
        c.SetPongHandler(
            func(string) error {
                log.Println("pong")
                c.SetReadDeadline(time.Now().Add(readPeriod));
                return nil
            },
        )
        for {
            _, _, err := c.ReadMessage()
            if err != nil {
                log.Println("disconnect:", userId)
                stopTxSource()
                return
            }
            select {
            case <- isStopped.Done(): return
            case isOk <- struct{}{}:
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
