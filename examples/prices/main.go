package main

import (
	"fmt"

	"github.com/jaztec/go-binance"
)

func main() {
	b, err := binance.NewAPICaller(binance.APIConfig{})
	if err != nil {
		panic(err)
	}

	p, err := b.AvgPrice("ETHBTC")
	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}
