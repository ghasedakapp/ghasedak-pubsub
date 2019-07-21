package pkg

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
)

var (
	logOnce sync.Once
	logInst *Log
)

type Log struct {
	*logrus.Logger
}

func NewLog() *Log {
	return &Log{}
}

func GetLogger() *Log {
	logOnce.Do(func() {
		logInst = NewLog()
	})
	return logInst
}

func (l *Log) Initialize(level string) {
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
}
