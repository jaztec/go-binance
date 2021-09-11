package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jaztec/go-binance/model"
)

const (
	myTradesPath = "/api/v3/myTrades"
)

func init() {
	requireSignature(myTradesPath)
}

// MyTrades returns trades performed by a user in a time window
func (a *api) MyTrades(symbol string, startTime, endTime int64, limit int) (t []model.UserTrade, err error) {
	if symbol == "" {
		return t, NoSymbolProvided
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

	body, err := a.Request(http.MethodGet, myTradesPath, q)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, fmt.Errorf("encountered error while unmarshaling '%s' into model.Trade", body)
	}

	return t, nil
}
