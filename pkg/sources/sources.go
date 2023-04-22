package sources

type EventHandler interface {
	ProcessLine(string)
}
