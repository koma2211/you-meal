package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
)

const (
	categoriesTable   = "categories"
	mealsTable        = "meals"
	clientsTable      = "clients"
	ordersTable       = "orders"
	orderedMealsTable = "ordered_meals"
	deliveriesTable   = "deliveries"
	selfPickupsTable  = "self_pickups"
)

type Categorier interface {
	GetBurgersCategoryByPage(ctx context.Context, limit, offset int) (entities.Category, error)
	GetBurgerIngredientsById(ctx context.Context, burgerId int) ([]entities.Ingredient, error)
	GetNumberOfPagesByBurgers(ctx context.Context, limit int) (int, error)
	CheckExistenceImage(ctx context.Context, burgerId int, fileName string) (bool, error)
}

type Customer interface {
	AddClientInfo(ctx context.Context, tx pgx.Tx, phoneNumber string, name string) (int, error)
	PlaceAnOrder(ctx context.Context, tx pgx.Tx, clientId int, totalPrice float64, receivingAt time.Time) (int, error)
	AddOrderedMeals(ctx context.Context, tx pgx.Tx, orderId int, orders []entities.OrderedMeals) error
	AddDelivery(ctx context.Context, tx pgx.Tx, orderId int, address string, floor int) error
	AddSelfPickups(ctx context.Context, tx pgx.Tx, orderId int) error
	TotalAmountOfOrders(ctx context.Context, tx pgx.Tx, orders []entities.OrderedMeals) (float64, error)
	CheckClientExistence(ctx context.Context, tx pgx.Tx, phoneNumber string) (bool, error)
	GetClientIDByPhoneNumber(ctx context.Context, tx pgx.Tx, phoneNumber string) (int, error)
}

type Tr interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}

type Repository struct {
	Categorier
	Customer
	Tr
}

func NewRepository(
	db *pgxpool.Pool,
	logger *logger.Logger,
) *Repository {
	return &Repository{
		Categorier: NewCategoryRepository(db, logger),
		Customer:   NewOrderRepository(db, logger),
		Tr:         NewTransactionManager(db),
	}
}
