package synchronizer

import (
	"log"

	"swap.io-agent/src/blockchain/ethereum/transactionFormatter"
	"swap.io-agent/src/blockchain/networks"
	"swap.io-agent/src/env"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
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
	var formatter *transactionFormatter.TransactionFormatter
	err = reg.FetchService(&formatter)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitialiseSynchronizer(SynchronizerConfig{
			Api:       networkApi,
			Formatter: formatter,
			Store:     transactionStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
