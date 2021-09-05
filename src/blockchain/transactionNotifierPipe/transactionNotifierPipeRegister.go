package transactionNotifierPipe

import (
	"log"

	"swap.io-agent/src/blockchain/indexer"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var indexer *indexer.Indexer
	err := reg.FetchService(&indexer)
	if err != nil {
		log.Panicln(err)
	}

	var subscribersStore *redisStore.RedisDb
	err = reg.FetchService(&subscribersStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeTransactionNotifierPipe(TransactionNotifierPipeConfig{
			Input:            indexer.NewTransactions,
			SubscribersStore: subscribersStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
