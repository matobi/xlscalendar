package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

// Row of cells
type Row struct {
	Cells []string
}

func readXLS(xlPath string, sheetIndex int) ([]Row, error) {
	xlFile, err := xlsx.OpenFile(xlPath)
	if err != nil {
		return nil, err
	}
	if len(xlFile.Sheets) <= sheetIndex {
		return nil, fmt.Errorf("sheet number not in xls; sheet=%d; file=%s", sheetIndex, xlPath)
	}
	fmt.Printf("xls sheet name = %s\n", xlFile.Sheets[sheetIndex].Name)
	rows := []Row{}
	for _, r := range xlFile.Sheets[sheetIndex].Rows {
		cells := []string{}
		for _, c := range r.Cells {
			cells = append(cells, strings.TrimSpace(c.String()))
		}
		rows = append(rows, Row{Cells: cells})
	}
	return rows, nil
}

var regexpFilenameYear = regexp.MustCompile(`(?:\D)(\d{4})(?:\D|$)`)
var regexpFilenameWeek = regexp.MustCompile(`(?:\D)(\d{2})(?:\D|$)`)

func parseFilenameWeek(s string) int {
	return parseFilenameInt(s, regexpFilenameWeek)
}
func parseFilenameYear(s string) int {
	return parseFilenameInt(s, regexpFilenameYear)
}
func parseFilenameInt(s string, r *regexp.Regexp) int {
	parts := r.FindStringSubmatch(s)
	if parts == nil {
		return -1
	}
	n, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		return -1
	}
	return int(n)
}
