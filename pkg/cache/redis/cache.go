package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/redis/go-redis/v9"
)

type Cacher interface {
	SetStruct(ctx context.Context, key string, val interface{}) error
	GetStruct(ctx context.Context, key string, valType interface{}) error
	SetVal(ctx context.Context, key string, val interface{}) error
	GetVal(ctx context.Context, key string) (any, error)
}

type Cache struct {
	rdb      *redis.Client
	cacheTTL time.Duration
}

func NewCache(
	rdb *redis.Client,
	cacheTTL time.Duration,
) *Cache {
	return &Cache{
		rdb:      rdb,
		cacheTTL: cacheTTL,
	}
}

func (c *Cache) SetStruct(ctx context.Context, key string, val interface{}) error {
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = c.rdb.Set(ctx, key, body, c.cacheTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetStruct(ctx context.Context, key string, valType interface{}) error {
	result, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return entities.ErrCacheEmpty
		}

		return err
	}

	err = json.Unmarshal([]byte(result), &valType)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) SetVal(ctx context.Context, key string, val interface{}) error {
	err := c.rdb.Set(ctx, key, val, c.cacheTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetVal(ctx context.Context, key string) (any, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, entities.ErrCacheEmpty
		}
		return nil, err
	}

	return val, nil
}
