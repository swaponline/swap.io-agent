package ethercsan

import (
	"strings"
	"swap.io-agent/src/common/Set"
)

const TransferType = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func GetAllSpendAddressFromLogs(
	logs []TransactionLog,
	transaction *BlockTransaction,
	miner string,
) []string {
	buf := Set.New()
	buf.Add(miner)
	buf.Add(transaction.From)
	buf.Add(transaction.To)
	for _, value := range logs {
		if len(value.Topics) == 3 && value.Topics[0] == TransferType {
			buf.Add(strings.Replace(value.Topics[1], "000000000000000000000000", "", 1))
			buf.Add(strings.Replace(value.Topics[2], "000000000000000000000000", "", 1))
		}
	}
	return buf.Keys()
}

func AllSpendAddressesTransaction(
	apiKey string,
	transaction *BlockTransaction,
) ([]string, int) {
	return make([]string, 0), 0
}