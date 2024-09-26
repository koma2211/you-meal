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

func (cr *CategoryRepository) GetBurgersCount(ctx context.Context, limit int) (int, error) {
	var totalPages int 

	query := "SELECT CEIL(COUNT(*)::decimal / $1) AS total_pages FROM meals m INNER JOIN categories c ON m.category_id = c.id WHERE c.title = 'Бургеры'"
	if err := cr.db.QueryRow(ctx, query, limit).Scan(&totalPages); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}

	return totalPages, nil  
}