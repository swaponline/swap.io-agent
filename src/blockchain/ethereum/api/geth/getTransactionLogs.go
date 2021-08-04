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

func (e *Geth) GetTransactionLogs(
	hash string,
) (*api.TransactionLogs, int) {
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
		log.Println(err)
		return nil, api.RequestError
	}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, api.ParseBodyError
	}

	var resBody GetTransactionLogsResponse
	if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, api.ParseBodyError
	}

	return &resBody.Result, api.RequestSuccess
}