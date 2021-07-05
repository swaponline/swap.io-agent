package ethercsan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetTransactionLogs(
	apiKey string,
	hash string,
) (*GetTransactionLogsResponse, int) {
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%v&apikey=%v",
			hash,
			apiKey,
		),
	)
	if err != nil {
		return nil, RequestError
	}

	var resBody GetTransactionLogsResponse
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, ParseBodyError
	}

	return &resBody, RequestSuccess
}