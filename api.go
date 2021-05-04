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
}

type API interface {
	Account() (ai model.AccountInfo, err error)
	AllOrders(symbol string, startTime, endTime int64, limit int) ([]model.UserOrder, error)
	AvgPrice(symbol string) (model.AvgPrice, error)
	Depth(symbol string, limit int) (model.Order, error)
	ExchangeInfo() (model.ExchangeInfo, error)
	Order(symbol string, side OrderSide, orderType OrderType, params OrderParams) (model.OrderResponse, error)
	OrderTest(symbol string, side OrderSide, orderType OrderType, params OrderParams) (model.OrderResponse, error)
	Ticker24h(symbol string) ([]model.TickerStatistics, error)
	TickerPrice(symbol string) ([]model.Price, error)

	Streamer() Streamer
}

type api struct {
	cfg          APIConfig
	checker      weightChecker
	logger       Logger
	streamer     Streamer
	exchangeInfo *model.ExchangeInfo
}

func (a *api) exchangeInfoValid() bool {
	if a.exchangeInfo == nil {
		return false
	}
	return true
}

func (a *api) Streamer() Streamer {
	return a.streamer
}

func NewAPI(cfg APIConfig, logger Logger) (API, error) {
	if logger == nil {
		return nil, errors.New("api expects an instantiated logger")
	}
	a := &api{
		cfg: cfg,
		checker: weightChecker{
			allowed: true,
			weight:  0,
		},
		logger: logger,
	}

	a.streamer = newStreamer(a, logger)

	return a, nil
}
