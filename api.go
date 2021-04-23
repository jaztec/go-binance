package binance

import (
	"encoding/json"
	"gitlab.jaztec.info/checkers/checkers/model"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseApi = "https://api.binance.com"

	tickerPath = "/api/v3/ticker/price"

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

func (wc weightChecker) checkResponse(response *http.Response) *BinanceAPIError {
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
		return &BinanceTooMuchCalls
	}
	if response.StatusCode == http.StatusTeapot {
		fn(response)
		return &BinanceBlocked
	}
	return nil
}

type APIConfig struct {
	Key    string
	Secret string
}

type API interface {
	Prices(symbol string) ([]model.BinancePrice, error)
}

type api struct {
	cfg     APIConfig
	checker weightChecker
}

func (a *api) Prices(symbol string) ([]model.BinancePrice, error) {
	var q url.Values
	if symbol != "" {
		q = url.Values{}
		q.Set("symbol", symbol)
	}

	body, err := a.doRequest(http.MethodGet, tickerPath, q, nil)
	if err != nil {
		return nil, err
	}

	var list []model.BinancePrice
	if symbol != "" {
		var p model.BinancePrice
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		list = []model.BinancePrice{p}
	} else {
		err = json.Unmarshal(body, &list)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func NewAPI(cfg APIConfig) API {
	return &api{
		cfg: cfg,
		checker: weightChecker{
			allowed: true,
			weight:  0,
		},
	}
}
