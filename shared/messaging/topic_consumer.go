package messaging

import (
	"context"
	"encoding/json"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

type TopicConsumer struct {
	kf      *kafka.KafkaClient
	connMgr *ConnectionManager
	topics  []string
}

func NewTopicConsumer(kf *kafka.KafkaClient, connMgr *ConnectionManager, topics []string) *TopicConsumer {
	return &TopicConsumer{
		kf:      kf,
		connMgr: connMgr,
		topics:  topics,
	}
}

func (tc *TopicConsumer) Consume(ctx context.Context) error {
	return tc.kf.Consumer.SubscribeAndConsume(ctx, tc.topics,
		func(ctx context.Context, msg *ckafka.Message) error {
			var kfMsg contracts.KafkaMessage
			if err := json.Unmarshal(msg.Value, &kfMsg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}

			entityID := kfMsg.EntityID

			var payload any
			if kfMsg.Data != nil {
				if err := json.Unmarshal(kfMsg.Data, &payload); err != nil {
					log.Printf("Failed to unmarshal payload: %v", err)
					return err
				}
			}

			clientMsg := contracts.WSMessage{
				Type: *msg.TopicPartition.Topic,
				Data: payload,
			}

			return tc.connMgr.SendMessage(entityID, clientMsg)
		},
	)
}
