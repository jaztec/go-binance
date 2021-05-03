package binance_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.jaztec.info/checkers/checkers/services/binance"
)

const (
	apiKey    = "aap"
	apiSecret = "noot"
)

func testServer(expectedPath string, expectedQueryParts map[string]struct{}, status int, response []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()
		if status != http.StatusOK {
			w.WriteHeader(status)
			_, _ = w.Write(response)
			return
		}
		Expect(r.URL.Path).To(Equal(expectedPath))
		for k, _ := range r.URL.Query() {
			_, ok := expectedQueryParts[k]
			Expect(ok).To(Equal(true))
		}
		_, _ = w.Write(response)
	}))
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

	Context("call API endpoints", func() {
		newAPI := func(uri string) binance.API {
			a, _ := binance.NewAPI(binance.APIConfig{
				Key:           apiKey,
				Secret:        apiSecret,
				BaseURI:       uri,
				BaseStreamURI: strings.ReplaceAll(uri, "ws", "http"),
			}, testLogger{})
			return a
		}

		Context("should call account data", func() {
			It("should work on success", func() {
				v := url.Values{}
				v.Set("timestamp", "0001")
				v.Set("signature", "aSignature")
				res := loadFixture("account_data")
				ts := testServer("/api/v3/account", map[string]struct{}{
					"timestamp": {},
					"signature": {},
				}, http.StatusOK, res)
				defer ts.Close()

				a := newAPI(ts.URL)

				ai, err := a.UserAccount()
				Expect(err).To(BeNil(), "calling UserAccount should not return error")

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
				}, http.StatusInternalServerError, []byte("{}"))
				defer ts.Close()

				a := newAPI(ts.URL)

				_, err := a.UserAccount()
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(model.Error{
					Code: 0,
					Msg:  "",
				}))
			})
		})
	})
})
