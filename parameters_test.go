package binance_test

import (
	"testing"

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

	Context("should function as required", func() {
		It("should encode in order of adding", func() {
			Expect(p.Encode()).To(Equal("mies=merel&aap=noot"))
		})

		It("should override a set parameter", func() {
			p.Set("mies", "boom")
			Expect(p.Encode()).To(Equal("mies=boom&aap=noot"))
		})
	})

	Context("have some limits on running time", func() {
		Measure("it should complete quickly", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				Expect(p.Encode()).To(Equal("mies=merel&aap=noot"))
			})

			Ω(runtime.Microseconds()).Should(BeNumerically("<", 100), "Encode() shouldn't take too long.")
			b.RecordValue("microseconds run", float64(runtime.Microseconds()))
		}, 10)
	})
})

func BenchmarkNewParameters(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		binance.NewParameters(10)
	}
}

func BenchmarkParameters_Encode(b *testing.B) {
	b.ReportAllocs()
	p := binance.NewParameters(5)
	p.Set("app", "noot")
	p.Set("mies", "mees")
	p.Set("boom", "roos")
	p.Set("vis", "kip")
	p.Set("maan", "pet")

	for n := 0; n < b.N; n++ {
		p.Encode()
	}
}

func BenchmarkParameters_Set(b *testing.B) {
	b.ReportAllocs()
	p := binance.NewParameters(2)
	for n := 0; n < b.N; n++ {
		p.Set("app", "noot")
		p.Set("mies", "mees")
	}
}
