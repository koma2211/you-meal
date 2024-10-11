package cacherepository

import (
	"context"
	"time"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Categorier interface {
	GetBurgersCategoryByPage(ctx context.Context, page int) (any, error)
	GetNumberOfPagesByBurgers(ctx context.Context) (any, error)
	SetBurgersCategoryByPage(ctx context.Context, key int, burgersCategory entities.Category) error
	SetNumberOfPagesByBurgers(ctx context.Context, burgersCount int) error
}

type CacheRepository struct {
	Categorier
}

func NewCacheRepository(
	rdb *redis.Client,
	logger *logger.Logger,
	cacheCategoryTTL time.Duration,
) *CacheRepository {
	return &CacheRepository{
		Categorier: NewCacheCategory(rdb, logger, cacheCategoryTTL),
	}
}
