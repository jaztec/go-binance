package binance

import (
	"encoding/json"
	"fmt"
	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
	"net/http"
	"strconv"
	"time"
)

const accountPath = "/api/v3/account"

func init() {
	requireSignature(accountPath)
}

func (a *api) UserAccount() (ai model.AccountInfo, err error) {
	q := Parameters{}
	q.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))

	body, err := a.doRequest(http.MethodGet, accountPath, q)
	if err != nil {
		return ai, err
	}

	err = json.Unmarshal(body, &ai)
	if err != nil {
		return ai, fmt.Errorf("encountered error while unmarshaling '%s' into model.AccountInfo", body)
	}

	return ai, nil
}
