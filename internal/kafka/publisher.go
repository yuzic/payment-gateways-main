//go:generate mockgen -source publisher.go -destination mocks/publisher.go -package mocks

package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sony/gobreaker"
)

type kafkaPublisher struct {
	writer         *kafka.Writer
	circuitBreaker *gobreaker.CircuitBreaker
}

type KafkaPublisher interface {
	PublishTransaction(ctx context.Context, transactionID string, message []byte, dataFormat string) error
	Close() error
}

var (
	defaultCircuitBreakerSettings = gobreaker.Settings{
		Name:        "KafkaPublisher",
		MaxRequests: 1,
		Interval:    5 * time.Second,
		Timeout:     3 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	}
)

func NewPublisher(kafkaURL string) KafkaPublisher {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		BatchTimeout:           100 * time.Millisecond,
		Async:                  true,
	}

	return &kafkaPublisher{
		writer:         writer,
		circuitBreaker: gobreaker.NewCircuitBreaker(defaultCircuitBreakerSettings),
	}
}

func (p *kafkaPublisher) PublishTransaction(ctx context.Context, transactionID string, message []byte, dataFormat string) error {
	if p.writer == nil {
		return fmt.Errorf("kafka writer not initialized")
	}

	topic, err := getTopic(dataFormat)
	if err != nil {
		return fmt.Errorf("topic resolution failed: %w", err)
	}

	// Execute with circuit breaker protection
	_, err = p.circuitBreaker.Execute(func() (interface{}, error) {
		msg := kafka.Message{
			Key:   []byte(transactionID),
			Value: message,
			Topic: topic,
		}

		err := p.writer.WriteMessages(ctx, msg)
		if err != nil {
			log.Printf("Kafka publish error: %v", err)
			return nil, fmt.Errorf("kafka write failed: %w", err)
		}

		log.Printf("Successfully published message to topic: %s", topic)
		return nil, nil
	})

	return err
}

func (p *kafkaPublisher) Close() error {
	if p.writer != nil {
		if err := p.writer.Close(); err != nil {
			return fmt.Errorf("error closing kafka writer: %w", err)
		}
	}
	return nil
}

func getTopic(dataFormat string) (string, error) {
	switch dataFormat {
	case "application/json":
		return "transactions.json", nil
	case "text/xml", "application/xml":
		return "transactions.xml", nil
	default:
		return "", fmt.Errorf("unsupported data format '%s' - allowed formats: JSON, XML", dataFormat)
	}
}
