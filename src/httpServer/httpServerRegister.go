package httpServer

import (
	"log"

	"swap.io-agent/src/blockchain/synchronizer"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var synchronizer *synchronizer.Synchronizer
	err := reg.FetchService(&synchronizer)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeServer(HttpServerConfig{
			Synhronizer: synchronizer,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
