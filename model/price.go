package model

import "strconv"

type Price struct {
	Symbol   string `json:"symbol"`
	RawPrice string `json:"price"`
}

func (bp Price) Price() float64 {
	f, err := strconv.ParseFloat(bp.RawPrice, 32)
	if err != nil {
		panic(err)
	}
	return f
}
