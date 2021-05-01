package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"

	"github.com/go-kit/kit/log"

	"github.com/gorilla/websocket"
)

type SubscribeType string

const (
	BaseStreamURI = "wss://stream.binance.com:9443"

	pongWait   = 5 * time.Minute
	pingPeriod = 5 * time.Minute

	Subscribe   SubscribeType = "SUBSCRIBE"
	Unsubscribe SubscribeType = "UNSUBSCRIBE"
)

type subscriberMap map[string][]chan model.StreamData

type StreamerConfig struct {
	API           API
	BaseStreamURI string
}

type Streamer interface {
	AccountData(ctx context.Context) (<-chan model.StreamData, error)
	KlineData(ctx context.Context, symbols []string, interval string) (<-chan model.KlineData, error)
}

type streamer struct {
	api     *api
	logger  log.Logger
	streams []*stream
}

type subscribeMessage struct {
	Method SubscribeType `json:"method"`
	Params []string      `json:"params"`
	ID     uint64        `json:"lastID"`
}

type stream struct {
	conn          *websocket.Conn
	subscriptions uint16
	writes        chan []byte
	lastID        uint64
	subscribers   subscriberMap
}

func (s *stream) unsubscribe(params []string) error {
	atomic.AddUint64(&s.lastID, 1)

	msg := subscribeMessage{
		Method: Unsubscribe,
		Params: params,
		ID:     s.lastID,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	s.writes <- b

	for _, param := range params {
		// TODO This should be composed better, only close actual channels that are requesting unsubscribe
		if list, ok := s.subscribers[param]; ok {
			for _, ch := range list {
				close(ch)
			}
			delete(s.subscribers, param)
		}
	}

	return nil
}

func (s *stream) subscribe(params []string) (<-chan model.StreamData, error) {
	atomic.AddUint64(&s.lastID, 1)

	msg := subscribeMessage{
		Method: Subscribe,
		Params: params,
		ID:     s.lastID,
	}

	reads := make(chan model.StreamData, 5)
	for _, param := range params {
		if _, ok := s.subscribers[param]; !ok {
			s.subscribers[param] = make([]chan model.StreamData, 0, 1)
		}
		s.subscribers[param] = append(s.subscribers[param], reads)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	s.writes <- b

	return reads, nil
}

func (s *stream) readPump(logger log.Logger) {
	defer func() {
		err := s.conn.Close()
		if err != nil {
			_ = logger.Log("close", "readPump", "error", err)
		}
	}()
	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			return
		}

		var sd model.StreamData
		if err = json.Unmarshal(msg, &sd); err != nil {
			_ = logger.Log("read", "error", "msg", err.Error())
			continue
		}

		list, ok := s.subscribers[sd.Stream]
		if !ok {
			continue
		}
		for _, ch := range list {
			ch <- sd
		}
	}
}

func (s *stream) writePump(ctx context.Context, logger log.Logger) {
	_ = s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error {
		return s.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	t := time.NewTicker(pingPeriod)
	defer t.Stop()

	for {
		select {
		case msg := <-s.writes:
			if err := s.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				_ = logger.Log("write", "writePump", "error", err)
				return
			}
		case _ = <-ctx.Done():
			_ = logger.Log("writePump", "close signal")
			return
		case <-t.C:
			if err := s.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (s *streamer) keepAlive(ctx context.Context, path string, interval time.Duration) {
	uri := fmt.Sprintf("%s/%s", s.api.cfg.BaseStreamURI, path)
	go func(ctx context.Context, uri string, interval time.Duration) {
		tC := time.Tick(interval)
		for {
			select {
			case <-tC:
				_, _ = s.api.doRequest(http.MethodPut, uri, nil)
			case <-ctx.Done():
				return
			}

		}
	}(ctx, uri, interval)
}

func (s *streamer) stream(ctx context.Context) (*stream, error) {
	if len(s.streams) > 0 {
		c := s.streams[0]
		return c, nil
	}

	fullURI := fmt.Sprintf("%s/stream", s.api.cfg.BaseStreamURI)
	d := &websocket.Dialer{}
	_ = s.logger.Log("msg", "starting stream", "uri", fullURI)
	conn, _, err := d.Dial(fullURI, nil)
	if err != nil {
		return nil, err
	}

	st := &stream{
		conn:          conn,
		subscriptions: 0,
		writes:        make(chan []byte, 5),
		subscribers:   make(subscriberMap),
	}
	s.streams = append(s.streams, st)

	go st.readPump(s.logger)
	go st.writePump(ctx, s.logger)

	return st, nil
}

func newStreamer(a *api, logger log.Logger) Streamer {
	return &streamer{
		api:    a,
		logger: logger,
	}
}
