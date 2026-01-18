// Package messaging provides Redpanda (Kafka-compatible) client configuration.
// Redpanda is a ZooKeeper-free, high-performance streaming platform.
package messaging

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// RedpandaConfig holds Redpanda connection configuration
type RedpandaConfig struct {
	// Brokers is a list of Redpanda broker addresses
	Brokers []string
	// Username for SASL authentication
	Username string
	// Password for SASL authentication
	Password string
	// TLSEnabled enables TLS connections
	TLSEnabled bool
	// ConsumerGroup for consumer group coordination
	ConsumerGroup string
	// ClientID for client identification
	ClientID string
	// Producer configuration
	Producer ProducerConfig
	// Consumer configuration
	Consumer ConsumerConfig
}

// ProducerConfig holds producer-specific configuration
type ProducerConfig struct {
	// BatchMaxBytes is the maximum size of a batch
	BatchMaxBytes int32
	// LingerMs is the time to wait for more messages before sending
	LingerMs int
	// Compression type (none, gzip, snappy, lz4, zstd)
	Compression string
	// Acks required (none, leader, all)
	Acks string
	// MaxRetries before failing
	MaxRetries int
}

// ConsumerConfig holds consumer-specific configuration
type ConsumerConfig struct {
	// AutoOffsetReset (earliest, latest)
	AutoOffsetReset string
	// MaxPollRecords per poll
	MaxPollRecords int
	// SessionTimeout for consumer group
	SessionTimeout time.Duration
}

// DefaultRedpandaConfig returns default configuration
func DefaultRedpandaConfig() RedpandaConfig {
	return RedpandaConfig{
		Brokers:  []string{"localhost:9092"},
		ClientID: "omniroute-service",
		Producer: ProducerConfig{
			BatchMaxBytes: 1048576, // 1MB
			LingerMs:      5,
			Compression:   "snappy",
			Acks:          "all",
			MaxRetries:    3,
		},
		Consumer: ConsumerConfig{
			AutoOffsetReset: "earliest",
			MaxPollRecords:  500,
			SessionTimeout:  30 * time.Second,
		},
	}
}

// NewRedpandaClient creates a new Redpanda client
func NewRedpandaClient(cfg RedpandaConfig) (*kgo.Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),

		// Producer settings
		kgo.ProducerBatchMaxBytes(cfg.Producer.BatchMaxBytes),
		kgo.ProducerLinger(time.Duration(cfg.Producer.LingerMs) * time.Millisecond),
		kgo.RequiredAcks(parseAcks(cfg.Producer.Acks)),
		kgo.ProducerBatchCompression(parseCompression(cfg.Producer.Compression)),
		kgo.RequestRetries(cfg.Producer.MaxRetries),

		// Consumer settings
		kgo.ConsumeResetOffset(parseOffset(cfg.Consumer.AutoOffsetReset)),
		kgo.SessionTimeout(cfg.Consumer.SessionTimeout),

		// Retry settings
		kgo.RetryBackoffFn(func(attempt int) time.Duration {
			return time.Duration(attempt*100) * time.Millisecond
		}),
	}

	// Add consumer group if specified
	if cfg.ConsumerGroup != "" {
		opts = append(opts, kgo.ConsumerGroup(cfg.ConsumerGroup))
	}

	// Add SASL authentication
	if cfg.Username != "" {
		mechanism := scram.Auth{
			User: cfg.Username,
			Pass: cfg.Password,
		}.AsSha256Mechanism()
		opts = append(opts, kgo.SASL(mechanism))
	}

	// Add TLS if enabled
	if cfg.TLSEnabled {
		opts = append(opts, kgo.DialTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
		}))
	}

	return kgo.NewClient(opts...)
}

// Event represents a message event
type Event struct {
	Topic     string
	Key       string
	Value     []byte
	Headers   map[string]string
	Timestamp time.Time
	Partition int32
	Offset    int64
}

// EventHandler handles incoming events
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

// EventProducer publishes events to Redpanda
type EventProducer struct {
	client *kgo.Client
}

// NewEventProducer creates a new event producer
func NewEventProducer(client *kgo.Client) *EventProducer {
	return &EventProducer{client: client}
}

// Publish publishes a single event
func (p *EventProducer) Publish(ctx context.Context, topic, key string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
		Headers: []kgo.RecordHeader{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}

	results := p.client.ProduceSync(ctx, record)
	return results.FirstErr()
}

// PublishAsync publishes an event asynchronously
func (p *EventProducer) PublishAsync(ctx context.Context, topic, key string, event interface{}, callback func(error)) {
	data, err := json.Marshal(event)
	if err != nil {
		callback(fmt.Errorf("marshal: %w", err))
		return
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	}

	p.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
		callback(err)
	})
}

// PublishBatch publishes multiple events
func (p *EventProducer) PublishBatch(ctx context.Context, topic string, events []Event) error {
	records := make([]*kgo.Record, len(events))
	for i, event := range events {
		records[i] = &kgo.Record{
			Topic: topic,
			Key:   []byte(event.Key),
			Value: event.Value,
		}
	}

	results := p.client.ProduceSync(ctx, records...)
	return results.FirstErr()
}

// EventConsumer consumes events from Redpanda
type EventConsumer struct {
	client   *kgo.Client
	handlers map[string]EventHandler
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(client *kgo.Client) *EventConsumer {
	return &EventConsumer{
		client:   client,
		handlers: make(map[string]EventHandler),
	}
}

// Subscribe subscribes to topics with a handler
func (c *EventConsumer) Subscribe(topics []string, handler EventHandler) {
	c.client.AddConsumeTopics(topics...)
	for _, topic := range topics {
		c.handlers[topic] = handler
	}
}

// Start starts consuming events
func (c *EventConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fetches := c.client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return nil
		}

		fetches.EachError(func(topic string, partition int32, err error) {
			fmt.Printf("fetch error topic %s partition %d: %v\n", topic, partition, err)
		})

		fetches.EachRecord(func(record *kgo.Record) {
			handler, ok := c.handlers[record.Topic]
			if !ok {
				return
			}

			event := Event{
				Topic:     record.Topic,
				Key:       string(record.Key),
				Value:     record.Value,
				Timestamp: record.Timestamp,
				Partition: record.Partition,
				Offset:    record.Offset,
			}

			if err := handler.Handle(ctx, event); err != nil {
				fmt.Printf("handler error: %v\n", err)
			}
		})

		// Commit offsets after processing
		if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
			fmt.Printf("commit error: %v\n", err)
		}
	}
}

// Close closes the consumer
func (c *EventConsumer) Close() {
	c.client.Close()
}

// Helper functions

func parseAcks(acks string) kgo.Acks {
	switch acks {
	case "none":
		return kgo.NoAck()
	case "leader":
		return kgo.LeaderAck()
	default:
		return kgo.AllISRAcks()
	}
}

func parseOffset(offset string) kgo.Offset {
	switch offset {
	case "latest":
		return kgo.NewOffset().AtEnd()
	default:
		return kgo.NewOffset().AtStart()
	}
}

func parseCompression(compression string) kgo.CompressionCodec {
	switch compression {
	case "gzip":
		return kgo.GzipCompression()
	case "snappy":
		return kgo.SnappyCompression()
	case "lz4":
		return kgo.Lz4Compression()
	case "zstd":
		return kgo.ZstdCompression()
	default:
		return kgo.NoCompression()
	}
}
