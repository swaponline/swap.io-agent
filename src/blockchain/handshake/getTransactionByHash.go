package handshake

import "swap.io-agent/src/blockchain"

func (a *Api) GetTransactionByHash(
	hash string,
) (
	*blockchain.Transaction,
	int,
) {
	nodeTx, err := a.nodeApi.GetTransactionByHash(hash)
	if err != blockchain.ApiRequestSuccess {
		return nil, err
	}
	nodeBlock, err := a.nodeApi.GetBlockByIndex(nodeTx.Height)
	if err != blockchain.ApiRequestSuccess {
		return nil, err
	}

	rewardTx, minerAddress := a.formatter.GetRewardTx(nodeBlock)

	if nodeTx.Fee == 0 {
		return rewardTx, blockchain.ApiRequestSuccess
	} else {
		tx := a.formatter.FormatTransaction(nodeTx, minerAddress)
		return tx, blockchain.ApiRequestSuccess
	}
}
