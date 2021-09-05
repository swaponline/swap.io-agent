package indexer

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/levelDbStore"
)

type Indexer struct {
	api               blockchain.IBlockchinApi
	transactionsStore levelDbStore.ITransactionsStore
	isSynchronize     chan struct{}
	NewTransactions   chan *blockchain.Transaction
}
type IndexerConfig struct {
	Api               blockchain.IBlockchinApi
	TransactionsStore *levelDbStore.TransactionsStore
}

func InitializeIndexer(config IndexerConfig) *Indexer {
	indexer := &Indexer{
		api:               config.Api,
		transactionsStore: config.TransactionsStore,
		isSynchronize:     make(chan struct{}),
		NewTransactions:   make(chan *blockchain.Transaction),
	}

	go indexer.RunScanner()
	<-indexer.isSynchronize

	return indexer
}

func (*Indexer) Start() {}
func (*Indexer) Stop() error {
	return nil
}
func (*Indexer) Status() error {
	return nil
}
