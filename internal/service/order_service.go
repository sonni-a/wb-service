package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sonni-a/wb-service/internal/models"
	"github.com/sonni-a/wb-service/internal/repository"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
}

type OrderService struct {
	repo  repository.OrderRepo
	cache Cache
}

var _ OrderServiceInterface = (*OrderService)(nil)

func NewOrderService(repo repository.OrderRepo, cache Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := s.repo.InsertOrder(ctx, order); err != nil {
		return err
	}
	s.cache.Set(order.OrderUID, order)
	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	if order, ok := s.cache.Get(orderUID); ok {
		return order, nil
	}

	order, err := s.repo.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, fmt.Errorf("order %s not found: %w", orderUID, err)
		}
		return nil, err
	}

	s.cache.Set(orderUID, order)
	return order, nil
}

func (s *OrderService) LoadCache(ctx context.Context) error {
	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		return err
	}

	for _, o := range orders {
		fullOrder, err := s.repo.GetOrder(ctx, o.OrderUID)
		if err != nil {
			continue
		}
		s.cache.Set(o.OrderUID, fullOrder)
	}

	return nil
}
