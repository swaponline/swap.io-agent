package transactionNotifierPipe

import (
	"log"
	ethereum "swap.io-agent/src/blockchain/ethereum/blockchainIndexer"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var indexer *ethereum.BlockchainIndexer
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
			Input: indexer.NewTransactions,
			SubscribersStore: subscribersStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
