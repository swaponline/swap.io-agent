package ethereum

import (
	"log"
	"time"

	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/blockchain/ethereum/api"
)

func (indexer *BlockchainIndexer) RunScanner() {
	allIndexedTransactions := 0
	currentBlock := indexer.transactionsStore.GetLastTransactionBlock() + 1

	for {
		block, err := indexer.api.GetBlockByIndex(
			currentBlock,
		)
		if err == api.NotExistBlockError {
			break
		}
		if err != api.RequestSuccess {
			log.Println("get block request err", err)
			continue
		}
		log.Println(currentBlock, "|", block.Number, "- ok")

		blockTrace := make(chan struct{})
		go func() {
			_, err := indexer.api.GetBlockTraceByIndex(block.Number)
			if err != api.RequestSuccess {
				log.Panicln("ERROR", block.Number)
			}
			close(blockTrace)
		}()
		//transactions        := indexer.formatBlockTransactions(block)
		<-blockTrace

		indexedTransactions := make(map[string][]string)
		//indexingTransactions(indexedTransactions, transactions)
		allIndexedTransactions += len(block.Transactions)
		log.Println("indexed transactions -", len(block.Transactions))
		log.Println("all indexed transactions -", allIndexedTransactions)

		indexer.writeIndexedTransactionToStore(
			indexedTransactions, currentBlock,
		)

		currentBlock += 1
	}

	close(indexer.isSynchronize)
	log.Println("***blockchain synchronized***")

	for {
		block, err := indexer.api.GetBlockByIndex(
			currentBlock,
		)
		if err == api.NotExistBlockError {
			<-time.After(time.Millisecond * 500)
			continue
		}
		if err != api.RequestSuccess {
			log.Println("get block request err", err)
			continue
		}
		log.Println(currentBlock, "|", block.Number, "- ok")

		transactions := indexer.formatBlockTransactions(block)

		indexedTransactions := make(map[string][]string)
		indexingTransactions(indexedTransactions, transactions)
		allIndexedTransactions += len(transactions)
		log.Println("indexed transactions -", len(transactions))
		log.Println("all indexed transactions -", allIndexedTransactions)

		indexer.writeIndexedTransactionToStore(
			indexedTransactions, currentBlock,
		)
		for indexer.transactionsStore.Flush() != nil {
		}

		for _, transaction := range transactions {
			indexer.NewTransactions <- transaction
		}

		currentBlock += 1
	}
}

func (indexer *BlockchainIndexer) writeIndexedTransactionToStore(
	indexedTransactions map[string][]string,
	indexBlock int,
) {
	for {
		err := indexer.transactionsStore.WriteLastIndexedTransactions(
			indexedTransactions, indexBlock,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		break
	}
}
func (indexer *BlockchainIndexer) formatBlockTransactions(
	block *api.Block,
) []*blockchain.Transaction {
	transactions := make([]*blockchain.Transaction, 0)
	for _, blockTx := range block.Transactions {
		for {
			transaction, err := indexer.formatter.FormatTransaction(&blockTx, block)
			if err != nil {
				log.Println(err)
				continue
			}
			transactions = append(
				transactions,
				transaction,
			)
			break
		}
	}

	return transactions
}

func indexingTransactions(
	indexedTransactions map[string][]string,
	transactions []*blockchain.Transaction,
) {
	// address -> transactions
	for _, transaction := range transactions {
		for _, address := range transaction.AllSpendAddresses {
			if len(address) == 0 {
				log.Println(transaction.AllSpendAddresses)
				log.Println(
					transaction.Hash,
					transaction.From,
					transaction.To,
				)
				continue
			}
			indexedTransactions[address] = append(
				indexedTransactions[address],
				transaction.Hash,
			)
		}
	}
}
