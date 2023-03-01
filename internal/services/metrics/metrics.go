package metrics

import (
	"fmt"
	"time"

	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/numeric"
	"github.com/unnmdnwb3/dora/internal/utils/times"
)

// MovingAverages calculates the moving averages for a given slice of totals.
func MovingAverages(totals *[]int, window int) (*[]float64, error) {
	if len(*totals) == 0 {
		return nil, fmt.Errorf("no data provided to calculate moving averages")
	}

	offset := window - 1
	totalsInWindow := 0

	for index := 0; index < offset; index++ {
		totalsInWindow += (*totals)[index]
	}

	movingAverages := make([]float64, len(*totals)-offset)

	for index := offset; index < len(*totals); index++ {
		totalsInWindow += (*totals)[index]
		val := float64(totalsInWindow) / float64(window)
		movingAverages[index-offset] = numeric.Round(val, 2)
		totalsInWindow -= (*totals)[index-offset]
	}

	return &movingAverages, nil
}

// MovingAveragesRatio calculates the moving averages for a given two slices of totals.
func MovingAveragesRatio(numerators *[]int, denominators *[]int, window int) (*[]float64, error) {
	if len(*numerators) == 0 || len(*denominators) == 0 {
		return nil, fmt.Errorf("no data provided to calculate moving averages")
	}

	offset := window - 1
	numeratorsInWindow := 0
	denominatorsInWindow := 0

	for index := 0; index < offset; index++ {
		numeratorsInWindow += (*numerators)[index]
		denominatorsInWindow += (*denominators)[index]
	}

	movingAverages := make([]float64, len(*numerators)-offset)

	for index := offset; index < len(*numerators); index++ {
		numeratorsInWindow += (*numerators)[index]
		denominatorsInWindow += (*denominators)[index]

		if denominatorsInWindow == 0 {
			movingAverages[index-offset] = 0
		} else {
			val := (float64(numeratorsInWindow) / float64(denominatorsInWindow))
			movingAverages[index-offset] = numeric.Round(val, 2)
		}

		numeratorsInWindow -= (*numerators)[index-offset]
		denominatorsInWindow -= (*denominators)[index-offset]
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

// CompleteIncidentsPerDays returns a slice of the number of incidents per day,
// since provided IncidentsPerDays only account for the dates that any incidents were found.
func CompleteIncidentsPerDays(incidentsPerDays *[]models.IncidentsPerDay, dates *[]time.Time) (*[]int, *[]int, error) {
	if len(*dates) == 0 {
		return nil, nil, fmt.Errorf("no dates provided")
	}

	dailyIncidents := make([]int, len(*dates))
	dailyDurations := make([]int, len(*dates))

	curr := 0

	for index, date := range *dates {
		sumIncidents := 0
		sumDuration := 0
		for j := curr; j < len(*incidentsPerDays); j++ {
			if (*incidentsPerDays)[j].Date == date {
				sumIncidents += (*incidentsPerDays)[j].TotalIncidents
				sumDuration += int((*incidentsPerDays)[j].TotalDuration)
				curr++
			} else {
				break
			}
		}

		dailyIncidents[index] = sumIncidents
		dailyDurations[index] = sumDuration
	}

	return &dailyIncidents, &dailyDurations, nil
}

// CompletePipelineRunsPerDays returns a slice of the number of pipeline runs per day,
// since provided PipelineRunsPerDays only account for the dates that any pipeline runs were found.
func CompletePipelineRunsPerDays(pipelineRunsPerDays *[]models.PipelineRunsPerDay, dates *[]time.Time) (*[]int, error) {
	if len(*dates) == 0 {
		return nil, fmt.Errorf("no dates provided")
	}

	dailyPipelineRuns := make([]int, len(*dates))

	curr := 0
	for i, date := range *dates {
		sum := 0
		for j := curr; j < len(*pipelineRunsPerDays); j++ {
			if (*pipelineRunsPerDays)[j].Date == date {
				sum += (*pipelineRunsPerDays)[j].TotalPipelineRuns
				curr++
			} else {
				break
			}
		}

		dailyPipelineRuns[i] = sum
	}

	return &dailyPipelineRuns, nil
}

// CompleteChangesPerDays returns a slice of the number of changes per day,
// since provided ChangesPerDays only account for the dates that any changes were found.
func CompleteChangesPerDays(changesPerDays *[]models.ChangesPerDay, dates *[]time.Time) (*[]int, *[]int, error) {
	if len(*dates) == 0 {
		return nil, nil, fmt.Errorf("no dates provided")
	}

	dailyChanges := make([]int, len(*dates))
	dailyLeadTimes := make([]int, len(*dates))

	curr := 0

	for i, date := range *dates {
		sumChanges := 0
		sumLeadTimes := 0.0
		for j := curr; j < len(*changesPerDays); j++ {
			if (*changesPerDays)[j].Date == date {
				sumChanges += (*changesPerDays)[j].TotalChanges
				sumLeadTimes += (*changesPerDays)[j].TotalLeadTime
				curr++
			} else {
				break
			}
		}

		dailyChanges[i] = sumChanges
		dailyLeadTimes[i] = int(sumLeadTimes)
	}

	return &dailyChanges, &dailyLeadTimes, nil
}
