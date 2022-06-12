package bq

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/trichner/elog"
	"github.com/trichner/elog/util"
	"time"
)

type BigQueryEventLogger struct {
	client *bigquery.Client
	table  *bigquery.Table
}

type bigQueryRow struct {
	ID        string
	Kind      string
	SessionID string
	Timestamp time.Time
	Payload   string
}

type BiqQueryEventLoggerOptions struct {
	BilledProject  string
	DatasetProject string
	DatasetID      string
	TableID        string
}

func NewBiqQueryEventLogger(options *BiqQueryEventLoggerOptions) (elog.EventLogger, error) {

	var err error

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, options.BilledProject)
	if err != nil {
		return nil, fmt.Errorf("failed to create biquery client: %w", err)
	}

	ds := client.DatasetInProject(options.DatasetProject, options.DatasetID)
	if _, err := ds.Metadata(ctx); err != nil {
		return nil, fmt.Errorf("cannot find dataset %q in project %q: %w", options.DatasetID, options.DatasetProject, err)
	}

	tb := ds.Table(options.TableID)
	_, err = tb.Metadata(ctx)
	if err != nil {
		err := tb.Create(ctx, &bigquery.TableMetadata{
			Schema: []*bigquery.FieldSchema{
				{
					Name:     "ID",
					Required: true,
					Type:     bigquery.StringFieldType,
				},
				{
					Name:     "Kind",
					Required: true,
					Type:     bigquery.StringFieldType,
				},
				{
					Name:     "SessionID",
					Required: true,
					Type:     bigquery.StringFieldType,
				},
				{
					Name:     "Timestamp",
					Required: true,
					Type:     bigquery.TimestampFieldType,
				},
				{
					Name:     "Payload",
					Required: false,
					Type:     bigquery.StringFieldType,
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("table %q does not exist and cannot be created: %w", options.TableID, err)
		}
	}
	return &BigQueryEventLogger{client: client, table: tb}, nil
}

func (b *BigQueryEventLogger) Close() error {
	return b.client.Close()
}

func (b *BigQueryEventLogger) Log(e *elog.Event) error {

	ctx := context.Background()

	payload, err := util.MarshalPayloadToJson(e.Payload)
	if err != nil {
		return err
	}

	inserter := b.table.Inserter()
	err = inserter.Put(ctx, []bigQueryRow{{
		ID:        e.ID,
		Kind:      e.Kind,
		SessionID: e.SessionID,
		Timestamp: e.Timestamp,
		Payload:   payload,
	}})
	if err != nil {
		return fmt.Errorf("cannot insert row into bigquery: %w", err)
	}
	return nil
}
