package ethercsan

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
)

func (e *Etherscan) GetBlockByIndex(index int) (*nodeApi.Block, int) {
	log.Println("get block", index, "0x"+strconv.FormatInt(int64(index), 16))
	res, err := http.Get(
		fmt.Sprintf(
			"%v/api?boolean=true&apikey=%v&tag=%v&action=eth_getBlockByNumber&module=proxy",
			e.baseUrl,
			e.apiKey,
			"0x"+strconv.FormatInt(int64(index), 16),
		),
	)
	if err != nil {
		return nil, blockchain.ApiRequestError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		return nil, blockchain.ApiParseBodyError
	}
	var resError apiError
	var resData getBlockResponse

	// insert in error struct
	if err = json.Unmarshal(resBody, &resError); err == nil {
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

	if err = json.Unmarshal(resBody, &resData); err != nil {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}
	if &resData.Result == nil {
		return nil, blockchain.ApiNotExistBlockError
	}

	return &resData.Result, blockchain.ApiRequestSuccess
}
