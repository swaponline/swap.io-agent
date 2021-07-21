package ethercsan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetTransactionByHash(
	apiKey string,
	hash string,
) (*BlockTransaction, int) {
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByHash&apikey=%v&txhash=%v",
			apiKey, hash,
		),
	)
	if err != nil {
		return nil, RequestError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ParseBodyError
	}

	var resData getTransactionByHashResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return nil,RequestError
	}

	return &resData.Result, RequestSuccess
}