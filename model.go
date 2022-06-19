package elog

import (
	"github.com/google/uuid"
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

// EventBuilder helps to assemble events.
// The builder will create a universally unique ID (UUID) for the event ID if
// none is provided. Furthermore, the Timestamp will be set to 'now' in case it was
// not provided.
type EventBuilder interface {
	WithID(id string) EventBuilder
	WithKind(kind string) EventBuilder
	WithSessionID(session string) EventBuilder
	WithTimeStamp(ts time.Time) EventBuilder
	WithPayload(payload any) EventBuilder

	Build() *Event
}

func NewEventBuilder() EventBuilder {
	return &eventBuilder{e: &Event{}}
}

type eventBuilder struct {
	e *Event
}

func (e *eventBuilder) WithID(id string) EventBuilder {
	e.e.ID = id
	return e
}

func (e *eventBuilder) WithKind(kind string) EventBuilder {
	e.e.Kind = kind
	return e
}

func (e *eventBuilder) WithSessionID(session string) EventBuilder {
	e.e.SessionID = session
	return e
}

func (e *eventBuilder) WithTimeStamp(ts time.Time) EventBuilder {
	e.e.Timestamp = ts
	return e
}

func (e *eventBuilder) WithPayload(payload any) EventBuilder {
	e.e.Payload = payload
	return e
}

func (e *eventBuilder) Build() *Event {
	if e.e.ID == "" {
		e.e.ID = uuid.New().String()
	}
	if e.e.Timestamp.IsZero() {
		e.e.Timestamp = time.Now()
	}
	return e.e
}
