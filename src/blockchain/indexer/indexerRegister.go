package indexer

import (
	"log"

	"swap.io-agent/src/blockchain/networks"
	"swap.io-agent/src/env"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/queueEvents"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/subscribersManager"
)

func IndexerRegister(reg *serviceRegistry.ServiceRegistry) {
	var queueEvents *queueEvents.QueueEvents
	err := reg.FetchService(&queueEvents)
	if err != nil {
		log.Panicln(err)
	}

	var subscribeManager *subscribersManager.SubscribesManager
	err = reg.FetchService(&subscribeManager)
	if err != nil {
		log.Panicln(err)
	}

	var networks *networks.Networks
	err = reg.FetchService(&networks)
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
			QueueEvents:       queueEvents,
			SubscribesManager: subscribeManager,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
