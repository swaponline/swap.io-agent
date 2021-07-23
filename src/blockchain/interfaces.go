package blockchain

import (
	"context"
	"swap.io-agent/src/blockchain/ethereum/api"
)

type Broadcaster interface {
	broadcast(hex string) error
}
type Notifier interface {
	notify(ctx context.Context, address string) error
}
type Formatter interface {
	FormatTransaction(
		blockTransaction *api.BlockTransaction,
		block *api.Block,
	) (*Transaction, error)
	FormatTransactionFromHash(hash string) (*Transaction, error)
}