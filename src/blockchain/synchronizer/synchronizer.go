package synchronizer

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/common/functions"
	"swap.io-agent/src/env"
	"swap.io-agent/src/levelDbStore"
)

type Synchronizer struct {
	Api       blockchain.IBlockchinApi
	formatter blockchain.IFormatter
	store     levelDbStore.ITransactionsStore
	//sendedTransactions map[string]Set.Set
}
type SynchronizerConfig struct {
	Api       blockchain.IBlockchinApi
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
) (*blockchain.CursorTransactions, error) {
	cursor, err := s.store.GetCursorFromAddress(address)
	if err != nil {
		return nil, err
	}

	return s.GetCursorData(cursor)
}
func (s *Synchronizer) GetCursorData(
	cursor string,
) (*blockchain.CursorTransactions, error) {
	cursorHashes, err := s.store.GetCursorTransactionHashes(cursor)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	e, _ := json.Marshal(cursorHashes)
	log.Println(string(e))

	txs := make([]*blockchain.Transaction, 0)
	for _, hash := range cursorHashes.Hashes {
		tx, err := s.Api.GetTransactionByHash(hash)
		if err != blockchain.ApiRequestSuccess {
			return nil, fmt.Errorf("err - %v", err)
		}
		txs = append(txs, tx)
	}
	log.Println(cursor)

	return &blockchain.CursorTransactions{
		Cursor:       cursor,
		NextCursor:   cursorHashes.NextCursor,
		Transactions: txs,
	}, nil
}
func (s *Synchronizer) SynchronizeAddress(
	userId string,
	address string,
	startTime int,
	endTime int,
) ([]*blockchain.Transaction, error) {
	transactionsHash, err := s.store.GetAddressTransactionsHash(
		address,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}

	transactions := make([]*blockchain.Transaction, len(transactionsHash))
	err = functions.ForWidthBreaks(
		len(transactionsHash),
		env.BLOCKCHAIN_REQUESTS_LIMIT,
		time.Second,
		func(wg *sync.WaitGroup, step int) error {
			transaction, err := s.formatter.FormatTransactionFromHash(
				transactionsHash[step],
			)
			if err != nil {
				return err
			}
			transactions[step] = transaction

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (*Synchronizer) Start() {}
func (*Synchronizer) Stop() error {
	return nil
}
func (*Synchronizer) Status() error {
	return nil
}
