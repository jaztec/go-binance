package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaztec/go-binance"
)

type logger struct{}

func (*logger) Log(vals ...interface{}) error {
	log.Println(vals...)
	return nil
}

func main() {
	key := flag.String("key", "false", "The API key provided by Binance")
	secret := flag.String("secret", "false", "The secret that belongs to the API key")
	flag.Parse()

	b, err := binance.NewAPI(binance.APIConfig{
		Key:    *key,
		Secret: *secret,
	}, &logger{})
	if err != nil {
		panic(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	// use lowercase for streams
	ch, err := b.Stream().(binance.StreamCaller).UserDataStream(ctx)
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
