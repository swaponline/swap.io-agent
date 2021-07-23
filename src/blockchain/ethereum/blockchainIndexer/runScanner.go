package ethereum

import (
	"log"
	"strconv"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/api"
	"swap.io-agent/src/blockchain/ethereum/api/ethercsan"
	"swap.io-agent/src/env"
	"sync"
	"time"
)

type indexedBlock struct {
	transactions *map[string][]string
	index int
	timestamp int
}

func (indexer *BlockchainIndexer) RunScanner() {
	isSynchronize := false
	requestsStepLen := env.BLOCKCHAIN_REQUESTS_LIMIT

	for {
		waits := new(sync.WaitGroup)
		waits.Add(requestsStepLen)

		bufIndexedBlocks := make([]indexedBlock, requestsStepLen)
		lockerChange := new(sync.Mutex)
		for t:=0; t<requestsStepLen; t++ {
			go func(blockIndex int, ItemIndexInBufIndexedBlocks int) {
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
				blockTimestamp, errConv := strconv.ParseInt(block.Timestamp, 0, 64)
				if errConv != nil {
					log.Panicln(
						"block timestamp invalid",
						errConv,
					)
				}

				lockerChange.Lock()
				<-time.After(time.Second)
				transactions, fErr := formattedBlockTransactions(
					indexer.formatter,
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
				indexedTransactions := make(map[string][]string)
				indexingTransactions(&indexedTransactions, transactions)
				bufIndexedBlocks[ItemIndexInBufIndexedBlocks] = indexedBlock{
					transactions: &indexedTransactions,
					index: blockIndex,
					timestamp: int(blockTimestamp),
				}
				lockerChange.Unlock()

				waits.Done()
			}(indexer.transactionsStore.GetLastTransactionBlock()+1+t, t)
		}
		// pending all done requests
		waits.Wait()

		for _, indexedBlock := range bufIndexedBlocks {
			err := indexer.transactionsStore.WriteLastIndexedBlockTransactions(
				indexedBlock.transactions,
				indexedBlock.index,
				indexedBlock.timestamp,
			)
			if err != nil {
				log.Panicf("not write block transaction %v", err)
			}
			log.Printf(
				"last block indexed - %v",
				indexer.transactionsStore.GetLastTransactionBlock(),
			)
		}

		//pending
		<- time.After(time.Second)
		// blockchain synchronize stop synchronize
		if isSynchronize {
			break
		}
	}

	log.Println("Blockchain synchronize ***")
	log.Println(
		"last block - indexed",
		indexer.transactionsStore.GetLastTransactionBlock(),
	)

	for {
		nextBlock := indexer.transactionsStore.GetLastTransactionBlock() + 1
		block, err := ethercsan.GetBlockByIndex(
			indexer.apiKey,
			nextBlock,
		)
		switch err {
			case ethercsan.RequestSuccess: {
				blockTimestamp, errConv := strconv.Atoi(block.Timestamp)
				if errConv != nil {
					log.Panicln(
						"block timestamp invalid",
						errConv,
					)
				}

				transactions, fErr := formattedBlockTransactions(
					indexer.formatter,
					block.Transactions,
					block,
					requestsStepLen,
				)
				if fErr != nil {
					log.Panicf("block transaction not formatted err - %v", err)
				}
				indexedTransactions := make(map[string][]string)
				indexingTransactions(&indexedTransactions, transactions)

				err := indexer.transactionsStore.WriteLastIndexedBlockTransactions(
					&indexedTransactions,
					nextBlock,
					blockTimestamp,
				)
				if err != nil {
					log.Panicf("not write block transaction %v", err)
				}
				log.Printf(
					"block indexed - %v",
					indexer.transactionsStore.GetLastTransactionBlock(),
				)
			}
			case ethercsan.NotExistBlockError: {}
			default: log.Panicln(err, "error code ethercsan request ethercsan.GetBlockByIndex")
		}
		<-time.After(time.Second * 1)
	}
}

func formattedBlockTransactions(
	formatter blockchain.Formatter,
	transactions []api.BlockTransaction,
	block *api.Block,
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
				transaction, fError := formatter.FormatTransaction(
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
