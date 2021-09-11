package transactionFormatter

import (
	"strconv"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
	"swap.io-agent/src/blockchain/handshake/nodeApi/fullNodeApi"
	"swap.io-agent/src/blockchain/journal"
)

const HSD = "HSD"

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

	rewardTx, minerAddress := tf.GetRewardTx(nodeBlock)
	txs = append(txs, rewardTx)

	for _, nodeTx := range nodeBlock.Txs {
		tx := tf.FormatTransaction(
			&nodeTx,
			minerAddress,
		)
		txs = append(txs, tx)
	}

	block.Transactions = txs

	return &block
}
func (tf *TransactionFormatter) GetRewardTx(
	block *nodeApi.Block,
) (*blockchain.Transaction, string) {
	minerAddress, allFee, blockReward, rewardTx := GetBlockMinderData(block)
	tx := blockchain.Transaction{
		Hash: rewardTx.Hash,
	}

	journal := journal.New(HSD)
	journal.Add(HSD, blockchain.Spend{
		Wallet: minerAddress,
		Value:  strconv.Itoa(blockReward),
		Label:  blockchain.SPEND_LABEL_TRANSFER,
	})
	journal.Add(HSD, blockchain.Spend{
		Wallet: blockchain.BLOCK_REWARD_CREATER_ADDRESS,
		Value:  strconv.Itoa(-blockReward),
		Label:  blockchain.SPEND_LABEL_BLOCK_REWARD,
	})
	if allFee > 0 {
		journal.Add(HSD, blockchain.Spend{
			Wallet: minerAddress,
			Value:  strconv.Itoa(allFee),
			Label:  blockchain.SPEND_LABEL_FEES,
		})
	}

	tx.Journal = journal.GetSpends()
	tx.AllSpendAddresses = journal.GetSpendsAddress()

	return &tx, minerAddress
}
func (tf *TransactionFormatter) FormatTransaction(
	nodeTx *nodeApi.Transaction,
	minerAddress string,
) *blockchain.Transaction {
	// todo: add check reward tx
	tx := blockchain.Transaction{
		Hash: nodeTx.Hash,
	}

	journal := journal.New(HSD)
	AddSpendsToJournal(nodeTx, journal, minerAddress)

	tx.Journal = journal.GetSpends()
	tx.AllSpendAddresses = journal.GetSpendsAddress()

	return &tx
}
