package model

import "strconv"

type BinancePrice struct {
	Symbol   string `json:"symbol"`
	RawPrice string `json:"price"`
}

func (bp BinancePrice) Price() float64 {
	f, err := strconv.ParseFloat(bp.RawPrice, 32)
	if err != nil {
		panic(err)
	}
	return f
}
