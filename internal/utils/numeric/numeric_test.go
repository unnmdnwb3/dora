package numeric_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/unnmdnwb3/dora/internal/utils/numeric"
)

func TestNumeric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "numeric Suite")
}

var _ = Describe("utils.numeric", func() {
	var _ = When("Round", func() {
		It("rounds a float64 to a given precision.", func() {
			val := 1.23456789
			precision := 2
			rounded := numeric.Round(val, precision)
			Expect(rounded).To(Equal(1.23))
		})
	})
})
