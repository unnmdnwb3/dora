package times_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/utils/times"
)

func TestTime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "times Suite")
}

var _ = Describe("utils.time", func() {
	var _ = When("Date", func() {
		It("returns true if the time are of the same date.", func() {
			ts1 := "2019-10-09T09:11:20.861Z"
			t1, err := time.Parse(time.RFC3339, ts1)
			Expect(err).To(BeNil())

			ts2 := "2019-10-09T09:12:20.861Z"
			t2, err := time.Parse(time.RFC3339, ts2)
			Expect(err).To(BeNil())

			same := times.SameDay(t1, t2)
			Expect(same).To(BeTrue())
		})

		It("returns false if the time are of the same date.", func() {
			ts1 := "2019-10-09T09:11:20.861Z"
			t1, err := time.Parse(time.RFC3339, ts1)
			Expect(err).To(BeNil())

			ts2 := "2019-10-10T09:12:20.861Z"
			t2, err := time.Parse(time.RFC3339, ts2)
			Expect(err).To(BeNil())

			same := times.SameDay(t1, t2)
			Expect(same).To(BeFalse())
		})
	})

	var _ = When("Date", func() {
		It("returns a time with the time set to 00:00:00.", func() {
			ts := "2019-10-09T09:11:20.861Z"
			t, err := time.Parse(time.RFC3339, ts)
			Expect(err).To(BeNil())

			d := times.Date(t)
			Expect(d.Year()).To(Equal(2019))
			Expect(d.Month()).To(Equal(time.October))
			Expect(d.Day()).To(Equal(9))
			Expect(d.Hour()).To(Equal(0))
			Expect(d.Minute()).To(Equal(0))
			Expect(d.Second()).To(Equal(0))
			Expect(d.Nanosecond()).To(Equal(0))
		})
	})

	var _ = When("Day", func() {
		It("returns the day of a date in the format YYYY-MM-DD.", func() {
			ts := "2019-10-09T09:11:20.861Z"
			t, err := time.Parse(time.RFC3339, ts)
			Expect(err).To(BeNil())

			d := times.Day(t)
			Expect(d).To(Equal("2019-10-09"))
		})
	})
})
