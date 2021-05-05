package binance

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"gitlab.jaztec.info/checkers/checkers/services/binance/model"

	"github.com/gorilla/websocket"
)

type SubscribeType string

const (
	BaseStreamURI = "wss://stream.binance.com:9443"

	pongPeriod = 2 * time.Minute

	Subscribe   SubscribeType = "SUBSCRIBE"
	Unsubscribe SubscribeType = "UNSUBSCRIBE"
)

type subscriberMap map[string][]chan model.StreamData
type channelList []string

func (cl channelList) Len() int           { return len(cl) }
func (cl channelList) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
func (cl channelList) Less(i, j int) bool { return cl[i] < cl[j] }
func (cl channelList) IndexOf(s string) int {
	for n, el := range cl {
		if el == s {
			return n
		}
	}
	return -1
}

type StreamerConfig struct {
	API           API
	BaseStreamURI string
}

// SubscribeMessage is a representation of the Binance subscribe and unsubscribe
// messages data structure.
type SubscribeMessage struct {
	Method SubscribeType `json:"method"`
	Params []string      `json:"params"`
	ID     uint64        `json:"id"`
}

type stream struct {
	conn        *websocket.Conn
	channels    channelList
	writes      chan []byte
	lastID      uint64
	subscribers subscriberMap
	logger      Logger
	closed      chan struct{}
	newClosed   chan chan struct{}
}

// reset allows the user to provide a new connection to the
// stream that will continue work where it stopped. The Binance
// API does a standard disconnect after 24h.
func (s *stream) reset(ctx context.Context, conn *websocket.Conn) error {
	newC := make(chan struct{})
	s.newClosed <- newC
	s.closed = newC

	_ = s.conn.Close()
	s.conn = conn
	go s.readPump()
	go s.writePump(ctx)

	msg := SubscribeMessage{
		Method: Subscribe,
		Params: s.channels,
		ID:     s.lastID,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	s.writes <- b

	return nil
}

func (s *stream) unsubscribe(params []string) error {
	atomic.AddUint64(&s.lastID, 1)

	msg := SubscribeMessage{
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
		if list, ok := s.subscribers[param]; ok {
			for _, ch := range list {
				close(ch)
			}
			delete(s.subscribers, param)
		}
		if n := s.channels.IndexOf(param); n > -1 {
			// remove channel but keep order intact
			s.channels = append(s.channels[:n], s.channels[n+1:]...)
		}
	}

	return nil
}

func (s *stream) subscribe(params []string) (<-chan model.StreamData, error) {
	atomic.AddUint64(&s.lastID, 1)
	_ = s.logger.Log("subscribe", strings.Join(params, ", "))

	newParams := make([]string, 0, len(params))
	reads := make(chan model.StreamData, 5)
	for _, param := range params {
		if _, ok := s.subscribers[param]; !ok {
			s.subscribers[param] = make([]chan model.StreamData, 0, 1)
			newParams = append(newParams, param)
		}
		s.subscribers[param] = append(s.subscribers[param], reads)
	}

	// keep track of channels we connect on
	s.channels = append(s.channels, params...)
	sort.Sort(s.channels)

	if len(newParams) > 0 {
		msg := SubscribeMessage{
			Method: Subscribe,
			Params: newParams,
			ID:     s.lastID,
		}

		b, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		s.writes <- b
	}

	return reads, nil
}

func (s *stream) readPump() {
	defer func() {
		err := s.conn.Close()
		if err != nil {
			_ = s.logger.Log("close", "readPump", "error", err)
		}
	}()
	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			_ = s.logger.Log("method", "readPump", "error", err.Error())
			close(s.closed)
			return
		}

		var sd model.StreamData
		if err = json.Unmarshal(msg, &sd); err != nil {
			_ = s.logger.Log("read", "error", "msg", err.Error())
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

func (s *stream) writePump(ctx context.Context) {
	t := time.NewTicker(pongPeriod)
	defer t.Stop()

	for {
		select {
		case msg := <-s.writes:
			if err := s.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				_ = s.logger.Log("write", string(msg), "error", err)
				return
			}
		case <-s.closed:
			// when top stream closes we exit too, reset will start new procedures
			return
		case _ = <-ctx.Done():
			_ = s.logger.Log("writePump", "close signal")
			return
		case <-t.C:
			if err := s.conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
