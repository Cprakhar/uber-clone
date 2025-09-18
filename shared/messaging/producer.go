package messaging

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	pr *kafka.Producer
}

// NewProducer creates a confluent producer with safe defaults.
func NewProducer(brokers []string) (*Producer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":                     strings.Join(brokers, ","),
		"security.protocol":                     "PLAINTEXT",
		"acks":                                  "all",
		"enable.idempotence":                    true,
		"max.in.flight.requests.per.connection": 1,
		"retries":                               5,
		"linger.ms":                             5,
		"batch.size":                            32 * 1024, // 32KB
		"compression.type":                      "zstd",
	}

	pr, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range pr.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed to topic %s [%d] at offset %v: %v\n",
						*ev.TopicPartition.Topic,
						ev.TopicPartition.Partition,
						ev.TopicPartition.Offset,
						ev.TopicPartition.Error,
					)
				} else {
					log.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*ev.TopicPartition.Topic,
						ev.TopicPartition.Partition,
						ev.TopicPartition.Offset,
					)
				}
			}
		}
	}()

	return &Producer{pr: pr}, nil
}

func (p *Producer) SendMessage(topic string, payload []byte, entityID string) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(entityID),
		Value:          payload,
	}

	return p.pr.Produce(msg, nil)
}

func (p *Producer) SendMessageAndWait(topic string, payload []byte, entityID string, timeout time.Duration) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(entityID),
		Value:          payload,
	}

	if err := p.pr.Produce(msg, deliveryChan); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	select {
	case ev := <-deliveryChan:
		m := ev.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
		}
		return nil
	case <-time.After(timeout):
		return context.DeadlineExceeded
	}
}

// Close flushes and shuts down the producer and listeners.
func (p *Producer) Close() {
	if p.pr != nil {
		p.pr.Flush(15 * 1000)
		p.pr.Close()
	}
}
