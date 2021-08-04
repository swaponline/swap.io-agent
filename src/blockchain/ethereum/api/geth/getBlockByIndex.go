package geth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"swap.io-agent/src/blockchain/ethereum/api"
)


func (e *Geth) GetBlockByIndex(index int) (*api.Block,int) {
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
					"id":1}
				`,
				"0x"+strconv.FormatInt(int64(index), 16),
			),
		),
	)
	if err != nil {return nil, api.RequestError
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, api.ParseBodyError
	}
	var resError apiError
	var resData getBlockResponse

	// insert in error struct
	if err = json.Unmarshal(resBody, &resError); err == nil {
		if resError.Result == "Max rate limit reached" {
			return nil, api.RequestLimitError
		}
		// if error parsed width empty filed then block not exit
		if resError.Result  == "" &&
			resError.Status  == "" &&
			resError.Message == "" {
			return nil, api.NotExistBlockError
		}
		return nil, api.RequestError
	}

	if err = json.Unmarshal(resBody, &resData); err != nil {
		log.Println(err)
		return nil, api.ParseBodyError
	}
	if &resData.Result == nil {
		return nil, api.NotExistBlockError
	}

	return &resData.Result, api.RequestSuccess
}