package ethercsan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"swap.io-agent/src/blockchain/ethereum/api"
)

func (e *Etherscan) GetTransactionByHash(
	hash string,
) (*api.BlockTransaction, int) {
	res, err := http.Get(
		fmt.Sprintf(
			"%v/api?module=proxy&action=eth_getTransactionByHash&apikey=%v&txhash=%v",
			e.baseUrl,
			e.apiKey,
			hash,
		),
	)
	if err != nil {
		return nil, api.RequestError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		return nil, api.ParseBodyError
	}

	var resData getTransactionByHashResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return nil, api.RequestError
	}

	return &resData.Result, api.RequestSuccess
}
