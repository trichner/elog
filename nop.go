package elog

type noOpEventLogger struct {
}

func (n *noOpEventLogger) Log(e *Event) error {
	return nil
}

func NewNoOpEventLogger() EventLogger {
	return &noOpEventLogger{}
}
