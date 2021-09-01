package indexer

import (
	"log"

	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/serviceRegistry"
)

func IndexerRegister(reg *serviceRegistry.ServiceRegistry) {
	var transactionStore *levelDbStore.TransactionsStore
	err := reg.FetchService(&transactionStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeIndexer(IndexerConfig{
			//Api:               api,
			TransactionsStore: transactionStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
