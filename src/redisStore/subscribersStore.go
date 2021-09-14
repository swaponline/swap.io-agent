package redisStore

import (
	"context"

	"swap.io-agent/src/common/Set"
)

var ctx = context.Background()

func (db *RedisDb) GetSubscribersFromAddresses(addresses []string) []string {
	subscribers := Set.New()
	for _, address := range addresses {
		subscribers.Adds(
			db.GetAddressSubscribers(address),
		)
	}
	return subscribers.Keys()
}
func (db *RedisDb) GetAddressSubscribers(address string) []string {
	return db.client.SMembers(ctx, address).Val()
}
func (db *RedisDb) SubscribeUserToAddress(userId string, address string) error {
	isAddUserAddress := db.client.SAdd(ctx, userId, address).Err()
	if isAddUserAddress != nil {
		return isAddUserAddress
	}
	isAddAddressUser := db.client.SAdd(ctx, address, userId).Err()
	if isAddAddressUser != nil {
		//redis transaction not supporting many work in goroutines
		db.client.SRem(ctx, userId, address)
		return isAddAddressUser
	}

	return nil
}
func (db *RedisDb) ClearAllUserSubscriptions(userId string) error {
	subscriptions := db.client.SMembers(ctx, userId).Val()
	for _, address := range subscriptions {
		err := db.client.SRem(ctx, address, userId).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
