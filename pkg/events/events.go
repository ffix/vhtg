package events

type EventType int

const (
	ServerStartType EventType = iota
	ServerStopType
	PlayerLoggedInType
	PlayerLoggedOutType
	NewServerSessionStartType
	RandomEventType
	NoEventType
)

type Event struct {
	Type    EventType
	Message string
}

func NewEvent(eventType EventType, msg string) Event {
	return Event{Type: eventType, Message: msg}

}

func NoEvent() Event {
	return NewEvent(NoEventType, "")
}

func ServerStartEvent(msg string) Event {
	return NewEvent(ServerStartType, msg)
}

func ServerStopEvent(msg string) Event {
	return NewEvent(ServerStopType, msg)
}

func PlayerLoggedInEvent(msg string) Event {
	//msg := fmt.Sprintf("%s logged in", playerName)
	return NewEvent(PlayerLoggedInType, msg)
}

func PlayerLoggedOutEvent(msg string) Event {
	//msg := fmt.Sprintf("%s logged out", playerName)
	return NewEvent(PlayerLoggedOutType, msg)
}

func NewServerSessionStartEvent(msg string) Event {
	return NewEvent(NewServerSessionStartType, msg)
}

func RandomEventEvent(msg string) Event {
	return NewEvent(RandomEventType, msg)
}
