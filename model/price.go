package model

import "strconv"

// Price holds information about a symbols price
type Price struct {
	Symbol   string `json:"symbol"`
	RawPrice string `json:"price"`
}

// Price function returns the string price value as a float
func (p Price) Price() float64 {
	f, err := strconv.ParseFloat(p.RawPrice, 32)
	if err != nil {
		// we do not handle this error
		return 0
	}
	return f
}

// PriceCollection is a alias for a slice of prices with Sort interface
// implemented
type PriceCollection []Price

func (pc PriceCollection) Len() int           { return len(pc) }
func (pc PriceCollection) Swap(i, j int)      { pc[i], pc[j] = pc[j], pc[i] }
func (pc PriceCollection) Less(i, j int) bool { return pc[i].Price() < pc[j].Price() }

// AvgPrice holds the average price over the provided minutes
type AvgPrice struct {
	Mins  int    `json:"mins"`
	Price string `json:"price"`
}
