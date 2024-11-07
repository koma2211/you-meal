package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
)

const (
	prepareMealExistence = "prepare_meal_existence"
	prepareMealPrice     = "prepare_meal_price"
)

type OrderRepository struct {
	db     *pgxpool.Pool
	logger *logger.Logger
}

func NewOrderRepository(
	db *pgxpool.Pool,
	logger *logger.Logger,
) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger,
	}
}

func (or *OrderRepository) AddOrderedMeals(ctx context.Context, tx pgx.Tx, orderId int, orders []entities.OrderedMeals) error {
	query := fmt.Sprintf("INSERT INTO %s (order_id, meal_id, quantity) VALUES ($1, $2, $3);", orderedMealsTable)

	for _, order := range orders {
		if _, err := tx.Exec(ctx, query, orderId, order.ID, order.Quantity); err != nil {
			or.logger.ErrorLog.Err(err).Msg(err.Error())
			return err
		}
	}

	return nil
}

func (or *OrderRepository) AddClientInfo(ctx context.Context, tx pgx.Tx, phoneNumber string, name string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (phone_number, name) VALUES ($1, $2) RETURNING id;", clientsTable)
	if err := tx.QueryRow(ctx, query, phoneNumber, name).Scan(&id); err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}

	return id, nil
}

func (or *OrderRepository) CheckClientExistence(ctx context.Context, tx pgx.Tx, phoneNumber string) (bool, error) {
	var exists bool

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE phone_number = $1);", clientsTable)
	if err := tx.QueryRow(ctx, query, phoneNumber).Scan(&exists); err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return false, err
	}

	return exists, nil
}

func (or *OrderRepository) GetClientIDByPhoneNumber(ctx context.Context, tx pgx.Tx, phoneNumber string) (int, error) {
	var clientId int
	query := fmt.Sprintf("SELECT id FROM %s WHERE phone_number = $1;", clientsTable)
	if err := tx.QueryRow(ctx, query, phoneNumber).Scan(&clientId); err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err
	}

	return clientId, nil
}

func (or *OrderRepository) PlaceAnOrder(ctx context.Context, tx pgx.Tx, clientId int, totalPrice float64, receivingAt time.Time) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (client_id, total_price, receiving_at) VALUES ($1, $2, $3) RETURNING id;", ordersTable)
	if err := tx.QueryRow(ctx, query, clientId, totalPrice, receivingAt).Scan(&id); err != nil {
		or.logger.ErrorLog.Err(err).Msg(fmt.Sprintf("%s:%x", err.Error(), receivingAt))
		return 0, err
	}

	return id, nil
}

func (or *OrderRepository) AddDelivery(ctx context.Context, tx pgx.Tx, orderId int, address string, floor int) error {
	query := fmt.Sprintf("INSERT INTO %s (order_id, address, floor) VALUES ($1, $2, $3);", deliveriesTable)
	_, err := tx.Exec(ctx, query, orderId, address, floor)
	if err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (or *OrderRepository) AddSelfPickups(ctx context.Context, tx pgx.Tx, orderId int) error {
	query := fmt.Sprintf("INSERT INTO %s (order_id) VALUES ($1);", selfPickupsTable)
	_, err := tx.Exec(ctx, query, orderId)
	if err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (or *OrderRepository) PrepareCheckMealExistenceManager(ctx context.Context, tx pgx.Tx) (*pgconn.StatementDescription, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1);", mealsTable)
	stmt, err := tx.Prepare(ctx, prepareMealExistence, query)
	if err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return stmt, nil
}

func (or *OrderRepository) CheckMealExistence(ctx context.Context, tx pgx.Tx, stmtName string, mealId int) (bool, error) {
	var exists bool
	if err := tx.QueryRow(ctx, stmtName, mealId).Scan(&exists); err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return false, err
	}

	return exists, nil
}

func (or *OrderRepository) PrepareGetMealPriceByMealIdManager(ctx context.Context, tx pgx.Tx) (*pgconn.StatementDescription, error) {
	query := fmt.Sprintf("SELECT price FROM %s WHERE id = $1;", mealsTable)
	stmt, err := tx.Prepare(ctx, prepareMealPrice, query)
	if err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return nil, err
	}

	return stmt, nil
}

func (or *OrderRepository) GetMealPriceByMealId(ctx context.Context, tx pgx.Tx, stsmtName string, mealId int) (float64, error) {
	var price float64
	err := tx.QueryRow(ctx, stsmtName, mealId).Scan(&price)
	if err != nil {
		or.logger.ErrorLog.Err(err).Msg(err.Error())
		return 0, err 
	}

	return price, nil 
}
