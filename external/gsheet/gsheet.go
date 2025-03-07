package gsheet

import (
	"errors"
	"net/http"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

var delays = []int{10, 30, 90}

func GetContent(srv *sheets.Service, url, sheetRange string) ([][]interface{}, error) {
	info, err := FromURL(url)
	if err != nil {
		return nil, err
	}

	content, err := GetSheetContent(srv, info.SheetID, sheetRange)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func GetSheetContent(srv *sheets.Service, sheetID, sheetRange string) ([][]interface{}, error) {
	batchGetCall := srv.Spreadsheets.Values.BatchGet(sheetID)
	if sheetRange != "" {
		batchGetCall = batchGetCall.Ranges(sheetRange)
	}

	for _, d := range delays {
		resp, err := batchGetCall.Do()
		if err != nil {
			var batchErr *googleapi.Error
			if errors.As(err, &batchErr) && batchErr.Code == http.StatusTooManyRequests {
				// When too many request, sleep delay sec and try again once
				time.Sleep(time.Second * time.Duration(d))
				continue
			}

			return nil, err
		}

		if len(resp.ValueRanges) == 0 {
			return nil, errors.New("no sheets found in the spreadsheet")
		}

		return resp.ValueRanges[0].Values, nil
	}

	return nil, errors.New("failed all the retry attempts")
}

func GetSheetName(srv *sheets.Service, sheetURL string) (string, error) {
	sheetInfo, err := FromURL(sheetURL)
	if err != nil {
		return "", err
	}
	spreadsheet, err := srv.Spreadsheets.Get(sheetInfo.SheetID).Do()
	if err != nil {
		return "", err
	}

	if len(spreadsheet.Sheets) == 0 {
		return "", errors.New("no sub sheet found")
	}

	for _, s := range spreadsheet.Sheets {
		if s.Properties.SheetId == sheetInfo.GID {
			return s.Properties.Title, nil
		}
	}
	sid := spreadsheet.Sheets[0].Properties.Title
	return sid, err
}
