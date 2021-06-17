package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"swap.io-agent/src/auth"
	"swap.io-agent/src/runApp"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {
	err := runApp.LoadConfig()
	if err != nil {panic(err)}

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
		RequestChecker: auth.AuthenticationSocketConnect,
	})
	server.OnConnect("/", func(s socketio.Conn) error {
		url := s.URL()
		id, _ := auth.VerifyAccessToken(
			url.Query().Get("token"),
		)
		log.Printf("connect: %v", id)
		return nil
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/getToken", func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if token, err := auth.GenerateAccessToken(0); err == nil {
			io.WriteString(writer, token)
		}
	})

	log.Printf("Serving at localhost:%s...", os.Getenv("PORT"))
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", os.Getenv("PORT")),
			nil,
		),
	)
}