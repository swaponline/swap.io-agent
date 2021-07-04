package ethereum

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
	"swap.io-agent/src/blockchainHandlers/ethereum/api/ethercsan"
	"swap.io-agent/src/levelDbStore"
	"sync"
	"time"
)

func (indexer *BlockchainIndexer) RunScanner() {
	isSynchronize := false
	requestsStepLen := 4
	for {
		waits := new(sync.WaitGroup)
		waits.Add(requestsStepLen)

		lastGetBlockIndex := indexer.lastBlock
		buf := make(map[string][]string)
		lockerChange := new(sync.Mutex)
		for t:=1; t<=requestsStepLen; t++ {
			go func(blockIndex int) {
				block, err := ethercsan.GetBlockByIndex(
					indexer.apiKey,
					blockIndex,
				)
				if err == ethercsan.NotExistBlockError {
					// if notExistBlockErr then scan synchronize
					lockerChange.Lock()
					if !isSynchronize {
						isSynchronize = true
						close(indexer.isSynchronize)
					}
					lockerChange.Unlock()
					waits.Done()
					return
				} else if err != ethercsan.RequestSuccess {
					log.Panicln(err, "error code ethercsan")
				}

				lockerChange.Lock()
				<-time.After(time.Second)
				err = setAllSpendAddressesForTransactions(
					indexer.apiKey,
					block.Transactions,
					requestsStepLen,
				)
				if err != ethercsan.RequestSuccess {
					log.Panicln(
						"not indexing all spend transactions(contracts) errors",
						err,
					)
				}
				indexingTransactions(&buf, block.Transactions)
				if lastGetBlockIndex < blockIndex {
					lastGetBlockIndex = blockIndex
				}
				lockerChange.Unlock()
				waits.Done()
			}(indexer.lastBlock + t)
		}
		// pending all done requests
		waits.Wait()

		// last block not change when blockchain synchronize
		if indexer.lastBlock != lastGetBlockIndex {
			err := writeIndexTransactionToDb(
				indexer.db,
				&buf,
				lastGetBlockIndex,
			)
			if err != nil {
				log.Panicf("not write block transaction %v", err)
			}

			indexer.lastBlock = lastGetBlockIndex
			log.Printf("last block indexed - %v", indexer.lastBlock)
		}

		//pending
		<- time.After(time.Second)
		// blockchain synchronize stop synchronize
		if isSynchronize {
			break
		}
	}

	log.Println("Blockchain synchronize ***")
	log.Println("last block - indexed", indexer.lastBlock)

	for {
		nextBlock := indexer.lastBlock + 1
		block, err := ethercsan.GetBlockByIndex(
			indexer.apiKey,
			nextBlock,
		)
		err = setAllSpendAddressesForTransactions(
			indexer.apiKey,
			block.Transactions,
			requestsStepLen,
		)
		if err != ethercsan.RequestSuccess {
			log.Panicf(
				"not indexing all spend transactions(contracts) errors",
			)
		}
		switch err {
			case ethercsan.RequestSuccess: {
				indexedTransactions := make(map[string][]string)
				indexingTransactions(&indexedTransactions, block.Transactions)
				err := writeIndexTransactionToDb(
					indexer.db,
					&indexedTransactions,
					nextBlock,
				)
				if err != nil {
					log.Panicf("not write block transaction %v", err)
				}
				indexer.lastBlock = nextBlock
				log.Printf("block indexed - %v", indexer.lastBlock)
			}
			case ethercsan.NotExistBlockError: {}
			default: log.Panicln(err, "error code ethercsan")
		}
		<-time.After(time.Second * 7)
	}
}

func setAllSpendAddressesForTransactions(
	apiKey string,
	transactions []ethercsan.BlockTransaction,
	requestLimitSecond int,
) int {
	err := ethercsan.RequestSuccess
	for t:=0; t<len(transactions); t+=requestLimitSecond {
		wg := new(sync.WaitGroup)
		steps := requestLimitSecond
		if len(transactions) - t < requestLimitSecond {
			steps = len(transactions) - t
		}
		wg.Add(steps)
		for r:=0; r<steps; r++ {
			go func(index int) {
				addresses, reqErr := ethercsan.AllSpendAddressesTransaction(
					apiKey,
					&transactions[index],
				)
				if reqErr != ethercsan.RequestSuccess {
					err = reqErr
				}

				transactions[index].AllSpendAddresses = addresses
				wg.Done()
			}(t+r)
		}
		wg.Wait()
		<-time.After(time.Second)
	}

	if err != ethercsan.RequestSuccess {
		return err
	}

	return ethercsan.RequestSuccess
}

func indexingTransactions(
	buf *map[string][]string,
	transactions []ethercsan.BlockTransaction,
) {
	// address -> transactions
	bufValue := *buf
	for _, transaction := range transactions {
		for _, address := range transaction.AllSpendAddresses {
			bufValue[address] = append(bufValue[address], transaction.Hash)
		}
	}
}

func writeIndexTransactionToDb(
	db *leveldb.DB,
	indexedTransactions *map[string][]string,
	indexBlock int,
) error {
	bdTransaction, err := db.OpenTransaction()
	if err != nil {
		return err
	}
	for address, transactions := range *indexedTransactions {
		// push to back address transaction
		err = levelDbStore.ArrayStringPush(
			bdTransaction, address, transactions,
		)
		if err != nil {
			bdTransaction.Discard()
			return err
		}
	}
	// update lastBlock
	err = bdTransaction.Put(
		lastBlockKey,
		[]byte(strconv.Itoa(indexBlock)),
		nil,
	)
	if err != nil {
		bdTransaction.Discard()
		return err
	}

	// commit transaction
	err = bdTransaction.Commit()
	if err != nil {
		bdTransaction.Discard()
		return err
	}

	return nil
}