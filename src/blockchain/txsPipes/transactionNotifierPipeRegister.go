package txsPipes

import (
	"log"

	"swap.io-agent/src/blockchain/indexer"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/subscribersManager"
)

func NewTxsNotifierPipeRegister(reg *serviceRegistry.ServiceRegistry) {
	var indexer *indexer.Indexer
	err := reg.FetchService(&indexer)
	if err != nil {
		log.Panicln(err)
	}

	var subscribersManager *subscribersManager.SubscribesManager
	err = reg.FetchService(&subscribersManager)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeNewTxsNotifierPipe(NewTxsNotifierPipeConfig{
			Input:              indexer.NewTransactions,
			SubscribersManager: subscribersManager,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
