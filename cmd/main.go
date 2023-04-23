package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ffix/vhtg/pkg/eventhandler"
	"github.com/ffix/vhtg/pkg/events"
	"github.com/ffix/vhtg/pkg/notifications"
	"github.com/ffix/vhtg/pkg/queue"
	"github.com/ffix/vhtg/pkg/sources"
)

type sendQueue interface {
	AddTask(events.Event, time.Time)
	WaitAndExit()
}

type processor interface {
	Process()
}

const (
	Telegram = iota
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//logger.SetLevel(logrus.InfoLevel)
	logger.SetLevel(logrus.DebugLevel)

	var testMode, dockerMode bool

	pflag.BoolVar(&testMode, "test", false, "Enable test mode (don't send messages to Telegram)")
	pflag.BoolVar(&dockerMode, "docker", false, "Integrate with docker and read data from docker socket instead of stdin")
	pflag.Parse()

	var sendQ sendQueue

	if !testMode {
		telegram, err := notifications.NewTelegramClient(logger)
		if err != nil {
			logger.Fatalf("Failed to initialize Telegram client client: %s", err.Error())
		}

		sendQ = queue.NewTaskQueue(
			func(task *queue.Task) error {
				//time.Sleep(1000 * time.Millisecond)
				//logger.Warn("Calling task with a payload: %v", task.Payload)
				//return fmt.Errorf("aaa")
				err := telegram.SendMessage(task.Payload.Message, task.Payload.Type == events.NewServerSessionStartType)
				if err != nil {
					logger.Warnf("Failed to send Telegram message: %s", err.Error())
				}
				return err
			},
			[]int{Telegram},
			logger,
		)
		defer sendQ.WaitAndExit()
	}

	password := os.Getenv("VALHEIM_SERVER_PASSWORD")
	eventHandler := eventhandler.New(logger, sendQ, password)

	var proc processor
	if dockerMode {
		proc = sources.NewDockerProcessor(eventHandler, "valheim.id", logger)
	} else {
		proc = sources.NewStdinProcessor(eventHandler)
	}
	proc.Process()
}
