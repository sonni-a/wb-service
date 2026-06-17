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
	"github.com/sonni-a/wb-service/internal/validator"
)

type Consumer struct {
	reader    *kafka.Reader
	dlqWriter *kafka.Writer
	svc       service.OrderServiceInterface
}

func NewConsumer(brokers []string, topic, groupID string, svc service.OrderServiceInterface) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	dlqWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    "orders-dlq",
		Balancer: &kafka.LeastBytes{},
	})

	return &Consumer{
		reader:    r,
		dlqWriter: dlqWriter,
		svc:       svc,
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

		if err := validator.ValidateOrder(&order); err != nil {
			log.Printf("Invalid order (%s): %v", order.OrderUID, err)
			metrics.KafkaProcessingErrorsTotal.Inc()
			c.sendToDLQ(ctx, m, "validation failed")
			continue
		}

		if err := c.svc.CreateOrder(ctx, &order); err != nil {
			if errors.Is(err, service.ErrOrderAlreadyExists) {
				log.Printf("Order %s already exists, skipping duplicate message", order.OrderUID)
				metrics.KafkaMessagesProcessedTotal.Inc()
				continue
			}

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
	payload := map[string]interface{}{
		"original_key":   string(msg.Key),
		"original_value": string(msg.Value),
		"reason":         reason,
		"time":           time.Now(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal DLQ payload: %v", err)
		return
	}

	if err := c.dlqWriter.WriteMessages(ctx, kafka.Message{
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
	if err := c.dlqWriter.Close(); err != nil {
		return err
	}
	return c.reader.Close()
}
