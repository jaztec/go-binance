package binance

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	BaseStreamURI = "wss://stream.binance.com:9443"

	pongWait   = 60 * time.Second
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
	log.Printf("Starting stream on %s", fullURI)
	conn, _, err := d.Dial(fullURI, nil)
	if err != nil {
		return nil, nil, err
	}

	var reads = make(chan []byte, 5)
	var writes = make(chan []byte, 5)
	go pingPong(ctx, conn)
	go readPump(conn, reads)
	go writePump(ctx, conn, writes)

	return reads, writes, nil
}

func pingPong(ctx context.Context, conn *websocket.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	t := time.NewTicker(pingPeriod)
	defer t.Stop()

	for {
		select {
		case _ = <-ctx.Done():
			return
		case <-t.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func readPump(conn *websocket.Conn, ch chan []byte) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Encountered error while closing socket: %s", err)
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

func writePump(ctx context.Context, conn *websocket.Conn, ch chan []byte) {
	for {
		select {
		case msg := <-ch:
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Received error while writing to connection, write pump closing: %s", err)
				return
			}
		case _ = <-ctx.Done():
			log.Printf("Closing socket writePump after close signal")
			return
		}
	}
}
