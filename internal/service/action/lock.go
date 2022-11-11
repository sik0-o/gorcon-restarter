package action

type Lock struct{}

func NewLock() Lock {
	return Lock{}
}

func (a Lock) Command() string {
	return "#lock"
}
