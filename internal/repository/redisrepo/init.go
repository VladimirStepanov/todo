package redisrepo

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(redisAddr string, redisPort string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisAddr, redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		return nil, err
	}

	return client, nil
}
