package ethereum

import (
	"log"
	"swap.io-agent/src/blockchain/ethereum/api/ethercsan"
	"swap.io-agent/src/blockchain/ethereum/transactionFormatter"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/serviceRegistry"
)

func BlockchainIndexerRegister(reg *serviceRegistry.ServiceRegistry) {
	var api *ethercsan.Etherscan
	err := reg.FetchService(&api)
	if err != nil {
		log.Panicln(err)
	}

	var formatter *transactionFormatter.TransactionFormatter
	err = reg.FetchService(&formatter)
	if err != nil {
		log.Panicln(err)
	}

	var transactionStore *levelDbStore.TransactionsStore
	err = reg.FetchService(&transactionStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeIndexer(BlockchainIndexerConfig{
			Api: api,
			TransactionsStore: transactionStore,
			Formatter: formatter,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}