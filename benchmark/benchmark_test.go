package benchmark_test

import (
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gray.net/tool-container-benchmark/benchmark"
)

var _ = Describe("Benchmark", func() {
	var ctrl *gomock.Controller // repo *mock_repo.MockInterface
	// b    *benchmark.Benchmark

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		// b = &benchmark.Benchmark{}
	})

	AfterEach(func() { ctrl.Finish() })

	Describe("GenerateRandomString", func() {
		It("returns string of random characters of a given length", func() {
			length := 30
			s := benchmark.GenerateRandomString(length)
			Expect(len(s)).To(Equal(length))
		})
		It("returned string does not contain numbers", func() {
			length := 10
			nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

			s := benchmark.GenerateRandomString(length)
			Expect(strings.Split(s, "")).ToNot(ContainElements(nums))
		})
	})
})
