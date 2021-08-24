package httpServer

import (
	"encoding/json"
	"net/http"
)

//TODO: add validations invalid cursor

func (server *HttpServer) InitializeCursorTxsEndoints() {
	http.HandleFunc(
		"/getCursorTransactions/:address",
		func(rw http.ResponseWriter, r *http.Request) {
			address := r.URL.Query().Get("address")
			if len(address) == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("invalid address"))
				return
			}

			cursorData, err := server.synhronizer.GetAddressFirstCursorData(
				address,
			)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			if data, err := json.Marshal(cursorData); err == nil {
				rw.Write(data)
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
				rw.WriteHeader(http.StatusInternalServerError)
			}
			if data, err := json.Marshal(cursorData); err == nil {
				rw.Write(data)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
		},
	)
}
