package db

import (
	"fmt"

	"github.com/CharlesChou03/_git/block-chain.git/config"
	"github.com/go-redis/redis"
)

type RedisTutorDB struct {
	DB *redis.Client
}

var (
	RedisDB *RedisTutorDB
)

func SetupRedisDB() *RedisTutorDB {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       0, // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return &RedisTutorDB{DB: client}
}

func (r *RedisTutorDB) Close() {
	r.DB.Close()
}

func (r *RedisTutorDB) DeleteCacheByKey(key string) {
	r.DB.Del(key)
}

func (r *RedisTutorDB) SetValueToCache(key string, value []byte) {
	// currentTime := time.Now().Unix()
	// obj := make(map[string]interface{})
	// obj["value"] = value
	// obj["createdAt"] = currentTime
	// r.DB.HMSet(key, obj)
	r.DB.Set(key, value, 0)
}

func (r *RedisTutorDB) GetValueFromCache(key string, duration int64) string {
	// vals, _ := r.DB.HMGet(key, "value", "createdAt").Result()
	// if len(vals) == 0 {
	// 	return false, ""
	// }
	// value := vals[0]
	// createdAt := vals[1]
	// if value == nil || createdAt == nil {
	// 	return false, ""
	// }
	// currentTime := time.Now().Unix()
	// if currentTime-createdAt.(int64) > duration {
	// 	return true, value.(string)
	// }
	// return false, value.(string)
	val, _ := r.DB.Get(key).Result()
	return val
}

// func (r *RedisTutorDB) GetValueFromCache(key string) map[string]interface {

// 	// vals, _ := r.DB.HMGet(key, "value", "createdAt").Result()
// 	// val, _ := r.DB.Get(key).Result()
// 	obj := make(map[string]interface{})
// 	// obj["value"] = vals[0]
// 	// obj["createdAt"] = vals[1]
// 	return obj
// }
