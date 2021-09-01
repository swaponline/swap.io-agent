package ethercsan

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
)

func (e *Etherscan) GetTransactionLogs(
	hash string,
) (*nodeApi.TransactionLogs, int) {
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
		return nil, blockchain.ApiRequestError
	}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}

	var resBody GetTransactionLogsResponse
	if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, blockchain.ApiParseBodyError
	}

	return &resBody.Result, blockchain.ApiRequestSuccess
}
