package times

import (
	"fmt"
	"strconv"
	"time"
)

// Date returns a time with the time set to 00:00:00.
func Date(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.UTC().Location())
}

// Day returns the day of a date in the format YYYY-MM-DD.
func Day(date time.Time) string {
	return date.Format("2006-01-02")
}

// SameDay returns true if the two times are on the same day.
func SameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

// Duration returns the duration from a Prometheus step.
func Duration(step string) (time.Duration, error) {
	if len(step) < 2 {
		return 0, fmt.Errorf("invalid format: %s", step)
	}

	count, err := strconv.Atoi(step[:len(step)-1])
	if err != nil {
		return 0, err
	}
	if count < 1 {
		return 0, fmt.Errorf("invalid step count: %d", count)
	}

	unit := step[len(step)-1:]
	switch unit {
	case "s":
		return time.Second * time.Duration(count), nil
	case "m":
		return time.Minute * time.Duration(count), nil
	case "h":
		return time.Hour * time.Duration(count), nil
	case "d":
		return time.Hour * 24 * time.Duration(count), nil
	case "w":
		return time.Hour * 24 * 7 * time.Duration(count), nil
	case "y":
		return time.Hour * 24 * 365 * time.Duration(count), nil
	default:
		return 0, fmt.Errorf("invalid step unit: %s", unit)
	}
}
