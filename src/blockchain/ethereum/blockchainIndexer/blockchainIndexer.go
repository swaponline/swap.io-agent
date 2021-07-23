package ethereum

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/env"
	"swap.io-agent/src/levelDbStore"
)

type BlockchainIndexer struct {
	formatter         blockchain.Formatter
	transactionsStore levelDbStore.ITransactionsStore
	apiKey            string
	isSynchronize     chan struct{}
	newTransactions   chan blockchain.Transaction
}

type BlockchainIndexerConfig struct {
	Formatter blockchain.Formatter
	TransactionsStore levelDbStore.ITransactionsStore
}

func InitializeIndexer(config BlockchainIndexerConfig) *BlockchainIndexer {
	bi := &BlockchainIndexer{
		transactionsStore: config.TransactionsStore,
		apiKey: env.ETHERSCAN_API_KEY,
		isSynchronize: make(chan struct{}),
		newTransactions: make(chan blockchain.Transaction),
	}
	bi.RunScanner()

	return bi
}

func (_ *BlockchainIndexer) Start()  {}
func (_ *BlockchainIndexer) Stop() error {
	return nil
}
func (_ *BlockchainIndexer) Status() error {
	return nil
}