package socketServer

import (
	"log"
	"swap.io-agent/src/redisStore"
	. "swap.io-agent/src/serviceRegistry"
)

func Register(reg *ServiceRegistry) {
	var dbService *redisStore.RedisDb
	err := reg.FetchService(&dbService)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeServer(Config{
			dbService,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}