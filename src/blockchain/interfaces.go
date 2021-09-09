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
	) (*CursorTransactions, error)
	GetCursorData(
		cursor string,
	) (*CursorTransactions, error)
}
type ISubscribeManager interface {
	SubscribeUserToAddress(
		userId string,
		address string,
	) error
	ClearAllUserSubscriptions(
		userId string,
	) error
}
