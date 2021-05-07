package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaztec/go-binance"
)

type logger struct{}

func (*logger) Log(vals ...interface{}) error {
	log.Println(vals)
	return nil
}

func main() {
	b, err := binance.NewAPI(binance.APIConfig{}, &logger{})
	if err != nil {
		panic(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ch, err := b.Streamer().Kline(ctx, []string{"ETHBTC"}, "1m")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				fmt.Println(msg)
			case <-ctx.Done():
				return
			}
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)

	<-c
}
