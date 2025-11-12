package main

import (
	"context"
	"log"
	"time"

	"github.com/sonni-a/wb-service/internal/faker"
	"github.com/sonni-a/wb-service/internal/kafka"
)

func main() {
	faker.Init()

	ctx := context.Background()
	topic := "orders"

	for i := 0; i < 5; i++ {
		order := faker.GenerateFakeOrder()

		if err := kafka.SendOrder(ctx, topic, &order); err != nil {
			log.Println("Failed to send order:", err)
		} else {
			log.Println("Order sent successfully:", order.OrderUID)
		}

		time.Sleep(time.Second)
	}
}
