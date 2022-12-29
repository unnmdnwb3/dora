package trigger_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrigger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "services.trigger Suite")
}
