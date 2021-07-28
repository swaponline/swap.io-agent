package synchronizer

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/common/Set"
	"swap.io-agent/src/common/functions"
	"swap.io-agent/src/env"
	"swap.io-agent/src/levelDbStore"
	"sync"
	"time"
)

type Synchronizer struct {
	formatter blockchain.IFormatter
	store     levelDbStore.ITransactionsStore
	sendedTransactions map[string]Set.Set
}
type SynchronizerConfig struct {
	Formatter blockchain.IFormatter
	Store     levelDbStore.ITransactionsStore
}

func InitialiseSynchronizer(config SynchronizerConfig) *Synchronizer {
	return &Synchronizer{
		formatter: config.Formatter,
		store: config.Store,
	}
}

func (s *Synchronizer) SynchronizeAddress(
	userId string,
	address string,
	startTime int,
	endTime int,
)([]*blockchain.Transaction, error) {
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

	return transactions,nil
}

func (_ *Synchronizer) Start() {}
func (_ *Synchronizer) Stop() error {
	return nil
}
func (_ *Synchronizer) Status() error {
	return nil
}