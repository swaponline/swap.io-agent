package httpHandler

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type HttpHandler struct {
	server *http.Server
}

func InitializeServer() *HttpHandler {
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", os.Getenv("PORT")),
	}
	return &HttpHandler{
		server: server,
	}
}

func (httpHandler *HttpHandler) Start() {
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
