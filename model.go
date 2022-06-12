package elog

import (
	"time"
)

type Event struct {
	ID        string
	Kind      string
	SessionID string
	Timestamp time.Time
	Payload   any
}

type EventLogger interface {
	Log(e *Event) error
}
