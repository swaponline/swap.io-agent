package blockchain

import "swap.io-agent/src/redisStore"

type SubscribeManager struct {
	synchroniser struct{}
	subscribersStore redisStore.ISubscribersStore
	formatter struct{}
}

type SubscribeManagerConfig struct {
	synchroniser struct{}
	subscribersStore redisStore.ISubscribersStore
	formatter struct{}
}

func InitializeSubscribeManager(config SubscribeManagerConfig) SubscribeManager {
	return SubscribeManager{
		synchroniser: config.synchroniser,
		subscribersStore: config.subscribersStore,
		formatter: config.formatter,
	}
}

func (sm *SubscribeManager) SubscribeUserToAddress(
	userId string,
	address string,
) error {
	return sm.subscribersStore.SubscribeUserToAddress(userId, address)
}
func (sm *SubscribeManager) ClearAllUserSubscriptions(userId string) error {
	return sm.subscribersStore.ClearAllUserSubscriptions(userId)
}

func (sm *SubscribeManager) Start() {}
func (sm *SubscribeManager) Stop() error {
	return nil
}
func (sm *SubscribeManager) Status() error {
	return nil
}