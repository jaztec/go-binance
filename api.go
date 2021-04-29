package binance

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
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
	Prices(symbol string) ([]model.Price, error)
	OrderBook(symbol string, limit int) (model.Order, error)
	UserOrderBook(symbol string, startTime, endTime int64, limit int) ([]model.UserOrder, error)
	StartUserDataStream(ctx context.Context) error
	UserAccount() (ai model.AccountInfo, err error)
}

//symbol	STRING	YES
//orderId	LONG	NO
//startTime	LONG	NO
//endTime	LONG	NO
//limit	INT	NO	Default 500; max 1000.
//recvWindow	LONG	NO	The value cannot be greater than 60000
//timestamp	LONG	YES

type api struct {
	cfg     APIConfig
	checker weightChecker
	logger  log.Logger
}

func NewAPI(cfg APIConfig, logger log.Logger) API {
	return &api{
		cfg: cfg,
		checker: weightChecker{
			allowed: true,
			weight:  0,
		},
		logger: logger,
	}
}
