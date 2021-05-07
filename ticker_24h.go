package binance

import (
	"encoding/json"
	"net/http"

	"github.com/jaztec/go-binance/model"
)

const ticker24hPath = "/api/v3/ticker/24hr"

func (a *api) Ticker24h(symbol string) (ts []model.TickerStatistics, err error) {
	var q Parameters
	if symbol != "" {
		q = NewParameters(1)
		q.Set("symbol", symbol)
	}

	body, err := a.doRequest(http.MethodGet, ticker24hPath, q)
	if err != nil {
		return nil, err
	}

	if symbol != "" {
		var p model.TickerStatistics
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		ts = append(ts, p)
	} else {
		err = json.Unmarshal(body, &ts)
		if err != nil {
			return nil, err
		}
	}

	return
}
