package ethercsan

import (
	"encoding/json"
	"fmt"
	"io"
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
		return nil, api.RequestError
	}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
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
