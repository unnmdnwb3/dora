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
			val := (float64(numeratorsInWindow) / float64(denominatorsInWindow)) * 100
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

// CompleteChangesPerDays returns a slice of the number of changes per day,
// since provided ChangesPerDays only account for the dates that any changes were found.
func CompleteChangesPerDays(changesPerDays *[]models.ChangesPerDay, dates *[]time.Time) (*[]int, *[]int, error) {
	if len(*changesPerDays) == 0 {
		return nil, nil, fmt.Errorf("no change aggregates provided")
	}
	if len(*dates) == 0 {
		return nil, nil, fmt.Errorf("no dates provided")
	}
	if len(*dates) < len(*changesPerDays) {
		return nil, nil, fmt.Errorf("more changes per day than dates provided")
	}

	dailyChanges := make([]int, len(*dates))
	dailyLeadTimes := make([]int, len(*dates))
	curr := 0
	for index, date := range *dates {
		if curr < len(*changesPerDays) && date == (*changesPerDays)[curr].Date {
			dailyChanges[index] = (*changesPerDays)[curr].TotalChanges
			dailyLeadTimes[index] = int((*changesPerDays)[curr].TotalLeadTime)
			curr++
		} else {
			dailyChanges[index] = 0
			dailyLeadTimes[index] = 0
		}
	}

	return &dailyChanges, &dailyLeadTimes, nil
}
