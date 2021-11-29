package transactionFormatter

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
	"swap.io-agent/src/blockchain/handshake/nodeApi/fullNodeApi"
	"swap.io-agent/src/blockchain/journal"
	"swap.io-agent/src/config"
)

const HSN = "HSN"

type TransactionFormatter struct {
	api *fullNodeApi.FullNodeApi
}
type TransactionFormatterConfig struct {
	Api *fullNodeApi.FullNodeApi
}

func InitializeTransactionFormatter(
	config TransactionFormatterConfig,
) *TransactionFormatter {
	return &TransactionFormatter{
		api: config.Api,
	}
}

func (tf *TransactionFormatter) FormatBlock(
	nodeBlock *nodeApi.Block,
) *blockchain.Block {
	block := blockchain.Block{}
	txs := make([]*blockchain.Transaction, 0)

	minerAddress := GetBlockMinderAddress(nodeBlock)

	for _, nodeTx := range nodeBlock.Txs {
		tx := tf.FormatTransaction(
			&nodeTx,
			nodeBlock,
			minerAddress,
		)
		txs = append(txs, tx)
	}

	block.Transactions = txs

	return &block
}

func (tf *TransactionFormatter) FormatTransaction(
	nodeTx *nodeApi.Transaction,
	nodeBlock *nodeApi.Block,
	minerAddress string,
) *blockchain.Transaction {
	tx := blockchain.Transaction{
		Hash: nodeTx.Hash,
	}

	journal := journal.New(config.BLOCKCHAIN)

	if nodeTx.Fee == 0 {
		AddMinerRewardToJournal(nodeTx, nodeBlock, journal)
	} else {
		AddSpendsToJournal(nodeTx, journal, minerAddress)
	}

	tx.Journal = journal.GetSpends()
	tx.AllSpendAddresses = journal.GetSpendsAddress()

	return &tx
}
