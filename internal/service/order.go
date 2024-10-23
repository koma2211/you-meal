package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/koma2211/you-meal/pkg/validate"
)

type OrderService struct {
	orderRepo repository.Customer
	trRepo    repository.Tr
	valid     validate.Validator
	logger    *logger.Logger
}

func NewOrderService(
	orderRepo repository.Customer,
	trRepo repository.Tr,
	valid validate.Validator,
	logger *logger.Logger,
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

	clientId, err := os.orderRepo.AddClientInfo(ctx, tx, client.PhoneNumber, client.Name)
	if err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	totalPrice, err := os.orderRepo.TotalAmountOfOrders(ctx, tx, client.Orders)
	if err != nil {
		if err := os.trRepo.Rollback(ctx, tx); err != nil {
			os.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}

		os.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	orderId, err := os.orderRepo.PlaceAnOrder(ctx, tx, clientId, totalPrice)
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

	if client.IsDelivery {
		err = os.orderRepo.AddDelivery(ctx, tx, orderId, client.Address, client.Floor)
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
