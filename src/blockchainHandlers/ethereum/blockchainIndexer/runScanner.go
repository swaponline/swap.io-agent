package ethereum

import (
	"log"
	"strconv"
	"swap.io-agent/src/blockchainHandlers/ethereum/api/ethercsan"
	"swap.io-agent/src/levelDbStore"
	"sync"
)


func (indexer *BlockchainIndexer) RunScanner() {
	isSynchronize   := false
	requestsStepLen := 5
	for {
		waits := new(sync.WaitGroup)
		waits.Add(requestsStepLen)

		lastGetBlockIndex := indexer.lastBlock
		buf := make(map[string][]string)
		lockerChange := new(sync.Mutex)
		for t:=1; t<=requestsStepLen; t++ {
			go func(upIndexBlock int) {
				blockIndex := indexer.lastBlock + upIndexBlock
				block, err := ethercsan.GetBlockByIndex(
					indexer.apiKey,
					blockIndex,
				)
				if err == ethercsan.NotExistBlockError {
					isSynchronize = true
					close(indexer.isSynchronize)
				} else if err != ethercsan.RequestSuccess {
					log.Panicf("not get block by index request etherscan")
				}

				for _, transaction := range block.Transactions {
					log.Println(transaction)
					lockerChange.Lock()
					buf[transaction.From] = append(buf[transaction.From], transaction.Hash)
					buf[transaction.To]   = append(buf[transaction.To], transaction.To)
					if lastGetBlockIndex < blockIndex {
						lastGetBlockIndex = blockIndex
					}
					lockerChange.Unlock()
				}
				waits.Done()
			}(t)
		}
		waits.Wait()

		transaction, err := indexer.db.OpenTransaction()
		for address, blockTransactions := range buf {
			err = levelDbStore.ArrayStringPush(
				transaction, address, blockTransactions,
			)
			if err != nil {
				transaction.Discard()
				log.Panicln(err)
			}
		}

		err = indexer.db.Put(
			lastBlockKey,
			[]byte(strconv.Itoa(lastGetBlockIndex)),
			nil,
		)
		if err != nil ||
		   transaction.Commit() != nil {
			transaction.Discard()
			log.Panicln(err)
		}
	}
}