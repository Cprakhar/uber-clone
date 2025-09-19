package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Consumer wraps a Kafka consumer.
type Consumer struct {
	cr *kafka.Consumer // Kafka consumer instance
}

// NewConsumer creates a confluent consumer with safe defaults.
func newConsumer(brokers []string, groupID string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(brokers, ","),
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
		"session.timeout.ms": 6000,
	}

	cr, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{cr: cr}, nil
}

// MessageHandler defines the function signature for processing Kafka messages.
type MessageHandler func(context.Context, *kafka.Message) error

// subscribeAndConsume subscribes to the given topic and processes messages using the provided handler.
func (c *Consumer) SubscribeAndConsume(ctx context.Context, topics []string, handler MessageHandler) error {
	if err := c.cr.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			e := c.cr.Poll(100)
			if e == nil {
				continue
			}
			switch ev := e.(type) {
			case *kafka.Message:
				if err := handler(ctx, ev); err != nil {
					log.Printf("Error handling message: %v", err)
					continue
				}
				if ev.Headers != nil {
					log.Printf("Headers: %v\n", ev.Headers)
				}
				log.Printf("Message on %v: %s\n", ev.TopicPartition, string(ev.Value))
				if _, err := c.cr.CommitMessage(ev); err != nil {
					log.Printf("Failed to commit message: %v", err)
				}
			case kafka.Error:
				log.Printf("Kafka error: %v, code: %v\n", ev, ev.Code())
			default:
				log.Printf("Ignored event: %v\n", ev)
			}
		}
	}
}

func (c *Consumer) Subscribe(topics []string) error {
	return c.cr.SubscribeTopics(topics, nil)
}

// Close shuts down the consumer.
func (c *Consumer) Close() {
	if c.cr != nil {
		c.cr.Close()
	}
}
