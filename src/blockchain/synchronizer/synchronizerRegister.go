package synchronizer

import (
	"log"
	"swap.io-agent/src/blockchain/ethereum/transactionFormatter"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var transactionStore *levelDbStore.TransactionsStore
	err := reg.FetchService(&transactionStore)
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
		 	Formatter: formatter,
		 	Store: transactionStore,
		 }),
	)
	if err != nil {
		log.Panicln(err)
	}
}