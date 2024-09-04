package cache

import (
	"geolocation-service/config"

	"github.com/go-redis/redis"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	// Создаем клиента Redis
	rclient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConf.Host + ":" + cfg.RedisConf.Port,
	})
	return rclient

}
