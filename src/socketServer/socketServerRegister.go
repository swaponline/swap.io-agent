package socketServer

import (
	"log"

	"swap.io-agent/src/blockchain/synchronizer"
	"swap.io-agent/src/blockchain/transactionNotifierPipe"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var notifyTransactionPipe *transactionNotifierPipe.TransactionNotifierPipe
	err := reg.FetchService(&notifyTransactionPipe)
	if err != nil {
		log.Panicln(err)
	}

	var synchroniser *synchronizer.Synchronizer
	err = reg.FetchService(&synchroniser)
	if err != nil {
		log.Panicln(err)
	}

	var subscribeManager *redisStore.RedisDb
	err = reg.FetchService(&subscribeManager)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeServer(Config{
			synchronizer:     synchroniser,
			subscribeManager: subscribeManager,
			onNotifyUsers:    notifyTransactionPipe.Out,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
