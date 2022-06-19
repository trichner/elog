package elog_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trichner/elog"
	"github.com/trichner/elog/csv"
	"strings"
	"testing"
	"time"
)

func TestEventLogger(t *testing.T) {

	var writer strings.Builder

	logger, err := csv.NewCsvEventLogger(&writer)
	if err != nil {
		panic(err)
	}

	id := "17"
	session := "my_session"

	ts, _ := time.Parse(time.RFC3339, "2022-08-01T23:22:14.000Z")

	event := elog.NewEventBuilder().
		WithID(id).
		WithTimeStamp(ts).
		WithSessionID(session).
		WithKind("user_login").
		WithPayload(map[string]string{
			"user":   "Alice",
			"remote": "42.0.4.23",
		}).
		Build()

	logger.Log(event)

	actual := writer.String()

	expected := `id,session,timestamp,kind,payload
17,my_session,2022-08-01T23:22:14Z,user_login,"{""remote"":""42.0.4.23"",""user"":""Alice""}"
`

	assert.Equal(t, expected, actual)
}
