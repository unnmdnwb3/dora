package metrics_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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
})
