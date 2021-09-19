package blockchain

import (
	"swap.io-agent/src/blockchain/ethereum/nodeApi"
)

type IBlockchinApi interface {
	GetBlockByIndex(index int) (*Block, int)
	GetTransactionByHash(hash string) (*Transaction, int)
}
type IFormatter interface {
	FormatTransaction(
		blockTransaction *nodeApi.BlockTransaction,
		block *nodeApi.Block,
	) (*Transaction, error)
	FormatTransactionFromHash(hash string) (*Transaction, error)
}
type ISynchronizer interface {
	GetAddressFirstCursorData(
		address string,
	) (*CursorTransactions, int)
	GetCursorData(
		cursor string,
	) (*CursorTransactions, int)
}
