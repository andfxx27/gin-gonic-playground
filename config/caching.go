package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func NewRedisConnection(addr string, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
		Protocol: 2,
	})

	log.Info().Msg("Connected to redis")

	return redisClient
}
