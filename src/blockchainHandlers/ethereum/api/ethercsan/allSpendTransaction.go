package ethercsan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"swap.io-agent/src/common/Set"
)

const transferType = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func AllSpendTransaction(apiKey string, transactionHash string) ([]string, int) {
	res, err := http.Get(
		fmt.Sprintf(
			"https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&apikey=%v",
			apiKey,
		),
	)
	if err != nil {
		return nil, RequestError
	}

	var resBody allSpendTransactionResponse
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, ParseBodyError
	}

	buf := Set.New()
	for _, value := range resBody.Result.Logs {
		if len(value.Topics) > 0 && value.Topics[0] == transferType {
			buf.Add(strings.Replace(value.Topics[1], "000000000000000000000000", "", 1))
			buf.Add(strings.Replace(value.Topics[2], "000000000000000000000000", "", 1))
		}
	}
	return buf.Keys(), RequestSuccess
}