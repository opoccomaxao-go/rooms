package proto

const (
	CommandAuthRequired uint16 = iota + 1
	CommandAuth
	CommandAuthSuccess
	CommandRoomCreate
	CommandRoomCreateSuccess
	CommandRoomCreateError
	CommandRoomFinished
	CommandStop
	CommandStats
)
