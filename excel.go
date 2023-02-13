package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// ReadExcel converts an Excel spreadsheet's "Sheet1" page into an array of
// values based on the first row being the column names.
func ReadExcel(fn, sheet string) ([]map[string]string, error) {
	out := []map[string]string{}

	if sheet == "" {
		sheet = "Sheet1"
	}

	f, err := excelize.OpenFile(fn)
	if err != nil {
		return out, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("ERR: %s", err.Error())
		}
	}()

	// Get column names from first row
	columns := make([]string, 0)
	for i := 'A'; i <= 'Z'; i++ {
		cell, err := f.GetCellValue(sheet, fmt.Sprintf("%c1", i))
		if err != nil {
			log.Printf("ERR: %s", err.Error())
			break
		}
		if cell == "" {
			break
		}
		columns = append(columns, cell)
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return out, err
	}
	for k, row := range rows {
		if k == 0 {
			log.Printf("INFO: Skipping first row")
			continue
		}
		r := map[string]string{}
		for k2, col := range row {
			if k2 > len(columns) {
				break
			}
			r[columns[k2]] = col
		}
		out = append(out, r)
	}

	return out, nil
}
