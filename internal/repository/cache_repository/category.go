package cacherepository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type CacheCategory struct {
	cache            *redis.Client
	logger           *logger.Logger
	categoryCacheTTL time.Duration
}

func NewCacheCategory(
	cache *redis.Client,
	logger *logger.Logger,
	categoryCacheTTL time.Duration,
) *CacheCategory {
	return &CacheCategory{
		cache:            cache,
		logger:           logger,
		categoryCacheTTL: categoryCacheTTL,
	}
}

func (cc *CacheCategory) GetBurgersCategoryByPage(ctx context.Context, page int) (any, error) {
	result, err := cc.cache.Get(ctx, generateBurgerPageKey(page)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, entities.ErrEmptyBurgers
		}
		cc.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	var input entities.Category
	err = json.Unmarshal([]byte(result), &input)
	if err != nil {
		cc.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return input, nil
}

func (cc *CacheCategory) GetNumberOfPagesByBurgers(ctx context.Context) (any, error) {
	result, err := cc.cache.Get(ctx, generateBurgerKeyCount()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, entities.ErrEmptyBurgers
		}
		cc.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	burgersCount, err := strconv.Atoi(result)
	if err != nil {
		cc.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return burgersCount, nil
}

func (cc *CacheCategory) SetBurgersCategoryByPage(ctx context.Context, key int, burgersCategory entities.Category) error {
	body, err := json.Marshal(burgersCategory)
	if err != nil {
		cc.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	return cc.cache.Set(ctx, generateBurgerPageKey(key), string(body), cc.categoryCacheTTL).Err()
}

func (cc *CacheCategory) SetNumberOfPagesByBurgers(ctx context.Context, burgersCount int) error {
	return cc.cache.Set(ctx, generateBurgerKeyCount(), burgersCount, cc.categoryCacheTTL).Err()
}

func generateBurgerKeyCount() string {
	return "burgers-count"
}

func generateBurgerPageKey(page int) string {
	return fmt.Sprintf("burgers:page:%d", page)
}

