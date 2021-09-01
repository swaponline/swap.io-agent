package ethercsan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"swap.io-agent/src/blockchain"
)

func (e *Etherscan) GetCurrentIndexBlock() (int, int) {
	res, err := http.Get(
		fmt.Sprintf(
			"%v/api?module=proxy&action=eth_blockNumber&apikey=%v",
			e.baseUrl,
			e.apiKey,
		),
	)
	if err != nil {
		return 0, blockchain.ApiRequestError
	}

	var currentBlockInfo getCurrentBlockIndexResponse
	if err = json.NewDecoder(res.Body).Decode(&currentBlockInfo); err != nil {
		return 0, blockchain.ApiParseBodyError
	}

	currentBlockId, err := strconv.ParseInt(
		strings.TrimPrefix(currentBlockInfo.Result, "0x"),
		16,
		64,
	)
	if err != nil {
		return 0, blockchain.ApiParseIndexError
	}

	return int(currentBlockId), blockchain.ApiRequestSuccess
}
