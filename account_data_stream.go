package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
)

const userDataStreamPath = "/api/v3/userDataStream"

func (s *streamer) UserDataStream(ctx context.Context) (<-chan model.StreamData, error) {
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

	return reads, nil
}
