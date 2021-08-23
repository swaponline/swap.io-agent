package synchronizer

import "swap.io-agent/src/blockchain"

type CursorTransactions struct {
	Cursor string
	NextCursor string
	Transactions []*blockchain.Transaction
}