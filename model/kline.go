package model

// Kline /Candlestick data for a symbol
type Kline struct {
	StartTime                int64  `json:"t"`
	CloseTime                int64  `json:"T"`
	Symbol                   string `json:"s"`
	Interval                 string `json:"i"`
	FirstTradeID             int    `json:"f"`
	LastTradeID              int    `json:"L"`
	OpenPrice                string `json:"o"`
	ClosePrice               string `json:"c"`
	HighPrice                string `json:"h"`
	LowPrice                 string `json:"l"`
	BaseAssetVolume          string `json:"v"`
	NumberOfTrades           int    `json:"n"`
	Closed                   bool   `json:"x"`
	QuoteAssetVolume         string `json:"q"`
	TakerBuyBaseAssetVolume  string `json:"V"`
	TakerBuyQuoteAssetVolume string `json:"Q"`
}

// KlineData wrapper for a symbol
type KlineData struct {
	Type   string `json:"e"`
	Time   int64  `json:"E"`
	Symbol string `json:"s"`
	Kline  Kline  `json:"k"`
}
