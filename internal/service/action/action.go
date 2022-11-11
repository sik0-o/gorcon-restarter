package action

type Action interface {
	Command() string // Command returns RCON command of action
}
