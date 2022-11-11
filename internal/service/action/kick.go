package action

import (
	"fmt"
	"strings"
)

type Kick struct {
	playerId int
	message  string
}

func NewKick(playerId int, message ...string) Kick {
	return Kick{
		playerId: playerId,
		message:  strings.Join(message, "\n"),
	}
}

func (a Kick) Command() string {
	var cmd string
	if a.playerId < 0 {
		// kick all
		cmd = "#kick -1"
	} else {
		cmd = fmt.Sprintf("kick %d %s", a.playerId, a.message)
		cmd = strings.TrimRight(cmd, " ")
	}

	return cmd
}
