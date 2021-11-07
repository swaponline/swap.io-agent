package handshake

import (
	"swap.io-agent/src/blockchain"
	transactionFormatter "swap.io-agent/src/blockchain/handshake/formatter"
)

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

	minerAddress := transactionFormatter.GetBlockMinderAddress(nodeBlock)
	tx := a.formatter.FormatTransaction(nodeTx, nodeBlock, minerAddress)

	return tx, blockchain.ApiRequestSuccess
}
