package metrics

import (
	"fmt"
	"time"

	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/utils/times"
)

// CalculateMovingAverages calculates the moving averages for a given slice of totals.
func CalculateMovingAverages(totals *[]int, window int) (*[]float64, error) {
	if len(*totals) == 0 {
		return nil, fmt.Errorf("no data provided to calculate moving averages")
	}
	if len(*totals) != (window*2)-1 {
		return nil, fmt.Errorf("number of data provided to calculate moving averages does not match window")
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

// CalculateMovingAveragesRatio calculates the moving averages for a given two slices of totals.
func CalculateMovingAveragesRatio(numerators *[]int, denominators *[]int, window int) (*[]float64, error) {
	if len(*numerators) == 0 || len(*denominators) == 0 {
		return nil, fmt.Errorf("no data provided to calculate moving averages")
	}
	if len(*numerators) != (window*2)-1 || len(*denominators) != (window*2)-1 {
		return nil, fmt.Errorf("number of data provided to calculate moving averages does not match window")
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
			movingAverages[index-offset] = (float64(numeratorsInWindow) / float64(denominatorsInWindow)) * 100
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
	if len(*incidentsPerDays) == 0 {
		return nil, nil, fmt.Errorf("no incidents aggregates provided")
	}
	if len(*dates) == 0 {
		return nil, nil, fmt.Errorf("no dates provided")
	}
	if len(*dates) < len(*incidentsPerDays) {
		return nil, nil, fmt.Errorf("more incidents per day than dates provided")
	}

	dailyIncidents := make([]int, len(*dates))
	dailyDurations := make([]int, len(*dates))
	curr := 0
	for index, date := range *dates {
		if curr < len(*incidentsPerDays) && date == (*incidentsPerDays)[curr].Date {
			dailyIncidents[index] = (*incidentsPerDays)[curr].TotalIncidents
			dailyDurations[index] = int((*incidentsPerDays)[curr].TotalDuration)
			curr++
		} else {
			dailyIncidents[index] = 0
			dailyDurations[index] = 0
		}
	}

	return &dailyIncidents, &dailyDurations, nil
}

// CompletePipelineRunsPerDays returns a slice of the number of pipeline runs per day,
// since provided PipelineRunsPerDays only account for the dates that any pipeline runs were found.
func CompletePipelineRunsPerDays(pipelineRunsPerDays *[]models.PipelineRunsPerDay, dates *[]time.Time) (*[]int, error) {
	if len(*pipelineRunsPerDays) == 0 {
		return nil, fmt.Errorf("no pipeline runs aggregates provided")
	}
	if len(*dates) == 0 {
		return nil, fmt.Errorf("no dates provided")
	}
	if len(*dates) < len(*pipelineRunsPerDays) {
		return nil, fmt.Errorf("more pipeline runs per day than dates provided")
	}

	dailyPipelineRuns := make([]int, len(*dates))
	curr := 0
	for index, date := range *dates {
		if curr < len(*pipelineRunsPerDays) && date == (*pipelineRunsPerDays)[curr].Date {
			dailyPipelineRuns[index] = (*pipelineRunsPerDays)[curr].TotalPipelineRuns
			curr++
		} else {
			dailyPipelineRuns[index] = 0
		}
	}

	return &dailyPipelineRuns, nil
}
