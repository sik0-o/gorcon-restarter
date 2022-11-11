package action

type Unlock struct{}

func NewUnlock() Lock {
	return Lock{}
}

func (a Unlock) Command() string {
	return "#unlock"
}
