package main

// Finds all training times (day, start- and end-time) in xls.

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Find all training times in xls
func parseSpots(rows []Row) ([]Spot, error) {
	currDay := -1
	spots := []Spot{}
	for y := range rows {
		if len(rows[y].Cells) == 0 {
			continue
		}
		rowText := strings.ToLower(strings.Join(rows[y].Cells, " ; ")) // merge all cells in a row

		if day, err := findDayInWeek(rowText); day >= 0 || err != nil {
			if err != nil {
				return nil, fmt.Errorf("line %d; %s", y, err.Error())
			}
			currDay = day
		}

		durations, err := findDurations(rowText)
		if err != nil {
			return nil, fmt.Errorf("line %d; %s", y, err.Error())
		}
		switch len(durations) {
		case 0:
			continue // no timecodes on this line
		case 2:
			spot := Spot{WeekDay: currDay, Start: durations[0], End: durations[1], Row: y}
			spots = append(spots, spot)
		default:
			return nil, fmt.Errorf("line %d; Expects 0 or 2 timecodes per line, found %d", y, len(durations))
		}
	}
	return spots, nil
}

// searches for day names in an xls row.
func findDayInWeek(line string) (int, error) {
	days := getDaysName()
	day := -1
	for i := range days {
		if strings.Contains(line, days[i]) {
			if day >= 0 {
				// found more than one day on a single row. We don't know which one to use.
				return -1, fmt.Errorf("multiple days on same row; %s; %s", days[day], days[i])
			}
			day = i
		}
	}
	return day, nil
}

var regexpTime = regexp.MustCompile(`(\d\d[\.\:]\d\d)`)

// find al time codes in an xls row. A time code is a string likw "NN:NN".
func findDurations(line string) ([]time.Duration, error) {
	durations := []time.Duration{}
	timecodes := regexpTime.FindAllString(line, -1)
	if timecodes == nil {
		return durations, nil
	}

	for _, timecode := range timecodes {
		d, err := clockToDuration(timecode)
		if err != nil {
			return durations, err
		}
		durations = append(durations, d)
	}
	return durations, nil
}

// parses string "NN:NN" to a duration (hour and minute)
func clockToDuration(s string) (time.Duration, error) {
	empty := time.Duration(0)
	if len(s) != 5 {
		return empty, fmt.Errorf("bad time format; %s", s)
	}
	hour, errHour := strconv.ParseInt(s[0:2], 10, 32)
	minute, errMinute := strconv.ParseInt(s[3:5], 10, 32)
	if errHour != nil || errMinute != nil {
		return empty, fmt.Errorf("bad time format; %s", s)
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 60 {
		return empty, fmt.Errorf("bad time format; %s", s)
	}
	return time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute, nil
}
