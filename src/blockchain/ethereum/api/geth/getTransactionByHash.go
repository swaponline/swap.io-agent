package geth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"swap.io-agent/src/blockchain/ethereum/api"
)

func (e *Geth) GetTransactionByHash(
	hash string,
) (*api.BlockTransaction, int) {
	/*
	   curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
	   "method": "eth_getTransactionByHash",
	   "params": ["0x7c9b1e9bdc3560195bd6cbe0b72f85ebc67c4edfc9333b87462be2f91e6aa872"],
	   "jsonrpc": "2.0",
	   "id": "2"
	   }'
	*/
	res, err := http.Post(
		e.baseUrl,
		"application/json",
		strings.NewReader(
			fmt.Sprintf(
				`{
					"jsonrpc":"2.0",
					"method":"eth_getTransactionByHash",
					"params":["%v"],
					"id":1
				}`,
				hash,
			),
		),
	)
	if err != nil {
		return nil, api.RequestError
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		return nil, api.ParseBodyError
	}

	var resBody getTransactionByHashResponse
	if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, api.RequestError
	}

	return &resBody.Result, api.RequestSuccess
}
