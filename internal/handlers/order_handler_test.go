package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/sonni-a/wb-service/internal/models"
	"github.com/sonni-a/wb-service/internal/repository"
	"github.com/sonni-a/wb-service/internal/service/mock_service"
)

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_service.NewMockOrderServiceInterface(ctrl)
	handler := NewOrderHandler(mockSvc)

	order := models.Order{OrderUID: "123"}
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
		Return(nil, repository.ErrOrderNotFound)

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
