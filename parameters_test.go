package binance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.jaztec.info/checkers/checkers/services/binance"
)

var _ = Describe("Parameters", func() {
	It("Should encode in order of adding", func() {
		p := binance.Parameters{}
		p.Set("mies", "merel")
		p.Set("aap", "noot")

		Expect(p.Encode()).To(Equal("mies=merel&aap=noot"))
	})
})
