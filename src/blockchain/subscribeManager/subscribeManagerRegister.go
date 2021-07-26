package subscribeManager

import (
	"log"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var subscribersStore *redisStore.RedisDb
	err := reg.FetchService(&subscribersStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitializeSubscribeManager(SubscribeManagerConfig{
			subscribersStore: subscribersStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}