package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sonni-a/wb-service/internal/models"
	"github.com/sonni-a/wb-service/internal/service"
)

type Consumer struct {
	reader *kafka.Reader
	svc    service.OrderServiceInterface
}

func NewConsumer(brokers []string, topic, groupID string, svc service.OrderServiceInterface) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		reader: r,
		svc:    svc,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	log.Println("Kafka consumer starting...")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		testCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		_, err := c.reader.ReadMessage(testCtx)
		cancel()

		if err == nil {
			c.reader.SetOffset(kafka.FirstOffset)
			log.Println("Kafka consumer ready.")
			break
		}

		log.Println("Kafka not ready, retrying...")
		time.Sleep(2 * time.Second)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopped")
			log.Println("Kafka consumer exited gracefully")
			return nil
		default:
		}

		msgCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		m, err := c.reader.ReadMessage(msgCtx)
		cancel()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Println("Kafka read timeout, retrying...")
				continue
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Invalid JSON: %v", err)
			continue
		}

		if err := order.Validate(); err != nil {
			log.Printf("Invalid order (%s): %v", order.OrderUID, err)
			continue
		}

		if err := c.svc.CreateOrder(ctx, &order); err != nil {
			log.Printf("Failed to save order %s: %v", order.OrderUID, err)
			continue
		}

		log.Printf("Order %s saved successfully", order.OrderUID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
