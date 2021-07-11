package ethereum

import ethereum "swap.io-agent/src/blockchain/ethereum/blockchainIndexer"

type TransactionNotifier struct {
	activeAddresses struct{}
	indexer *ethereum.BlockchainIndexer
}