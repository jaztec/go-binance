package model

// SymbolInfo provides data about a symbol
type SymbolInfo struct {
	Symbol                     string        `json:"symbol"`
	Status                     string        `json:"status"`
	BaseAsset                  string        `json:"baseAsset"`
	BaseAssetPrecision         int           `json:"baseAssetPrecision"`
	QuoteAsset                 string        `json:"quoteAsset"`
	QuotePrecision             int           `json:"quotePrecision"`
	QuoteAssetPrecision        int           `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int           `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int           `json:"quoteCommissionPrecision"`
	OrderTypes                 []string      `json:"orderTypes"`
	IcebergAllowed             bool          `json:"icebergAllowed"`
	OcoAllowed                 bool          `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool          `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool          `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool          `json:"isMarginTradingAllowed"`
	Filters                    []interface{} `json:"filters"`
	Permissions                []string      `json:"permissions"`
}

// ExchangeInfo holds data about the exchange and internals
type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int64         `json:"serverTime"`
	RateLimits      []interface{} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []SymbolInfo  `json:"symbols"`
}
