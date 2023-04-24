package eventhandler

import (
	"time"

	"github.com/ffix/vhtg/pkg/events"
)

type notifier interface {
	AddTask(payload events.Event, expiry time.Time)
}

type logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type stringFinder interface {
	FindStringSubmatch(s string) []string
	SubexpNames() []string
}
