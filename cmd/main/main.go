// Package main starts the Demo Order service.
// The service loads configuration, initializes PostgreSQL, Kafka consumer,
// builds all dependencies, restores cache from DB, starts HTTP API,
// and performs graceful shutdown on OS signals.

// @title Demo Order Service API
// @version 1.0
// @description Demo service that receives orders from Kafka, stores them in PostgreSQL,
// @description and exposes HTTP API with in-memory caching.
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/sonni-a/wb-service/docs"

	"github.com/sonni-a/wb-service/internal/config"
	"github.com/sonni-a/wb-service/internal/db"
	"github.com/sonni-a/wb-service/internal/handlers"
	"github.com/sonni-a/wb-service/internal/kafka"
	"github.com/sonni-a/wb-service/internal/repository"
	"github.com/sonni-a/wb-service/internal/service"
	"github.com/sonni-a/wb-service/internal/shutdown"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	fmt.Println("Starting demo service...")

	cfg := config.Load()

	pool, err := db.NewPool(cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer pool.Close()

	orderRepo := repository.NewOrderRepository(pool)
	cache := service.NewMemoryCache(100)
	orderSvc := service.NewOrderService(orderRepo, cache)
	orderHandler := handlers.NewOrderHandler(orderSvc)

	if err := orderSvc.LoadCache(context.Background()); err != nil {
		log.Println("Failed to load cache:", err)
	}

	consumer := kafka.NewConsumer(
		[]string{cfg.KafkaBrokers},
		"orders",
		"order-service-group",
		orderSvc,
	)
	defer consumer.Close()

	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	go func() {
		if err := consumer.Consume(consumerCtx); err != nil {
			log.Println("Kafka consumer error:", err)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("/swagger/doc.json", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("internal/web/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("internal/web/js"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/web/index.html")
	})

	mux.HandleFunc("/order", orderHandler.CreateOrder)
	mux.HandleFunc("/order/", orderHandler.GetOrderByUID)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		log.Println("HTTP server started at :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	shutdown.GracefulShutdown(srv, consumerCancel)
	if err := consumer.Close(); err != nil {
		log.Println("Error closing Kafka consumer:", err)
	}
}
