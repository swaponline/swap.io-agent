package ethereum

import (
	"os"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/levelDbStore"
)

var lastBlockKey = []byte("lastBlock")
var dbpath = "./blockchainIndexes/ethereum"

type BlockchainIndexer struct {
	transactionsStore levelDbStore.ITransactionsStore
	apiKey            string
	isSynchronize     chan struct{}
	newTransactions   chan blockchain.Transaction
}

type BlockchainIndexerConfig struct {
	TransactionsStore levelDbStore.ITransactionsStore
}

func InitializeIndexer(config BlockchainIndexerConfig) *BlockchainIndexer {
	return &BlockchainIndexer{
		transactionsStore: config.TransactionsStore,
		apiKey: os.Getenv("ETHERSCAN_API_KEY"),
		isSynchronize: make(chan struct{}),
		newTransactions: make(chan blockchain.Transaction),
	}
}