package utils_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/utils"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utils Suite")
}

var _ = Describe("mongodb.Service", func() {
	var _ = When("Date", func() {
		It("returns true if the time are of the same date.", func() {
			ts1 := "2019-10-09T09:11:20.861Z"
			t1, err := time.Parse(time.RFC3339, ts1)
			Expect(err).To(BeNil())

			ts2 := "2019-10-09T09:12:20.861Z"
			t2, err := time.Parse(time.RFC3339, ts2)
			Expect(err).To(BeNil())

			same := utils.SameDay(t1, t2)
			Expect(same).To(BeTrue())
		})

		It("returns false if the time are of the same date.", func() {
			ts1 := "2019-10-09T09:11:20.861Z"
			t1, err := time.Parse(time.RFC3339, ts1)
			Expect(err).To(BeNil())

			ts2 := "2019-10-10T09:12:20.861Z"
			t2, err := time.Parse(time.RFC3339, ts2)
			Expect(err).To(BeNil())

			same := utils.SameDay(t1, t2)
			Expect(same).To(BeFalse())
		})
	})
})
