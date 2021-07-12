package httpServer

import (
	"io"
	"net/http"
	"swap.io-agent/src/auth"
)

type HttpServer struct {}

func InitializeServer() *HttpServer {
	httpServer := HttpServer{}

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/getToken", func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if token, err := auth.GenerateAccessToken("0"); err == nil {
			io.WriteString(writer, token)
		}
	})

	return &httpServer
}

func (_ *HttpServer) Start()  {}
func (httpServer *HttpServer) Stop() error {
	return nil
}
func (_ *HttpServer) Status() error {
	return nil
}
