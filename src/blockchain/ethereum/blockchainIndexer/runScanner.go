package ethereum

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/api/ethercsan"
	"swap.io-agent/src/blockchain/ethereum/transactionFormater"
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
				transactions, fErr := formattedBlockTransactions(
					indexer.apiKey,
					block.Transactions,
					block,
					requestsStepLen,
				)
				if fErr != nil {
					log.Panicln(
						"not indexing all spend transactions(contracts) errors",
						fErr,
					)
				}
				indexingTransactions(&buf, transactions)
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
		if err != ethercsan.RequestSuccess {
			log.Panicf(
				"not indexing all spend transactions(contracts) errors",
			)
		}
		transactions, fErr := formattedBlockTransactions(
			indexer.apiKey,
			block.Transactions,
			block,
			requestsStepLen,
		)
		if fErr != nil {
			log.Panicf("block transaction not formatted err - %v", err)
		}

		switch err {
			case ethercsan.RequestSuccess: {
				indexedTransactions := make(map[string][]string)
				indexingTransactions(&indexedTransactions, transactions)
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

func formattedBlockTransactions(
	apiKey string,
	transactions []ethercsan.BlockTransaction,
	block *ethercsan.Block,
	requestLimitSecond int,
) ([]blockchain.Transaction, error) {
	var err error
	formattedTransactions := make([]blockchain.Transaction, 0)
	for t:=0; t<len(transactions); t+=requestLimitSecond {
		wg := new(sync.WaitGroup)
		steps := requestLimitSecond
		if len(transactions) - t < requestLimitSecond {
			steps = len(transactions) - t
		}
		wg.Add(steps)
		for r:=0; r<steps; r++ {
			go func(index int) {
				transaction, fError := transactionFormater.FormatTransaction(
					apiKey,
					&transactions[index],
					block,
				)
				defer wg.Done()
				if fError != nil {
					err = fError
					return
				}
				formattedTransactions = append(formattedTransactions, *transaction)
			}(t+r)
		}
		wg.Wait()
		<-time.After(time.Second)
	}

	return formattedTransactions, err
}

func indexingTransactions(
	buf *map[string][]string,
	transactions []blockchain.Transaction,
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