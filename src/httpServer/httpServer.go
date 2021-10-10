package httpServer

import (
	"io"
	"net/http"
	"swap.io-agent/src/subscribersManager"

	"swap.io-agent/src/auth"
	"swap.io-agent/src/blockchain"
)

type HttpServer struct {
	synhronizer blockchain.ISynchronizer
}
type HttpServerConfig struct {
	Synhronizer blockchain.ISynchronizer
	SubscribersManager *subscribersManager.SubscribesManager
}

func InitializeServer(config HttpServerConfig) *HttpServer {
	httpServer := HttpServer{
		synhronizer: config.Synhronizer,
	}

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/getToken", func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if token, err := auth.GenerateAccessToken("0"); err == nil {
			io.WriteString(writer, token)
		}
	})
	httpServer.InitializeCursorTxsEndoints()
	httpServer.InitializeSubscriptionsEndpoints(config.SubscribersManager)

	return &httpServer
}

func (*HttpServer) Start() {}
func (httpServer *HttpServer) Stop() error {
	return nil
}
func (*HttpServer) Status() error {
	return nil
}
