package handshake

import "swap.io-agent/src/blockchain"

func (a *Api) WatchMempoolTxs(out chan *blockchain.Transaction) error {
	return a.nodeApi.WatchMempoolTxs(out)
}