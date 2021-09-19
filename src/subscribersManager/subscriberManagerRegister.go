package subscribersManager

import (
	"log"

	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/redisStore"
	"swap.io-agent/src/serviceRegistry"
)

func Register(reg *serviceRegistry.ServiceRegistry) {
	var redisStore *redisStore.RedisDb
	err := reg.FetchService(&redisStore)
	if err != nil {
		log.Panicln(err)
	}

	var diskStore *levelDbStore.SubscribersStore
	err = reg.FetchService(&diskStore)
	if err != nil {
		log.Panicln(err)
	}

	err = reg.RegisterService(
		InitialiseSubscribersStore(SubscribesManagerConfig{
			MemoryStore: redisStore,
			DiskStore:   diskStore,
		}),
	)
	if err != nil {
		log.Panicln(err)
	}
}
