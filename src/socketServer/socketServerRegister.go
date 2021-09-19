package socketServer

import (
	"log"

	"swap.io-agent/src/blockchain/transactionNotifierPipe"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/subscribersManager"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var subscribeManager *subscribersManager.SubscribesManager
	err := reg.FetchService(&subscribeManager)
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
			subscribeManager: subscribeManager,
			onNotifyUsers:    notifyTransactionPipe.Out,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
