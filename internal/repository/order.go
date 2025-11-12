package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonni-a/wb-service/internal/models"
)

type OrderRepo interface {
	InsertOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, uid string) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
}

type OrderRepository struct {
	db *pgxpool.Pool
}

var _ OrderRepo = (*OrderRepository)(nil)

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) InsertOrder(ctx context.Context, order *models.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	log.Printf("Parsed order before insert: %+v", order)
	_, err = tx.Exec(ctx, InsertOrderQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	log.Printf("Parsed delivery before insert: %+v", order.Delivery)
	_, err = tx.Exec(ctx, InsertDeliveryQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("insert delivery: %w", err)
	}

	_, err = tx.Exec(ctx, InsertPaymentQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, InsertItemQuery,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	order := &models.Order{}

	err := r.db.QueryRow(ctx, GetOrderWithJoinsQuery, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order with joins: %w", err)
	}

	order.Delivery.OrderUID = order.OrderUID
	order.Payment.OrderUID = order.OrderUID

	rows, err := r.db.Query(ctx, GetItemsQuery, orderUID)
	if err != nil {
		return nil, fmt.Errorf("get items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		item.OrderUID = order.OrderUID
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	rows, err := r.db.Query(ctx, GetAllOrdersQuery)
	if err != nil {
		return nil, fmt.Errorf("query all orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := models.Order{}
		if err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		); err != nil {
			return nil, fmt.Errorf("scan order row: %w", err)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}
