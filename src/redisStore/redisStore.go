package redisStore

import (
	"github.com/go-redis/redis/v8"
	"swap.io-agent/src/env"
)

type RedisDb struct {
	client *redis.Client
}

func InitializeDB() (RedisDb, error) {
	db := RedisDb{}

	db.client = redis.NewClient(&redis.Options{
		Addr: env.REDIS_ADDR,
		Password: env.REDIS_PASSWORD,
		DB: env.REDIS_DB,
	})

	return db,nil
}

func (db *RedisDb) Start() {}
func (db *RedisDb) Stop() error {
	return db.client.Close()
}
func (db *RedisDb) Status() error {
	return nil
}