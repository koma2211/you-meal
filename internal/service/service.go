package service

import (
	"context"
	"time"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	"github.com/koma2211/you-meal/pkg/cache/redis"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/koma2211/you-meal/pkg/validate"
)

type Categorier interface {
	GetCategories(ctx context.Context) (*entities.Menu, error)
	GetMealsByCategoryID(ctx context.Context, categoryId, limit, page int) (entities.MealCategory, error)
	GetMealPageCountByCategoryId(ctx context.Context, categoryId, limit int) (int, error)
}

type Customer interface {
	AddOrder(ctx context.Context, client entities.Client) error
}

type Service struct {
	Categorier
	Customer
}

func NewService(
	repo *repository.Repository,
	cacher redis.Cacher,
	logger *logger.Logger,
	valid validate.Validator,
	recievingTTL time.Duration,
) *Service {
	return &Service{
		Categorier: NewCategoryService(repo.Categorier, repo.Tr, cacher, logger),
		Customer:   NewOrderService(repo.Customer, repo.Tr, valid, logger, recievingTTL),
	}
}
