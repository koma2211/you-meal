package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
)

const (
	prepareMealIngredients = "meal_ingredients_prepare"
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

func (cr *CategoryRepository) GetCategories(ctx context.Context) ([]entities.Category, error) {
	query := fmt.Sprintf("SELECT id, title FROM %s;", categoriesTable)

	rows, err := cr.db.Query(ctx, query)
	if err != nil {
		if pgx.ErrNoRows == err {
			return nil, entities.ErrEmptyCategories
		}
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(&category.ID, &category.Title); err != nil {
			cr.logger.ErrorLog.Err(err).Msg(err.Error())
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return categories, nil
}

func (cr *CategoryRepository) GetMealsByCategoryID(ctx context.Context, tx pgx.Tx, categoryId int, limit, offset int) ([]entities.Meal, error) {
	query := fmt.Sprintf("SELECT id, title, description, weight, calorie, price, image_path FROM %s WHERE category_id = $1 LIMIT $2 OFFSET $3;", mealsTable)
	rows, err := tx.Query(ctx, query, categoryId, limit, offset)
	if err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	defer rows.Close()

	var meals []entities.Meal
	for rows.Next() {
		var meal entities.Meal
		if err := rows.Scan(&meal.ID, &meal.Title, &meal.Description, &meal.Weight, &meal.Calorie, &meal.Price, &meal.ImagePath); err != nil {
			cr.logger.ErrorLog.Err(err).Msg(err.Error())
			return nil, err
		}

		meals = append(meals, meal)
	}

	if err := rows.Err(); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return meals, nil
}

func (cr *CategoryRepository) PrepareMealIngredientsStatement(ctx context.Context, tx pgx.Tx) (*pgconn.StatementDescription, error) {
	query := fmt.Sprintf("SELECT id, title FROM %s WHERE meal_id = $1;", ingredientsTable)
	stmt, err := tx.Prepare(ctx, prepareMealIngredients, query)
	if err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err 
	}

	return stmt, nil
}

func (cr *CategoryRepository) GetMealIngredientsByMealID(ctx context.Context, tx pgx.Tx, stmtName string, mealId int) ([]entities.Ingredient, error) {
	var ingredients []entities.Ingredient

	rows, err := tx.Query(ctx, stmtName, mealId)
	if err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
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

	if err := rows.Err(); err != nil {
		cr.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return ingredients, nil
}


func (cr *CategoryRepository) GetNumberOfPagesMealByCategoryId(ctx context.Context, categoryId, limit int) (int, error) {
	var totalPages int

	query := "SELECT CEIL(COUNT(*)::decimal / $1) AS total_pages FROM meals WHERE category_id = $2;"
	if err := cr.db.QueryRow(ctx, query, limit, categoryId).Scan(&totalPages); err != nil {
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
