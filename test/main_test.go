package test

import (
	"fmt"
	pb "ghasedak-pubsub/api/proto/src"
	"ghasedak-pubsub/api/rpc"
	pubsub2 "ghasedak-pubsub/internal/pubsub"
	"ghasedak-pubsub/pkg"
	"ghasedak-pubsub/pkg/pubsub"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"os"
	"testing"
	time "time"
)

var PubClient pb.PublisherClient
var SubClient pb.SubscriberClient
var PulsarPS pubsub.PubSub
var Log *pkg.Logger

func setup() {
	rand.Seed(time.Now().Unix())
	conf := initConfig()
	Log = pkg.NewLog(conf.Log.Level)
	pubsubAdapter := pubsub2.NewAdapter(
		Log,
		conf.PubSubType,
		conf.Kafka.Host,
		conf.Kafka.Port,
		conf.Pulsar.Host,
		conf.Pulsar.Port)
	rpc.NewGrpc(Log, pubsubAdapter, ":5050")
	time.Sleep(500 * time.Millisecond)
	initGrpcClient(":5050")
}

func teardown() {

}

func TestMain(m *testing.M) {
	fmt.Println("Initializing integration test...")
	setup()
	r := m.Run()
	teardown()
	os.Exit(r)
}

func initConfig() *pkg.Config {
	configPath := fmt.Sprintf("%s/config/test.yaml", os.Getenv("PWD"))
	fmt.Println("Config path is ", configPath)
	return pkg.NewConfig(configPath)
}

func initGrpcClient(address string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	PubClient = pb.NewPublisherClient(conn)
	SubClient = pb.NewSubscriberClient(conn)

}
