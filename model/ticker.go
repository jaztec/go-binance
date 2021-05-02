package model

type Ticker struct {
	MessageType            string `json:"e"`
	EventTime              int64  `json:"E"`
	Symbol                 string `json:"s"`
	PriceChange            string `json:"p"`
	PriceChangePercent     string `json:"P"`
	WeightedAveragePrice   string `json:"w"`
	FirstPrice             string `json:"x"`
	LastPrice              string `json:"c"`
	LastQuantity           string `json:"Q"`
	BestBidPrice           string `json:"b"`
	BestBidQuantity        string `json:"B"`
	BestAskPrice           string `json:"a"`
	BestAskQuantity        string `json:"A"`
	OpenPrice              string `json:"o"`
	HighPrice              string `json:"h"`
	LowPrice               string `json:"l"`
	TotalTradedBaseVolume  string `json:"v"`
	TotalTradedQuoteVolume string `json:"q"`
	StatisticsOpenTime     int64  `json:"O"`
	StatisticsCloseTime    int64  `json:"C"`
	FirstTradeID           int    `json:"F"`
	LastTradeID            int    `json:"L"`
	NumberOfTrades         int    `json:"n"`
}
