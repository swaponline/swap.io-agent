package httpServer

import (
	"log"
	"swap.io-agent/src/blockchain/networks"
	"swap.io-agent/src/config"
	"swap.io-agent/src/subscribersManager"

	"swap.io-agent/src/blockchain/synchronizer"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var s *synchronizer.Synchronizer
	err := reg.FetchService(&s)
	if err != nil {
		log.Panicln(err)
	}

	var sm *subscribersManager.SubscribesManager
	err = reg.FetchService(&sm)
	if err != nil {
		log.Panicln(err)
	}

	var n *networks.Networks
	err = reg.FetchService(&n)
	if err != nil {
		log.Panicln(err)
	}
	blockchainApi, ok := (*n)[config.BLOCKCHAIN]
	if !ok {
		log.Panicln(ok, "not found blockchain api")
	}

	err = reg.RegisterService(
		InitializeServer(HttpServerConfig{
			Synhronizer:        s,
			SubscribersManager: sm,
			BlockchainApi:      blockchainApi,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
