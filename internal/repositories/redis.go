package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func RedisInitialize(port, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("host.docker.internal:%s", port),
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to redis server: %v", err)
	}

	return client
}
