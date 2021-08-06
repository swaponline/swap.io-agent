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

func (e *Geth) GetTransactionByHash(
	hash string,
) (*api.BlockTransaction, int) {
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
	if err != nil {
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