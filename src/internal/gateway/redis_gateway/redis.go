package redis_gateway

import "github.com/redis/go-redis/v9"

func NewRedisClient(addr, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0, // use default DB
	})

	return rdb
}
