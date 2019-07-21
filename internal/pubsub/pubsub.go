package pubsub

import (
	"ghasedak-pubsub/pkg"
	"ghasedak-pubsub/pkg/pubsub"
	"sync"
)

var (
	adapterOnce sync.Once
	adapterInst *Adapter
)

type Adapter struct {
	pubsub.PubSub
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func GetAdapter() *Adapter {
	adapterOnce.Do(func() {
		adapterInst = NewAdapter()
	})
	return adapterInst
}

func (p *Adapter) Initialize(kafkaHost string, kafkaPort int32, pulsarHost string, pulsarPort int32) {
	if pkg.GetConfig().PubSub == "kafka" {
		p.PubSub = pubsub.GetKafka().Initialize(kafkaHost, kafkaPort)
	} else if pkg.GetConfig().PubSub == "pulsar" {
		p.PubSub = pubsub.GetPulsar().Initialize(pulsarHost, pulsarPort)
	}
}
