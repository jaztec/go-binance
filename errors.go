package binance

import "fmt"

type BinanceAPIError struct {
	msg string
}

func (bae BinanceAPIError) Error() string {
	return fmt.Sprintf("msg=%s", bae.msg)
}

var (
	BinanceTooMuchCalls = BinanceAPIError{"too much calls to API"}
	BinanceBlocked      = BinanceAPIError{"IP ban active"}
	BinanceAtTimeout    = BinanceAPIError{"API in timeout now"}
)
