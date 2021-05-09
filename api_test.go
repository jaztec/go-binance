package binance_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/jaztec/go-binance/model"

	"github.com/jaztec/go-binance"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	apiKey    = "aap"
	apiSecret = "noot"
)

func testServer(expectedPath string, expectedQueryParts map[string]struct{}, status int, response []byte, responseHeaders map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()
		if responseHeaders != nil {
			for k, v := range responseHeaders {
				w.Header().Set(k, v)
			}
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			_, _ = w.Write(response)
			return
		}
		Expect(r.URL.Path).To(Equal(expectedPath))
		for k, _ := range r.URL.Query() {
			_, ok := expectedQueryParts[k]
			Expect(ok).To(Equal(true), fmt.Sprintf("Expected %s to be present", k))
		}
		_, _ = w.Write(response)
	}))
}

func newAPI(uri string) binance.API {
	defer GinkgoRecover()

	a, err := binance.NewAPI(binance.APIConfig{
		Key:           apiKey,
		Secret:        apiSecret,
		BaseURI:       uri,
		BaseStreamURI: strings.ReplaceAll(uri, "ws", "http"),
	}, testLogger{})
	Expect(err).To(BeNil())
	return a
}

func loadFixture(name string) []byte {
	defer GinkgoRecover()
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/data/%s.json", name))
	Expect(err).To(BeNil())
	return b
}

var _ = Describe("Api", func() {
	Context("create new API", func() {
		var a binance.API
		BeforeEach(func() {
			var err error
			a, err = binance.NewAPI(binance.APIConfig{
				Key:           apiKey,
				Secret:        apiSecret,
				BaseURI:       "http://mies.mees",
				BaseStreamURI: "ws://mies.mees",
			}, testLogger{})
			Expect(err).To(BeNil())
		})

		It("should create a new API", func() {
			Expect(a).ToNot(BeNil())

			_, err := binance.NewAPI(binance.APIConfig{
				Key:           apiKey,
				Secret:        apiSecret,
				BaseURI:       "http://mies.mees",
				BaseStreamURI: "ws://mies.mees",
			}, nil)
			Expect(err).ToNot(BeNil(), "Error should not be nil when no logger was provided")
		})

		It("should have an instance of a Streamer", func() {
			Expect(a.Streamer()).ToNot(BeNil())
		})
	})

	Context("API weight results must be respected", func() {
		Context("Should halt on warning", func() {
			It("should respect API limits", func() {
				ts := testServer("/api/v3/avgPrice", map[string]struct{}{
					"symbol": {},
				}, http.StatusTooManyRequests, nil, map[string]string{"Retry-After": "30"})
				defer ts.Close()

				a := newAPI(ts.URL)
				_, err := a.AvgPrice("ETHBTC")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(binance.TooMuchCalls.Error()))

				_, err = a.AvgPrice("ETHBTC")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(binance.AtTimeout.Error()))
			})
		})

		Context("Should halt on block", func() {
			It("should respect API limits", func() {
				ts := testServer("/api/v3/avgPrice", map[string]struct{}{
					"symbol": {},
				}, http.StatusTeapot, nil, map[string]string{"Retry-After": "30"})
				defer ts.Close()

				a := newAPI(ts.URL)
				_, err := a.AvgPrice("ETHBTC")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(binance.Blocked.Error()))

				_, err = a.AvgPrice("ETHBTC")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(binance.AtTimeout.Error()))
			})
		})
	})

	Context("call API endpoints", func() {

		Context("should call account data", func() {
			It("should work on success", func() {
				v := url.Values{}
				v.Set("timestamp", "0001")
				v.Set("signature", "aSignature")
				res := loadFixture("account_data")
				ts := testServer("/api/v3/account", map[string]struct{}{
					"timestamp": {},
					"signature": {},
				}, http.StatusOK, res, nil)
				defer ts.Close()

				a := newAPI(ts.URL)

				ai, err := a.Account()
				Expect(err).To(BeNil(), "calling Account should not return error")

				Expect(ai.MakerCommission).To(Equal(10))
				Expect(ai.UpdateTime).To(Equal(1619000000000))
				Expect(ai.Balances).To(HaveLen(1))
				Expect(ai.Balances[0].Asset).To(Equal("BTC"))
				Expect(ai.Balances[0].Free).To(Equal("10.00000000"))
				Expect(ai.Balances[0].Locked).To(Equal("5.00000000"))
			})

			It("should work on error", func() {
				v := url.Values{}
				v.Set("timestamp", "0001")
				v.Set("signature", "aSignature")
				ts := testServer("/api/v3/account", map[string]struct{}{
					"timestamp": {},
					"signature": {},
				}, http.StatusInternalServerError, []byte("{}"), nil)
				defer ts.Close()

				a := newAPI(ts.URL)

				_, err := a.Account()
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(model.Error{
					Code: 0,
					Msg:  "",
				}))
			})
		})

		Context("Work on AvgPrice", func() {
			var ts *httptest.Server
			var a binance.API

			BeforeEach(func() {
				res := loadFixture("avg_price_data")
				ts = testServer("/api/v3/avgPrice", map[string]struct{}{
					"symbol": {},
				}, http.StatusOK, res, nil)

				a = newAPI(ts.URL)
			})

			AfterEach(func() {
				ts.Close()
			})

			It("should work on regular call", func() {
				av, err := a.AvgPrice("ETHBTC")
				Expect(err).To(BeNil(), "calling AvgPrice should not return error")

				Expect(av.Mins).To(Equal(5))
				Expect(av.Price).To(Equal("0.06656334"))
			})

			It("should detect missing symbol", func() {
				_, err := a.AvgPrice("")
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(binance.NoSymbolProvided))
			})
		})
	})
})
