package service

import (
	"context"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	"github.com/koma2211/you-meal/pkg/logger"
)

type CategoryService struct {
	categoryRepo repository.Categorier
	logger       *logger.Logger
}

func NewCategoryService(
	categoryRepo repository.Categorier,
	logger *logger.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (cs *CategoryService) GetBurgersByPage(ctx context.Context, limit, page int) ([]entities.Burger, error) {
	return nil, nil
}
