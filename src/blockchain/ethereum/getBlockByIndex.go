package ethereum

import (
	"swap.io-agent/src/blockchain"
)

func (a *Api) GetBlockByIndex(index int) (*blockchain.Block, int) {
	apiBlock, err := a.nodeApi.GetBlockByIndex(index)
	if err != blockchain.ApiRequestSuccess {
		return nil, err
	}
	block := blockchain.Block{}
	for _, apiTransaction := range apiBlock.Transactions {
		if transaction, err := a.formatter.FormatTransaction(
			&apiTransaction,
			apiBlock,
		); err == nil {
			block.Transactions = append(
				block.Transactions,
				transaction,
			)
		} else {
			return nil, blockchain.ApiRequestError
		}
	}
	return &block, blockchain.ApiRequestSuccess
}
