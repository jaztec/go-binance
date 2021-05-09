package binance_test

import (
	"net/http"

	. "github.com/onsi/gomega"

	"github.com/jaztec/go-binance"
	"github.com/jaztec/go-binance/model"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Order", func() {
	Context("Should insert orders into the Binance system", func() {
		It("should add a buying order", func() {
			res := loadFixture("buy_order_ack_data")
			ts := testServer(
				"/api/v3/order",
				map[string]struct{}{
					"symbol":           {},
					"side":             {},
					"type":             {},
					"timeInForce":      {},
					"quantity":         {},
					"price":            {},
					"stopPrice":        {},
					"newOrderRespType": {},
					"timestamp":        {},
					"signature":        {},
				},
				http.StatusOK,
				res,
				nil,
			)
			defer ts.Close()

			a := newAPI(ts.URL)
			ap, err := a.Order("DOGEUSDT", model.Buy, model.TakeProfitLimit, binance.OrderParams{
				TimeInForce:      model.GoodTilCanceled,
				Quantity:         50,
				QuoteOrderQty:    0,
				Price:            0.495,
				StopPrice:        0.50,
				IcebergQty:       0,
				NewOrderRespType: model.Ack,
				RecvWindow:       0,
			})
			Expect(err).To(BeNil())
			_, ok := ap.(*model.OrderResponseAck)
			Expect(ok).To(BeTrue())
		})
	})
})
