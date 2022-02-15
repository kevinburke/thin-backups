package main

import (
	"flag"
	"fmt"
	"log"
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

func Add(t time.Time, d SuperDuration) time.Time {
	if d.Duration != 0 {
		return t.Add(d.Duration)
	}
	return time.Date(t.Year(), t.Month(), t.Day()+d.Days, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func buildRanges(cfg rangeConfig) ([]Range, error) {
	now := time.Now()
	ranges := make([]Range, 0)
	if cfg.EverythingUntil != "" {
		e, err := parseDuration(cfg.EverythingUntil)
		if err != nil {
			return nil, fmt.Errorf("invalid EverythingUntil value %q: %v", cfg.EverythingUntil, err)
		}
		ranges = append(ranges, Range{
			Start:    now,
			End:      Add(now, e),
			Preserve: Everything,
		})
	}
	if cfg.DailyUntil != "" {
		e, err := parseDuration(cfg.DailyUntil)
		if err != nil {
			return nil, fmt.Errorf("invalid DailyUntil value %q: %v", cfg.DailyUntil, err)
		}
		var start time.Time
		if len(ranges) > 0 {
			start = ranges[len(ranges)-1].End
		} else {
			start = now
		}
		end := Add(now, e)
		if start.After(end) {
			return nil, fmt.Errorf("invalid DailyUntil value: start (%v) is after end (%v)", start, end)
		}
		ranges = append(ranges, Range{
			Start:    start,
			End:      end,
			Preserve: Daily,
		})
	}
	return ranges, nil
}

type SuperDuration struct {
	Duration time.Duration
	Days     int
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
			pfx := strings.TrimSpace(d[:len(d)-len(customUnits[i])])
			num, err := strconv.ParseInt(pfx, 10, 64)
			if err != nil {
				return SuperDuration{}, fmt.Errorf("could not parse number in %q as a duration: %w", d, err)
			}
			switch customUnits[i] {
			case "d", "day", "days":
				return SuperDuration{Days: int(num)}, nil
			case "w", "week", "weeks":
				return SuperDuration{Days: int(num) * 7}, nil
			case "mo", "month", "months":
				return SuperDuration{Days: int(num) * 30}, nil
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
	ranges, err := buildRanges(rangeConfig{
		EverythingUntil: *everythingUntil,
		DailyUntil:      *dailyUntil,
		WeeklyUntil:     *weeklyUntil,
		MonthlyUntil:    *monthlyUntil,
	})
	if err != nil {
		log.Fatal(err)
	}
	dir := flag.Arg(1)
	if dir == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "please provide a directory to read files from\n")
		flag.Usage()
	}
	// names := os.Readdirnames(dir)
	fmt.Printf("ranges: %#v\n", ranges)

	_ = extension
	_ = everythingUntil
	_ = dailyUntil
	_ = weeklyUntil
	_ = monthlyUntil
}
