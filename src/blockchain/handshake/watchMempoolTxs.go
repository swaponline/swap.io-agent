package handshake

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/handshake/nodeApi"
)

func (a *Api) WatchMempoolTxs(out chan *blockchain.Transaction) error {
	nonFormatTxChan := make(chan *nodeApi.Transaction)
	go func() {
		a.nodeApi.WatchMempoolTxs(nonFormatTxChan)
	}()
	for {
		nonFormatTx := <-nonFormatTxChan
		out <- a.formatter.FormatTransaction(nonFormatTx, nil, "")
	}
}
