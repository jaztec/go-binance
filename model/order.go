package model

// Orders holds information from the depth endpoint of the Binance API
type Orders struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// OrderResponse interface exposes the fields all order response types
// hold for easy passing along inside the program
type OrderResponse interface {
	// Symbol returns the symbol for this order
	Symbol() string
	// OrderID of this order
	OrderID() int
	// OrderListID for this order
	OrderListID() int
	// ClientOrderID for this order
	ClientOrderID() string
	// TransactionTime of the order
	TransactionTime() int64
}

// UserOrder holds data about a single order from the user order book
type UserOrder struct {
	Symbol              string `json:"symbol"`
	OrderID             int    `json:"orderId"`
	OrderListID         int    `json:"orderListId"`
	ClientOrderID       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

// OrderResponseAck holds the fields for the ACK order response
type OrderResponseAck struct {
	Sym          string `json:"symbol"`
	Order        int    `json:"orderId"`
	OrderList    int    `json:"orderListId"`
	ClientOrder  string `json:"clientOrderId"`
	TransactTime int64  `json:"transactTime"`
}

// Symbol returns the symbol for this order
func (r OrderResponseAck) Symbol() string { return r.Sym }

// OrderID of this order
func (r OrderResponseAck) OrderID() int { return r.Order }

// OrderListID for this order
func (r OrderResponseAck) OrderListID() int { return r.OrderList }

// ClientOrderID for this order
func (r OrderResponseAck) ClientOrderID() string { return r.ClientOrder }

// TransactionTime of the order
func (r OrderResponseAck) TransactionTime() int64 { return r.TransactTime }

// OrderResponseResult holds the fields for the RESPONSE order response
type OrderResponseResult struct {
	Sym                 string `json:"symbol"`
	Order               int    `json:"orderId"`
	OrderList           int    `json:"orderListId"`
	ClientOrder         string `json:"clientOrderId"`
	TransactTime        int64  `json:"transactTime"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
}

// Symbol returns the symbol for this order
func (r OrderResponseResult) Symbol() string { return r.Sym }

// OrderID of this order
func (r OrderResponseResult) OrderID() int { return r.Order }

// OrderListID for this order
func (r OrderResponseResult) OrderListID() int { return r.OrderList }

// ClientOrderID for this order
func (r OrderResponseResult) ClientOrderID() string { return r.ClientOrder }

// TransactionTime of the order
func (r OrderResponseResult) TransactionTime() int64 { return r.TransactTime }

// Fill holds details about how a order was filled
type Fill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}

// OrderResponseFull holds the fields for the full order response type
type OrderResponseFull struct {
	Sym                 string `json:"symbol"`
	Order               int    `json:"orderId"`
	OrderList           int    `json:"orderListId"`
	ClientOrder         string `json:"clientOrderId"`
	TransactTime        int64  `json:"transactTime"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	Fills               []Fill `json:"fills"`
}

// Symbol returns the symbol for this order
func (r OrderResponseFull) Symbol() string { return r.Sym }

// OrderID of this order
func (r OrderResponseFull) OrderID() int { return r.Order }

// OrderListID for this order
func (r OrderResponseFull) OrderListID() int { return r.OrderList }

// ClientOrderID for this order
func (r OrderResponseFull) ClientOrderID() string { return r.ClientOrder }

// TransactionTime of the order
func (r OrderResponseFull) TransactionTime() int64 { return r.TransactTime }
