package synchronizer

import (
	"log"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/levelDbStore"
)

type Synchronizer struct {
	Api       blockchain.IBlockchainApi
	formatter blockchain.IFormatter
	store     levelDbStore.ITransactionsStore
}
type SynchronizerConfig struct {
	Api       blockchain.IBlockchainApi
	Formatter blockchain.IFormatter
	Store     levelDbStore.ITransactionsStore
}

func InitialiseSynchronizer(config SynchronizerConfig) *Synchronizer {
	return &Synchronizer{
		Api:       config.Api,
		formatter: config.Formatter,
		store:     config.Store,
	}
}

func (s *Synchronizer) GetAddressFirstCursorData(
	address string,
) (*blockchain.CursorTransactions, int) {
	cursor, err := s.store.GetCursorFromAddress(address)
	if err != nil {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}

	return s.GetCursorData(cursor)
}
func (s *Synchronizer) GetCursorData(
	cursor string,
) (*blockchain.CursorTransactions, int) {
	cursorHashes, err := s.store.GetCursorTransactionHashes(cursor)
	if err != nil {
		log.Println(err)
		return nil, blockchain.ApiParseBodyError
	}

	txs := make([]*blockchain.Transaction, 0)
	for _, hash := range cursorHashes.Hashes {
		tx, err := s.Api.GetTransactionByHash(hash)
		if err != blockchain.ApiRequestSuccess {
			return nil, err
		}

		txs = append(txs, tx)
	}

	return &blockchain.CursorTransactions{
		Cursor:       cursor,
		NextCursor:   cursorHashes.NextCursor,
		Transactions: txs,
	}, blockchain.ApiRequestSuccess
}

func (*Synchronizer) Start() {}
func (*Synchronizer) Stop() error {
	return nil
}
func (*Synchronizer) Status() error {
	return nil
}
