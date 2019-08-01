package cmd

import (
	"ghasedak-pubsub/api/rpc"
	"ghasedak-pubsub/internal/pubsub"
	"ghasedak-pubsub/pkg"
)

func initialize() {
	conf := pkg.NewConfig("")
	logger := pkg.NewLog(conf.Log.Level)
	pubsubAdapter := pubsub.NewAdapter(
		logger,
		conf.PubSubType,
		conf.Kafka.Host,
		conf.Kafka.Port,
		conf.Pulsar.Host,
		conf.Pulsar.Port)
	rpc.NewGrpc(logger, pubsubAdapter, ":5050")
}

func Main() {
	initialize()
	pkg.Wait()
}
