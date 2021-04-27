package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
	"net/http"
	"time"
)

const userDataStreamPath = "/api/v3/userDataStream"

func (a *api) StartUserDataStream(ctx context.Context) error {
	res, err := a.doRequest(http.MethodPost, userDataStreamPath, nil)
	if err != nil {
		return err
	}

	var key model.ListenKey
	err = json.Unmarshal(res, &key)
	if err != nil {
		return err
	}

	p := Parameters{}
	p.Set("streams", key.ListenKey)
	reads, _, err := a.stream(ctx, p)
	if err != nil {
		return err
	}
	a.keepAlive(ctx, p.Encode(), time.Minute*30)

	for {
		select {
		case msg := <-reads:
			fmt.Println(string(msg))
		case <-ctx.Done():
			return nil
		}
	}
}
