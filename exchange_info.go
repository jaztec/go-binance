package binance

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaztec/go-binance/model"
)

const exchangeInfoPath = "/api/v3/exchangeInfo"

func (a *api) ExchangeInfo() (ei model.ExchangeInfo, err error) {
	body, err := a.Request(http.MethodGet, exchangeInfoPath, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ei)
	if err != nil {
		return ei, fmt.Errorf("encountered error while unmarshaling '%s' into model.ExchangeInfo", body)
	}

	// update internal exchange information as well
	a.exchangeInfo = &ei

	return
}
