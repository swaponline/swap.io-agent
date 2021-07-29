package main

import (
	"log"
	"swap.io-agent/src/blockchain/ethereum/api/ethercsan"
	ethereum "swap.io-agent/src/blockchain/ethereum/blockchainIndexer"
	"swap.io-agent/src/blockchain/ethereum/transactionFormatter"
	"swap.io-agent/src/blockchain/subscribeManager"
	"swap.io-agent/src/blockchain/synchronizer"
	"swap.io-agent/src/blockchain/transactionNotifierPipe"
	"swap.io-agent/src/env"
	"swap.io-agent/src/httpHandler"
	"swap.io-agent/src/httpServer"
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
	"swap.io-agent/src/socketServer"
)

func main() {
	registry := serviceRegistry.NewServiceRegistry()

	err := env.InitializeConfig()
	if err != nil {
		log.Panicln(err.Error())
	}

	db, err := redisStore.InitializeDB()
	if err != nil {
		log.Panicf("redisStore not initialize, err: %v", err)
	}
	err = registry.RegisterService(&db)
	if err != nil {
		log.Panicln(err.Error())
	}

	transactionStore, err := levelDbStore.InitialiseTransactionStore(
		levelDbStore.TransactionsStoreConfig{
			Name: env.BLOCKCHAIN,
			DefaultScannedBlocks: env.BLOCKCHAIN_DEFAULT_SCANNED_BLOCK,
		},
	)
	if err != nil {
		log.Panicln(err.Error())
	}
	err = registry.RegisterService(transactionStore)
	if err != nil {
		log.Panicln(err.Error())
	}

	api := ethercsan.InitializeEthercsan()
	err = registry.RegisterService(
		api,
	)
	if err != nil {
		log.Panicln(err.Error())
	}

	transactionFormatter.Register(registry)

	ethereum.BlockchainIndexerRegister(registry)

	synchronizer.Register(registry)

	transactionNotifierPipe.Register(registry)

	subscribeManager.Register(registry)

	socketServer.Register(registry)

	httpServerEntity := httpServer.InitializeServer()
	err = registry.RegisterService(httpServerEntity)
	if err != nil {
		log.Panicln(err.Error())
	}

	httpHandlerEntity := httpHandler.InitializeServer()
	err = registry.RegisterService(httpHandlerEntity)
	if err != nil {
		log.Panicln(err.Error())
	}

	registry.StartAll()

	<-make(chan struct{})
}