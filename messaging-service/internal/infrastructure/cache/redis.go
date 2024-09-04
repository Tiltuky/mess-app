package cache

import "github.com/go-redis/redis"

func NewRedisClient(host, port string) *redis.Client {
	// Создаем клиента Redis
	rclient := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})
	return rclient

}
