package export

// code from google

import (
	"fmt"
	"log"
	"net/http"

	"github.com/athletifit/social-network-insights/models"
	"github.com/athletifit/social-network-insights/sheet"
	sheets "google.golang.org/api/sheets/v4"
)

// SheetExporter represents a google sheet exporter.
type SheetExporter struct {
	SheetClient *http.Client
}

// NewSheetExporter returns a google sheet exporter.
func NewSheetExporter() Exporter {
	sc := sheet.GetSheetClient()
	return SheetExporter{
		SheetClient: sc,
	}
}

// Export is our main method. Writes the data to export.
func (se SheetExporter) Export(document Document) {
	srv, err := sheets.New(se.SheetClient)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return
	}

	var s *sheets.Spreadsheet
	create := true
	if create {
		s, err = se.createSheet(srv, document)
		if err != nil {
			log.Fatalf("Unable to write Sheet: %v", err)
			return
		}
	}

	// Prints the url of the exported doc.
	fmt.Println("Available at: " + s.SpreadsheetUrl)

}

func (se SheetExporter) createSheet(srv *sheets.Service, document Document) (*sheets.Spreadsheet, error) {
	sheetsToWrite := se.getSheets(document.dataSets)

	rb := &sheets.Spreadsheet{
		Sheets: sheetsToWrite,
		Properties: &sheets.SpreadsheetProperties{
			Title: document.name,
		},
	}

	s, err := srv.Spreadsheets.Create(rb).Do()
	if err != nil {
		fmt.Printf("Err creating spreadsheet: %+v ", err)
		return nil, err
	}

	requests := make([]*sheets.Request, 0, 1)
	for i, ds := range document.dataSets {
		r := se.getFilterView(s.Sheets[i].Properties.SheetId, len(ds.Data))
		req := &sheets.Request{AddFilterView: r}
		requests = append(requests, req)
	}

	bu := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	_, err = srv.Spreadsheets.BatchUpdate(s.SpreadsheetId, bu).Do()
	if err != nil {
		fmt.Printf("Err sorting update spreadsheet: %+v ", err)
		return nil, err
	}

	return s, nil
}

func (se SheetExporter) getSheets(dataSets []models.DataSet) []*sheets.Sheet {
	sheetsToWrite := make([]*sheets.Sheet, 0, 1)

	for _, d := range dataSets {
		s := se.getSheet(d)
		sheetsToWrite = append(sheetsToWrite, s)
	}

	return sheetsToWrite
}

func (se SheetExporter) getSheet(dataSet models.DataSet) *sheets.Sheet {
	rows := make([]*sheets.RowData, 0, 1)
	rows = append(rows, se.getHeaderRow())

	for _, d := range dataSet.Data {
		r := se.getUserRow(d)
		rows = append(rows, r)
	}

	gridData := []*sheets.GridData{&sheets.GridData{
		RowData: rows,
	}}

	return &sheets.Sheet{
		Properties: &sheets.SheetProperties{
			Title: dataSet.Title,
		},
		Data: gridData,
	}
}

func (se SheetExporter) getHeaderRow() *sheets.RowData {
	headers := []string{"Screen Name", "Followers Count", "Email", "Name", "URL", "Link"}
	cells := make([]*sheets.CellData, 0, len(headers))

	for _, header := range headers {
		cell := &sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: header,
			},
		}
		cells = append(cells, cell)
	}

	return &sheets.RowData{
		Values: cells,
	}
}

// getTwitterRow creates a sheet row with twitter data.
// may not belong here..until we use reflect to create a row out of any struct..?
func (se SheetExporter) getUserRow(u models.User) *sheets.RowData {

	cells := make([]*sheets.CellData, 0, 2)

	screenNameCell := getTextCell(u.ScreenName)
	cells = append(cells, screenNameCell)

	countCell := getNumberCell(u.FollowersCount)
	cells = append(cells, countCell)

	emailCell := getTextCell(u.Email)
	cells = append(cells, emailCell)

	nameCell := getTextCell(u.Name)
	cells = append(cells, nameCell)

	urlCell := getTextCell(u.URL)
	cells = append(cells, urlCell)

	linkCell := getTextCell("https://twitter.com/" + u.ScreenName)
	cells = append(cells, linkCell)

	return &sheets.RowData{
		Values: cells,
	}
}

func (se SheetExporter) getFilterView(sheetID int64, maxRow int) *sheets.AddFilterViewRequest {
	return &sheets.AddFilterViewRequest{
		Filter: &sheets.FilterView{
			Title: "Sorted Desc",
			Range: &sheets.GridRange{
				EndColumnIndex:   6,
				StartColumnIndex: 0,
				StartRowIndex:    0,
				EndRowIndex:      int64(maxRow) + 1,
				SheetId:          sheetID,
			},
			SortSpecs: []*sheets.SortSpec{
				&sheets.SortSpec{
					SortOrder: "DESCENDING",
				},
				&sheets.SortSpec{
					SortOrder: "DESCENDING",
				},
			},
		},
	}
}

func getTextCell(text string) *sheets.CellData {
	return &sheets.CellData{
		UserEnteredValue: &sheets.ExtendedValue{
			StringValue: text,
		},
	}
}

func getNumberCell(number int64) *sheets.CellData {
	return &sheets.CellData{
		UserEnteredValue: &sheets.ExtendedValue{
			NumberValue: float64(number),
		},
	}
}
