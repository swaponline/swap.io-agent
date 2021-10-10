package socketServer

import (
	"log"

	"swap.io-agent/src/blockchain/transactionNotifierPipe"
	"swap.io-agent/src/queueEvents"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/subscribersManager"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var queueEvents *queueEvents.QueueEvents
	err := reg.FetchService(&queueEvents)
	if err != nil {
		log.Panicln(err)
	}

	var usersManager *redisStore.RedisDb
	err = reg.FetchService(&usersManager)
	if err != nil {
		log.Panicln(err)
	}

	var subscribersManager *subscribersManager.SubscribesManager
	err = reg.FetchService(&subscribersManager)
	if err != nil {
		log.Panicln(err)
	}

	var notifyTransactionPipe *transactionNotifierPipe.TransactionNotifierPipe
	err = reg.FetchService(&notifyTransactionPipe)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeServer(Config{
			usersManager:     usersManager,
			subscribersManager: subscribersManager,
			queueEvents:      queueEvents,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
