package ethercsan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetBlockByIndex(apiKey string, index int) (*block,int) {
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?tag=%v&boolean=true&apikey=%v&action=eth_getBlockByNumber&module=proxy",
			"0x"+strconv.FormatInt(int64(index), 16),
			apiKey,
		),
	)
	if err != nil {return nil, RequestError}

	var blockInfo blockRes
	if err = json.NewDecoder(res.Body).Decode(&blockInfo); err != nil {
		return nil, ParseBodyError
	}

	return &blockInfo.Result, Success
}