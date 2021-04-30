package binance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/gorilla/websocket"
)

const (
	BaseStreamURI = "wss://stream.binance.com:9443"

	pongWait   = 5 * time.Minute
	pingPeriod = 5 * time.Minute
)

func (a *api) keepAlive(ctx context.Context, path string, interval time.Duration) {
	uri := fmt.Sprintf("%s/%s", a.cfg.BaseStreamURI, path)
	go func(ctx context.Context, uri string, interval time.Duration) {
		tC := time.Tick(interval)
		for {
			select {
			case <-tC:
				_, _ = a.doRequest(http.MethodPut, uri, nil)
			case <-ctx.Done():
				return
			}

		}
	}(ctx, uri, interval)
}

func (a *api) stream(ctx context.Context, p Parameters) (chan []byte, chan []byte, error) {
	fullURI := fmt.Sprintf("%s/stream?%s", a.cfg.BaseStreamURI, p.Encode())
	d := &websocket.Dialer{}
	_ = a.logger.Log("starting stream", fullURI)
	conn, _, err := d.Dial(fullURI, nil)
	if err != nil {
		return nil, nil, err
	}

	var reads = make(chan []byte, 5)
	var writes = make(chan []byte, 5)
	go readPump(a.logger, conn, reads)
	go writePump(a.logger, ctx, conn, writes)

	return reads, writes, nil
}

func readPump(logger log.Logger, conn *websocket.Conn, ch chan []byte) {
	defer func() {
		err := conn.Close()
		if err != nil {
			_ = logger.Log("close", "readPump", "error", err)
		}
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		ch <- msg
	}
}

func writePump(logger log.Logger, ctx context.Context, conn *websocket.Conn, ch chan []byte) {
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	t := time.NewTicker(pingPeriod)
	defer t.Stop()

	for {
		select {
		case msg := <-ch:
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				_ = logger.Log("write", "writePump", "error", err)
				return
			}
		case _ = <-ctx.Done():
			_ = logger.Log("writePump", "close signal")
			return
		case <-t.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
