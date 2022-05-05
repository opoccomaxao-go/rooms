package proto

const (
	CommandMasterAuthRequired uint16 = iota + 1
	CommandMasterAuthSuccess
	CommandMasterRoomCreate
)

const (
	CommandSessionAuth uint16 = iota + 1
	CommandSessionRoomCreated
	CommandSessionRoomError
	CommandSessionRoomFinished
	CommandSessionStats
)
