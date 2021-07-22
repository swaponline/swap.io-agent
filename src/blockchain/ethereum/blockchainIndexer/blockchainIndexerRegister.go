package ethereum

import (
	"log"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/serviceRegistry"
)

func BlockchainIndexerRegister(reg *serviceRegistry.ServiceRegistry) {
	var transactionStore *levelDbStore.TransactionsStore
	err := reg.FetchService(&transactionStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeIndexer(BlockchainIndexerConfig{
			TransactionsStore: transactionStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}