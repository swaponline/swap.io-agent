package ethercsan

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetBlockByIndex(apiKey string, index int) (*Block,int) {
	log.Println("0x"+strconv.FormatInt(int64(index), 16))
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?tag=%v&boolean=true&apikey=%v&action=eth_getBlockByNumber&module=proxy",
			"0x"+strconv.FormatInt(int64(index), 16),
			apiKey,
		),
	)
	if err != nil {return nil, RequestError}

	var reqData blockRes
	// todo add switch check time limit error and parse error
	if err = json.NewDecoder(res.Body).Decode(&reqData); err != nil {
		log.Println(err)
		return nil, ParseBodyError
	}

	return &reqData.Result, RequestSuccess
}