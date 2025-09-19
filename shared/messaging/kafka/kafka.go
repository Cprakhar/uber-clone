package kafka

// KafkaClient is a wrapper around Kafka producer and consumer.
type KafkaClient struct {
	Producer *Producer // Kafka producer
	Consumer *Consumer // Kafka consumer
}

// NewKafkaClient creates a new KafkaClient with the given brokers and group ID.
func NewKafkaClient(brokers []string, groupID string) (*KafkaClient, error) {
	p, err := newProducer(brokers)
	if err != nil {
		return nil, err
	}

	c, err := newConsumer(brokers, groupID)
	if err != nil {
		p.Close()
		return nil, err
	}

	return &KafkaClient{
		Producer: p,
		Consumer: c,
	}, nil
}

func (kc *KafkaClient) Close() {
	if kc.Producer != nil {
		kc.Producer.Close()
	}
	if kc.Consumer != nil {
		kc.Consumer.Close()
	}
}
