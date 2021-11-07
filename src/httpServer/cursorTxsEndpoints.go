package httpServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"swap.io-agent/src/blockchain"
)

//TODO: add validations invalid cursor

func (server *HttpServer) InitializeCursorTxsEndoints() {
	http.HandleFunc(
		"/getFirstCursorTransactions",
		func(rw http.ResponseWriter, r *http.Request) {
			address := r.URL.Query().Get("address")
			if len(address) == 0 || strings.Contains(address, "|") {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("invalid address"))
				return
			}

			log.Printf("get cursor address %v", address)
			cursorData, err := server.synhronizer.GetAddressFirstCursorData(
				address,
			)
			if err == blockchain.ApiNotExist {
				log.Printf("not found cursor for address %v", address)
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte(
					fmt.Sprintf("not found cursor for address %v", address),
				))
				return
			}
			if err != blockchain.ApiRequestSuccess {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(strconv.Itoa(err)))
				return
			}
			if data, err := json.Marshal(cursorData); err == nil {
				rw.Header().Set("Content-Type", "application/json; charset=utf-8")
				rw.Write(data)
				return
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
			if err == blockchain.ApiNotExist {
				log.Printf("not found cursor data for %v", cursor)
				rw.WriteHeader(http.StatusNotFound)
				rw.Write([]byte(
					fmt.Sprintf("not found cursor data for %v", cursor),
				))
				return
			}
			if err != blockchain.ApiRequestSuccess {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(strconv.Itoa(err)))
				return
			}
			if data, err := json.Marshal(cursorData); err == nil {
				rw.Header().Set("Content-Type", "application/json; charset=utf-8")
				rw.Write(data)
				return
			} else {
				log.Println(err)
			}

			rw.WriteHeader(http.StatusInternalServerError)
		},
	)
}
