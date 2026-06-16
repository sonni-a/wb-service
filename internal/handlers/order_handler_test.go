package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/sonni-a/wb-service/internal/models"
	"github.com/sonni-a/wb-service/internal/service"
	"github.com/sonni-a/wb-service/internal/service/mock_service"
)

func validTestOrder() models.Order {
	orderUID := "550e8400-e29b-41d4-a716-446655440000"
	track := "WBIL12345678"
	empty := ""

	return models.Order{
		OrderUID:          orderUID,
		TrackNumber:       track,
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: &empty,
		CustomerID:        "customer-1",
		DeliveryService:   "meest",
		ShardKey:          "1",
		SmID:              1,
		OofShard:          "1",
		Delivery: models.Delivery{
			OrderUID: orderUID,
			Name:     "John Doe",
			Phone:    "+12345678901",
			Zip:      "12345",
			City:     "Moscow",
			Address:  "Red Square 1",
			Region:   "Moscow",
			Email:    "john@example.com",
		},
		Payment: models.Payment{
			OrderUID:     orderUID,
			Transaction:  orderUID,
			RequestID:    &empty,
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1000,
			PaymentDt:    1700000000,
			Bank:         "AlphaBank",
			DeliveryCost: 100,
			GoodsTotal:   900,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				OrderUID:    orderUID,
				ChrtID:      123456,
				TrackNumber: track,
				Price:       500,
				RID:         "rid-1",
				Name:        "T-Shirt",
				Sale:        0,
				Size:        "M",
				TotalPrice:  500,
				NmID:        1000001,
				Brand:       "Brand",
				Status:      202,
			},
		},
	}
}

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockOrderServiceInterface(ctrl)
	handler := NewOrderHandler(mockSvc)

	order := validTestOrder()
	body, _ := json.Marshal(order)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockSvc.EXPECT().
		CreateOrder(gomock.Any(), &order).
		Return(nil)

	handler.CreateOrder(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestOrderHandler_CreateOrder_InvalidOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockOrderServiceInterface(ctrl)
	handler := NewOrderHandler(mockSvc)

	order := models.Order{OrderUID: "123"}
	body, _ := json.Marshal(order)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateOrder(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid order")
	}
}

func TestOrderHandler_CreateOrder_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockOrderServiceInterface(ctrl)
	handler := NewOrderHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader([]byte("{bad json")))
	w := httptest.NewRecorder()

	handler.CreateOrder(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid JSON")
	}
}

func TestOrderHandler_GetOrderByUID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockOrderServiceInterface(ctrl)
	handler := NewOrderHandler(mockSvc)

	order := &models.Order{OrderUID: "abc"}

	mockSvc.EXPECT().
		GetOrder(context.Background(), "abc").
		Return(order, nil)

	req := httptest.NewRequest(http.MethodGet, "/order/abc", nil)
	w := httptest.NewRecorder()

	handler.GetOrderByUID(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var got models.Order
	_ = json.NewDecoder(resp.Body).Decode(&got)

	if got.OrderUID != "abc" {
		t.Fatalf("wrong order returned")
	}
}

func TestOrderHandler_GetOrderByUID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockOrderServiceInterface(ctrl)
	handler := NewOrderHandler(mockSvc)

	mockSvc.EXPECT().
		GetOrder(context.Background(), "zzz").
		Return(nil, fmt.Errorf("zzz: %w", service.ErrOrderNotFound))

	req := httptest.NewRequest(http.MethodGet, "/order/zzz", nil)
	w := httptest.NewRecorder()

	handler.GetOrderByUID(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found")
	}
}

func TestOrderHandler_GetOrderByUID_BadURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := NewOrderHandler(nil) // svc не нужен, ошибка до вызова

	req := httptest.NewRequest(http.MethodGet, "/order", nil)
	w := httptest.NewRecorder()

	handler.GetOrderByUID(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400")
	}
}
