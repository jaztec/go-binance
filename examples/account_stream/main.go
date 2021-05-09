package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaztec/go-binance"
)

func main() {
	key := flag.String("key", "false", "The API key provided by Binance")
	secret := flag.String("secret", "false", "The secret that belongs to the API key")
	flag.Parse()

	b, err := binance.NewAPICaller(binance.APIConfig{
		Key:    *key,
		Secret: *secret,
	})
	if err != nil {
		panic(err)
	}

	caller, ok := b.(binance.APICaller)
	if !ok {
		panic("Not an APICaller")
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ch, err := caller.StreamCaller().UserDataStream(ctx)
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
