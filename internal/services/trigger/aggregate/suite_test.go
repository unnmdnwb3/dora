package aggregate_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAggregate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "services.trigger.aggregate Suite")
}
