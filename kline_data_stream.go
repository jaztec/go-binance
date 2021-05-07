package binance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jaztec/go-binance/model"
)

func (s *streamer) Kline(ctx context.Context, symbols []string, interval string) (<-chan model.KlineData, error) {
	params := make([]string, 0, len(symbols))
	for _, s := range symbols {
		params = append(params, fmt.Sprintf("%s@kline_%s", s, interval))
	}
	st, err := s.stream(ctx)
	if err != nil {
		return nil, err
	}

	reads, err := st.subscribe(params)
	if err != nil {
		return nil, err
	}

	readStream := make(chan model.KlineData)
	go func() {
		for {
			select {
			case msg := <-reads:
				var k model.KlineData
				_ = json.Unmarshal(msg.Data, &k)
				readStream <- k
			case <-ctx.Done():
				return
			}
		}
	}()

	return readStream, nil
}
