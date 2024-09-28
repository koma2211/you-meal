package cacherepository

import (
	"context"
	"time"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Categorier interface {
	GetBurgersByPage(ctx context.Context, page int) (any, error)
	GetNumberOfPagesByBurgers(ctx context.Context) (any, error)
	SetBurgersByPage(ctx context.Context, key int, burgers []entities.Burger) error
	SetNumberOfPagesByBurgers(ctx context.Context, burgersCount int) error
	SetImagePath(ctx context.Context, imagePath string) error
	CheckImageExists(ctx context.Context, imagePath string) error
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
