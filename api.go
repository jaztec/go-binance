package binance

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jaztec/go-binance/model"
)

const (
	// BaseAPIURI contains the default Binance API location
	BaseAPIURI = "https://api.binance.com"
	// APIKeyHeaderName is the header name Binance API expects the API token to be
	APIKeyHeaderName = "X-MBX-APIKEY"
)

type weightChecker struct {
	allowed bool
	weight  int
}

func (wc *weightChecker) deactivate(seconds int) {
	wc.allowed = false
	go func(seconds int) {
		tC := time.Tick(time.Second * time.Duration(seconds))
		<-tC
		wc.allowed = true
	}(seconds)
}

func (wc *weightChecker) checkResponse(response *http.Response) error {
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
		return TooMuchCalls
	}
	if response.StatusCode == http.StatusTeapot {
		fn(response)
		return Blocked
	}
	return nil
}

// APIConfig lets us setup the values required by the adapter
type APIConfig struct {
	// API Key to be used in the communication
	Key string
	// Secret attached to the API Key
	Secret string
	// BaseURI of the API. Will default to BaseAPIURI
	BaseURI string
	// BaseStreamURI for the websocket API. Will default to BaseStreamURI
	BaseStreamURI string
}

// API interface exposes all the available (implemented) endpoints to the Binance REST API. The Streamer can be
// used to start streams onto websockets.
type API interface {
	// Request makes an actual request to the Binance API. It will check the API rate limits
	// and deactivate the API for the period defined by the API when in violation. The function
	// will return the raw body of the result on success or an error on failure.
	Request(method, path string, params Parameters) ([]byte, error)

	// Stream returns a Streamer
	Stream() Streamer
}

// APICaller exposes readily implemented calls for the Binance REST API
type APICaller interface {
	API

	// Account information
	Account() (ai model.AccountInfo, err error)
	// AllOrders for a symbol from the user account
	AllOrders(symbol string, startTime, endTime int64, limit int) ([]model.UserOrder, error)
	// AvgPrice of a symbol
	AvgPrice(symbol string) (model.AvgPrice, error)
	// Depth endpoint on Binance API
	Depth(symbol string, limit int) (model.Orders, error)
	// ExchangeInfo as set by Binance
	ExchangeInfo() (model.ExchangeInfo, error)
	// Order to put into the Binance system
	Order(symbol string, side model.OrderSide, orderType model.OrderType, params OrderParams) (model.OrderResponse, error)
	// OrderTest will validate an order but not put it into the system
	OrderTest(symbol string, side model.OrderSide, orderType model.OrderType, params OrderParams) (model.OrderResponse, error)
	// Ticker24h will retrieve information about a symbol for a 24H period. WARNING, heavy penalty
	Ticker24h(symbol string) ([]model.TickerStatistics, error)
	// TickerPrice returns price information about a symbol
	TickerPrice(symbol string) ([]model.Price, error)

	// StreamCaller returns a stream with readily implemented functions
	StreamCaller() StreamCaller
}

type api struct {
	cfg          APIConfig
	checker      *weightChecker
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

func (a *api) Stream() Streamer {
	return a.streamer
}

func (a *api) StreamCaller() StreamCaller {
	return a.streamer.(StreamCaller)
}

// NewAPI will return a new API interface fully setup to run
func NewAPI(cfg APIConfig, logger Logger) (API, error) {
	if logger == nil {
		return nil, errors.New("api expects an instantiated logger")
	}
	if cfg.BaseURI == "" {
		cfg.BaseURI = BaseAPIURI
	}
	if cfg.BaseStreamURI == "" {
		cfg.BaseStreamURI = BaseStreamURI
	}
	a := &api{
		cfg: cfg,
		checker: &weightChecker{
			allowed: true,
			weight:  0,
		},
		logger: logger,
	}

	a.streamer = newStreamer(a, logger)

	return a, nil
}
