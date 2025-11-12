package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sonni-a/wb-service/internal/models"
	"github.com/sonni-a/wb-service/internal/repository/mock_repository"
)

func TestOrderService_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOrderRepo(ctrl)
	cache := NewMemoryCache(2)

	service := NewOrderService(mockRepo, cache)

	ctx := context.Background()
	order := &models.Order{OrderUID: "123"}

	mockRepo.EXPECT().GetOrder(ctx, "123").Return(order, nil)

	got, err := service.GetOrder(ctx, "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != order {
		t.Fatalf("expected order %v, got %v", order, got)
	}

	got2, err := service.GetOrder(ctx, "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got2 != order {
		t.Fatalf("expected order %v, got %v", order, got2)
	}
}

func TestOrderService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOrderRepo(ctrl)
	cache := NewMemoryCache(2)
	service := NewOrderService(mockRepo, cache)

	ctx := context.Background()
	order := &models.Order{OrderUID: "abc"}

	mockRepo.EXPECT().InsertOrder(ctx, order).Return(nil)

	err := service.CreateOrder(ctx, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cached, ok := cache.Get("abc")
	if !ok {
		t.Fatalf("expected order in cache")
	}
	if cached != order {
		t.Fatalf("cached order mismatch")
	}
}

func TestMemoryCache_Eviction(t *testing.T) {
	cache := NewMemoryCache(2)

	order1 := &models.Order{OrderUID: "1"}
	order2 := &models.Order{OrderUID: "2"}
	order3 := &models.Order{OrderUID: "3"}

	cache.Set("1", order1)
	cache.Set("2", order2)
	cache.Set("3", order3)

	if _, ok := cache.Get("1"); ok {
		t.Fatalf("expected order1 to be evicted")
	}
	if _, ok := cache.Get("2"); !ok {
		t.Fatalf("expected order2 to remain")
	}
	if _, ok := cache.Get("3"); !ok {
		t.Fatalf("expected order3 to remain")
	}
}
