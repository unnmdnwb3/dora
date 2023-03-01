package metrics_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unnmdnwb3/dora/internal/services/metrics"
)

var _ = Describe("services.metrics.metrics", func() {

	var _ = When("MovingAverages", func() {
		It("returns a list of MovingAverages.", func() {
			deploymentsPerDay := []int{1, 2, 3, 4, 5}

			movingAverages, err := metrics.MovingAverages(&deploymentsPerDay, 3)
			Expect(err).To(BeNil())
			Expect(len(*movingAverages)).To(Equal(3))
			Expect((*movingAverages)[0]).To(Equal(2.0))
			Expect((*movingAverages)[1]).To(Equal(3.0))
			Expect((*movingAverages)[2]).To(Equal(4.0))
		})
	})

	var _ = When("MovingAveragesRatio", func() {
		It("returns a list of MovingAverages when given two slices.", func() {
			numerators := []int{1, 2, 2, 0, 3}
			denominators := []int{5, 5, 10, 1, 9}

			movingAverages, err := metrics.MovingAveragesRatio(&numerators, &denominators, 3)
			Expect(err).To(BeNil())
			Expect(len(*movingAverages)).To(Equal(3))
			Expect((*movingAverages)).To(Equal([]float64{0.25, 0.25, 0.25}))
		})
	})
})
