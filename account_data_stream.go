package binance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jaztec/go-binance/model"
)

const userDataStreamPath = "/api/v3/userDataStream"

func (s *streamer) UserDataStream(ctx context.Context) (<-chan model.UserAccountUpdate, error) {
	res, err := s.api.doRequest(http.MethodPost, userDataStreamPath, nil)
	if err != nil {
		return nil, err
	}

	var key model.ListenKey
	err = json.Unmarshal(res, &key)
	if err != nil {
		return nil, err
	}
	st, err := s.stream(ctx)
	if err != nil {
		return nil, err
	}

	reads, err := st.subscribe([]string{key.ListenKey})
	if err != nil {
		return nil, err
	}

	p := NewParameters(1)
	p.Set("listenKey", key.ListenKey)
	path := fmt.Sprintf("%s?%s", userDataStreamPath, p.Encode())
	s.keepAlive(ctx, path, time.Minute*30)

	ch := make(chan model.UserAccountUpdate, 5)
	go func() {
		for {
			select {
			case msg := <-reads:
				m, err := findModel(msg.Data)
				if err != nil {
					_ = s.logger.Log("stream", "account_data_stream", "error", err.Error())
					continue
				}
				ch <- m
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

func findModel(in []byte) (model.UserAccountUpdate, error) {
	cmp := func(b []byte, update model.AccountUpdateType) bool {
		return bytes.Contains(b, []byte(fmt.Sprintf("\"e\":\"%s\"", string(update))))
	}
	if cmp(in, model.OutboundAccountPositionType) {
		out := model.OutboundAccountPosition{}
		err := json.Unmarshal(in, &out)
		return out, err
	}
	if cmp(in, model.BalanceUpdateType) {
		out := model.BalanceUpdate{}
		err := json.Unmarshal(in, &out)
		return out, err
	}
	if cmp(in, model.ExecutionReportType) {
		out := model.ExecutionReport{}
		err := json.Unmarshal(in, &out)
		return out, err
	}
	return nil, errors.New("no AccountUpdateType matched")
}
