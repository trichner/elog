package gsheets

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var scopes = []string{
	"https://www.googleapis.com/auth/drive",
	"https://www.googleapis.com/auth/drive.file",
	"https://www.googleapis.com/auth/drive.readonly",
	"https://www.googleapis.com/auth/spreadsheets",
	"https://www.googleapis.com/auth/spreadsheets.readonly",
}

type SheetValueAppender interface {
	AppendValues(spreadsheetId string, sheetId int64, data [][]string) error
}

type SheetService interface {
	SheetValueAppender
	GetFirstSheet(spreadsheetId string) (*Sheet, error)
}

type sheetsService struct {
	service *sheets.Service
}

type SpreadSheet struct {
	Id string
}

type Sheet struct {
	Id    int64
	Title string
	Index int64
}

func NewDefaultService() (SheetService, error) {

	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		return nil, err
	}

	service, err := sheets.NewService(context.Background(), option.WithCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return &sheetsService{service: service}, nil
}

func NewService(service *sheets.Service) (SheetService, error) {
	return &sheetsService{service: service}, nil
}

func (s *sheetsService) AppendValues(spreadsheetId string, sheetId int64, data [][]string) error {

	sheetName, err := s.getSheetName(spreadsheetId, sheetId)
	if err != nil {
		return fmt.Errorf("unable to append data, spreadsheet='%s' sheetId='%d': %w", spreadsheetId, sheetId, err)
	}

	_range := fmt.Sprintf("'%s'!A:A", sheetName)
	values := toValues(data)
	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	_, err = s.service.Spreadsheets.Values.Append(spreadsheetId, _range, valueRange).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()

	if err != nil {
		return fmt.Errorf("unable to append data, spreadsheet='%s' sheetId='%d': %w", spreadsheetId, sheetId, err)
	}

	return nil
}

func (s *sheetsService) GetFirstSheet(spreadsheetId string) (*Sheet, error) {
	return s.getSheetByIndex(spreadsheetId, 0)
}

func (s *sheetsService) getSheetByIndex(spreadsheetId string, index int64) (*Sheet, error) {

	sheets, err := s.getSheets(spreadsheetId)
	if err != nil {
		return nil, fmt.Errorf("cannot read spreadsheet: %w", err)
	}

	for _, sheet := range sheets {
		if sheet.Index == index {
			return sheet, nil
		}
	}
	return nil, fmt.Errorf("cannot find sheet with index 0")
}

func (s *sheetsService) getSheets(spreadsheetId string) ([]*Sheet, error) {
	resp, err := s.service.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot read spreadsheet: %w", err)
	}

	sheets := make([]*Sheet, 0, len(resp.Sheets))
	for _, sheet := range resp.Sheets {
		sheets = append(sheets, &Sheet{
			Id:    sheet.Properties.SheetId,
			Title: sheet.Properties.Title,
			Index: sheet.Properties.Index,
		})
	}
	return sheets, nil
}

func (s *sheetsService) getSheetName(spreadsheetId string, sheetId int64) (string, error) {
	allSheets, err := s.getSheets(spreadsheetId)
	if err != nil {
		return "", err
	}

	for _, sheet := range allSheets {
		if sheet.Id == sheetId {
			return sheet.Title, nil
		}
	}
	return "", fmt.Errorf("sheet does not exist, spreadsheet='%s' sheetId='%d'", spreadsheetId, sheetId)
}

func toValues(data [][]string) [][]interface{} {
	values := make([][]interface{}, len(data))
	for i, row := range data {
		values[i] = make([]interface{}, len(row))
		for j, cell := range row {
			values[i][j] = cell
		}
	}
	return values
}
