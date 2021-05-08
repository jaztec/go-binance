package model

// OrderType reflects the available Binance order types
type OrderType string

const (
	// Limit order
	Limit OrderType = "LIMIT"
	// Market order
	Market OrderType = "MARKET"
	// StopLoss order
	StopLoss OrderType = "STOP_LOSS"
	// StopLossLimit order
	StopLossLimit OrderType = "STOP_LOSS_LIMIT"
	// TakeProfit order
	TakeProfit OrderType = "TAKE_PROFIT"
	// TakeProfitLimit order
	TakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
	// LimitMaker order
	LimitMaker OrderType = "LIMIT_MAKER"
)

// OrderResponseType reflects the Binance order response types
type OrderResponseType string

const (
	// Ack order response
	Ack OrderResponseType = "ACK"
	// Result order response
	Result OrderResponseType = "RESULT"
	// Full order response
	Full OrderResponseType = "FULL"
)

// OrderSide reflects the available Binance order sides
type OrderSide string

const (
	// Buy order
	Buy OrderSide = "BUY"
	// Sell order
	Sell OrderSide = "SELL"
)

// TimeInForce reflects the Binance enums how long an order should stay in place
type TimeInForce string

const (
	// GoodTilCanceled order
	GoodTilCanceled TimeInForce = "GTC"
	// ImmediateOrCancel order
	ImmediateOrCancel TimeInForce = "IOC"
	// FillOrKill order
	FillOrKill TimeInForce = "FOK"
)

// ExecutionType enumerates different execution types an order can have
type ExecutionType string

const (
	// New - The order has been accepted into the engine.
	New ExecutionType = "NEW"
	// Canceled - The order has been canceled by the user.
	Canceled ExecutionType = "CANCELED"
	// Replaced - (currently unused)
	Replaced ExecutionType = "REPLACED"
	// Rejected - The order has been rejected and was not processed. (This is never
	// pushed into the User Data Stream)
	Rejected ExecutionType = "REJECTED"
	// Trade - Part of the order or all of the order's quantity has filled.
	Trade ExecutionType = "TRADE"
	// Expired - The order was canceled according to the order type's rules (e.g.
	// LIMIT FOK orders with no fill, LIMIT IOC or MARKET orders that partially
	// fill) or by the exchange, (e.g. orders canceled during liquidation, orders
	// canceled during maintenance)
	Expired ExecutionType = "EXPIRED"
)

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
