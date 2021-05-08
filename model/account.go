package model

import (
	"strconv"
)

// AccountUpdateType defines the messages that can be send via
// the user data stream
type AccountUpdateType string

const (
	// OutboundAccountPositionType is sent any time an account balance has changed and
	// contains the assets that were possibly changed by the event that generated
	// the balance change.
	OutboundAccountPositionType AccountUpdateType = "outboundAccountPosition"

	// BalanceUpdateType shows deposits or withdrawals from the account and transfer
	// of funds between accounts
	BalanceUpdateType AccountUpdateType = "balanceUpdate"

	// ExecutionReportType shows updates on orders
	ExecutionReportType AccountUpdateType = "executionReport"
)

// UserAccountUpdate interface allows multiple typse of messages to be
// bundled inside a single stream
type UserAccountUpdate interface {
	// Type returns the type of message this account update is
	Type() AccountUpdateType
}

// OutboundAccountPosition is sent any time an account balance has changed and
// contains the assets that were possibly changed by the event that generated
// the balance change.
type OutboundAccountPosition struct {
	EventType      string    `json:"e"`
	EventTime      int64     `json:"E"`
	LastUpdateTime int64     `json:"u"`
	Balances       []Balance `json:"B"`
}

// Type returns the type of message this account update is
func (oap OutboundAccountPosition) Type() AccountUpdateType {
	return OutboundAccountPositionType
}

// BalanceUpdate shows deposits or withdrawals from the account and transfer
// of funds between accounts
type BalanceUpdate struct {
	EventType    string `json:"e"`
	EventTime    int64  `json:"E"`
	Asset        string `json:"a"`
	BalanceDelta string `json:"d"`
	ClearTime    int64  `json:"T"`
}

// Type returns the type of message this account update is
func (bu BalanceUpdate) Type() AccountUpdateType {
	return BalanceUpdateType
}

// ExecutionReport shows updates on orders
type ExecutionReport struct {
	EventType                string        `json:"e"`
	EventTime                int64         `json:"E"`
	Symbol                   string        `json:"s"`
	ClientOrderID            string        `json:"c"`
	Side                     OrderSide     `json:"S"`
	OrderType                OrderType     `json:"o"`
	TIF                      TimeInForce   `json:"f"`
	OrderQuantity            string        `json:"q"`
	OrderPrice               string        `json:"p"`
	StopPrice                string        `json:"P"`
	IcebergQuantity          string        `json:"F"`
	OrderListID              int           `json:"g"`
	OriginalClientOrderID    interface{}   `json:"C"`
	CurrentExecutionType     ExecutionType `json:"x"`
	CurrentOrderStatus       string        `json:"X"`
	OrderRejectReason        string        `json:"r"`
	OrderID                  int           `json:"i"`
	LastExecutedQuantity     string        `json:"l"`
	CumulativeFilledQuantity string        `json:"z"`
	LastExecutedPrice        string        `json:"L"`
	CommissionAmount         string        `json:"n"`
	CommissionAsset          interface{}   `json:"N"`
	TransactionTime          int64         `json:"T"`
	TradeID                  int           `json:"t"`
	OnOrderBook              bool          `json:"w"`
	MakerSide                bool          `json:"m"`
	OrderCreationTime        int64         `json:"O"`
	CumulativeQuoteQuantity  string        `json:"Z"`
	LastQuoteQuantity        string        `json:"Y"`
	QuoteOrderQuantity       string        `json:"Q"`
}

// Type returns the type of message this account update is
func (er ExecutionReport) Type() AccountUpdateType {
	return ExecutionReportType
}

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
