package redisStore

import "context"

type IAddressInfoService interface {
	AddUser(id int) error
	RemoveUser(id int) error
}

var ctx = context.Background()

func (db *RedisDb) AddAddressUser(address string, userId int) error {
	return db.client.SAdd(ctx, address, userId).Err()
}
func (db *RedisDb) RemoveAddressUser(address string, userId int) error {
	return db.client.SRem(ctx, address, userId).Err()
}

func (db *RedisDb) AddUser(id int) error {
	return db.client.SAdd(ctx, "users", id).Err()
}
func (db *RedisDb) RemoveUser(id int) error {
	return db.client.SRem(ctx, "users", id).Err()
}