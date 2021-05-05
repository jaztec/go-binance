package binance_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.jaztec.info/checkers/checkers/services/binance"
)

var (
	stopTestServer = false
	stopMux        = sync.Mutex{}
)

type testStreamServer struct {
	it           uint32
	wasStopped   bool
	withStopper  uint32
	afterStopped chan binance.SubscribeMessage
}

func (s *testStreamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer GinkgoRecover()
	accountFn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/userDataStream" {
			k := model.ListenKey{ListenKey: "userDataStreamAllowed"} // anything works here
			b, _ := json.Marshal(k)
			_, _ = w.Write(b)
		}
	}

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

		stopMux.Lock()
		s.it++
		if stopTestServer {
			s.afterStopped <- parsed
		}
		if s.withStopper > 0 && s.withStopper == s.it {
			stopTestServer = true
			stopMux.Unlock()
			return
		}
		stopMux.Unlock()
	}
}

func startTestStreamServer(withStopper uint32, afterStopped chan binance.SubscribeMessage) *httptest.Server {
	return httptest.NewServer(&testStreamServer{
		it:           0,
		wasStopped:   false,
		withStopper:  withStopper,
		afterStopped: afterStopped,
	})
}

var _ = Describe("Streamer", func() {
	Context("Create an API with a Streamer", func() {
		Context("Should connect to stream and subscribe to channels", func() {
			var a binance.API
			var start = func(withStopper uint32, afterStopped chan binance.SubscribeMessage) {
				s := startTestStreamServer(withStopper, afterStopped)
				var err error
				a, err = binance.NewAPI(binance.APIConfig{
					Key:           apiKey,
					Secret:        apiSecret,
					BaseURI:       s.URL,
					BaseStreamURI: strings.ReplaceAll(s.URL, "http", "ws"),
				}, testLogger{})
				Expect(err).To(BeNil())
			}

			BeforeEach(func() {
				start(0, nil)
			})

			It("should have created a streamer", func() {
				Expect(a.Streamer()).ToNot(BeNil())
			})

			It("should call Kline function", func() {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()
				_, err := a.Streamer().Kline(ctx, []string{"BTCETH"}, "5m")
				Expect(err).To(BeNil())
			})

			It("should call UserDataStream function", func() {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()
				_, err := a.Streamer().UserDataStream(ctx)
				Expect(err).To(BeNil())
			})

			It("should call TickerArr function", func() {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()
				_, err := a.Streamer().TickerArr(ctx)
				Expect(err).To(BeNil())
			})

			It("should keep the stream alive", func() {
				afterStopped := make(chan binance.SubscribeMessage, 1)
				start(3, afterStopped)

				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()
				var err error

				_, err = a.Streamer().Kline(ctx, []string{"BTCETH"}, "5m")
				Expect(err).To(BeNil())
				_, err = a.Streamer().UserDataStream(ctx)
				Expect(err).To(BeNil())
				_, err = a.Streamer().TickerArr(ctx)
				Expect(err).To(BeNil())

				////TODO This tests actually works sometimes but other times it blocks
				Expect(<-afterStopped).To(Equal(binance.SubscribeMessage{
					Method: binance.Subscribe,
					Params: []string{
						"!ticker@arr",
						"BTCETH@kline_5m",
						"userDataStreamAllowed",
					},
					ID: 4,
				}))
			})
		})
	})
})
