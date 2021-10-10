package httpServer

import (
	"log"
	"net/http"

	"swap.io-agent/src/auth"
    "swap.io-agent/src/subscribersManager"
)

func (*HttpServer) InitializeSubscriptionsEndpoints(
    subscribersManager *subscribersManager.SubscribesManager,
) {
	http.HandleFunc(
		"/subscribe",
		func(rw http.ResponseWriter, r *http.Request) {
            userId, err := auth.AuthenticationRequest(r)
            if err != nil {
                rw.WriteHeader(http.StatusUnauthorized)
                return
            }

			address := r.URL.Query().Get("address")
			if len(address) == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("not address field"))
				return
			}

            err = subscribersManager.SubscribeUserToAddress(
                userId,
                address,
                true,
            )
            if err != nil {
                rw.WriteHeader(http.StatusInternalServerError)
                return
            }

			log.Printf("user %v subscribe %v", userId, address)

            rw.WriteHeader(http.StatusOK)
		},
	)
    http.HandleFunc(
		"/unsubscribe",
		func(rw http.ResponseWriter, r *http.Request) {
            userId, err := auth.AuthenticationRequest(r)
            if err != nil {
                rw.WriteHeader(http.StatusUnauthorized)
                return
            }

			address := r.URL.Query().Get("address")
			if len(address) == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("not address field"))
				return
			}

            err = subscribersManager.UnsubscribeUserToAddress(
                userId,
                address,
            )
            if err != nil {
                rw.WriteHeader(http.StatusInternalServerError)
                return
            }

			log.Printf("user %v unsubscribe %v", userId, address)

            rw.WriteHeader(http.StatusOK)
		},
	)
}
