package model

import "strconv"

// Balance holds per symbol asset details
type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// Total of a asset
func (b Balance) Total() float64 {
	f, err := strconv.ParseFloat(b.Free, 32)
	if err != nil {
		panic(err)
	}
	l, err := strconv.ParseFloat(b.Locked, 32)
	if err != nil {
		panic(err)
	}
	return f + l
}

// AccountInfo holds full list of account details
type AccountInfo struct {
	MakerCommission  int       `json:"makerCommission"`
	TakerCommission  int       `json:"takerCommission"`
	BuyerCommission  int       `json:"buyerCommission"`
	SellerCommission int       `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       int       `json:"updateTime"`
	AccountType      string    `json:"accountType"`
	Balances         []Balance `json:"balances"`
	Permissions      []string  `json:"permissions"`
}
