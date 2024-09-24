package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
)

type CategoryRepository struct {
	db     *pgxpool.Pool
	logger *logger.Logger
}

func NewCategoryRepository(
	db *pgxpool.Pool,
	logger *logger.Logger,
) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
	}
}

func (cr *CategoryRepository) GetBurgersByPage(ctx context.Context, limit, page int) ([]entities.Burger, error) {
	// Расчет offset
	// offset := (page - 1) * limit
	return nil, nil
}
