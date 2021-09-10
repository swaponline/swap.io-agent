package httpServer

import (
	"encoding/json"
	"log"
	"net/http"
)

//TODO: add validations invalid cursor

func (server *HttpServer) InitializeCursorTxsEndoints() {
	http.HandleFunc(
		"/getFirstCursorTransactions",
		func(rw http.ResponseWriter, r *http.Request) {
			address := r.URL.Query().Get("address")
			if len(address) == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("invalid address"))
				return
			}

			log.Printf("get cursor address %v", address)
			cursorData, err := server.synhronizer.GetAddressFirstCursorData(
				address,
			)
			if err != nil {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if data, err := json.Marshal(cursorData); err == nil {
				rw.Write(data)
			} else {
				log.Println(err)
			}
			rw.WriteHeader(http.StatusInternalServerError)
		},
	)
	http.HandleFunc(
		"/getCursorTransactions",
		func(rw http.ResponseWriter, r *http.Request) {
			cursor := r.URL.Query().Get("cursor")
			if len(cursor) == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("invalid cursor"))
				return
			}

			cursorData, err := server.synhronizer.GetCursorData(
				cursor,
			)
			if err != nil {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if data, err := json.Marshal(cursorData); err == nil {
				rw.Write(data)
				return
			} else {
				log.Println(err)
			}

			rw.WriteHeader(http.StatusInternalServerError)
		},
	)
}
