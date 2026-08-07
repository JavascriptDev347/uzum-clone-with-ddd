package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: host + ":" + port})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
