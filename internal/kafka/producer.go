package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sonni-a/wb-service/internal/models"
)

func SendOrder(ctx context.Context, topic string, order *models.Order) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer func() {
		if err := w.Close(); err != nil {
			log.Println("kafka producer close error:", err)
		}
	}()

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: data,
		Time:  time.Now(),
	}

	for i := 0; i < 5; i++ {
		if err := w.WriteMessages(ctx, msg); err != nil {
			log.Printf("Kafka write error, retry %d/5: %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[Kafka Producer] Sent order %s", order.OrderUID)
		return nil
	}

	return err
}
