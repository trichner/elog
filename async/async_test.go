package async

import (
	"github.com/stretchr/testify/assert"
	"github.com/trichner/elog"
	"testing"
)

type mockEventLogger struct {
	calls []*elog.Event
}

func (m *mockEventLogger) Log(e *elog.Event) error {
	m.calls = append(m.calls, e)
	return nil
}

func TestAsyncLogger_Publish(t *testing.T) {

	m := &mockEventLogger{}
	a := NewAsyncLogger(m)

	a.Start()

	e := &elog.Event{}
	err := a.Log(e)

	a.Shutdown()

	assert.NoError(t, err)
	assert.Equal(t, e, m.calls[0])
}
