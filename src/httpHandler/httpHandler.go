package httpHandler

import (
	"fmt"
	"log"
	"net/http"
	"swap.io-agent/src/env"
)

type HttpHandler struct {
	server *http.Server
}

func InitializeServer() *HttpHandler {
	server := &http.Server{
		Addr: fmt.Sprintf(":%v", env.PORT),
	}
	return &HttpHandler{
		server: server,
	}
}

func (httpHandler *HttpHandler) Start() {
	log.Printf("Http handle:%v", env.PORT)

	err := httpHandler.server.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
}
func (httpHandler *HttpHandler) Stop() error {
	return httpHandler.server.Close()
}
func (httpHandler *HttpHandler) Status() error {
	return nil
}
