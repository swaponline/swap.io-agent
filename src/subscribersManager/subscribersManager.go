package main

import (
	"swap.io-agent/src/levelDbStore"
	"swap.io-agent/src/redisStore"
)

type SubscribesManager struct {
	memoryStore redisStore.ISubscribersStore
	diskStore   levelDbStore.SubscribersStore
}

type SubscribesManagerConfig struct {
	MemoryStore redisStore.ISubscribersStore
	DiskStore   levelDbStore.SubscribersStore
}

func InitialiseSubscribersStore(config SubscribesManagerConfig) *SubscribesManager {
	return &SubscribesManager{
		memoryStore: config.MemoryStore,
		diskStore:   config.DiskStore,
	}
}

func (s *SubscribesManager) GetCountSubcribers(
	userId string,
) (int, error) {
	return s.diskStore.GetCountSubcribers(userId)
}

func (s *SubscribesManager) GetSubscribersFromAddresses(
	addresses []string,
) []string {
	return s.memoryStore.GetSubscribersFromAddresses(addresses)
}

func (s *SubscribesManager) SubscribeUserToAddress(
	userId string, address string,
) error {
	err := s.diskStore.AddSubscription(userId, address)
	if err != nil {
		return err
	}
	return s.memoryStore.AddSubscription(userId, address)
}
func (s *SubscribesManager) SubscribeUserToAddresses(
	userId string, addresses []string,
) int {
	writed := 0
	for _, address := range addresses {
		err := s.SubscribeUserToAddress(userId, address)
		if err != nil {
			return len(addresses) - writed
		}
		writed += 1
	}

	return 0
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
func (s *SubscribesManager) UnsubscribeUserToAddresses(
	userId string, addresses []string,
) int {
	deleted := 0
	for _, address := range addresses {
		err := s.UnsubscribeUserToAddress(userId, address)
		if err != nil {
			return len(addresses) - deleted
		}
		deleted += 1
	}

	return 0
}
func (s *SubscribesManager) ClearAllUserSubscriptions(
	userId string,
) error {
	//!!! no clear disk store
	return s.memoryStore.RemoveSubscriptions(userId)
}
