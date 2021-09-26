package subscribersManager

import (
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/redisStore"
)

type SubscribesManager struct {
	memoryStore *redisStore.RedisDb
	diskStore   *levelDbStore.SubscribersStore
}

type SubscribesManagerConfig struct {
	MemoryStore *redisStore.RedisDb
	DiskStore   *levelDbStore.SubscribersStore
}

func InitialiseSubscribersStore(config SubscribesManagerConfig) *SubscribesManager {
	return &SubscribesManager{
		memoryStore: config.MemoryStore,
		diskStore:   config.DiskStore,
	}
}

func (s *SubscribesManager) LoadAllSubscriptions() error {
	allSubscriptions, err := s.diskStore.GetAllSubscriptions()
	if err != nil {
		return err
	}
	for userId, subscriptions := range allSubscriptions {
		for _, subscription := range subscriptions {
			err := s.memoryStore.AddSubscription(userId, subscription)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SubscribesManager) GetSubscribersFromAddresses(
	addresses []string,
) []string {
	return s.memoryStore.GetSubscribersFromAddresses(addresses)
}

func (s *SubscribesManager) SubscribeUserToAddress(
	userId string, address string, isWriteToDisk bool,
) error {
	err := s.diskStore.AddSubscription(userId, address)
	if err != nil {
		return err
	}

	return s.memoryStore.AddSubscription(userId, address)
}
func (s *SubscribesManager) UnsubscribeUserToAddress(
	userId string, address string,
) error {
	err := s.diskStore.RemoveSubscription(userId, address)
	if err != nil {
		return err
	}

	return s.memoryStore.RemoveSubscription(userId, address)
}

func (*SubscribesManager) Start() {}
func (*SubscribesManager) Stop() error {
	return nil
}
func (*SubscribesManager) Status() error {
	return nil
}
