package test

import (
	"fmt"
	"ghasedak-pubsub/pkg"
	"ghasedak-pubsub/pkg/pubsub"
	"os"
	"testing"
)

func setup() {
	configPath := fmt.Sprintf("%s/test/config/test.yaml", os.Getenv("PWD"))
	fmt.Println("Config path is ", configPath)
	pkg.InitConfig(configPath)
	pubsub.New()
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	teardown()
	os.Exit(r)
}
