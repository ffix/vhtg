package eventhandler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type LogEventHandler func(map[string]string) string
type LogEvent struct {
	Pattern string
	Handler LogEventHandler
}

type Character struct {
	Name   string
	Online bool
}

type EventHandler struct {
	characters map[int]*Character
	logger     logger
	notifier   notifier
}

func New(logger logger, notifier notifier) *EventHandler {
	e := EventHandler{logger: logger, notifier: notifier}
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

func (e *EventHandler) currentlyOnlineString(exclude string) string {
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
		return " They are the only player online."
	}
	return fmt.Sprintf(" Currently online players: %s.", joined)

}

func (e *EventHandler) dungeonDBStartHandler(_ map[string]string) string {
	e.initOnlinePlayers()
	return "The server has started."

}

func (e *EventHandler) characterConnected(matches map[string]string) string {
	characterID, err := parseStringIDToInt(matches["id"])
	if err != nil {
		return ""
	}

	if characterID == 0 {
		// A Player has died
		return ""
	}

	character := e.characters[characterID]
	if character != nil && character.Online == true {
		return ""
	}

	character = &Character{
		Name:   matches["name"],
		Online: true,
	}

	e.characters[characterID] = character

	return fmt.Sprintf(
		"A player named %s has entered the server.%s\n",
		matches["name"],
		e.currentlyOnlineString(matches["name"]),
	)

}

//func (e *EventHandler) charactedDied(matches map[string]string) string {
//	// don't do anything if character has died
//	return ""
//}

//func (e *EventHandler) gotHandshakeHandler(matches map[string]string) string {
//	// fixme
//	e.logger.Warn("Not handled")
//	return ""
//}

//func (e *EventHandler) peerWrongPasswordHandler(matches map[string]string) {
//	fmt.Printf("Wrong password from peer %s.\n", matches["peer_id"])
//
//}

//func (e *EventHandler) closingSocketHandler(matches map[string]string) string {
//	// fixme
//	e.logger.Warn("Not handled")
//	return ""
//}

func (e *EventHandler) serverShutDownComplete(_ map[string]string) string {
	return "The server has stopped."
}

func (e *EventHandler) newSessionHandler(matches map[string]string) string {
	return fmt.Sprintf(
		"A session with ID %s and join code %s has started.\n",
		matches["session"],
		matches["join_code"],
	)

}
func (e *EventHandler) playerDisconnected(matches map[string]string) string {
	//fmt.Println(matches)
	characterID, err := parseStringIDToInt(matches["id"])
	if err != nil {
		return ""
	}
	character, ok := e.characters[characterID]
	if !ok {
		return ""
	}

	if !character.Online {
		return ""
	}
	character.Online = false

	return fmt.Sprintf(
		"The player %s disconnected.%s\n",
		character.Name,
		e.currentlyOnlineString(character.Name),
	)

}

func (e *EventHandler) ProcessLine(line string) {
	logEvents := []LogEvent{
		{
			Pattern: `DungeonDB Start`,
			Handler: e.dungeonDBStartHandler,
		},
		//{
		//	Pattern: `Got character ZDOID from (?P<name>.+) : (0:0)`,
		//	Handler: e.charactedDied,
		//},
		{
			Pattern: `Got character ZDOID from (?P<name>[\w.-]+) : (?P<id>-{0,1}\d+):`,
			Handler: e.characterConnected,
		},
		//{
		//	Pattern: `Got handshake from client (?P<client_id>\d+)`,
		//	Handler: e.gotHandshakeHandler,
		//},
		//{
		//	Pattern: `Peer (?P<peer_id>\d+)( has wrong password)`,
		//	Handler: e.peerWrongPasswordHandler,
		//},
		//{
		//	Pattern: `Closing socket (?P<socket_id>[0-9]+)`,
		//	Handler: e.closingSocketHandler,
		//},
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

			message := logEvent.Handler(namedMatches)

			if message != "" {
				e.logger.Info(message)
				if e.notifier != nil {
					err := e.notifier.SendMessage(message)
					if err != nil {
						e.logger.Warnf("Failed to send a message: %w", err)
					}
				}
			}
			break
		}
	}
}

func parseStringIDToInt(id string) (int, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse string ID to integer: %w", err)
	}
	return intID, nil
}
