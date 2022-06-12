package gsheets

import (
	"fmt"
	"github.com/trichner/elog"
	"github.com/trichner/elog/row"
)

type sheetAppender struct {
	sheetValueAppender SheetValueAppender
	spreadsheetId      string
	sheetId            int64
}

func NewEventLogger(ss SheetValueAppender, spreadsheetId string, sheetId int64) (elog.EventLogger, error) {
	appender := &sheetAppender{
		sheetValueAppender: ss,
		spreadsheetId:      spreadsheetId,
		sheetId:            sheetId,
	}
	return row.NewRowEventLogger(appender)
}

func NewLoggerWithFirstSheet(ss SheetService, spreadsheetId string) (elog.EventLogger, error) {
	sheet, err := ss.GetFirstSheet(spreadsheetId)
	if err != nil {
		return nil, fmt.Errorf("cannot find first sheet of spreadsheet: %w", err)
	}
	sheetId := sheet.Id
	appender := &sheetAppender{
		sheetValueAppender: ss,
		spreadsheetId:      spreadsheetId,
		sheetId:            sheetId,
	}
	return row.NewRowEventLogger(appender)
}

func (s *sheetAppender) Append(row []string) error {

	data := [][]string{row}
	err := s.sheetValueAppender.AppendValues(s.spreadsheetId, s.sheetId, data)
	if err != nil {
		return fmt.Errorf("cannot google sheet row: %w", err)
	}
	return nil
}
