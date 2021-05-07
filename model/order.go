package model

type Order struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type OrderResponse interface {
	Symbol() string
	OrderID() int
	OrderListID() int
	ClientOrderID() string
	TransactionTime() int64
}

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

type OrderResponseAck struct {
	Sym          string `json:"symbol"`
	Order        int    `json:"orderId"`
	OrderList    int    `json:"orderListId"`
	ClientOrder  string `json:"clientOrderId"`
	TransactTime int64  `json:"transactTime"`
}

func (r OrderResponseAck) Symbol() string         { return r.Sym }
func (r OrderResponseAck) OrderID() int           { return r.Order }
func (r OrderResponseAck) OrderListID() int       { return r.OrderList }
func (r OrderResponseAck) ClientOrderID() string  { return r.ClientOrder }
func (r OrderResponseAck) TransactionTime() int64 { return r.TransactTime }

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

func (r OrderResponseResult) Symbol() string         { return r.Sym }
func (r OrderResponseResult) OrderID() int           { return r.Order }
func (r OrderResponseResult) OrderListID() int       { return r.OrderList }
func (r OrderResponseResult) ClientOrderID() string  { return r.ClientOrder }
func (r OrderResponseResult) TransactionTime() int64 { return r.TransactTime }

type Fill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}

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

func (r OrderResponseFull) Symbol() string         { return r.Sym }
func (r OrderResponseFull) OrderID() int           { return r.Order }
func (r OrderResponseFull) OrderListID() int       { return r.OrderList }
func (r OrderResponseFull) ClientOrderID() string  { return r.ClientOrder }
func (r OrderResponseFull) TransactionTime() int64 { return r.TransactTime }
