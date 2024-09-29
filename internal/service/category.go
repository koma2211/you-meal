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

func (cs *CategoryService) GetBurgersByPage(ctx context.Context, limit, page int) ([]entities.Burger, error) {
	values, err := cs.categoryCacheRepo.GetBurgersByPage(ctx, page)
	if err == nil {
		return values.([]entities.Burger), nil 
	}

	if err != entities.ErrEmptyBurgers {
		return nil, err
	}

	// Расчет offset
	offset := (page - 1) * limit
	burgers, err := cs.categoryRepo.GetBurgersByPage(ctx, limit, offset)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	err = cs.categoryCacheRepo.SetBurgersByPage(ctx, page, burgers)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return burgers, nil
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

func (cs *CategoryService) CheckExistenceImage(ctx context.Context, burgerId int, fileName string) error {
	err := cs.categoryCacheRepo.CheckImageExists(ctx, fileName)
	if err == nil {
		return nil
	}
	
	if err != entities.ErrImageNotExists {
		return err
	}

	exists, err := cs.categoryRepo.CheckExistenceImage(ctx, burgerId, fileName)
	if err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return err 
	}

	if !exists {
		return entities.ErrImageNotExists
	}

	if err := cs.categoryCacheRepo.SetImagePath(ctx, fileName); err != nil {
		cs.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	return nil
}
