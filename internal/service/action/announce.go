package action

type Announce struct {
	message string
}

func NewAnnounce(message string) Announce {
	return Announce{
		message: message,
	}
}

func (a Announce) Command() string {
	// Use Say action
	return NewSay(a.message, -1).Command()
}
