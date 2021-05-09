package main

import (
	"flag"
	"fmt"
	"log"

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

	b, err := binance.NewAPICaller(binance.APIConfig{
		Key:    *key,
		Secret: *secret,
	}, &logger{})
	if err != nil {
		panic(err)
	}

	caller, ok := b.(binance.APICaller)
	if !ok {
		panic("Not an APICaller")
	}

	a, err := caller.Account()
	if err != nil {
		panic(err)
	}

	fmt.Println(a)
}
