package state

type StateOfAddr int

const (
	READY_READY StateOfAddr = iota
	READY_NOT
	NOT_READY
	NOT_NOT
)
