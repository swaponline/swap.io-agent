package redisStore

import "context"

type IAddressInfoService interface {
	AddUser(id int) error
	RemoveUser(id int) error
}

var ctx = context.Background()

func (db *RedisDb) AddressIsActive(address string) (bool, error) {
	isExist, err := db.client.Exists(ctx, address).Result()
	return isExist > 0, err
}
func (db *RedisDb) AddAddressUser(address string, userId int) error {
	return db.client.SAdd(ctx, address, userId).Err()
}
func (db *RedisDb) RemoveAddressUser(address string, userId int) error {
	return db.client.SRem(ctx, address, userId).Err()
}

func (db *RedisDb) UserIsActive(id int) (bool,error) {
	return db.client.SIsMember(ctx, "users", id).Result()
}
func (db *RedisDb) AddUser(id int) error {
	return db.client.SAdd(ctx, "users", id).Err()
}
func (db *RedisDb) RemoveUser(id int) error {
	return db.client.SRem(ctx, "users", id).Err()
}