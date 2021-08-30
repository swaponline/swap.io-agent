package geth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"

	"swap.io-agent/src/blockchain/ethereum/api"
)

/*
curl --location --request POST 'localhost:8545/' --header 'Content-Type: application/json' --data-raw '{
"method": "eth_getBlockByNumber",
"params": ["0xa68a4e", true],
"jsonrpc": "2.0",
"id": "2"
}'
*/
func (e *Geth) GetBlockByIndex(index int) (*api.Block, int) {
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
		return nil, api.RequestError
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		return nil, api.ParseBodyError
	}
	var resError apiError
	var resBody getBlockResponse

	// insert in error struct
	if err = json.Unmarshal(resBodyBytes, &resError); err == nil {
		log.Println(err, string(resBodyBytes))
		if resError.Result == "Max rate limit reached" {
			return nil, api.RequestLimitError
		}
		// if error parsed width empty filed then block not exit
		if resError.Result == "" &&
			resError.Status == "" &&
			resError.Message == "" {
			return nil, api.NotExistBlockError
		}
		return nil, api.RequestError
	}

	if err = json.Unmarshal(resBodyBytes, &resBody); err != nil {
		log.Println(err, string(resBodyBytes))
		return nil, api.ParseBodyError
	}
	if &resBody.Result == nil {
		log.Println(string(resBodyBytes), "result = null")
		return nil, api.NotExistBlockError
	}

	return &resBody.Result, api.RequestSuccess
}
