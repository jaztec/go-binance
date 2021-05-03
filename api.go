package binance

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
)

const (
	BaseApiURI = "https://api.binance.com"

	APIKeyHeaderName = "X-MBX-APIKEY"
)

type weightChecker struct {
	allowed bool
	weight  int
}

func (wc weightChecker) deactivate(seconds int) {
	wc.allowed = false
	go func(seconds int) {
		tC := time.Tick(time.Second * time.Duration(seconds))
		<-tC
		wc.allowed = true
	}(seconds)
}

func (wc weightChecker) checkResponse(response *http.Response) *APIError {
	fn := func(response *http.Response) {
		retry := response.Header.Get("Retry-After")
		if retry == "" {
			return
		}
		i, err := strconv.Atoi(retry)
		if err != nil {
			panic(err)
		}
		wc.deactivate(i)
	}
	if response.StatusCode == http.StatusTooManyRequests {
		fn(response)
		return &TooMuchCalls
	}
	if response.StatusCode == http.StatusTeapot {
		fn(response)
		return &Blocked
	}
	return nil
}

type APIConfig struct {
	Key           string
	Secret        string
	BaseURI       string
	BaseStreamURI string
	Client        *http.Client
}

type API interface {
	Prices(symbol string) ([]model.Price, error)
	OrderBook(symbol string, limit int) (model.Order, error)
	UserOrderBook(symbol string, startTime, endTime int64, limit int) ([]model.UserOrder, error)
	UserAccount() (ai model.AccountInfo, err error)

	Streamer() Streamer
}

type api struct {
	cfg      APIConfig
	checker  weightChecker
	logger   Logger
	streamer Streamer
}

func (a *api) Streamer() Streamer {
	return a.streamer
}

func NewAPI(cfg APIConfig, logger Logger) (API, error) {
	if logger == nil {
		return nil, errors.New("expect an instantiated logger")
	}
	a := &api{
		cfg: cfg,
		checker: weightChecker{
			allowed: true,
			weight:  0,
		},
		logger: logger,
	}

	if a.cfg.Client == nil {
		a.cfg.Client = a.client()
	}

	a.streamer = newStreamer(a, logger)

	return a, nil
}
