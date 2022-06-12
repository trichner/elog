# Event Logger

This package provides a simple interface for event logging to
generic backends as well as a solid set of backend implementations.

### Available Storage Backends

- `csv`: logging events in CSV format, easy for development and log collection
- `bq`: logging to Google BigQuery, great for production
- `gsheets`: logging to a Google Sheet, great for low-volume production

### Backend Decorators

- `async`: a decorator to make logging asynchronous
- `row`: a decorator to make logging to row based backends easier

## Installation

```bash
go install github.com/trichner/elog
```

## Usage

```go
package main

func main() {
	// create a new logger with the CSV backend
	logger, err := NewCsvEventLogger(os.Stdout)

	if err != nil {
		panic(err)
	}

	// create a new event
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

	// ...and log it
	err = logger.Log(event)
	if err != nil {
		panic(err)
	}
}
```
