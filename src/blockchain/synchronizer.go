package blockchain

import (
	"log"
	"os"
	"strconv"
	"swap.io-agent/src/blockchain/ethereum/transactionFormater"
	"swap.io-agent/src/common/functions"
	"swap.io-agent/src/levelDbStore"
	"sync"
	"time"
)

type Synchronizer struct {
	apiKey string
	store levelDbStore.ITransactionsStore
}
type SynchronizerConfig struct {
	apiKey string
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
	transactions := make([]*Transaction, len(transactionsHash))
	if err != nil {
		return nil, err
	}
	requestsStepSize, err := strconv.ParseInt(
		os.Getenv("BLOCKCHAIN_REQUESTS_LIMIT"),
		0,
		64,
	)
	if err != nil {
		log.Panicln("set BLOCKCHAIN_REQUESTS_LIMIT in env")
	}
	err = functions.ForWidthBreaks(
		len(transactionsHash),
		int(requestsStepSize),
		time.Second,
		func(wg *sync.WaitGroup, step int) error {
			transaction, err := transactionFormater.FormatTransactionFromHash(
				s.apiKey,
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

func Start() {}
func Stop() error {
	return nil
}
func Status() error {
	return nil
}