package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Strategy int

const (
	Everything Strategy = iota
	Daily
	Weekly
	Monthly
)

type Range struct {
	Start, End time.Time
	Preserve   Strategy
}

type rangeConfig struct {
	EverythingUntil string
	DailyUntil      string
	WeeklyUntil     string
	MonthlyUntil    string
}

func buildRanges(cfg rangeConfig) []Range {
	if cfg.EverythingUntil != "" {
	}
	return nil
}

type SuperDuration struct {
	Duration time.Duration
	Days     int64
}

var customUnits = []string{
	"d",
	"day",
	"days",
	"w",
	"week",
	"weeks",
	"mo",
	"month",
	"months",
}

func parseDuration(d string) (SuperDuration, error) {
	d = strings.TrimSpace(d)
	found := false
	for i := range customUnits {
		if strings.HasSuffix(d, customUnits[i]) {
			fmt.Println("d", d, "custom", customUnits[i])
			pfx := strings.TrimSpace(d[:len(d)-len(customUnits[i])])
			fmt.Printf("pfx: %q\n", pfx)
			num, err := strconv.ParseInt(pfx, 10, 64)
			if err != nil {
				return SuperDuration{}, fmt.Errorf("could not parse number in %q as a duration: %w", d, err)
			}
			switch customUnits[i] {
			case "d", "day", "days":
				return SuperDuration{Days: num}, nil
			case "w", "week", "weeks":
				return SuperDuration{Days: num * 7}, nil
			case "m", "month", "months":
				return SuperDuration{Days: num * 30}, nil
			default:
				panic(fmt.Sprintf("unknown custom unit %q", customUnits[i]))
			}
		}
	}
	if !found {
		return SuperDuration{}, fmt.Errorf("did not find unit attached to duration %q", d)
	}
	return SuperDuration{}, nil
}

func main() {
	extension := flag.String("extension", "", "Only check files matching this extension ('.tar.gz' or 'zst')")
	everythingUntil := flag.String("everything-until", "", "Save everything within this duration of time")
	dailyUntil := flag.String("daily-until", "", "Save one file per day within this duration of time")
	weeklyUntil := flag.String("weekly-until", "", "Save one file per week within this duration of time")
	monthlyUntil := flag.String("monthly-until", "", "Save one file per month within this duration of time")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `backupdeleter [flags] <directory>

Print out a list of files to delete, one per line, matching the instructions 
provided by the command line flags.

`)
		flag.PrintDefaults()
	}
	flag.Parse()
	ranges := buildRanges(rangeConfig{
		EverythingUntil: *everythingUntil,
		DailyUntil:      *dailyUntil,
		WeeklyUntil:     *weeklyUntil,
		MonthlyUntil:    *monthlyUntil,
	})
	fmt.Printf("%#v\n", ranges)

	_ = extension
	_ = everythingUntil
	_ = dailyUntil
	_ = weeklyUntil
	_ = monthlyUntil
}
