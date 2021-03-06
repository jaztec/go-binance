package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jaztec/go-binance/model"
)

const accountPath = "/api/v3/account"

func init() {
	requireSignature(accountPath)
}

func (a *api) Account() (ai model.AccountInfo, err error) {
	q := NewParameters(1)
	q.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))

	body, err := a.Request(http.MethodGet, accountPath, q)
	if err != nil {
		return ai, err
	}

	err = json.Unmarshal(body, &ai)
	if err != nil {
		return ai, fmt.Errorf("encountered error while unmarshaling '%s' into model.AccountInfo", body)
	}

	return ai, nil
}
