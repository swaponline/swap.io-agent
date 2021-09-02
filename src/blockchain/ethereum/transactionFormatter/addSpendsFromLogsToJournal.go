package transactionFormatter

import (
	"strings"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
	"swap.io-agent/src/blockchain/journal"
)

const TransferType = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func AddSpendsFromLogsToJournal(
	logs []nodeApi.TransactionLog,
	journal *journal.Journal,
) {
	for _, value := range logs {
		if len(value.Topics) == 3 && value.Topics[0] == TransferType && IsNotEmptyVal(value.Data) {
			fromTransfer := strings.Replace(value.Topics[1], "000000000000000000000000", "", 1)
			toTransfer := strings.Replace(value.Topics[2], "000000000000000000000000", "", 1)

			if len(fromTransfer) > 0 {
				journal.Add(value.Address, blockchain.Spend{
					Wallet: fromTransfer,
					Value:  `-` + value.Data,
				})
			}
			if len(toTransfer) > 0 {
				journal.Add(value.Address, blockchain.Spend{
					Wallet: toTransfer,
					Value:  value.Data,
				})
			}
		}
	}
}
