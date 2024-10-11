package repository

import (
	"context"
	"fmt"

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

func (cr *CategoryRepository) GetBurgersCategoryByPage(ctx context.Context, limit, offset int) (entities.Category, error) {
	var category entities.Category

	query := fmt.Sprintf("SELECT id, title FROM %s WHERE title = 'Бургеры'", categoriesTable)

	err := cr.db.QueryRow(ctx, query).Scan(&category.ID, &category.Title)
	if err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.Category{}, err
	}

	query = fmt.Sprintf("SELECT id, title, description, weight, calorie, price, image_path FROM %s WHERE category_id = $1 ORDER BY id LIMIT $2 OFFSET $3;", mealsTable)
	rows, err := cr.db.Query(ctx, query, category.ID, limit, offset)
	if err != nil {
		if pgx.ErrNoRows == err {
			return entities.Category{}, entities.ErrEmptyBurgers
		}
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.Category{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var burgerMeal entities.Meal
		if err := rows.Scan(&burgerMeal.ID, &burgerMeal.Title, &burgerMeal.Description, &burgerMeal.Weight, &burgerMeal.Calorie, &burgerMeal.Price, &burgerMeal.ImagePath); err != nil {
			cr.logger.ErrorLog.Err(err).Msg(err.Error())
			return entities.Category{}, err
		}

		category.Meals = append(category.Meals, burgerMeal)
	}

	if err := rows.Err(); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.Category{}, err
	}

	return category, nil
}

func (cr *CategoryRepository) GetBurgerIngredientsById(ctx context.Context, burgerId int) ([]entities.Ingredient, error) {
	var ingredients []entities.Ingredient
	query := "SELECT id, title FROM ingredients WHERE meal_id = $1 ORDER BY id;"

	rows, err := cr.db.Query(ctx, query, burgerId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var ingredient entities.Ingredient
		if err := rows.Scan(&ingredient.Id, &ingredient.Title); err != nil {
			cr.logger.ErrorLog.Err(err).Msg(err.Error())
			return nil, err
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
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
