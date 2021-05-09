[![Build Status](https://www.travis-ci.com/jaztec/go-binance.svg?branch=main)](https://www.travis-ci.com/jaztec/go-binance)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/jaztec/go-binance/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaztec/go-binance)](https://goreportcard.com/report/github.com/jaztec/go-binance)
[![codecov](https://codecov.io/gh/jaztec/go-binance/branch/main/graph/badge.svg?token=VYL719M4RA)](https://codecov.io/gh/jaztec/go-binance)

# Go Binance API SDK

## Goal 

This library aims to provide some clean endpoints allowing a user to interact seamlessly and easily with the 
Binance V3 REST and Websocket API's.

### Examples

You can find examples inside the `examples` directory in the repository.

### Interfaces

The SDK consists of a couple of interfaces. The most important are the `API` and `Streamer`. These interfaces expose
the fundamentals to talk with the Binance API. These interfaces however are also as clean as possible.

To actually make use of the implemented function calls the `APICaller` and `StreamCaller` interfaces have been created.

### Missing methods

The SDK for the moment only exposes a couple of endpoints used for my own applications. However, you can easily use the
`Request` method to formulate your own request or implement a function and send in a pull request.

```go

package main

import (
	"fmt"
	"github.com/jaztec/go-binance"
	"net/http"
)

func main() {
	a, err := binance.NewAPI(binance.APIConfig{})
	if err != nil {
		panic(err)
	}

	p := binance.NewParameters(1)
	p.Set("symbol", "ETHBTC")

	res, err := a.Request(http.MethodGet, "/api/v3/avgPrice", p)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

```