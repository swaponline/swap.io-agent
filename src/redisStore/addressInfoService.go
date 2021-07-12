package redisStore

import (
	"context"
	"swap.io-agent/src/common/Set"
)

var ctx = context.Background()

func (db *RedisDb) GetSubscribersFromAddresses(addresses []string) []string {
	uniqueUsers := Set.New()
	for _, address := range addresses {
		uniqueUsers.Adds(
			db.GetAddressSubscriber(address),
		)
	}
	return uniqueUsers.Keys()
}
func (db *RedisDb) GetAddressSubscriber(address string) []string {
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