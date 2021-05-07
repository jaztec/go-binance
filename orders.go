package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jaztec/go-binance/model"
)

const (
	depthPath     = "/api/v3/depth"
	allOrdersPath = "/api/v3/allOrders"
	orderPath     = "/api/v3/order"
	orderTestPath = "/api/v3/order/test"
)

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

func init() {
	requireSignature(allOrdersPath, orderPath, orderTestPath)
}

func (a *api) AllOrders(symbol string, startTime, endTime int64, limit int) (uo []model.UserOrder, err error) {
	if symbol == "" {
		return uo, NoSymbolProvided
	}
	q := NewParameters(5)
	q.Set("symbol", symbol)
	if startTime != 0 {
		q.Set("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime != 0 {
		q.Set("endTime", strconv.FormatInt(endTime, 10))
	}
	if limit != 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	q.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))

	body, err := a.doRequest(http.MethodGet, allOrdersPath, q)
	if err != nil {
		return uo, err
	}

	err = json.Unmarshal(body, &uo)
	if err != nil {
		return uo, fmt.Errorf("encountered error while unmarshaling '%s' into model.UserOrder", body)
	}

	return uo, nil
}

func (a *api) Depth(symbol string, limit int) (o model.Orders, err error) {
	if symbol == "" {
		return o, NoSymbolProvided
	}
	q := NewParameters(2)
	q.Set("symbol", symbol)
	if limit != 0 {
		q.Set("limit", strconv.Itoa(limit))
	}

	body, err := a.doRequest(http.MethodGet, depthPath, q)
	if err != nil {
		return o, err
	}

	err = json.Unmarshal(body, &o)
	if err != nil {
		return o, fmt.Errorf("encountered error while unmarshaling '%s' into model.Orders", body)
	}

	return o, nil
}

// OrderParams hold all optional parameters for a new order. Some parameters
// may still be enforced depending on the OrderType
type OrderParams struct {
	TimeInForce      TimeInForce
	Quantity         float64
	QuoteOrderQty    float64
	Price            float64
	StopPrice        float64
	IcebergQty       float64
	NewOrderRespType OrderResponseType
	RecvWindow       int64
}

func (a *api) Order(symbol string, side OrderSide, orderType OrderType, params OrderParams) (model.OrderResponse, error) {
	return a.doOrder(orderPath, symbol, side, orderType, params)
}

func (a *api) OrderTest(symbol string, side OrderSide, orderType OrderType, params OrderParams) (model.OrderResponse, error) {
	return a.doOrder(orderTestPath, symbol, side, orderType, params)
}

func (a *api) doOrder(path string, symbol string, side OrderSide, orderType OrderType, params OrderParams) (model.OrderResponse, error) {
	if err := checkOrderParams(orderType, params); err != nil {
		return nil, err
	}

	p := NewParameters(11)
	p.Set("symbol", symbol)
	p.Set("side", string(side))
	p.Set("type", string(orderType))

	addOrderParams(p, params)

	p.Set("timestamp", strconv.FormatInt(time.Now().Unix()*1000, 10))

	res, err := a.doRequest(http.MethodPost, path, p)
	if err != nil {
		return nil, err
	}

	if params.NewOrderRespType != "" {
		i := orderResponse(params.NewOrderRespType)
		err = json.Unmarshal(res, i)
		if err != nil {
			return nil, err
		}
		return i.(model.OrderResponse), nil
	}
	i := model.OrderResponseAck{}

	err = json.Unmarshal(res, &i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func orderResponse(t OrderResponseType) interface{} {
	switch t {
	case Ack:
		return &model.OrderResponseAck{}
	case Result:
		return &model.OrderResponseResult{}
	case Full:
		return &model.OrderResponseFull{}
	default:
		return &model.OrderResponseAck{}
	}
}

func addOrderParams(p Parameters, params OrderParams) {
	if params.TimeInForce != "" {
		p.Set("timeInForce", string(params.TimeInForce))
	}
	if params.Quantity != 0.0 {
		p.Set("quantity", fmt.Sprintf("%.8f", params.Quantity))
	}
	if params.QuoteOrderQty != 0.0 {
		p.Set("quoteOrderQty", fmt.Sprintf("%.8f", params.QuoteOrderQty))
	}
	if params.Price != 0.0 {
		p.Set("price", fmt.Sprintf("%.8f", params.Price))
	}
	if params.StopPrice != 0.0 {
		p.Set("stopPrice", fmt.Sprintf("%.8f", params.StopPrice))
	}
	if params.IcebergQty != 0.0 {
		p.Set("icebergQty", fmt.Sprintf("%.8f", params.IcebergQty))
	}
	if params.NewOrderRespType != "" {
		p.Set("newOrderRespType", string(params.NewOrderRespType))
	}
	if params.RecvWindow != 0 {
		p.Set("recvWindow", strconv.FormatInt(params.RecvWindow, 10))
	}
}

func checkOrderParams(ot OrderType, params OrderParams) error {
	check := func(n []string) error {
		missing := make([]string, 0, len(n))
		v := reflect.ValueOf(params)
		for _, p := range n {
			f := v.FieldByName(p)
			if f.IsZero() {
				missing = append(missing, p)
			}
		}

		if len(missing) > 0 {
			return fmt.Errorf("required: %s", strings.Join(missing, ", "))
		}

		return nil
	}

	switch ot {
	case Limit:
		return check([]string{"TimeInForce", "Quantity", "Price"})
	case Market:
		err1 := check([]string{"Quantity"})
		err2 := check([]string{"QuoteOrderQty"})
		if err1 != nil && err2 != nil {
			return errors.New("required: Quantity or QuoteOrderQty")
		}
	case LimitMaker, StopLoss, TakeProfit:
		return check([]string{"Quantity", "StopPrice"})
	case StopLossLimit, TakeProfitLimit:
		return check([]string{"TimeInForce", "Quantity", "Price", "StopPrice"})
	}
	return nil
}
