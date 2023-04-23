package eventhandler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ffix/vhtg/pkg/events"
)

type LogEventHandler func(map[string]string) events.Event
type LogEvent struct {
	Pattern string
	Handler LogEventHandler
}

type Character struct {
	Name   string
	Online bool
}

type EventHandler struct {
	up         bool
	characters map[int]*Character
	logger     logger
	notifier   notifier
	password   string
}

var gameEvents = map[string]string{
	"army_eikthyr":  "Eikthyr rallies the creatures of the forest",
	"army_theelder": "The forest is moving...",
	"army_bonemass": "A foul smell from the swamp",
	"army_moder":    "A cold wind blows from the mountains",
	"army_goblin":   "The horde is attacking",
	"foresttrolls":  "The ground is shaking",
	"blobs":         "A foul smell from the swamp",
	"skeletons":     "Skeleton Surprise",
	"surtlings":     "There's a smell of sulfur in the air",
	"wolves":        "You are being hunted",
	"bats":          "You stirred the cauldron",
	"army_gjall":    "What's up, Gjall!?",
	"army_seekers":  "They sought you out",
}

func New(logger logger, notifier notifier, password string) *EventHandler {
	e := EventHandler{logger: logger, notifier: notifier, password: password}
	e.initOnlinePlayers()
	return &e
}

func (e *EventHandler) initOnlinePlayers() {
	e.characters = make(map[int]*Character)
}

func (e *EventHandler) currentlyOnline() []string {
	var online []string
	for _, val := range e.characters {
		if val.Online {
			online = append(online, val.Name)
		}
	}
	return online
}

func (e *EventHandler) currentlyOnlineString(exclude string, empty string) string {
	online := e.currentlyOnline()
	if exclude != "" {
		var tmp []string
		for _, name := range online {
			if name == exclude {
				continue
			}
			tmp = append(tmp, name)
		}
		online = tmp
	}

	joined := strings.Join(online, ", ")
	if joined == "" {
		return fmt.Sprintf(" %s", empty)
	}
	return fmt.Sprintf(" Currently online players: %s.", joined)

}

func (e *EventHandler) dungeonDBStartHandler(_ map[string]string) events.Event {
	e.initOnlinePlayers()
	if e.up {
		return events.NoEvent()
	}
	e.up = true
	return events.ServerStartEvent("The server has started.")

}

func (e *EventHandler) characterConnected(matches map[string]string) events.Event {
	characterID, err := parseStringIDToInt(matches["id"])
	if err != nil {
		return events.NoEvent()
	}

	if characterID == 0 {
		// A Player has died
		return events.NoEvent()
	}

	character := e.characters[characterID]
	if character != nil && character.Online {
		return events.NoEvent()
	}

	character = &Character{
		Name:   matches["name"],
		Online: true,
	}

	e.characters[characterID] = character

	return events.PlayerLoggedInEvent(
		fmt.Sprintf(
			"A player named %s has entered the server.%s",
			matches["name"],
			e.currentlyOnlineString(matches["name"], "They are the only player online."),
		),
	)

}

func (e *EventHandler) serverShutDownComplete(_ map[string]string) events.Event {
	if !e.up {
		return events.NoEvent()
	}
	e.up = false
	return events.ServerStopEvent("The server has stopped.")
}

func (e *EventHandler) newSessionHandler(matches map[string]string) events.Event {
	var password string
	if e.password != "" {
		password = fmt.Sprintf(" Connect with password: %s.", e.password)
	}
	return events.NewServerSessionStartEvent(fmt.Sprintf(
		"Session ID %s, join code %s active.%s",
		matches["session"],
		matches["join_code"],
		password,
	))

}
func (e *EventHandler) playerDisconnected(matches map[string]string) events.Event {
	//fmt.Println(matches)
	characterID, err := parseStringIDToInt(matches["id"])
	if err != nil {
		return events.NoEvent()
	}
	character, ok := e.characters[characterID]
	if !ok {
		return events.NoEvent()
	}

	if !character.Online {
		return events.NoEvent()
	}
	character.Online = false

	return events.PlayerLoggedOutEvent(fmt.Sprintf(
		"The player %s disconnected.%s",
		character.Name,
		e.currentlyOnlineString(character.Name, "They were the only player online."),
	))

}

func (e *EventHandler) randomEvent(matches map[string]string) events.Event {
	gameEventDescription, ok := gameEvents[matches["event"]]
	if ok {
		return events.RandomEventEvent(gameEventDescription)
	}
	return events.RandomEventEvent("Unknown game event started...")
}

// func (e *EventHandler) ProcessLine(line string, eventTime *time.Time) {
func (e *EventHandler) ProcessLine(line string, eventTime *time.Time) {
	logEvents := []LogEvent{
		{
			Pattern: `DungeonDB Start`,
			Handler: e.dungeonDBStartHandler,
		},
		{
			Pattern: `Got character ZDOID from (?P<name>[\w.-]+) : (?P<id>-{0,1}\d+):`,
			Handler: e.characterConnected,
		},

		{
			Pattern: `Shutdown complete`,
			Handler: e.serverShutDownComplete,
		},
		{
			Pattern: `Session "(?P<session>.+)" with join code (?P<join_code>\d+) and IP (?P<ip>[\d.]+:\d+)`,
			Handler: e.newSessionHandler,
		},
		{
			Pattern: `Destroying abandoned non persistent zdo -{0,1}\d+:\d+ owner (?P<id>-{0,1}\d+)`,
			Handler: e.playerDisconnected,
		},
		{
			Pattern: `Random event set:(?P<event>\w+)`,
			Handler: e.randomEvent,
		},
	}

	for _, logEvent := range logEvents {
		compiledRegex, err := regexp.Compile(logEvent.Pattern)
		if err != nil {
			log.Fatalf("Error compiling regex pattern: %s", err)
		}

		matches := compiledRegex.FindStringSubmatch(line)
		if matches != nil {
			namedMatches := make(map[string]string)
			for i, name := range compiledRegex.SubexpNames() {
				if i != 0 && name != "" {
					namedMatches[name] = matches[i]
				}
			}

			e.logger.Debugf("Matched: %s", line)

			event := logEvent.Handler(namedMatches)

			if event.Type != events.NoEventType {
				e.logger.Info(event.Message)
				if e.notifier != nil {
					expiry := time.Now().Add(5 * time.Minute)
					if eventTime != nil {
						expiry = eventTime.Add(5 * time.Minute)
					}
					e.notifier.AddTask(event, expiry)
				}
			}
			break
		}
	}
}

func parseStringIDToInt(id string) (int, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string ID to integer: %w", err)
	}
	return intID, nil
}
