package ethereum

import (
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/api"
	"sync"
)

type indexedBlock struct {
	transactions *map[string][]string
	index int
	timestamp int
}

func (indexer *BlockchainIndexer) RunScanner() {
	//requestsLimit := env.BLOCKCHAIN_REQUESTS_LIMIT
	//
	//for {
	//	lockerChange := new(sync.Mutex)
	//	waits := new(sync.WaitGroup)
	//	waits.Add(requestsLimit)
	//}

	//isSynchronize := false
	//requestsStepLen := env.BLOCKCHAIN_REQUESTS_LIMIT
	//allScannedTransactions := 0
	//
	//for {
	//	waits := new(sync.WaitGroup)
	//	waits.Add(requestsStepLen)
	//
	//	bufIndexedBlocks := make([]*indexedBlock, requestsStepLen)
	//	lockerChange := new(sync.Mutex)
	//	transactionsScannedBeforeStep := allScannedTransactions
	//	for t:=0; t<requestsStepLen; t++ {
	//		go func(blockIndex int, ItemIndexInBufIndexedBlocks int) {
	//			defer waits.Done()
	//			block, err := indexer.api.GetBlockByIndex(
	//				blockIndex,
	//			)
	//			if err == api.NotExistBlockError {
	//				// if notExistBlockErr then scan synchronize
	//				lockerChange.Lock()
	//				if !isSynchronize {
	//					isSynchronize = true
	//				}
	//				lockerChange.Unlock()
	//				return
	//			} else if err != api.RequestSuccess {
	//				log.Panicln(err, "error code ethercsan")
	//			}
	//			blockTimestamp, errConv := strconv.ParseInt(block.Timestamp, 0, 64)
	//			if errConv != nil {
	//				log.Panicln(
	//					"block timestamp invalid",
	//					errConv,
	//				)
	//			}
	//
	//			lockerChange.Lock()
	//			//<-time.After(time.Second)
	//			allScannedTransactions+=len(block.Transactions)
	//			transactions, fErr := formattedBlockTransactions(
	//				indexer.formatter,
	//				block.Transactions,
	//				block,
	//				requestsStepLen,
	//			)
	//			if fErr != nil {
	//				log.Panicln(
	//					"not indexing all spend transactions(contracts) errors",
	//					fErr,
	//				)
	//			}
	//			indexedTransactions := make(map[string][]string)
	//			indexingTransactions(&indexedTransactions, transactions)
	//			bufIndexedBlocks[ItemIndexInBufIndexedBlocks] = &indexedBlock{
	//				transactions: &indexedTransactions,
	//				index: blockIndex,
	//				timestamp: int(blockTimestamp),
	//			}
	//			lockerChange.Unlock()
	//		}(indexer.transactionsStore.GetLastTransactionBlock()+1+t, t)
	//	}
	//	// pending all done requests
	//	waits.Wait()
	//
	//	for _, indexedBlock := range bufIndexedBlocks {
	//		if indexedBlock != nil {
	//			err := indexer.transactionsStore.WriteLastIndexedBlockTransactions(
	//				indexedBlock.transactions,
	//				indexedBlock.index,
	//				indexedBlock.timestamp,
	//			)
	//			if err != nil {
	//				log.Panicf("not write block transaction %v", err)
	//			}
	//			log.Printf(
	//				"last block indexed - %v",
	//				indexer.transactionsStore.GetLastTransactionBlock(),
	//			)
	//		}
	//	}
	//	log.Printf(
	//		"writed transactions - %v",
	//		allScannedTransactions - transactionsScannedBeforeStep,
	//	)
	//	log.Printf("all scanned transactions - %v", allScannedTransactions)
	//
	//	//pending
	//	//<- time.After(time.Second)
	//	// blockchain synchronize stop synchronize
	//	if isSynchronize {
	//		break
	//	}
	//}
	//
	//
	//log.Println("Blockchain synchronize ***")
	//log.Println(
	//	"last block - indexed",
	//	indexer.transactionsStore.GetLastTransactionBlock(),
	//)
	//close(indexer.isSynchronize)
	//
	//for {
	//	nextBlock := indexer.transactionsStore.GetLastTransactionBlock() + 1
	//	block, err := indexer.api.GetBlockByIndex(
	//		nextBlock,
	//	)
	//	switch err {
	//		case api.RequestSuccess: {
	//			blockTimestamp, errConv := strconv.ParseInt(
	//				block.Timestamp,
	//				0,
	//				64,
	//			)
	//			if errConv != nil {
	//				log.Panicln(
	//					"block timestamp invalid",
	//					errConv,
	//				)
	//			}
	//
	//			//<-time.After(time.Second)
	//
	//			transactions, fErr := formattedBlockTransactions(
	//				indexer.formatter,
	//				block.Transactions,
	//				block,
	//				requestsStepLen,
	//			)
	//			if fErr != nil {
	//				log.Panicf("block transaction not formatted err - %v", err)
	//			}
	//			indexedTransactions := make(map[string][]string)
	//			indexingTransactions(&indexedTransactions, transactions)
	//			err := indexer.transactionsStore.WriteLastIndexedBlockTransactions(
	//				&indexedTransactions,
	//				nextBlock,
	//				int(blockTimestamp),
	//			)
	//			if err != nil {
	//				log.Panicf("not write block transaction %v", err)
	//			}
	//			for _, transaction := range transactions {
	//				indexer.NewTransactions <- transaction
	//			}
	//			log.Printf(
	//				"block indexed - %v | size - %v",
	//				indexer.transactionsStore.GetLastTransactionBlock(),
	//				len(block.Transactions),
	//			)
	//		}
	//		case api.NotExistBlockError: {}
	//		default: log.Panicln(err, "error code ethercsan request ethercsan.GetBlockByIndex")
	//	}
	//	//<-time.After(time.Second * 1)
	//}
}

func formattedBlockTransactions(
	formatter blockchain.IFormatter,
	transactions []api.BlockTransaction,
	block *api.Block,
	requestLimitSecond int,
) ([]blockchain.Transaction, error) {
	var err error
	formattedTransactions := make([]blockchain.Transaction, 0)
	wg := new(sync.WaitGroup)
	wg.Add(len(transactions))
	for r:=0; r<len(transactions); r++ {
		go func(index int) {
			defer wg.Done()
			transaction, fError := formatter.FormatTransaction(
				&transactions[index],
				block,
			)
			if fError != nil {
				err = fError
				return
			}
			formattedTransactions = append(formattedTransactions, *transaction)
		}(r)
	}
	wg.Wait()

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
			if len(address) == 0 {
				continue
			}
			bufValue[address] = append(bufValue[address], transaction.Hash)
		}
	}
}
