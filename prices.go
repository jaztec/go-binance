package binance

import (
	"encoding/json"
	"net/http"

	"github.com/jaztec/go-binance/model"
)

const pricesPath = "/api/v3/ticker/price"

func (a *api) TickerPrice(symbol string) ([]model.Price, error) {
	var q Parameters
	if symbol != "" {
		q = NewParameters(1)
		q.Set("symbol", symbol)
	}

	body, err := a.doRequest(http.MethodGet, pricesPath, q)
	if err != nil {
		return nil, err
	}

	var list []model.Price
	if symbol != "" {
		var p model.Price
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		list = []model.Price{p}
	} else {
		err = json.Unmarshal(body, &list)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
