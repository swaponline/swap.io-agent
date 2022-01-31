package fullNodeApi

import (
	"log"
	"swap.io-agent/src/blockchain"
	"time"
)

func (fn *FullNodeApi) WatchMempoolTxs(out chan *blockchain.Transaction) error {
	buf := make(map[string]struct{})
	for {
		txsHashes, err := fn.GetMempool()
		if err != nil {
			log.Println("ERROR GET MEMPOOLTXS:", err)
			continue
		}
		for _, txHash := range txsHashes {
			if _, exist := buf[txHash]; !exist {
				for {
					nodeTx, err := fn.GetTransactionByHash(txHash)
					if err != blockchain.ApiRequestSuccess {
						log.Println("ERROR GET TRANSACTION BY HASH:", err)
						<-time.After(time.Second)
						continue
					}

					tx := fn.formatter.FormatTransaction(nodeTx, nil, "")
					out <- tx

					buf[txHash] = struct{}{}
					break
				}
			}
		}

		<-time.After(time.Second)
	}
}
