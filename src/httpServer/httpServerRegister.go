package httpServer

import (
	"log"
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

	err = reg.RegisterService(
		InitializeServer(HttpServerConfig{
			Synhronizer: s,
			SubscribersManager: sm,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
