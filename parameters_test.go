package binance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.jaztec.info/checkers/checkers/services/binance"
)

var _ = Describe("Parameters", func() {
	var p binance.Parameters

	BeforeEach(func() {
		p = binance.NewParameters(2)
		p.Set("mies", "merel")
		p.Set("aap", "noot")
	})

	It("should encode in order of adding", func() {
		Expect(p.Encode()).To(Equal("mies=merel&aap=noot"))
	})

	It("should override a set parameter", func() {
		p.Set("mies", "boom")
		Expect(p.Encode()).To(Equal("mies=boom&aap=noot"))
	})

	Measure("it should complete quickly", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			Expect(p.Encode()).To(Equal("mies=merel&aap=noot"))
		})

		Ω(runtime.Seconds()).Should(BeNumerically("<", 0.2), "Encode() shouldn't take too long.")
	}, 10)
})
