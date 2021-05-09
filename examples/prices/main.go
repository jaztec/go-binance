package main

import (
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
	b, err := binance.NewAPICaller(binance.APIConfig{}, &logger{})
	if err != nil {
		panic(err)
	}

	p, err := b.(binance.APICaller).AvgPrice("ETHBTC")
	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}
