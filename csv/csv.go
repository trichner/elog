package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/trichner/elog"
	"io"
	"time"
)

type csvEventLogger struct {
	w *csv.Writer
}

func (c *csvEventLogger) Log(e *elog.Event) error {

	payload, err := marshalPayload(e.Payload)
	if err != nil {
		return err
	}

	err = c.w.Write([]string{
		e.ID,
		e.SessionID,
		e.Timestamp.Format(time.RFC3339),
		e.Kind,
		payload,
	})
	if err != nil {
		return fmt.Errorf("cannot write csv row for event %+v: %w", e, err)
	}
	c.w.Flush()
	return nil
}

func NewCsvEventLogger(w io.Writer) (elog.EventLogger, error) {
	csvWriter := csv.NewWriter(w)
	err := csvWriter.Write([]string{
		"id", "session", "timestamp", "kind", "payload",
	})
	if err != nil {
		return nil, err
	}

	csvWriter.Flush()
	return &csvEventLogger{w: csvWriter}, nil
}

func marshalPayload(payload any) (string, error) {

	if payload == nil {
		return "", nil
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("cannot serialize event payload: %w", err)
	}
	return string(payloadBytes), nil
}
