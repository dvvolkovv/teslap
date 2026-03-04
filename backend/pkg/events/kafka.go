// Package events provides Kafka producer and consumer abstractions
// for TeslaPay event-driven communication.
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Topic constants matching the architecture specification.
const (
	TopicLedgerEvents         = "ledger.events"
	TopicPaymentEvents        = "payment.events"
	TopicCardEvents           = "card.events"
	TopicKYCEvents            = "kyc.events"
	TopicCryptoEvents         = "crypto.events"
	TopicAuditEvents          = "audit.events"
	TopicNotificationCommands = "notification.commands"
	TopicFraudSignals         = "fraud.signals"
)

// Event represents a TeslaPay domain event published to Kafka.
type Event struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Source        string          `json:"source"`
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	Data          json.RawMessage `json:"data"`
	Metadata      EventMetadata   `json:"metadata"`
	CreatedAt     time.Time       `json:"created_at"`
}

// EventMetadata holds correlation and causation information for tracing.
type EventMetadata struct {
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id,omitempty"`
	ActorID       string `json:"actor_id,omitempty"`
	ActorType     string `json:"actor_type,omitempty"`
	RequestID     string `json:"request_id,omitempty"`
}

// Producer publishes events to Kafka topics.
type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

// ProducerConfig holds configuration for the Kafka producer.
type ProducerConfig struct {
	Brokers      []string
	BatchSize    int
	BatchTimeout time.Duration
	RequiredAcks kafka.RequiredAcks
}

// DefaultProducerConfig returns production-grade defaults.
// Uses RequireAll acks for financial data integrity.
func DefaultProducerConfig(brokers []string) *ProducerConfig {
	return &ProducerConfig{
		Brokers:      brokers,
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
	}
}

// NewProducer creates a new Kafka producer.
func NewProducer(cfg *ProducerConfig, logger *zap.Logger) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		RequiredAcks: kafka.RequiredAcks(cfg.RequiredAcks),
		Async:        false, // Synchronous writes for financial events.
	}

	return &Producer{
		writer: w,
		logger: logger,
	}
}

// Publish sends an event to the specified Kafka topic.
// The aggregate ID is used as the message key to ensure ordering per aggregate.
func (p *Producer) Publish(ctx context.Context, topic string, event *Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(event.AggregateID),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.Type)},
			{Key: "event_id", Value: []byte(event.ID)},
			{Key: "correlation_id", Value: []byte(event.Metadata.CorrelationID)},
		},
		Time: event.CreatedAt,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error("failed to publish event",
			zap.String("topic", topic),
			zap.String("event_type", event.Type),
			zap.String("event_id", event.ID),
			zap.Error(err),
		)
		return fmt.Errorf("publish to %s: %w", topic, err)
	}

	p.logger.Debug("event published",
		zap.String("topic", topic),
		zap.String("event_type", event.Type),
		zap.String("event_id", event.ID),
		zap.String("aggregate_id", event.AggregateID),
	)
	return nil
}

// Close shuts down the producer, flushing any buffered messages.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Consumer reads events from a Kafka topic using a consumer group.
type Consumer struct {
	reader  *kafka.Reader
	logger  *zap.Logger
	handler EventHandler
}

// EventHandler is the interface that consumers implement to process events.
type EventHandler interface {
	HandleEvent(ctx context.Context, event *Event) error
}

// EventHandlerFunc is an adapter to allow the use of ordinary functions as EventHandlers.
type EventHandlerFunc func(ctx context.Context, event *Event) error

// HandleEvent calls the wrapped function.
func (f EventHandlerFunc) HandleEvent(ctx context.Context, event *Event) error {
	return f(ctx, event)
}

// ConsumerConfig holds configuration for the Kafka consumer.
type ConsumerConfig struct {
	Brokers  []string
	Topic    string
	GroupID  string
	MinBytes int
	MaxBytes int
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(cfg *ConsumerConfig, handler EventHandler, logger *zap.Logger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MinBytes: cfg.MinBytes,
		MaxBytes: cfg.MaxBytes,
	})

	return &Consumer{
		reader:  r,
		logger:  logger,
		handler: handler,
	}
}

// Start begins consuming messages in a blocking loop. It respects context cancellation
// for graceful shutdown.
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("consumer started",
		zap.String("topic", c.reader.Config().Topic),
		zap.String("group", c.reader.Config().GroupID),
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consumer shutting down")
			return c.reader.Close()
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil // Context cancelled, clean shutdown.
			}
			c.logger.Error("failed to fetch message", zap.Error(err))
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.logger.Error("failed to unmarshal event",
				zap.Error(err),
				zap.Int64("offset", msg.Offset),
			)
			// Commit the malformed message to avoid blocking the consumer.
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		if err := c.handler.HandleEvent(ctx, &event); err != nil {
			c.logger.Error("failed to handle event",
				zap.String("event_type", event.Type),
				zap.String("event_id", event.ID),
				zap.Error(err),
			)
			// In production, implement dead-letter queue or retry logic here.
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("failed to commit message",
				zap.Int64("offset", msg.Offset),
				zap.Error(err),
			)
		}
	}
}

// Close shuts down the consumer.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
