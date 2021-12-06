package httpServer

import (
	"io"
	"net/http"
	"swap.io-agent/src/blockchain"
)

func (*HttpServer) InitialisePushTxEndpoint(api blockchain.IBlockchainApi) {
	http.HandleFunc("/pushTx", func(writer http.ResponseWriter, request *http.Request) {
		txHex, err := io.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("incorrect txHex"))
			return
		}

		result, err := api.PushTx(string(txHex))
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write(result)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(result)
	})
}
