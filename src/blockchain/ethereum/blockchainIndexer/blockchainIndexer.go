package ethereum

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum"
	"swap.io-agent/src/levelDbStore"
)

type BlockchainIndexer struct {
	api               ethereum.IGeth
	formatter         blockchain.IFormatter
	transactionsStore levelDbStore.ITransactionsStore
	isSynchronize     chan struct{}
	NewTransactions   chan *blockchain.Transaction
}

type BlockchainIndexerConfig struct {
	Api ethereum.IGeth
	Formatter blockchain.IFormatter
	TransactionsStore levelDbStore.ITransactionsStore
}

func InitializeIndexer(config BlockchainIndexerConfig) *BlockchainIndexer {
	bi := &BlockchainIndexer{
		api: config.Api,
		formatter: config.Formatter,
		transactionsStore: config.TransactionsStore,
		isSynchronize: make(chan struct{}),
		NewTransactions: make(chan *blockchain.Transaction),
	}

	go bi.RunScanner()
	<-bi.isSynchronize

	return bi
}

func (*BlockchainIndexer) Start()  {}
func (*BlockchainIndexer) Stop() error {
	return nil
}
func (*BlockchainIndexer) Status() error {
	return nil
}