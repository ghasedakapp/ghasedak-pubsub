package pubsub

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

type KafkaPubSub struct {
	host      string
	consumers map[string]*kafka.Consumer
	producers map[string]*kafka.Producer
}

func NewKafka() *KafkaPubSub {
	return &KafkaPubSub{}
}

var (
	kafkaOnce sync.Once

	kafkaPubSub *KafkaPubSub
)

func GetKafka() *KafkaPubSub {
	kafkaOnce.Do(func() {
		kafkaPubSub = NewKafka()
	})

	return kafkaPubSub
}

//
//var pubSub *PubSub
//
///**
//Get singletone pubsub object
//*/
//func GetPubSub() *PubSub {
//	if pubSub == nil {
//		p := NewPulsar(Conf.Pulsar.Host, Conf.Pulsar.Port)
//		pubSub = &p
//		return pubSub
//	} else {
//		return pubSub
//	}
//}

func (p *KafkaPubSub) Initialize(host string, port int32) *KafkaPubSub {
	p.host = host
	p.consumers = make(map[string]*kafka.Consumer)
	p.producers = make(map[string]*kafka.Producer)
	return p
}

func (p *KafkaPubSub) CreateProducer(topic string) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": p.host,
	})
	if err != nil {
		return err
	}

	p.producers[topic] = producer

	return nil
}

func (p *KafkaPubSub) Publish(topic string, body []byte) (*MessageId, error) {
	producer, ok := p.producers[topic]
	if !ok {
		return nil, Errors.TopicNotFound
	}

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          body,
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return nil, m.TopicPartition.Error
	}

	close(deliveryChan)

	return &MessageId{[]byte(m.TopicPartition.Offset.String())}, err
}

func (p *KafkaPubSub) Subscribe(name string, topics []string) error {
	if val, ok := p.consumers[name]; ok {
		_ = val.Close()
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": p.host,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return err
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return err
	}

	p.consumers[name] = consumer
	return err
}

func (p *KafkaPubSub) Receive(subscriptionName string) (*Message, error) {
	consumer, ok := p.consumers[subscriptionName]
	if !ok {
		return nil, Errors.SubscriptionNotFound
	}

	msg, err := consumer.ReadMessage(-1)
	if err != nil {
		return nil, err
	}
	return &Message{Value: msg.Value, Timestamp: msg.Timestamp}, nil
}

func (p *KafkaPubSub) Ack(subscriptionName string, mid MessageId) error {
	//consumer, ok := p.consumers[subscriptionName]
	//if !ok {
	//	return Errors.SubscriptionNotFound
	//}

	//id := PulsarMessageId{mid.ID}
	//err := consumer.Commit()
	//if err != nil {
	//	return err
	//}
	return nil

}

func (p *KafkaPubSub) Close() error {
	panic("implement me")
}
