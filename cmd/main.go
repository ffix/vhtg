package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ffix/vhtg/pkg/eventhandler"
	"github.com/ffix/vhtg/pkg/notifications"
	"github.com/ffix/vhtg/pkg/sources"
)

type messageSender interface {
	SendMessage(msg string) error
}

type processor interface {
	Process()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(logrus.InfoLevel)
	//logger.SetLevel(logrus.DebugLevel)

	var testMode, dockerMode bool

	pflag.BoolVar(&testMode, "test", false, "Enable test mode (don't send messages to Telegram)")
	pflag.BoolVar(&dockerMode, "docker", false, "Integrate with docker and read data from docker socket instead of stdin")
	pflag.Parse()

	var sender messageSender
	var err error

	if !testMode {
		sender, err = notifications.NewTelegramClient()
		if err != nil {
			log.Fatalf("Failed to initialize messageSender client: %s", err)
		}
	}

	eventHandler := eventhandler.New(logger, sender)

	var proc processor
	if dockerMode {
		proc = sources.NewDockerProcessor(eventHandler)
	} else {
		proc = sources.NewStdinProcessor(eventHandler)
	}
	proc.Process()
}
