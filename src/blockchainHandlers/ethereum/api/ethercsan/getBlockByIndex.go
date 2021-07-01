package ethercsan

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func GetBlockByIndex(apiKey string, index int) (*Block,int) {
	log.Println("get block", index, "0x"+strconv.FormatInt(int64(index), 16))
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?tag=%v&boolean=true&apikey=%v&action=eth_getBlockByNumber&module=proxy",
			"0x"+strconv.FormatInt(int64(index), 16),
			apiKey,
		),
	)
	if err != nil {return nil, RequestError}

	reqBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ParseBodyError
	}
	var reqError apiError
	var reqData blockResponse

	// insert in error struct
	if err = json.Unmarshal(reqBody, &reqError); err == nil {
		if reqError.Result == "Max rate limit reached" {
			return nil, RequestLimitError
		}
		// if error parsed width empty filed then block not exit
		if reqError.Result == "" &&
		   reqError.Status == "" &&
		   reqError.Message == "" {
			return nil, NotExistBlockError
		}
		return nil, RequestError
	}

	if err = json.Unmarshal(reqBody, &reqData); err != nil {
		log.Println(err)
		return nil, ParseBodyError
	}

	return &reqData.Result, RequestSuccess
}