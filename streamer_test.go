package binance_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.jaztec.info/checkers/checkers/services/binance"
)

func testStreamServer() *httptest.Server {
	accountFn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/userDataStream" {
			k := model.ListenKey{ListenKey: "allowed"} // anything works here
			b, _ := json.Marshal(k)
			_, _ = w.Write(b)
		}
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()

		if r.Header.Get("Connection") != "Upgrade" {
			accountFn(w, r)
			return
		}

		var upgrader = websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		Expect(err).To(BeNil())
		defer c.Close()

		for {
			mt, msg, err := c.ReadMessage()
			Expect(err).To(BeNil())

			if mt != websocket.TextMessage {
				continue
			}
			var parsed binance.SubscribeMessage
			err = json.Unmarshal(msg, &parsed)
			Expect(err).To(BeNil())
		}
	}))
}

var _ = Describe("Streamer", func() {
	Context("Create an API with a Streamer", func() {
		Context("Should connect to stream and subscribe to channels", func() {
			var a binance.API

			BeforeEach(func() {
				s := testStreamServer()
				var err error
				a, err = binance.NewAPI(binance.APIConfig{
					Key:           apiKey,
					Secret:        apiSecret,
					BaseURI:       s.URL,
					BaseStreamURI: strings.ReplaceAll(s.URL, "http", "ws"),
				}, testLogger{})
				Expect(err).To(BeNil())
			})

			It("should have created a streamer", func() {
				Expect(a.Streamer()).ToNot(BeNil())
			})

			It("should call KlineData function", func() {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()
				_, err := a.Streamer().KlineData(ctx, []string{"BTCETH"}, "5m")
				Expect(err).To(BeNil())
			})

			It("should call AccountData function", func() {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()
				_, err := a.Streamer().AccountData(ctx)
				Expect(err).To(BeNil())
			})

			It("should call AllTicker function", func() {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()
				_, err := a.Streamer().AllTicker(ctx)
				Expect(err).To(BeNil())
			})
		})
	})
})
