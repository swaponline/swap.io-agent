package ethercsan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"swap.io-agent/src/blockchain/ethereum/api"
)

func (e *Etherscan) GetTransactionLogs(
	hash string,
) (*api.TransactionLogs, int) {
	res, err := http.Get(
		fmt.Sprintf(
			"%v/api?module=proxy&action=eth_getTransactionReceipt&apikey=%v&txhash=%v",
			e.baseUrl,
			e.apiKey,
			hash,
		),
	)
	if err != nil {
		log.Println(err)
		return nil, RequestError
	}

	var resBody GetTransactionLogsResponse
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		log.Println(err)
		return nil, ParseBodyError
	}

	return &resBody.Result, RequestSuccess
}