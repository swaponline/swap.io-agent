package httpServer

import (
	"encoding/json"
	"io"
	"net/http"
	"swap.io-agent/src/blockchain"
)

func InitialisePushTxEndpoint(api blockchain.IBlockchainApi) {
	http.HandleFunc("/pushTx", func(writer http.ResponseWriter, request *http.Request) {
		txHex, err := io.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("incorrect txHex"))
			return
		}

		result, err := api.PushTx(string(txHex))
		resultBytes, _ := json.Marshal(result)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write(resultBytes)
			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write(resultBytes)
	})
}
