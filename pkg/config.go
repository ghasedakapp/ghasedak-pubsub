package pkg

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strings"
)

var Conf ConfYaml

var defaultConf = []byte(`
grpc:
  address: 127.0.0.1:5050

postgres:
  host: 127.0.0.1
  port: 5432

pulsar:
  host:127.0.0.1
  port:6650

kafka:
  host:127.0.0.1
  port:6650

log:
  level: debug

`)

// ConfYaml is config structure.
type ConfYaml struct {
	GRPC     SectionGRPC     `yaml:"grpc"`
	Postgres SectionPostgres `yaml:"postgres"`
	Pulsar   SectionPulsar   `yaml:"pulsar"`
	Kafka    SectionKafka    `yaml:"kafka"`
	Log      SectionLog      `yaml:"log"`
}

// SectionGRPC is sub section of config.
type SectionGRPC struct {
	Address string `yaml:"address"`
}

// SectionPostgres is sub section of config.
type SectionPostgres struct {
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
}

// SectionPulsar is sub section of config.
type SectionPulsar struct {
	Host string `yaml:"host"`
	Port int32  `yaml:"host"`
}

// SectionKafka is sub section of config.
type SectionKafka struct {
	Host string `yaml:"host"`
	Port int32  `yaml:"host"`
}

// SectionLog is sub section of config.
type SectionLog struct {
	Level string `yaml:"level"`
}

// LoadConf load config from file and read in environment variables that match
func loadConf(confPath string) (ConfYaml, error) {
	var conf ConfYaml

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()         // read in environment variables that match
	viper.SetEnvPrefix("pubsub") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		// Search config in home directory with name ".gorush" (without extension).
		viper.AddConfigPath("/etc/pubsub/")
		viper.AddConfigPath("$HOME/.pubsub")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return conf, err
			}
		}
	}

	// Grpc
	conf.GRPC.Address = viper.GetString("grpc.address")

	// Postgres
	conf.Postgres.Host = viper.GetString("grpc.host")
	conf.Postgres.Port = viper.GetInt32("grpc.port")

	// Pulsar
	conf.Pulsar.Host = viper.GetString("pulsar.host")
	conf.Pulsar.Port = viper.GetInt32("pulsar.port")

	// Kafka
	conf.Kafka.Host = viper.GetString("kafka.host")
	conf.Kafka.Port = viper.GetInt32("kafka.port")

	// log
	conf.Log.Level = viper.GetString("log.level")

	return conf, nil
}

func InitConfig(path string) {
	var err error
	Conf, err = loadConf(path)
	if err != nil {
		log.Fatal(err)
	}
}
