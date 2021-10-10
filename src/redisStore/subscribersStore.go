package redisStore

import (
	"context"
	"swap.io-agent/src/common/Set"
	"swap.io-agent/src/config"
)

var ctx = context.Background()

func keyWidthPrefix(key string) string {
	return config.BLOCKCHAIN + "|" + key
}

const activeUsersKey = "activeUsers"

func (db *RedisDb) UserIsActive(userId string) (bool, error) {
	return db.client.SIsMember(ctx, keyWidthPrefix(activeUsersKey), userId).Result()
}
func (db *RedisDb) ActiveUser(userId string) error {
	return db.client.SAdd(ctx, keyWidthPrefix(activeUsersKey), userId).Err()
}
func (db *RedisDb) DeactiveUser(userId string) error {
	return db.client.SRem(ctx, keyWidthPrefix(activeUsersKey), userId).Err()
}

// todo: return data, err
func (db *RedisDb) GetSubscribersFromAddresses(addresses []string) []string {
	subscribers := Set.New()
	for _, address := range addresses {
		subscribers.Adds(
			db.GetAddressSubscribers(address),
		)
	}
	return subscribers.Keys()
}

// todo: return data, err
func (db *RedisDb) GetAddressSubscribers(address string) []string {
	return db.client.SMembers(ctx, keyWidthPrefix(address)).Val()
}

func (db *RedisDb) AddSubscription(userId string, address string) error {
	isAddUserAddress := db.client.SAdd(
		ctx,
		keyWidthPrefix(userId),
		address,
	).Err()
	if isAddUserAddress != nil {
		return isAddUserAddress
	}
	isAddAddressUser := db.client.SAdd(
		ctx,
		keyWidthPrefix(address),
		userId,
	).Err()
	if isAddAddressUser != nil {
		//redis transaction not supporting many work in goroutines
		db.client.SRem(
			ctx,
			keyWidthPrefix(userId),
			address,
		)
		return isAddAddressUser
	}

	return nil
}
func (db *RedisDb) RemoveSubscription(userId string, address string) error {
	err := db.client.SRem(ctx, keyWidthPrefix(address), userId).Err()
	if err != nil {
		return err
	}
	err = db.client.SRem(ctx, keyWidthPrefix(userId), address).Err()
	if err != nil {
		return err
	}

	return nil
}
func (db *RedisDb) RemoveSubscriptions(userId string) error {
	subscriptions := db.client.SMembers(ctx, keyWidthPrefix(userId)).Val()
	for _, address := range subscriptions {
		db.RemoveSubscription(userId, address)
	}
	return nil
}
