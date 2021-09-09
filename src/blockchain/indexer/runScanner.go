package indexer

import (
	"log"
	"time"

	"swap.io-agent/src/blockchain"
)

func (indexer *Indexer) RunScanner() {
	allIndexedTransactions := 0
	currentBlock := indexer.transactionsStore.GetLastTransactionBlock() + 1

	for {
		block, err := indexer.api.GetBlockByIndex(
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

		for err := indexer.transactionsStore.WriteLastIndexedTransactions(
			indexedTransactions,
			currentBlock,
		); err != nil; {
			log.Println("ERROR", err, indexedTransactions)
		}

		currentBlock += 1
	}

	close(indexer.isSynchronize)
	log.Println("***blockchain synchronized***")

	for {
		block, err := indexer.api.GetBlockByIndex(
			currentBlock,
		)
		if err == blockchain.ApiNotExist {
			<-time.After(time.Millisecond * 500)
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

		for err := indexer.transactionsStore.WriteLastIndexedTransactions(
			indexedTransactions,
			currentBlock,
		); err != nil; {
			log.Println("ERROR", err)
		}

		for _, transaction := range block.Transactions {
			indexer.NewTransactions <- transaction
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
