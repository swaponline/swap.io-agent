package geth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"swap.io-agent/src/blockchain"
	nodeApi "swap.io-agent/src/blockchain/ethereum/nodeApi"
)

/*
curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
"method": "eth_getBlockByNumber",
"params": ["0xa68a4e", true],
"jsonrpc": "2.0",
"id": "2"
}'
*/
func (e *Geth) GetBlockByIndex(index int) (*nodeApi.Block, int) {
	log.Println("get block", index, "0x"+strconv.FormatInt(int64(index), 16))
	res, err := http.Post(
		e.baseUrl,
		"application/json",
		strings.NewReader(
			fmt.Sprintf(
				`{
					"jsonrpc":"2.0",
					"method":"eth_getBlockByNumber",
                   "params":["%v", true],
					"id":1
				}`,
				"0x"+strconv.FormatInt(int64(index), 16),
			),
		),
	)
	if err != nil {
		log.Println(err)
		return nil, blockchain.ApiRequestError
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}
	var resError apiError
	var resBody getBlockResponse

	// insert in error struct
	if err = json.Unmarshal(resBodyBytes, &resError); err == nil {
		log.Println(err, string(resBodyBytes))
		if resError.Result == "Max rate limit reached" {
			return nil, blockchain.ApiRequestLimitError
		}
		// if error parsed width empty filed then block not exit
		if resError.Result == "" &&
			resError.Status == "" &&
			resError.Message == "" {
			return nil, blockchain.ApiNotExistBlockError
		}
		return nil, blockchain.ApiRequestError
	}

	if err = json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, blockchain.ApiParseBodyError
	}
	if &resBody.Result == nil {
		log.Println(string(resBodyBytes), "result = null")
		return nil, blockchain.ApiNotExistBlockError
	}

	return &resBody.Result, blockchain.ApiRequestSuccess
}
