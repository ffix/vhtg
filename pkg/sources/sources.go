package sources

import (
	"time"
)

type EventHandler interface {
	ProcessLine(string, *time.Time)
}
