package times

import "time"

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
