package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Schedule contains all info parsed from xls.
type Schedule struct {
	Spots     []Spot
	Groups    []Group
	Trainings []Training
}

// Spot day and start- end-time of training.
type Spot struct {
	WeekDay int
	Start   time.Duration
	End     time.Duration
	Row     int // row in xls
}

// Group a skating group.
type Group struct {
	Display string         // Display name
	Short   string         // Short name
	RegExp  *regexp.Regexp // regexp to match group name in xls.
	Column  int            // column in xls
}

// Training an entry in the calendar.
type Training struct {
	GroupIdx int    // Which group
	SpotIdx  int    // What time
	Cell     string // raw cell value from xls
	Type     string // type of training, for example "Is", "Mark"
	Trainers []string
	Title    string
	Desc     string
}

func parseSchedule(rows []Row) (Schedule, error) {
	schedule := Schedule{}
	var err error
	if schedule.Spots, err = parseSpots(rows); err != nil {
		return schedule, err
	}
	if schedule.Groups, err = parseGroups(rows[0].Cells); err != nil {
		return schedule, err
	}
	if schedule.Trainings, err = parseTrainings(rows, schedule.Spots, schedule.Groups); err != nil {
		return schedule, err
	}
	return schedule, nil
}

// Merges groups (from column header) and times (from row header)
// with training description (cell). Result is all training calendar entries.
func parseTrainings(rows []Row, spots []Spot, groups []Group) ([]Training, error) {
	trainings := []Training{}
	regexpLineBreaks := regexp.MustCompile(`[\r\n]+`)

	// Loop each groups column in xls
	for groupIdx, group := range groups {
		if group.Column < 0 {
			continue // this group was not found
		}
		// loop all rows with time info
		for spotIdx, spot := range spots {
			if len(rows) <= spot.Row || len(rows[spot.Row].Cells) <= group.Column {
				continue // this cell was not in the spreadsheet.
			}
			cell := rows[spot.Row].Cells[group.Column]
			if cell == "" {
				continue
			}
			// change line breaks to <br> so it can be added in csv.
			cell = regexpLineBreaks.ReplaceAllString(cell, "<br>")

			training := Training{
				Cell:     cell,
				GroupIdx: groupIdx,
				SpotIdx:  spotIdx,
			}

			// Sometimes group name is placed directly in traingin cell. (usually for "Skridskoskolan")
			if overrideGoupIdx := findGroup(groups, cell); overrideGoupIdx >= 0 {
				training.GroupIdx = overrideGoupIdx
			}
			trainings = append(trainings, training)
		}
	}
	return trainings, nil
}

// Format calendar entry description.
func formatDescription(schedule Schedule) {
	trainingTypes := getTrainingTypes()
	trainers := getTrainers()
	for i := range schedule.Trainings {
		t := &schedule.Trainings[i]
		t.Trainers = findNameMultiple(trainers, t.Cell)
		t.Type = findName(trainingTypes, t.Cell)

		trainersString := strings.Join(t.Trainers, ",")
		t.Title = fmt.Sprintf("%s; %s; %s", t.Type, schedule.Groups[t.GroupIdx].Short, trainersString)
		t.Desc = fmt.Sprintf(`träning: <b>%s</b><br/>grupp: <b>%s</b><br/>tränare: <b>%s</b><br/>%s<br/><i>xlsauto</i>`,
			t.Type, schedule.Groups[t.GroupIdx].Display, trainersString, t.Cell)
	}
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
