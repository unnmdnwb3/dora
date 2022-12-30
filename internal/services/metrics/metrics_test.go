package metrics_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/unnmdnwb3/dora/internal/models"
	"github.com/unnmdnwb3/dora/internal/services/metrics"
)

var _ = Describe("services.metrics.metrics", func() {

	var _ = When("CalculateMovingAverages", func() {
		It("returns a list of MovingAverages.", func() {
			deploymentsPerDay := []int{1, 2, 3, 4, 5}

			movingAverages, err := metrics.CalculateMovingAverages(&deploymentsPerDay, 3)
			Expect(err).To(BeNil())
			Expect(len(*movingAverages)).To(Equal(3))
			Expect((*movingAverages)[0]).To(Equal(2.0))
			Expect((*movingAverages)[1]).To(Equal(3.0))
			Expect((*movingAverages)[2]).To(Equal(4.0))
		})
	})

	var _ = When("CalculateMovingAveragesRatio", func() {
		It("returns a list of MovingAverages when given two slices.", func() {
			numerators := []int{1, 2, 2, 0, 3}
			denominators := []int{5, 5, 10, 1, 9}

			movingAverages, err := metrics.CalculateMovingAveragesRatio(&numerators, &denominators, 3)
			Expect(err).To(BeNil())
			Expect(len(*movingAverages)).To(Equal(3))
			Expect((*movingAverages)).To(Equal([]float64{25, 25, 25}))
		})
	})

	var _ = When("DatesBetween", func() {
		It("returns a list of dates between two dates.", func() {
			startDate, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")
			endDate, _ := time.Parse(time.RFC3339, "2020-02-10T15:29:50.092Z")

			dates, err := metrics.DatesBetween(startDate, endDate)
			Expect(err).To(BeNil())
			Expect(len(*dates)).To(Equal(7))
		})

		It("returns an error", func() {
			startDateAfterEndDate, _ := time.Parse(time.RFC3339, "2020-02-06T14:29:50.092Z")
			endDate, _ := time.Parse(time.RFC3339, "2020-02-04T14:29:50.092Z")

			_, err := metrics.DatesBetween(startDateAfterEndDate, endDate)
			Expect(err).To(Not(BeNil()))
		})
	})

	var _ = When("CompleteIncidentsPerDays", func() {
		It("returns the complete list of incidents and durations per day between two dates.", func() {
			deploymentID := primitive.NewObjectID()
			incidentsPerDays := []models.IncidentsPerDay{
				{
					DeploymentID:   deploymentID,
					Date:           time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 1,
					TotalDuration:  600,
				},
				{
					DeploymentID:   deploymentID,
					Date:           time.Date(2022, 12, 26, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 2,
					TotalDuration:  1200,
				},
				{
					DeploymentID:   deploymentID,
					Date:           time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 1,
					TotalDuration:  600,
				},
				{
					DeploymentID:   deploymentID,
					Date:           time.Date(2022, 12, 29, 0, 0, 0, 0, time.UTC),
					TotalIncidents: 2,
					TotalDuration:  1200,
				},
			}

			dates, err := metrics.DatesBetween(
				time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 12, 29, 0, 0, 0, 0, time.UTC),
			)
			Expect(err).To(BeNil())

			dailyIncidents, dailyDuration, err := metrics.CompleteIncidentsPerDays(&incidentsPerDays, dates)
			Expect(err).To(BeNil())
			Expect(len(*dailyIncidents)).To(Equal(6))
			Expect(len(*dailyDuration)).To(Equal(6))
		})
	})

	var _ = When("CompletePipelineRunsPerDays", func() {
		It("returns the complete list of pipeline runs per day between two dates.", func() {
			pipelineID := primitive.NewObjectID()
			date1, _ := time.Parse(time.RFC3339, "2020-02-04T00:00:00.000Z")
			date2, _ := time.Parse(time.RFC3339, "2020-02-05T00:00:00.000Z")
			date3, _ := time.Parse(time.RFC3339, "2020-02-06T00:00:00.000Z")

			dates := []time.Time{date1, date2, date3}
			persistedDailyPipelineRuns := []models.PipelineRunsPerDay{
				{
					PipelineID:        pipelineID,
					Date:              date1,
					TotalPipelineRuns: 1,
				},
				{
					PipelineID:        pipelineID,
					Date:              date3,
					TotalPipelineRuns: 1,
				},
			}

			completeDailyPipelineRuns, err := metrics.CompletePipelineRunsPerDays(&persistedDailyPipelineRuns, &dates)
			Expect(err).To(BeNil())
			Expect(len(*completeDailyPipelineRuns)).To(Equal(3))
			Expect((*completeDailyPipelineRuns)[0]).To(Equal(1))
			Expect((*completeDailyPipelineRuns)[1]).To(Equal(0))
			Expect((*completeDailyPipelineRuns)[2]).To(Equal(1))
		})
	})
})
