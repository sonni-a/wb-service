package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sonni-a/wb-service/internal/models"
	"github.com/sonni-a/wb-service/internal/repository"
	"github.com/sonni-a/wb-service/internal/service"
)

type OrderHandler struct {
	service service.OrderServiceInterface
}

func NewOrderHandler(svc service.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{service: svc}
}

// CreateOrder godoc
// @Summary      Create new order
// @Description  Accepts order JSON and stores it in PostgreSQL and cache
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      models.Order  true  "Order data"
// @Success      201    {object}  map[string]string
// @Failure      400    {string}  string  "bad request"
// @Failure      409    {string}  string  "order already exists"
// @Failure      500    {string}  string  "internal error"
// @Router       /order [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.CreateOrder(context.Background(), &order)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				http.Error(w, "order already exists", http.StatusConflict)
				return
			default:
				http.Error(w, "database error: "+pgErr.Message, http.StatusInternalServerError)
				return
			}
		}

		http.Error(w, "failed to insert order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetOrderByUID godoc
// @Summary      Get order by UID
// @Description  Returns order from cache or PostgreSQL
// @Tags         orders
// @Produce      json
// @Param        uid  path      string  true  "Order UID"
// @Success      200  {object}  models.Order
// @Failure      400  {string}  string  "missing or invalid order_uid"
// @Failure      404  {string}  string  "order not found"
// @Failure      500  {string}  string  "internal error"
// @Router       /order/{uid} [get]
func (h *OrderHandler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[2] == "" {
		http.Error(w, "missing order_uid", http.StatusBadRequest)
		return
	}
	orderUID := parts[2]

	order, err := h.service.GetOrder(context.Background(), orderUID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}
