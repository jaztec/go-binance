package binance_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

type testHandler struct {
	expectedPath       string
	expectedQueryParts map[string]struct{}
	status             int
	response           []byte
}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer GinkgoRecover()
	if h.status != http.StatusOK {
		w.WriteHeader(h.status)
		_, _ = w.Write(h.response)
		return
	}
	Expect(r.URL.Path).To(Equal(h.expectedPath))
	for k, _ := range r.URL.Query() {
		_, ok := h.expectedQueryParts[k]
		Expect(ok).To(Equal(true))
	}
	_, _ = w.Write(h.response)
}

func testServer(expectedPath string, expectedQueryParts map[string]struct{}, status int, res []byte) *httptest.Server {
	return httptest.NewServer(testHandler{
		expectedPath:       expectedPath,
		expectedQueryParts: expectedQueryParts,
		status:             status,
		response:           res,
	})
}

func generateSignature(secret, path string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(path))

	return hex.EncodeToString(h.Sum(nil))
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
		newAPI := func(client *http.Client, uri string) binance.API {
			a, _ := binance.NewAPI(binance.APIConfig{
				Key:           apiKey,
				Secret:        apiSecret,
				BaseURI:       uri,
				BaseStreamURI: strings.ReplaceAll(uri, "ws", "http"),
				Client:        client,
			}, testLogger{})
			return a
		}

		Context("should call account data", func() {
			It("should work on success", func() {
				v := url.Values{}
				v.Set("timestamp", "0001")
				v.Set("signature", generateSignature(apiSecret, "/api/v3/account?timestamp=0001"))
				res := loadFixture("account_data")
				ts := testServer("/api/v3/account", map[string]struct{}{
					"timestamp": {},
					"signature": {},
				}, http.StatusOK, res)
				defer ts.Close()

				a := newAPI(ts.Client(), ts.URL)

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
				v.Set("signature", generateSignature(apiSecret, "/api/v3/account?timestamp=0001"))
				ts := testServer("/api/v3/account", map[string]struct{}{
					"timestamp": {},
					"signature": {},
				}, http.StatusInternalServerError, []byte("{}"))
				defer ts.Close()

				a := newAPI(ts.Client(), ts.URL)

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
