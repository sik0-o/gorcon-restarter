package action

type Shutdown struct{}

func NewShutdown() Shutdown {
	return Shutdown{}
}

func (a Shutdown) Command() string {
	return "#shutdown"
}
