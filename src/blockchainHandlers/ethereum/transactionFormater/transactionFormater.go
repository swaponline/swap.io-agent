package transactionFormater

import (
	"errors"
	"fmt"
	"swap.io-agent/src/blockchainHandlers"
	"swap.io-agent/src/blockchainHandlers/ethereum/api/ethercsan"
)

func FormatTransaction(
	apiKey string,
	transaction *ethercsan.BlockTransaction,
	block ethercsan.Block,
) (*blockchainHandlers.Transaction, error) {
	transactionLogs, err := ethercsan.GetTransactionLogs(
		apiKey,
		transaction.Hash,
	)
	if err != ethercsan.RequestSuccess {
		return nil, errors.New(fmt.Sprintf(
			"not get transactionLogs error - %v", err,
		))
	}

}
