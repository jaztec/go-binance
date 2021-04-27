package binance_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBinance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Binance Suite")
}
