package ethercsan

import (
	"strconv"
	"strings"
	"swap.io-agent/src/blockchainHandlers"
	"swap.io-agent/src/blockchainHandlers/journal"
)

const TransferType = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func AddSpendsFromLogsToJournal(
	logs []TransactionLog,
	journal *journal.Journal,
) error {
	for _, value := range logs {
		if len(value.Topics) == 3 && value.Topics[0] == TransferType {
			fromTransfer := strings.Replace(value.Topics[1], "000000000000000000000000", "", 1)
			toTransfer   := strings.Replace(value.Topics[2], "000000000000000000000000", "", 1)
			valueTransfer, err := strconv.ParseInt(value.Data, 16, 64)
			if err != nil {
				return err
			}

			journal.Add(value.Address, blockchainHandlers.Spend{
				Wallet: fromTransfer,
				Value: -valueTransfer,
			})
			journal.Add(value.Address, blockchainHandlers.Spend{
				Wallet: toTransfer,
				Value: valueTransfer,
			})
		}
	}

	return nil
}

func AllSpendAddressesTransaction(
	apiKey string,
	transaction *BlockTransaction,
) ([]string, int) {
	return make([]string, 0), 0
}