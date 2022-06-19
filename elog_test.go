package elog_test

import (
	"github.com/trichner/elog"
	"github.com/trichner/elog/csv"
	"os"
	"time"
)

func ExampleEventLogger() {

	// create a new logger with the CSV backend
	logger, err := csv.NewCsvEventLogger(os.Stdout)
	if err != nil {
		panic(err)
	}

	event := elog.NewEventBuilder().
		WithKind("user_login").
		WithTimeStamp(time.Now()).
		WithPayload(map[string]string{
			"user":   "Alice",
			"remote": "42.0.4.23",
		}).
		Build()

	logger.Log(event)
}
