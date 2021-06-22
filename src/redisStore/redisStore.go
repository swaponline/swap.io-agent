package redisStore

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

type RedisDb struct {
	client *redis.Client
}

func InitializeDB() (RedisDb, error) {
	db := RedisDb{}

	redisDbIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return db, fmt.Errorf(
			"error(%v) parse REDIS_DB from env, expected int", err,
		)
	}

	db.client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: redisDbIndex,
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