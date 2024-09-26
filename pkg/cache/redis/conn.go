package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func RedisCacheMemoryConn(redisSource string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisSource)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
