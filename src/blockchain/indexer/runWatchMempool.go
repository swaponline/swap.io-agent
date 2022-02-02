package indexer

import (
	"log"
)

func (i *Indexer) RunWatchMempool() {
	err := i.api.WatchMempoolTxs(i.NewMempoolTransactions)
	if err != nil {
		log.Println("ERROR WATCHTXMEMPOOL", err)
	}
}
