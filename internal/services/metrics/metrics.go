package metrics

import (
	"fmt"
	"time"

	"github.com/unnmdnwb3/dora/internal/utils/times"
)

// CalculateMovingAverages calculates the moving averages for a given slice of totals.
func CalculateMovingAverages(totals *[]int, window int) (*[]float64, error) {
	if len(*totals) == 0 {
		return nil, fmt.Errorf("no deployments per day provided to calculate moving averages")
	}
	if len(*totals) != (window*2)-1 {
		return nil, fmt.Errorf("number of pipeline runs per day provided to calculate moving averages does not match window")
	}

	offset := window - 1
	totalsInWindow := 0
	for index := 0; index < offset; index++ {
		totalsInWindow += (*totals)[index]
	}

	movingAverages := make([]float64, len(*totals)-offset)
	for index := offset; index < len(*totals); index++ {
		totalsInWindow += (*totals)[index]
		movingAverages[index-offset] = float64(totalsInWindow) / float64(window)
		totalsInWindow -= (*totals)[index-offset]
	}

	return &movingAverages, nil
}

// DatesBetween returns a slice of dates between the start and end dates.
func DatesBetween(startDate time.Time, endDate time.Time) (*[]time.Time, error) {
	startDate = times.Date(startDate)
	endDate = times.Date(endDate)

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date is after end date")
	}

	dates := []time.Time{}
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		dates = append(dates, times.Date(date))
	}
	dates = append(dates, times.Date(endDate))

	return &dates, nil
}
