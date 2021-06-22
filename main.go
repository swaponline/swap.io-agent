package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"io"
	"log"
	"net/http"
	"os"
	"swap.io-agent/src/auth"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/runApp"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/socket"
)

func main() {
	err := runApp.LoadConfig()
	if err != nil {panic(err)}

	registry := serviceRegistry.NewServiceRegistry()

	db, err := redisStore.InitializeDB()
	if err != nil {
		log.Panicf("redisStore not initialize, err: %v", err)
	}
	err = registry.RegisterService(&db)
	if err != nil {
		log.Panicf(err.Error())
	}

	server := socketio.NewServer(&engineio.Options{
		Transports: socket.DefaultTransport,
		RequestChecker: auth.AuthenticationSocketConnect,
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		url := s.URL()
		id, _ := auth.DecodeAccessToken(
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

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/getToken", func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if token, err := auth.GenerateAccessToken(0); err == nil {
			io.WriteString(writer, token)
		}
	})
	http.Handle("/socket.io/", server)

	log.Printf("Serving at localhost:%s...", os.Getenv("PORT"))
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", os.Getenv("PORT")),
			nil,
		),
	)
}