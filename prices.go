package binance

import (
	"encoding/json"
	model2 "gitlab.jaztec.info/checkers/checkers/services/binance/model"
	"net/http"
)

const pricesPath = "/api/v3/ticker/price"

func (a *api) Prices(symbol string) ([]model2.BinancePrice, error) {
	var q Parameters
	if symbol != "" {
		q = Parameters{}
		q.Set("symbol", symbol)
	}

	body, err := a.doRequest(http.MethodGet, pricesPath, q)
	if err != nil {
		return nil, err
	}

	var list []model2.BinancePrice
	if symbol != "" {
		var p model2.BinancePrice
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		list = []model2.BinancePrice{p}
	} else {
		err = json.Unmarshal(body, &list)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
