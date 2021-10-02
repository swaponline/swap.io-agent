package indexer

import (
	"log"
	"strconv"
	"time"

	"swap.io-agent/src/blockchain"
)

func (i *Indexer) RunScanner() {
	allIndexedTransactions := 0
	currentBlock := i.transactionsStore.GetLastTransactionBlock() + 1

	for {
		block, err := i.api.GetBlockByIndex(
			currentBlock,
		)
		if err == blockchain.ApiNotExist {
			break
		}
		if err != blockchain.ApiRequestSuccess {
			log.Println("get block request err", err)
			continue
		}
		log.Println("block", currentBlock, "- ok")

		indexedTransactions := indexingTransactions(block.Transactions)

		allIndexedTransactions += len(block.Transactions)

		log.Println("indexed transactions -", len(block.Transactions))
		log.Println("all indexed transactions -", allIndexedTransactions)

		for err := i.writeExpectedTxsToQueueEvents(
			block.Transactions,
		); err != nil; {
			log.Println("ERROR", err)
		}

		for err := i.transactionsStore.WriteLastIndexedTransactions(
			indexedTransactions,
			currentBlock,
		); err != nil; {
			log.Println("ERROR", err, indexedTransactions)
		}

		currentBlock += 1
	}

	close(i.isSynchronize)
	log.Println("***blockchain synchronized***")

	// todo: refactoring
	for {
		block, err := i.api.GetBlockByIndex(
			currentBlock,
		)
		if err == blockchain.ApiNotExist {
			for err := i.writeExpectedTxsToQueueEvents(
				[]*blockchain.Transaction{{
					Hash:              "time" + strconv.Itoa(int(time.Now().Unix())),
					AllSpendAddresses: []string{"mi46vEy3EPcDx1PLMw7hgAhHqCWSBPnuMA"},
				}},
			); err != nil; {
				log.Println("ERROR", err)
			}
			<-time.After(time.Millisecond * 2000)
			continue
		}
		if err != blockchain.ApiRequestSuccess {
			log.Println("get block request err", err)
			continue
		}
		log.Println("block", currentBlock, "- ok")

		indexedTransactions := indexingTransactions(block.Transactions)

		allIndexedTransactions += len(block.Transactions)

		log.Println("indexed transactions -", len(block.Transactions))
		log.Println("all indexed transactions -", allIndexedTransactions)

		for err := i.writeExpectedTxsToQueueEvents(
			block.Transactions,
		); err != nil; {
			log.Println("ERROR", err)
		}

		for err := i.transactionsStore.WriteLastIndexedTransactions(
			indexedTransactions,
			currentBlock,
		); err != nil; {
			log.Println("ERROR", err)
		}

		currentBlock += 1
	}
}

func indexingTransactions(
	transactions []*blockchain.Transaction,
) map[string][]string {
	// address -> hashTx
	buf := make(map[string][]string)
	for _, transaction := range transactions {
		for _, address := range transaction.AllSpendAddresses {
			buf[address] = append(
				buf[address],
				transaction.Hash,
			)
		}
	}
	return buf
}

func (i *Indexer) getExpectedTxs(
	txs []*blockchain.Transaction,
) map[string][]*blockchain.Transaction {
	buf := make(map[string][]*blockchain.Transaction)
	for _, tx := range txs {
		subscribers := i.subscribesManager.GetSubscribersFromAddresses(
			tx.AllSpendAddresses,
		)
		for _, subscriber := range subscribers {
			buf[subscriber] = append(buf[subscriber], tx)
		}
	}
	return buf
}
func (i *Indexer) writeExpectedTxsToQueueEvents(
	txs []*blockchain.Transaction,
) error {
	expectedTxs := i.getExpectedTxs(txs)
	return i.queueEvents.WriteTxsEvents(expectedTxs)
}
