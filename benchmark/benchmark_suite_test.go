package benchmark_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRabbit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Benchmark Suite")
}
