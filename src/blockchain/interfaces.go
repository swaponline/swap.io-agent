package blockchain

import (
	"swap.io-agent/src/blockchain/ethereum/api"
)

type IFormatter interface {
	FormatTransaction(
		blockTransaction *api.BlockTransaction,
		block *api.Block,
	) (*Transaction, error)
	FormatTransactionFromHash(hash string) (*Transaction, error)
}
type ISynchronizer interface {
	SynchronizeAddress(
		userId string,
		address string,
		startTime int,
		endTime int,
	)([]*Transaction, error)
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