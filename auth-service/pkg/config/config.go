package config

import (
	"auth-service/internal/utils"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	RedisClient *redis.Client
}

func NewConfig() *Config {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Printf("Connected to Redis: %s", utils.GetEnv("REDIS_ADDR", "localhost:6379"))

	return &Config{
		RedisClient: redisClient,
	}
}
