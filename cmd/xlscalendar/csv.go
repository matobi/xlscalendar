package main

// Creates csv info from schedule. Filter out a single skating group per csv.

import (
	"fmt"
	"time"
)

const csvHeader = `Subject,Start Date,Start Time,End Date,End Time,All day event,Description`
const csvLine = `"%s",%s,%s,%s,%s,FALSE,"%s"`

func generateCSV(schedule Schedule, groupIdx int, year int, weekStart int, weekEnd int) []string {
	lines := []string{csvHeader}
	group := schedule.Groups[groupIdx]

	for week := weekStart; week <= weekEnd; week++ {
		mondayDate := firstDayOfISOWeek(year, week)
		for _, training := range schedule.Trainings {
			if training.GroupIdx != groupIdx {
				continue
			}
			spot := schedule.Spots[training.SpotIdx]
			lines = append(lines, formatLine(group, spot, training, mondayDate))
		}
	}
	return lines
}

func formatLine(group Group, spot Spot, training Training, mondayDate time.Time) string {
	day := mondayDate.AddDate(0, 0, spot.WeekDay)
	start := day.Add(spot.Start)
	end := day.Add(spot.End)
	const layoutDate = "01/02/2006"
	const layoutTime = "15:04 PM"
	return fmt.Sprintf(csvLine, training.Title,
		start.Format(layoutDate), start.Format(layoutTime),
		end.Format(layoutDate), end.Format(layoutTime),
		training.Desc,
	)
}
