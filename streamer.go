package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
	"gitlab.jaztec.info/checkers/checkers/services/binance/model"
)

// Streamer defines functions that are available in the Binance Websocket API.
type Streamer interface {
	// UserDataStream updates when user account changes have occured
	UserDataStream(ctx context.Context) (<-chan model.StreamData, error)
	// Kline data for a list of tokens
	Kline(ctx context.Context, symbols []string, interval string) (<-chan model.KlineData, error)
	// TickerArr changes to prices from the ticker API
	TickerArr(ctx context.Context) (chan []model.Ticker, error)
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
	_ = s.logger.Log("stream", "request", "current", len(s.streams))
	if len(s.streams) > 0 {
		_ = s.logger.Log("stream", "request", "returning", "existing")
		c := s.streams[0]
		return c, nil
	}
	_ = s.logger.Log("stream", "request", "returning", "new")

	conn, err := s.conn()
	if err != nil {
		return nil, err
	}

	st := &stream{
		id:          uniuri.New(),
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

func (s *streamer) resetStream(ctx context.Context, st *stream) error {
	// remove stream from list
	s.removeStream(st.id)

	// get a new stream running
	nst, err := s.stream(ctx)
	if err != nil {
		return err
	}

	// copy values from last stream to the new one
	nst.lastID = st.lastID + 1
	for k, v := range st.subscribers {
		nst.subscribers[k] = v
	}
	copy(nst.channels, st.channels)

	_ = s.logger.Log("resetting", strings.Join(nst.channels, ","))
	// subscribe to the channels the old stream was subscribed to
	// we purposely don't use subscribe method to keep subscriber map intact
	if len(nst.channels) > 0 {
		msg := SubscribeMessage{
			Method: Subscribe,
			Params: nst.channels,
			ID:     nst.lastID,
		}

		b, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		nst.writes <- b
	}

	return nil
}

func (s *streamer) removeStream(id string) {
	for i, st := range s.streams {
		if st.id == id {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
		}
	}
}

func (s *streamer) monitor(ctx context.Context, st *stream) {
	for {
		select {
		case <-st.closed:
			err := s.resetStream(ctx, st)
			if err != nil {
				_ = s.logger.Log("streamer", "monitor", "error resetting", err.Error())
			}
			return
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
