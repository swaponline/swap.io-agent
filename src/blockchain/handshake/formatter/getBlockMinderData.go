package transactionFormatter

import "swap.io-agent/src/blockchain/handshake/nodeApi"

func GetBlockMinderData(
	block *nodeApi.Block,
) (
	address string,
	allFee int,
	blockReward int,
	indexRewordTxInBlockTxs int,
) {
	rewardTxValue := 0
	for index, tx := range block.Txs {
		if tx.Fee == 0 {
			address = tx.Outputs[0].Address
			rewardTxValue = tx.Outputs[0].Value
			indexRewordTxInBlockTxs = index
		}
		allFee += tx.Fee
	}
	return address, allFee, rewardTxValue - allFee, indexRewordTxInBlockTxs
}
