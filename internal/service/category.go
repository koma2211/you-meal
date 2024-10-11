package service

import (
	"context"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	cacherepository "github.com/koma2211/you-meal/internal/repository/cache_repository"
	"github.com/koma2211/you-meal/pkg/logger"
)

type CategoryService struct {
	categoryRepo repository.Categorier
	categoryCacheRepo cacherepository.Categorier
	logger       *logger.Logger
}

func NewCategoryService(
	categoryRepo repository.Categorier,
	categoryCacheRepo cacherepository.Categorier,
	logger *logger.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		categoryCacheRepo: categoryCacheRepo,
		logger:       logger,
	}
}

func (cs *CategoryService) GetBurgersCategoryByPage(ctx context.Context, limit, page int) (entities.Category, error) {
	values, err := cs.categoryCacheRepo.GetBurgersCategoryByPage(ctx, page)
	if err == nil {
		return values.(entities.Category), nil 
	}

	if err != entities.ErrEmptyBurgers {
		return entities.Category{}, err
	}

	// Расчет offset
	offset := (page - 1) * limit
	burgersCategory, err := cs.categoryRepo.GetBurgersCategoryByPage(ctx, limit, offset)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.Category{}, err
	}

	burgersLen := len(burgersCategory.Meals)
	
	for i := 0; i < burgersLen; i++ {
		ingredients, err := cs.categoryRepo.GetBurgerIngredientsById(ctx, burgersCategory.Meals[i].ID)
		if err != nil {
			cs.logger.ErrorLog.Err(err).Msg(err.Error())
			return entities.Category{}, err
		}

		burgersCategory.Meals[i].Ingredients = ingredients
	}

	err = cs.categoryCacheRepo.SetBurgersCategoryByPage(ctx, page, burgersCategory)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return entities.Category{}, err
	}

	return burgersCategory, nil
}

func (cs *CategoryService) GetNumberOfPagesByBurgers(ctx context.Context, page int) (int, error) {
	// check burgers count exists. If not exists then require to get data from database...
	val, err := cs.categoryCacheRepo.GetNumberOfPagesByBurgers(ctx)
	if err == nil {
		return val.(int), nil
	}

	if err != entities.ErrEmptyBurgers {
		return 0, err
	}

	burgersCount, err := cs.categoryRepo.GetNumberOfPagesByBurgers(ctx, page)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err 
	}

	// added burgers count in redis cache if not exists...
	err = cs.categoryCacheRepo.SetNumberOfPagesByBurgers(ctx, burgersCount)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}
	
	return burgersCount, nil 
}

