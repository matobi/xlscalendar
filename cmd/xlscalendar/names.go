package main

// Defines names of trainers, skating groups and training types.
// todo: Move to config file.

import (
	"regexp"
	"strings"
)

func getDaysName() []string {
	return []string{"måndag", "tisdag", "onsdag", "torsdag", "fredag", "lördag", "söndag"}
}

// Definition of all skating groups
func initGroups() []Group {
	return []Group{
		Group{Display: "Elit", Short: "Elit", RegExp: regexp.MustCompile(`elit`)},
		Group{Display: "Advanced", Short: "Advanced", RegExp: regexp.MustCompile(`advanced`)},
		Group{Display: "1", Short: "Gr 1", RegExp: regexp.MustCompile(`gr(upp)? 1`)},
		Group{Display: "2", Short: "Gr 2", RegExp: regexp.MustCompile(`gr(upp)? 2`)},
		Group{Display: "3", Short: "Gr 3", RegExp: regexp.MustCompile(`gr(upp)? 3`)},
		Group{Display: "FY", Short: "Gr FY", RegExp: regexp.MustCompile(`gr(upp)? fy`)},
		Group{Display: "FÄ", Short: "Gr FÄ", RegExp: regexp.MustCompile(`gr(upp)? fä`)},
		Group{Display: "Skridskoskola", Short: "Skridskoskola", RegExp: regexp.MustCompile(`skridsko.*skola`)},
	}
}

// Name and all matching search words.
type Name struct {
	Display string
	Tags    []string // All matching search tags. Should be changed to a regexp.
}

func getTrainers() []Name {
	return []Name{
		Name{Display: "Vera", Tags: []string{"vera"}},
		Name{Display: "Erika", Tags: []string{"erika"}},
		Name{Display: "Nicole", Tags: []string{"nicole"}},
		Name{Display: "Kia", Tags: []string{"kia"}},
		Name{Display: "Alexandra", Tags: []string{"alexandra"}},
	}
}

func getTrainingTypes() []Name {
	return []Name{
		Name{Display: "Is", Tags: []string{"is"}},
		Name{Display: "Mark", Tags: []string{"mark"}},
		Name{Display: "Balett", Tags: []string{"balett"}},
	}
}

func findName(names []Name, s string) string {
	search := strings.ToLower(s)
	for i := range names {
		for _, tag := range names[i].Tags {
			if strings.Contains(search, tag) {
				return names[i].Display
			}
		}
	}
	return ""
}

func findNameMultiple(names []Name, s string) []string {
	search := strings.ToLower(s)
	found := []string{}
	for i := range names {
		for _, tag := range names[i].Tags {
			if strings.Contains(search, tag) {
				found = append(found, names[i].Display)
			}
		}
	}
	return found
}
