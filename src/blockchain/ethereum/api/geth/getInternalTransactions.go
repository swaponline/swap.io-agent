package geth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"swap.io-agent/src/blockchain/ethereum/api"
)

func (e *Geth) GetInternalTransaction(hash string) (*api.InteranlTransaction, int) {
	/*
	   curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
	   "method": "debug_traceTransaction",
	   "params": ["0x7c9b1e9bdc3560195bd6cbe0b72f85ebc67c4edfc9333b87462be2f91e6aa872", {"tracer": "callTracer"}],
	   "jsonrpc": "2.0",
	   "id": "2"
	   }'
	*/
	res, err := http.Post(
		e.baseUrl,
		"application/json",
		strings.NewReader(
			fmt.Sprintf(
				`
				{
					"method": "debug_traceTransaction",
					"params": ["%v", {"tracer": "callTracer"}],
					"jsonrpc": "2.0",
					"id": "2"
				}
				`,
				hash,
			),
		),
	)
	if err != nil {
		return nil, api.RequestError
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, api.ParseBodyError
	}
	var resBody GetInternalTransactionsResponse
	if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, api.RequestError
	}

	return &resBody.Result, api.RequestSuccess
}
