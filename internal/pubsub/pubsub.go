package pubsub

import (
	"ghasedak-pubsub/pkg"
	"ghasedak-pubsub/pkg/pubsub"
)

type Adapter struct {
	log *pkg.Logger
	pubsub.PubSub
}

func NewAdapter(log *pkg.Logger, pubsubType string, kafkaHost string, kafkaPort int32, pulsarHost string, pulsarPort int32) *Adapter {
	p := &Adapter{
		log: log,
	}
	if pubsubType == "kafka" {
		p.PubSub = pubsub.NewKafka(log, kafkaHost, kafkaPort)
	} else if pubsubType == "pulsar" {
		p.PubSub = pubsub.NewPulsar(log, pulsarHost, pulsarPort)
	}
	return p
}
