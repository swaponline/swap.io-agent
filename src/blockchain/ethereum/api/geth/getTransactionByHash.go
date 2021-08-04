package geth

import (
	"encoding/json"
	"fmt"
	"io"
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

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, api.ParseBodyError
	}

	var resData getTransactionByHashResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return nil, api.RequestError
	}

	return &resData.Result, api.RequestSuccess
}