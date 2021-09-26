package main

import (
	"log"

	"swap.io-agent/src/blockchain/ethereum/nodeApi/geth"
	"swap.io-agent/src/blockchain/ethereum/transactionFormatter"
	"swap.io-agent/src/blockchain/handshake"
	"swap.io-agent/src/blockchain/indexer"
	"swap.io-agent/src/blockchain/networks"
	"swap.io-agent/src/blockchain/synchronizer"
	"swap.io-agent/src/blockchain/transactionNotifierPipe"
	"swap.io-agent/src/env"
	"swap.io-agent/src/httpHandler"
	"swap.io-agent/src/httpServer"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/socketServer"
	"swap.io-agent/src/subscribersManager"
)

func main() {
	registry := serviceRegistry.NewServiceRegistry()

	err := env.InitializeConfig()
	if err != nil {
		log.Panicln(err)
	}

	handshake.Test()

	networks := networks.InitializeNetworks()
	err = registry.RegisterService(networks)
	if err != nil {
		log.Panicln(err)
	}

	subscribersStoreMemory, err := redisStore.InitializeDB()
	if err != nil {
		log.Panicf("redisStore not initialize, err: %v", err)
	}
	err = registry.RegisterService(&subscribersStoreMemory)
	if err != nil {
		log.Panicln(err)
	}

	subscribersStoreDisk, err := levelDbStore.InitialiseSubscribersStore()
	if err != nil {
		log.Panicf("Subscribers store not initialize, err: %v", err)
	}
	err = registry.RegisterService(subscribersStoreDisk)
	if err != nil {
		log.Panicln(err)
	}

	transactionStore, err := levelDbStore.InitialiseTransactionStore(
		levelDbStore.TransactionsStoreConfig{
			Name:                 env.BLOCKCHAIN,
			DefaultScannedBlocks: env.BLOCKCHAIN_DEFAULT_SCANNED_BLOCK,
		},
	)
	if err != nil {
		log.Panicln(err)
	}
	err = registry.RegisterService(transactionStore)
	if err != nil {
		log.Panicln(err)
	}

	api := geth.InitializeGeth()
	err = registry.RegisterService(
		api,
	)
	if err != nil {
		log.Panicln(err)
	}
	formatter := transactionFormatter.InitializeTransactionFormatter(transactionFormatter.TransactionFormatterConfig{
		Api: api,
	})
	err = registry.RegisterService(
		formatter,
	)
	if err != nil {
		log.Panicln(err)
	}

	synchronizer.Register(registry)

	indexer.IndexerRegister(registry)

	subscribersManager.Register(registry)

	transactionNotifierPipe.Register(registry)

	socketServer.Register(registry)

	httpServer.Register(registry)

	httpHandlerEntity := httpHandler.InitializeServer()
	err = registry.RegisterService(httpHandlerEntity)
	if err != nil {
		log.Panicln(err)
	}

	registry.StartAll()

	<-make(chan struct{})
}
