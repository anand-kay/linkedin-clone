package core

import (
	"github.com/anand-kay/linkedin-clone/config"

	"github.com/go-redis/redis/v8"
)

func newRedisDB(config *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: "",
		DB:       0,
	})
}
