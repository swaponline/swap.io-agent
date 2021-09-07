package transactionFormatter

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
	"swap.io-agent/src/blockchain/handshake/nodeApi/fullNodeApi"
)

const HSD = "HSD"

type TransactionFormatter struct {
	api fullNodeApi.FullNodeApi
}
type TransactionFormatterConfig struct {
	Api fullNodeApi.FullNodeApi
}

func InitializeTransactionFormatter(
	config TransactionFormatterConfig,
) *TransactionFormatter {
	return &TransactionFormatter{
		api: config.Api,
	}
}

func (tf *TransactionFormatter) FormatBlock(
	block *nodeApi.Block,
) (*blockchain.Block, error) {
	//addressMiner, allFee, blockReward, indexRewordTxInBlockTxs := GetBlockMinderData(block)
	return nil, nil
}

func (tf *TransactionFormatter) FormatTransactionFromHash(
	hash string,
) (*blockchain.Transaction, error) {

	return nil, nil
}
func (tf *TransactionFormatter) FormatTransaction(
	transaction *nodeApi.Transaction,
) (*blockchain.Transaction, error) {

	return nil, nil
}
func (tf *TransactionFormatter) FormatRewordTx() {

}
