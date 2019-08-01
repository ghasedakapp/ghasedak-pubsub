package pubsub

import (
	"context"
	"fmt"
	"ghasedak-pubsub/pkg"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
)

type PulsarPubSub struct {
	log       *pkg.Logger
	client    pulsar.Client
	consumers map[string]pulsar.Consumer
	producers map[string]pulsar.Producer
}

type PulsarMessageId struct {
	id []byte
}

func (pm *PulsarMessageId) Serialize() []byte {
	return pm.id
}

func NewPulsar(log *pkg.Logger, host string, port int32) *PulsarPubSub {
	p := &PulsarPubSub{log: log}
	// Instantiate a PubSub client
	client, err := pulsar.NewClient(pulsar.ClientOptions{

		URL: fmt.Sprintf("pulsar://%s:%d", host, port),
	})

	if err != nil {
		log.Fatal(err)
	}
	p.client = client
	p.consumers = make(map[string]pulsar.Consumer)
	p.producers = make(map[string]pulsar.Producer)
	return p
}

func (p *PulsarPubSub) CreateProducer(topic string) error {
	producer, err := p.client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	p.producers[topic] = producer
	return err
}

func (p *PulsarPubSub) Publish(topic string, body []byte) (*MessageId, error) {
	producer, ok := p.producers[topic]
	if !ok {
		return nil, Errors.TopicNotFound
	}
	ctx := context.Background()
	err := producer.Send(ctx, pulsar.ProducerMessage{Payload: body})
	return &MessageId{}, err
}

func (p *PulsarPubSub) Subscribe(name string, topic []string) error {
	if val, ok := p.consumers[name]; ok {
		_ = val.Close()
	}

	consumer, err := p.client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic[0],
		SubscriptionName: name,
		Type:             pulsar.Exclusive,
	})
	p.consumers[name] = consumer
	return err
}

func (p *PulsarPubSub) Receive(subscriptionName string) (*Message, error) {
	consumer, ok := p.consumers[subscriptionName]
	if !ok {
		return nil, Errors.SubscriptionNotFound
	}

	msg, err := consumer.Receive(context.Background())
	if err != nil {
		return nil, err
	}
	return &Message{Value: msg.Payload(), Timestamp: msg.PublishTime()}, nil
}

func (p *PulsarPubSub) Ack(subscriptionName string, mid MessageId) error {
	consumer, ok := p.consumers[subscriptionName]
	if !ok {
		return Errors.SubscriptionNotFound
	}

	id := PulsarMessageId{mid.ID}
	err := consumer.AckID(&id)
	if err != nil {
		return err
	}
	return nil

}

func (p *PulsarPubSub) Close() error {
	panic("implement me")
}
