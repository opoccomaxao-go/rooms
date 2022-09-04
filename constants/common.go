package constants

import "time"

const (
	DefaultAddress = ":22100"

	DefaultTimeout          = time.Second * 10
	DefaultTimeoutReconnect = time.Second * 10

	Version = "1"
)
