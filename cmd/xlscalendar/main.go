package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type args struct {
	xlPath     string
	year       int
	weekStart  int
	weekEnd    int
	sheetIndex int
	isDebug    bool
}

func main() {
	args := parseArgs()
	rows, err := readXLS(args.xlPath, args.sheetIndex)
	if err != nil {
		log.Fatalf("error parsing xls file; %+v\n", err)
	}

	schedule, err := parseSchedule(rows)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	formatDescription(schedule)
	if args.isDebug {
		printDebug(schedule)
	}

	// create a csv file for each skating group
	os.MkdirAll("./calendars", os.ModePerm)
	for i := range schedule.Groups {
		lines := generateCSV(schedule, i, args.year, args.weekStart, args.weekEnd)

		weeks := fmt.Sprintf("%d", args.weekStart)
		if args.weekStart != args.weekEnd {
			weeks = fmt.Sprintf("%dto%d", args.weekStart, args.weekEnd)
		}
		filename := fmt.Sprintf("./calendars/%d-%s-%s.csv", args.year, weeks, strings.ToLower(schedule.Groups[i].Display))

		if err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			log.Fatalf("failed to create calendar file; %+v", err)
		}
		fmt.Printf("%3d trainings; %s; %s\n", len(lines)-1, schedule.Groups[i].Display, filename)
	}
}

func parseArgs() args {
	a := args{}
	flag.StringVar(&a.xlPath, "xl", "", "path to xls file")
	flag.IntVar(&a.weekStart, "wa", -1, "start week number, for example 30")
	flag.IntVar(&a.weekEnd, "wz", -1, "end week number, for example 33")
	flag.IntVar(&a.year, "y", -1, "year, for exmple 2018")
	flag.IntVar(&a.sheetIndex, "sheet", 0, "xls sheet number, starting at 0")
	flag.BoolVar(&a.isDebug, "debug", false, "show debug info")
	flag.Parse()

	if a.xlPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if a.year <= 0 {
		// if year not specified, then try parse it from file name.
		a.year = parseFilenameYear(a.xlPath)
	}
	if a.year <= 0 {
		// If year not found, then assume current year
		a.year = time.Now().Year()
	}
	if a.weekStart <= 0 {
		// If start week not found, try parse it from file name
		a.weekStart = parseFilenameWeek(a.xlPath)
	}
	if a.weekEnd <= 0 {
		// weekEnd default equals to weekStart
		a.weekEnd = a.weekStart
	}

	fmt.Printf("args: %+v\n", a)

	if a.weekStart > a.weekEnd {
		// switch start/end week
		tmp := a.weekStart
		a.weekStart = a.weekEnd
		a.weekEnd = tmp
	}
	if a.year < 2017 || a.year > 2030 {
		log.Fatalf("unexpected year; %d\n", a.year)
	}
	if a.weekStart < 1 || a.weekStart > 53 {
		log.Fatalf("invalid start week; %d\n", a.weekStart)
	}
	if a.weekEnd < 1 || a.weekEnd > 53 {
		log.Fatalf("invalid end week; %d\n", a.weekEnd)
	}
	return a
}

func printDebug(schedule Schedule) {
	fmt.Printf("------ debug parsed values ------\n")
	dayNames := getDaysName()
	for _, spot := range schedule.Spots {
		fmt.Printf("training; %s %s-%s\n", dayNames[spot.WeekDay], formatDuration(spot.Start), formatDuration(spot.End))
	}
	for _, group := range schedule.Groups {
		column := "not found"
		if group.Column >= 0 {
			column = strconv.FormatInt(int64(group.Column), 10)
		}
		fmt.Printf("groups; %s; column=%s\n", group.Display, column)
	}
	for _, t := range schedule.Trainings {
		fmt.Printf("training; spot=%d; group=%s; %+v\n", t.SpotIdx, schedule.Groups[t.GroupIdx].Display, t)
	}
	fmt.Printf("------ end debug parsed values ------\n")
}
