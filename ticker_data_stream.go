package binance

import (
	"context"
	"encoding/json"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
)

func (s *streamer) TickerArr(ctx context.Context) (chan []model.Ticker, error) {
	st, err := s.stream(ctx)
	if err != nil {
		return nil, err
	}

	ch, err := st.subscribe([]string{"!ticker@arr"})
	if err != nil {
		return nil, err
	}

	reads := make(chan []model.Ticker)

	go func() {
		for {
			select {
			case msg := <-ch:
				var t []model.Ticker
				if err := json.Unmarshal(msg.Data, &t); err != nil {
					_ = s.logger.Log("method", "TickerArr", "error", err.Error())
					continue
				}

				reads <- t
			case <-ctx.Done():
				return
			}
		}
	}()

	return reads, nil
}
