package handshake

import "swap.io-agent/src/blockchain"

func (a *Api) GetBlockByIndex(index int) (*blockchain.Block, int) {
	nodeBlock, err := a.nodeApi.GetBlockByIndex(index)
	if err != blockchain.ApiRequestSuccess {
		return nil, err
	}
	block := a.formatter.FormatBlock(nodeBlock)

	return block, blockchain.ApiRequestSuccess
}
