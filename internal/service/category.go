package service

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	"github.com/koma2211/you-meal/pkg/cache/redis"
	"github.com/koma2211/you-meal/pkg/logger"
)

type CategoryService struct {
	categoryRepo repository.Categorier
	trRepo       repository.Tr
	cacher       redis.Cacher
	logger       *logger.Logger
}

func NewCategoryService(
	categoryRepo repository.Categorier,
	trRepo repository.Tr,
	cacher redis.Cacher,
	logger *logger.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		trRepo:       trRepo,
		cacher:       cacher,
		logger:       logger,
	}
}

func (cs *CategoryService) GetCategories(ctx context.Context) (*entities.Menu, error) {
	key := entities.GenerateCategoriesKey()

	menu := &entities.Menu{}
	err := cs.cacher.GetStruct(ctx, key, menu)
	if err == nil {
		return menu, nil
	}

	if err != entities.ErrCacheEmpty {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	categories, err := cs.categoryRepo.GetCategories(ctx)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	*menu = entities.Menu{
		Categories: categories,
	}

	err = cs.cacher.SetStruct(ctx, entities.GenerateCategoriesKey(), *menu)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return menu, nil
}

func (cs *CategoryService) GetMealsByCategoryID(ctx context.Context, categoryId, limit, page int) (entities.MealCategory, error) {
	if categoryId == 0 || page == 0 {
		return entities.MealCategory{}, entities.ErrValidateAddressOrFloor
	}

	tx, err := cs.trRepo.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.MealCategory{}, err
	}

	defer cs.trRepo.Rollback(ctx, tx)

	offset := (page - 1) * limit
	meals, err := cs.categoryRepo.GetMealsByCategoryID(ctx, tx, categoryId, limit, offset)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.MealCategory{}, err
	}

	stmt, err := cs.categoryRepo.PrepareMealIngredientsStatement(ctx, tx)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.MealCategory{}, err
	}

	for mealIndex := 0; mealIndex < len(meals); mealIndex++ {
		ingredients, err := cs.categoryRepo.GetMealIngredientsByMealID(ctx, tx, stmt.Name, meals[mealIndex].ID)
		if err != nil {
			cs.logger.ErrorLog.Err(err).Msg(err.Error())
			return entities.MealCategory{}, err
		}

		meals[mealIndex].Ingredients = ingredients
	}

	if err := cs.trRepo.Commit(ctx, tx); err != nil {
		return entities.MealCategory{}, err
	}

	mealCategory := entities.MealCategory{
		ID:    categoryId,
		Meals: meals,
	}

	return mealCategory, nil
}

func (cs *CategoryService) GetMealPageCountByCategoryId(ctx context.Context, categoryId, limit int) (int, error) {
	categoryKey := strconv.Itoa(categoryId)

	pagesCount, err := cs.cacher.GetVal(ctx, entities.GenerateNumberOfPagesKey(categoryKey))
	if err == nil {
		page, err := strconv.Atoi(pagesCount.(string))
		if err != nil {
			cs.logger.ErrorLog.Err(err).Msg(err.Error())
			return 0, err 
		}
		return page, nil
	}

	if err != entities.ErrCacheEmpty {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}

	mealPages, err := cs.categoryRepo.GetNumberOfPagesMealByCategoryId(ctx, categoryId, limit)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}

	err = cs.cacher.SetVal(ctx, entities.GenerateNumberOfPagesKey(strconv.Itoa(categoryId)), mealPages)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}

	return mealPages, nil
}
