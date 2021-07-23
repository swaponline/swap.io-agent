package blockchain

import (
	"swap.io-agent/src/common/functions"
	"swap.io-agent/src/env"
	"swap.io-agent/src/levelDbStore"
	"sync"
	"time"
)

type Synchronizer struct {
	apiKey string
	formatter Formatter
	store levelDbStore.ITransactionsStore
}
type SynchronizerConfig struct {
	apiKey string
	Formatter Formatter
	Store levelDbStore.ITransactionsStore
}

func InitialiseSynchronizer(config SynchronizerConfig) *Synchronizer {
	return &Synchronizer{
		apiKey: config.apiKey,
		store: config.Store,
	}
}

func (s *Synchronizer) SynchronizeAddress(
	address string,
	startTime int,
	endTime int,
)([]*Transaction, error) {
	transactionsHash, err := s.store.GetAddressTransactionsHash(
		address,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}

	transactions := make([]*Transaction, len(transactionsHash))
	err = functions.ForWidthBreaks(
		len(transactionsHash),
		env.BLOCKCHAIN_REQUESTS_LIMIT,
		time.Second,
		func(wg *sync.WaitGroup, step int) error {
			//transaction, err := s.formatter.FormatTransactionFromHash(
			//	s.apiKey,
			//	transactionsHash[step],
			//)
			//if err != nil {
			//	return err
			//}
			//transactions[step] = transaction

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return transactions,nil
}

func Start() {}
func Stop() error {
	return nil
}
func Status() error {
	return nil
}