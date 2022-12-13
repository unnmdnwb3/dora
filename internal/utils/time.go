package utils

import "time"

// SameDay returns true if the two times are on the same day
func SameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}
