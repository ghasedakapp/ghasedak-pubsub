package cmd

import (
	"ghasedak-pubsub/api/rpc"
	"ghasedak-pubsub/internal/pubsub"
	"ghasedak-pubsub/pkg"
)

func initialize() {
	pkg.GetConfig().Initialize("")
	pkg.GetLogger().Initialize(pkg.GetConfig().Log.Level)
	rpc.NewGrpc().Initialize(":5050")
	pubsub.GetAdapter().Initialize(
		pkg.GetConfig().Kafka.Host,
		pkg.GetConfig().Kafka.Port,
		pkg.GetConfig().Pulsar.Host,
		pkg.GetConfig().Pulsar.Port)
}

func Main() {
	initialize()
	pkg.Wait()
}
