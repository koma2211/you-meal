package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (cr *CategoryRepository) GetBurgersByPage(ctx context.Context, limit, offset int) ([]entities.Burger, error) {
	query := "SELECT m.id, m.title, m.description, m.weight, m.calorie, m.price, m.image_path FROM meals m INNER JOIN categories c ON m.category_id = c.id WHERE c.title = 'Бургеры' ORDER BY m.id LIMIT $1 OFFSET $2;"

	rows, err := cr.db.Query(ctx, query, limit, offset)
	if err != nil {
		if pgx.ErrNoRows == err {
			return nil, entities.ErrEmptyBurgers
		}
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}
	defer rows.Close()

	var burgers []entities.Burger
	for rows.Next() {
		var burger entities.Burger
		if err := rows.Scan(&burger.ID, &burger.Title, &burger.Description, &burger.Weight, &burger.Calorie, &burger.Price, &burger.ImagePath); err != nil {
			cr.logger.ErrorLog.Err(err).Msg(err.Error())
			return nil, err 
		}

		burgers = append(burgers, burger)
	}

	if err := rows.Err(); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return burgers, nil
}

func (cr *CategoryRepository) GetNumberOfPagesByBurgers(ctx context.Context, limit int) (int, error) {
	var totalPages int

	query := "SELECT CEIL(COUNT(*)::decimal / $1) AS total_pages FROM meals m INNER JOIN categories c ON m.category_id = c.id WHERE c.title = 'Бургеры';"
	if err := cr.db.QueryRow(ctx, query, limit).Scan(&totalPages); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}

	return totalPages, nil
}

func (cr *CategoryRepository) CheckExistenceImage(ctx context.Context, burgerId int, fileName string) (bool, error) {
	var result bool 
	query := "SELECT EXISTS (SELECT 1 FROM meals WHERE id = $1 AND image_path = $2);"
	if err := cr.db.QueryRow(ctx, query, burgerId, fileName).Scan(&result); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return false, err
	}

	return result, nil
}