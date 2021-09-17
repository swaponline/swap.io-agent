package subscribeManager

import (
	"swap.io-agent/src/redisStore"
)

type SubscribeManager struct {
	subscribersStore redisStore.ISubscribersStore
}

type SubscribeManagerConfig struct {
	subscribersStore redisStore.ISubscribersStore
}

func InitializeSubscribeManager(config SubscribeManagerConfig) *SubscribeManager {
	return &SubscribeManager{
		subscribersStore: config.subscribersStore,
	}
}

func (sm *SubscribeManager) SubscribeUserToAddress(
	userId string,
	address string,
) error {
	return sm.subscribersStore.AddSubscription(userId, address)
}
func (sm *SubscribeManager) ClearAllUserSubscriptions(userId string) error {
	return sm.subscribersStore.RemoveSubscriptions(userId)
}

func (sm *SubscribeManager) Start() {}
func (sm *SubscribeManager) Stop() error {
	return nil
}
func (sm *SubscribeManager) Status() error {
	return nil
}
