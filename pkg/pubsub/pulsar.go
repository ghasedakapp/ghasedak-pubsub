package pubsub

import (
	"context"
	"fmt"
	"ghasedak-pubsub/pkg"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
)

type PulsarPubSub struct {
	client    pulsar.Client
	consumers map[string]pulsar.Consumer
}

type PulsarMessageId struct {
	id []byte
}

func (pm *PulsarMessageId) Serialize() []byte {
	return pm.id
}

func NewPulsar(host string, port int32) PubSub {
	// Instantiate a PubSub client
	client, err := pulsar.NewClient(pulsar.ClientOptions{

		URL: fmt.Sprintf("pulsar://%s:%d", host, port),
	})

	if err != nil {
		pkg.Logger.Fatal(err)
	}

	return &PulsarPubSub{client: client}
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

func (p *PulsarPubSub) Subscribe(name string, topic string) error {
	if val, ok := p.consumers[name]; ok {
		_ = val.Close()
	}

	consumer, err := p.client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: name,
		Type:             pulsar.Exclusive,
	})
	p.consumers[name] = consumer
	return err
}

func (p *PulsarPubSub) Receive(subscriptionName string) (*Message, error) {
	consumer, ok := p.consumers[subscriptionName]
	if !ok {
		return nil, SubscriptionNotFound
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
		return SubscriptionNotFound
	}

	id := PulsarMessageId{mid.id}
	err := consumer.AckID(&id)
	if err != nil {
		return err
	}
	return nil

}
