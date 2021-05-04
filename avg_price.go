package binance

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
)

const avgPricePath = "/api/v3/avgPrice"

func (a *api) AvgPrice(symbol string) (ap model.AvgPrice, err error) {
	q := NewParameters(1)
	q.Set("symbol", symbol)

	body, err := a.doRequest(http.MethodGet, avgPricePath, q)
	if err != nil {
		return ap, err
	}

	err = json.Unmarshal(body, &ap)
	if err != nil {
		return ap, fmt.Errorf("encountered error while unmarshaling '%s' into model.AvgPrice", body)
	}

	return
}
