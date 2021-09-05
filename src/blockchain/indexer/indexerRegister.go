package indexer

import (
	"log"

	"swap.io-agent/src/blockchain/networks"
	"swap.io-agent/src/env"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/serviceRegistry"
)

func IndexerRegister(reg *serviceRegistry.ServiceRegistry) {
	var networks *networks.Networks
	err := reg.FetchService(&networks)
	if err != nil {
		log.Panicln(err)
	}
	networkApi := (*networks)[env.BLOCKCHAIN]

	var transactionStore *levelDbStore.TransactionsStore
	err = reg.FetchService(&transactionStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeIndexer(IndexerConfig{
			Api:               networkApi,
			TransactionsStore: transactionStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
