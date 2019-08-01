package pkg

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLog(level string) *Logger {
	l := &Logger{}
	l.Logger = logrus.New()

	l.Logger.SetFormatter(&logrus.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true,
	})

	l.Logger.SetOutput(os.Stdout)

	lev, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	l.Logger.SetLevel(lev)

	return l
}
