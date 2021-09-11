package model

// UserTrade shows the data about trade that has been performed
type UserTrade struct {
	Symbol          string `json:"symbol"`
	ID              int    `json:"id"`
	OrderID         int    `json:"orderId"`
	OrderListID     int    `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}
