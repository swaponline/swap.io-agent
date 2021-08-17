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

func (e *Geth) GetTransactionLogs(
	hash string,
) (*api.TransactionLogs, int) {
/*
curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
"method": "eth_getTransactionReceipt",
"params": ["0x383d86755b0288a9a4cece3f0f6c6d721924b1f75fcf402b7f5e39c0d8ac10f5"],
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
					"method":"eth_getTransactionReceipt",
					"params":["%v"],
					"id":1
				}`,
				hash,
			),
		),
	)
	if err != nil {
		log.Println(err)
		return nil, api.RequestError
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, api.ParseBodyError
	}

	var resBody GetTransactionLogsResponse
	if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, api.ParseBodyError
	}

	return &resBody.Result, api.RequestSuccess
}