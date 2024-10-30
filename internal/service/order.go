package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/koma2211/you-meal/pkg/validate"
)

type OrderService struct {
	orderRepo    repository.Customer
	trRepo       repository.Tr
	valid        validate.Validator
	logger       *logger.Logger
	recievingTTL time.Duration
}

func NewOrderService(
	orderRepo repository.Customer,
	trRepo repository.Tr,
	valid validate.Validator,
	logger *logger.Logger,
	recievingTTL time.Duration,
) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		trRepo:    trRepo,
		logger:    logger,
		valid:     valid,
	}
}

func (os *OrderService) AddOrder(ctx context.Context, client entities.Client) error {
	if ok := os.valid.ValidatePhoneNumber(client.PhoneNumber); !ok {
		return entities.ErrPhoneNumberNotValid
	}

	if len(client.Orders) == 0 {
		return entities.ErrEmptyOrder
	}

	tx, err := os.trRepo.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	exists, err := os.orderRepo.CheckClientExistence(ctx, tx, client.PhoneNumber)
	if err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	var clientId int
	if !exists {
		clientId, err = os.orderRepo.AddClientInfo(ctx, tx, client.PhoneNumber, client.Name)
		if err != nil {
			if err := os.trRepo.Rollback(ctx, tx); err != nil {
				os.logger.ErrorLog.Err(err).Msg(err.Error())
				return err
			}

			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}
	} else {
		clientId, err = os.orderRepo.GetClientIDByPhoneNumber(ctx, tx, client.PhoneNumber)
		if err != nil {
			if err := os.trRepo.Rollback(ctx, tx); err != nil {
				os.logger.ErrorLog.Err(err).Msg(err.Error())
				return err
			}

			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}
	}

	stmtExistence, err := os.orderRepo.PrepareCheckMealExistenceManager(ctx, tx)
	if err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	stmtPrice, err := os.orderRepo.PrepareGetMealPriceByMealIdManager(ctx, tx)
	if err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	var totalPrice float64
	for _, order := range client.Orders {
		exists, err := os.orderRepo.CheckMealExistence(ctx, tx, stmtExistence.Name, order.ID)
		if err != nil {
			if err := os.trRepo.Rollback(ctx, tx); err != nil {
				os.logger.ErrorLog.Err(err).Msg(err.Error())
				return err
			}

			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		if !exists {
			if err := os.trRepo.Rollback(ctx, tx); err != nil {
				os.logger.ErrorLog.Err(err).Msg(err.Error())
				return err
			}
			return entities.ErrMealNotExists
		}

		mealPrice, err := os.orderRepo.GetMealPriceByMealId(ctx, tx, stmtPrice.Name, order.ID)
		if err != nil {
			if err := os.trRepo.Rollback(ctx, tx); err != nil {
				os.logger.ErrorLog.Err(err).Msg(err.Error())
				return err
			}

			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		totalPrice += (mealPrice + float64(order.Quantity))
	}

	receivingAt := time.Now().Add(os.recievingTTL)

	orderId, err := os.orderRepo.PlaceAnOrder(ctx, tx, clientId, totalPrice, receivingAt)
	if err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	if err := os.orderRepo.AddOrderedMeals(ctx, tx, orderId, client.Orders); err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	if *client.IsDelivery {
		err = os.orderRepo.AddDelivery(ctx, tx, orderId, *client.Address, *client.Floor)
		if err != nil {
			if err := os.trRepo.Rollback(ctx, tx); err != nil {
				os.logger.ErrorLog.Err(err).Msg(err.Error())
				return err
			}

			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}
	} else {
		err = os.orderRepo.AddSelfPickups(ctx, tx, orderId)
		if err != nil {
			if err := os.trRepo.Rollback(ctx, tx); err != nil {
				os.logger.ErrorLog.Err(err).Msg(err.Error())
				return err
			}

			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}
	}

	if err := os.trRepo.Commit(ctx, tx); err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	return nil
}
