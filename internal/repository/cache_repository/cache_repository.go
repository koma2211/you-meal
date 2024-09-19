package cacherepository

import "github.com/redis/go-redis/v9"

type CacheRepository struct {

}

func NewCacheRepository(
	rdb *redis.Client,
) *CacheRepository {
	return &CacheRepository{}
}