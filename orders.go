package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
)

const (
	orderBookPath     = "/api/v3/depth"
	userOrderBookPath = "/api/v3/allOrders"
)

func init() {
	requireSignature(userOrderBookPath)
}

func (a *api) AllOrders(symbol string, startTime, endTime int64, limit int) (uo []model.UserOrder, err error) {
	if symbol == "" {
		return uo, NoSymbolProvided
	}
	q := NewParameters(5)
	q.Set("symbol", symbol)
	if startTime != 0 {
		q.Set("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime != 0 {
		q.Set("endTime", strconv.FormatInt(endTime, 10))
	}
	if limit != 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	q.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))

	body, err := a.doRequest(http.MethodGet, userOrderBookPath, q)
	if err != nil {
		return uo, err
	}

	err = json.Unmarshal(body, &uo)
	if err != nil {
		return uo, fmt.Errorf("encountered error while unmarshaling '%s' into model.UserOrder", body)
	}

	return uo, nil
}

func (a *api) Depth(symbol string, limit int) (o model.Order, err error) {
	if symbol == "" {
		return o, NoSymbolProvided
	}
	q := NewParameters(2)
	q.Set("symbol", symbol)
	if limit != 0 {
		q.Set("limit", strconv.Itoa(limit))
	}

	body, err := a.doRequest(http.MethodGet, orderBookPath, q)
	if err != nil {
		return o, err
	}

	err = json.Unmarshal(body, &o)
	if err != nil {
		return o, fmt.Errorf("encountered error while unmarshaling '%s' into model.Order", body)
	}

	return o, nil
}
