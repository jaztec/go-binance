package binance

import "fmt"

type BinanceAPIError struct {
	msg string
}

func (bae BinanceAPIError) Error() string {
	return fmt.Sprintf("msg=%s", bae.msg)
}

var (
	TooMuchCalls     = BinanceAPIError{"too much calls to API"}
	Blocked          = BinanceAPIError{"IP ban active"}
	AtTimeout        = BinanceAPIError{"API in timeout now"}
	NoSymbolProvided = BinanceAPIError{"no symbol provided"}
)
