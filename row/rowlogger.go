package row

import (
	"fmt"
	"github.com/trichner/elog"
	"github.com/trichner/elog/util"
	"time"
)

type RowAppender interface {
	Append(row []string) error
}

type rowEventLogger struct {
	appender RowAppender
}

func (s *rowEventLogger) Log(e *elog.Event) error {

	payload, err := util.MarshalPayloadToJson(e.Payload)
	if err != nil {
		return err
	}

	row := []string{
		e.ID,
		e.SessionID,
		e.Timestamp.Format(time.RFC3339),
		e.Kind,
		payload,
	}

	err = s.appender.Append(row)
	if err != nil {
		return fmt.Errorf("cannot append row for event %+v: %w", e, err)
	}
	return nil
}

func NewRowEventLogger(appender RowAppender) (elog.EventLogger, error) {
	return &rowEventLogger{appender: appender}, nil
}
