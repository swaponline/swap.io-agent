package ethereum

import ethereum "swap.io-agent/src/blockchainHandlers/ethereum/blockchainIndexer"

type TransactionNotifier struct {
	activeAddresses struct{}
	indexer *ethereum.BlockchainIndexer
}