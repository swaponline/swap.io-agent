package redisStore

import (
	"github.com/go-redis/redis/v8"
	"swap.io-agent/src/config"
)

type RedisDb struct {
	client *redis.Client
}

func InitializeDB() (RedisDb, error) {
	db := RedisDb{}

	db.client = redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})

	return db, nil
}

func (db *RedisDb) Start() {}
func (db *RedisDb) Stop() error {
	return db.client.Close()
}
func (db *RedisDb) Status() error {
	return nil
}
