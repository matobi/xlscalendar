package main

// Parses skating group info from xls.

import (
	"strings"
)

// Find all groups in xls.
// Currently assunes that all groups are defined in first row in xls.
func parseGroups(cells []string) ([]Group, error) {
	groups := initGroups()
	for i := range groups {
		groups[i].Column = -1
	}
	for columnIdx, cell := range cells {
		groupIdx := findGroup(groups, cell)
		if groupIdx >= 0 && groups[groupIdx].Column < 0 {
			groups[groupIdx].Column = columnIdx
		}
	}
	return groups, nil
}

// Search if a single cell contains a group name
func findGroup(groups []Group, cell string) int {
	cell = strings.ToLower(cell)
	for i := range groups {
		if groups[i].RegExp.MatchString(cell) {
			return i
		}
	}
	return -1
}
