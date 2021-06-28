package ethercsan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetCurrentIndexBlock(apiKey string) (int,int) {
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=%v",
			apiKey,
		),
	)
	if err != nil {return 0, RequestError}

	var currentBlockInfo currentBlockIndexRes
	if err = json.NewDecoder(res.Body).Decode(&currentBlockInfo); err != nil {
		return 0, ParseBodyError
	}

	currentBlockId, err := strconv.ParseInt(
		strings.TrimPrefix(currentBlockInfo.Result, "0x"),
		16,
		64,
	)
	if err != nil {return 0, ParseIndexError}

	return int(currentBlockId), Success
}