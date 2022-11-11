package action

import "fmt"

type Say struct {
	playerId int
	message  string
}

func NewSay(message string, playerId int) Say {
	return Say{
		playerId: playerId,
		message:  message,
	}
}

func (a Say) Command() string {
	return fmt.Sprintf("say %d %s", a.playerId, a.message)
}
