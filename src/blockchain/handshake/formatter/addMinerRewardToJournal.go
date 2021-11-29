package transactionFormatter

import (
	"strconv"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
	"swap.io-agent/src/blockchain/journal"
)

func AddMinerRewardToJournal(
	tx *nodeApi.Transaction,
	block *nodeApi.Block,
	journal *journal.Journal,
) {
	reward := tx.Outputs[0].Value
	for _, blockTx := range block.Txs {
		reward -= blockTx.Fee
	}

	journal.Add(HSN, blockchain.Spend{
		Wallet: blockchain.BLOCK_REWARD_CREATER_ADDRESS,
		Label:  blockchain.SPEND_LABEL_BLOCK_REWARD,
		Value:  strconv.Itoa(-reward),
	})
	journal.Add(HSN, blockchain.Spend{
		Wallet: tx.Outputs[0].Address,
		Label: blockchain.SPEND_LABEL_TRANSFER,
		Value: strconv.Itoa(reward),
	})
}
