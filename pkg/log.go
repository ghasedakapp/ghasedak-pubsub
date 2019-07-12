package pkg

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func InitLog() {
	Logger = logrus.New()

	Logger.SetFormatter(&logrus.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true,
	})

	Logger.SetOutput(os.Stdout)

	Logger.SetLevel(logrus.DebugLevel)
}
