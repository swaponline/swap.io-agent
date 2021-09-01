package ethercsan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
)

func (e *Etherscan) GetTransactionByHash(
	hash string,
) (*nodeApi.BlockTransaction, int) {
	res, err := http.Get(
		fmt.Sprintf(
			"%v/api?module=proxy&action=eth_getTransactionByHash&apikey=%v&txhash=%v",
			e.baseUrl,
			e.apiKey,
			hash,
		),
	)
	if err != nil {
		return nil, blockchain.ApiRequestError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		return nil, blockchain.ApiParseBodyError
	}

	var resData getTransactionByHashResponse
	if err := json.Unmarshal(resBody, &resData); err != nil {
		return nil, blockchain.ApiParseBodyError
	}

	return &resData.Result, blockchain.ApiRequestSuccess
}
