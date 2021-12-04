package indexer

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/queueEvents"
	"swap.io-agent/src/subscribersManager"
)

type Indexer struct {
	api               blockchain.IBBlockchainApi
	transactionsStore levelDbStore.ITransactionsStore
	queueEvents       *queueEvents.QueueEvents
	subscribesManager *subscribersManager.SubscribesManager
	isSynchronize     chan struct{}
	NewTransactions   chan *blockchain.Transaction
}
type IndexerConfig struct {
	Api               blockchain.IBBlockchainApi
	TransactionsStore *levelDbStore.TransactionsStore
	QueueEvents       *queueEvents.QueueEvents
	SubscribesManager *subscribersManager.SubscribesManager
}

func InitializeIndexer(config IndexerConfig) *Indexer {
	indexer := &Indexer{
		api:               config.Api,
		transactionsStore: config.TransactionsStore,
		queueEvents:       config.QueueEvents,
		subscribesManager: config.SubscribesManager,
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
