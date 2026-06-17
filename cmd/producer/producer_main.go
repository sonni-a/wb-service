package main

import (
	"context"
	"log"
	"time"

	"github.com/sonni-a/wb-service/internal/config"
	"github.com/sonni-a/wb-service/internal/kafka"
)

func main() {
	initFaker()

	cfg := config.Load()

	ctx := context.Background()
	topic := "orders"
	brokers := []string{cfg.KafkaBrokers}

	for i := 0; i < 5; i++ {
		order := generateFakeOrder()

		if err := kafka.SendOrder(ctx, brokers, topic, &order); err != nil {
			log.Println("Failed to send order:", err)
		} else {
			log.Println("Order sent successfully:", order.OrderUID)
		}

		time.Sleep(time.Second)
	}
}
