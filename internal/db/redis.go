package db

import (
	"fmt"
	"time"

	"github.com/CharlesChou03/_git/block-chain.git/config"
	"github.com/go-redis/redis"
)

type RedisBlockChainDB struct {
	DB *redis.Client
}

var (
	RedisDB *RedisBlockChainDB
)

func SetupRedisDB() *RedisBlockChainDB {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       0, // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return &RedisBlockChainDB{DB: client}
}

func (r *RedisBlockChainDB) Close() {
	r.DB.Close()
}

func (r *RedisBlockChainDB) MSetValueToCache(dataMap map[string][]byte, duration int64) {
	for key, value := range dataMap {
		r.SetValueToCache(key, value, duration)
	}
}

func (r *RedisBlockChainDB) SetValueToCache(key string, value []byte, duration int64) {
	r.DB.Set(key, value, time.Duration(duration))
}

func (r *RedisBlockChainDB) GetValueFromCache(key string) string {
	val, err := r.DB.Get(key).Result()
	if err != nil {
		return ""
	}
	return val
}
