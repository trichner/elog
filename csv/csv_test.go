package csv

import (
	"github.com/google/uuid"
	"github.com/trichner/elog"
	"os"
	"time"
)

func ExampleNewCsvEventLogger() {

	logger, err := NewCsvEventLogger(os.Stdout)

	if err != nil {
		panic(err)
	}

	event := &elog.Event{
		ID:        uuid.New().String(),
		Kind:      "user_login",
		SessionID: uuid.New().String(),
		Timestamp: time.Now(),
		Payload: map[string]string{
			"user":   "Alice",
			"remote": "42.0.0.42",
		},
	}

	err = logger.Log(event)
	if err != nil {
		panic(err)
	}
}
