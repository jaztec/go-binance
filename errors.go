package binance

import "fmt"

type APIError struct {
	msg string
}

func (bae APIError) Error() string {
	return fmt.Sprintf("msg=%s", bae.msg)
}

var (
	TooMuchCalls     = APIError{"too much calls to API"}
	Blocked          = APIError{"IP ban active"}
	AtTimeout        = APIError{"API in timeout now"}
	NoSymbolProvided = APIError{"no symbol provided"}
)
