package model

import "strconv"

type Price struct {
	Symbol   string `json:"symbol"`
	RawPrice string `json:"price"`
}

func (p Price) Price() float64 {
	f, err := strconv.ParseFloat(p.RawPrice, 32)
	if err != nil {
		panic(err)
	}
	return f
}

type PriceCollection []Price

func (pc PriceCollection) Len() int           { return len(pc) }
func (pc PriceCollection) Swap(i, j int)      { pc[i], pc[j] = pc[j], pc[i] }
func (pc PriceCollection) Less(i, j int) bool { return pc[i].Price() < pc[j].Price() }

type AvgPrice struct {
	Min   int    `json:"min"`
	Price string `json:"price"`
}
