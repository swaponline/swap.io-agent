package transactionFormatter

import "swap.io-agent/src/blockchain/handshake/nodeApi"

func GetBlockMinderData(
	block *nodeApi.Block,
) (
	address string,
	allFee int,
	blockReward int,
	rewordTx *nodeApi.Transaction,
) {
	for _, tx := range block.Txs {
		if tx.Fee == 0 {
			address = tx.Outputs[0].Address
			rewordTx = &tx
		}
		allFee += tx.Fee
	}
	return address, allFee, rewordTx.Outputs[0].Value - allFee, rewordTx
}
