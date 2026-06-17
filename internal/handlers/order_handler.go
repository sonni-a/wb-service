package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sonni-a/wb-service/internal/models"
	"github.com/sonni-a/wb-service/internal/service"
	"github.com/sonni-a/wb-service/internal/validator"
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
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateOrder(&order); err != nil {
		http.Error(w, "invalid order: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.CreateOrder(r.Context(), &order)
	if err != nil {
		if errors.Is(err, service.ErrOrderAlreadyExists) {
			http.Error(w, "order already exists", http.StatusConflict)
			return
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("database error on create order: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("failed to create order: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
	orderUID := r.PathValue("uid")
	if orderUID == "" {
		http.Error(w, "missing order_uid", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderUID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		log.Printf("failed to get order %s: %v", orderUID, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}
