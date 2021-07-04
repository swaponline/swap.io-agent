package ethercsan

const transferType = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
func AllSpendAddressesTransaction(
	apiKey string,
	transaction *BlockTransaction,
) ([]string, int) {
	return make([]string, 0), 0
}
//func GetTransactionLogs(
//	apiKey string,
//	hash string,
//) (*GetTransactionLogsResponse, int) {
//	res, err := http.Get(
//		fmt.Sprintf(
//			"https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%v&apikey=%v",
//			hash,
//			apiKey,
//		),
//	)
//	if err != nil {
//		return nil, RequestError
//	}
//
//	var resBody GetTransactionLogsResponse
//	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
//		return nil, ParseBodyError
//	}
//
//	buf := Set.New()
//	buf.Add(transaction.From)
//	buf.Add(transaction.To)
//	for _, value := range resBody.Result.Logs {
//		if len(value.Topics) == 3 && value.Topics[0] == transferType {
//			buf.Add(strings.Replace(value.Topics[1], "000000000000000000000000", "", 1))
//			buf.Add(strings.Replace(value.Topics[2], "000000000000000000000000", "", 1))
//		}
//	}
//	return buf.Keys(), RequestSuccess
//}