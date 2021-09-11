package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaztec/go-binance"
)

func main() {
	b, err := binance.NewAPICaller(binance.APIConfig{})
	if err != nil {
		panic(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ch, err := b.StreamCaller().Kline(ctx, []string{"ETHBTC"}, "1m")
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	<-c
}
