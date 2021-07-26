package ethereum

import (
	"log"
	"swap.io-agent/src/blockchain/ethereum/transactionFormatter"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/serviceRegistry"
)

func BlockchainIndexerRegister(reg *serviceRegistry.ServiceRegistry) {
	var formatter *transactionFormatter.TransactionFormatter
	err := reg.FetchService(&formatter)
	if err != nil {
		log.Panicln(err)
	}

	var transactionStore *levelDbStore.TransactionsStore
	err = reg.FetchService(&transactionStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeIndexer(BlockchainIndexerConfig{
			TransactionsStore: transactionStore,
			Formatter: formatter,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}