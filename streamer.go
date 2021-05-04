package binance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
)

type Streamer interface {
	AccountData(ctx context.Context) (<-chan model.StreamData, error)
	KlineData(ctx context.Context, symbols []string, interval string) (<-chan model.KlineData, error)
	AllTicker(ctx context.Context) (chan []model.Ticker, error)
}

type streamer struct {
	api     *api
	logger  Logger
	streams []*stream
}

func (s *streamer) keepAlive(ctx context.Context, path string, interval time.Duration) {
	go func(ctx context.Context, interval time.Duration) {
		tC := time.Tick(interval)
		for {
			select {
			case <-tC:
				_, _ = s.api.doRequest(http.MethodPut, path, nil)
			case <-ctx.Done():
				return
			}

		}
	}(ctx, interval)
}

func (s *streamer) stream(ctx context.Context) (*stream, error) {
	if len(s.streams) > 0 {
		c := s.streams[0]
		return c, nil
	}

	conn, err := s.conn()
	if err != nil {
		return nil, err
	}

	st := &stream{
		conn:        conn,
		channels:    make(channelList, 0, 5),
		writes:      make(chan []byte, 5),
		subscribers: make(subscriberMap),
		logger:      s.logger,
		closed:      make(chan struct{}),
	}
	s.streams = append(s.streams, st)

	go st.readPump()
	go st.writePump(ctx)
	go s.monitor(ctx, st)

	return st, nil
}

func (s *streamer) conn() (*websocket.Conn, error) {
	fullURI := fmt.Sprintf("%s/stream", s.api.cfg.BaseStreamURI)
	d := &websocket.Dialer{}
	_ = s.logger.Log("msg", "starting stream", "uri", fullURI)
	conn, _, err := d.Dial(fullURI, nil)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (s *streamer) monitor(ctx context.Context, st *stream) {
	for {
		ch := st.closed
		select {
		case nch := <-st.newClosed:
			fmt.Println("got new closed")
			ch = nch
		case <-ch:
			conn, err := s.conn()
			if err != nil {
				_ = s.logger.Log("streamer", "monitor", "error creating new connection", err.Error())
				return
			}
			err = st.reset(ctx, conn)
			if err != nil {
				_ = s.logger.Log("streamer", "monitor", "error resetting", err.Error())
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func newStreamer(a *api, logger log.Logger) Streamer {
	return &streamer{
		api:    a,
		logger: logger,
	}
}
