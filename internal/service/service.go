package service

import (
	"context"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	cacherepository "github.com/koma2211/you-meal/internal/repository/cache_repository"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/koma2211/you-meal/pkg/validate"
)

type Categorier interface {
	GetBurgersCategoryByPage(ctx context.Context, limit, page int) (entities.Category, error)
	GetNumberOfPagesByBurgers(ctx context.Context, page int) (int, error)
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
	cacheRepo *cacherepository.CacheRepository,
	logger *logger.Logger,
	valid validate.Validator,
) *Service {
	return &Service{
		Categorier: NewCategoryService(repo.Categorier, cacheRepo.Categorier, logger),
		Customer: NewOrderService(repo.Customer, repo.Tr, valid, logger),
	}
}
