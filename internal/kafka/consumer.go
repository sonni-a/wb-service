package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sonni-a/wb-service/internal/metrics"
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
			if err := c.reader.SetOffset(kafka.FirstOffset); err != nil {
				log.Printf("failed to set Kafka offset: %v", err)
			}
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
			metrics.KafkaProcessingErrorsTotal.Inc()
			c.sendToDLQ(ctx, m, "invalid JSON")
			continue
		}

		if err := order.Validate(); err != nil {
			log.Printf("Invalid order (%s): %v", order.OrderUID, err)
			metrics.KafkaProcessingErrorsTotal.Inc()
			c.sendToDLQ(ctx, m, "validation failed")
			continue
		}

		if err := c.svc.CreateOrder(ctx, &order); err != nil {
			log.Printf("Failed to save order %s: %v", order.OrderUID, err)
			metrics.KafkaProcessingErrorsTotal.Inc()
			c.sendToDLQ(ctx, m, "DB write failed")
			continue
		}

		log.Printf("Order %s saved successfully", order.OrderUID)
		metrics.KafkaMessagesProcessedTotal.Inc()
	}
}

func (c *Consumer) sendToDLQ(ctx context.Context, msg kafka.Message, reason string) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "orders-dlq",
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	payload := map[string]interface{}{
		"original_key":   string(msg.Key),
		"original_value": string(msg.Value),
		"reason":         reason,
		"time":           time.Now(),
	}
	data, _ := json.Marshal(payload)

	if err := w.WriteMessages(ctx, kafka.Message{
		Key:   msg.Key,
		Value: data,
		Time:  time.Now(),
	}); err != nil {
		log.Printf("Failed to send to DLQ: %v", err)
	} else {
		log.Printf("Message sent to DLQ (key=%s): %s", msg.Key, reason)
		metrics.KafkaDLQMessagesTotal.Inc()
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
