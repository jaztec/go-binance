package binance

import "fmt"

// APIError encapsulates some expected errors
type APIError struct {
	code int
	msg  string
}

// Satisfy the Error interface
func (bae APIError) Error() string {
	return fmt.Sprintf("msg=%s", bae.msg)
}

var (
	// TooMuchCalls to API, hold back to prevent a ban
	TooMuchCalls = APIError{msg: "too much calls to API"}
	// Blocked IP address for too many calls to the API
	Blocked = APIError{msg: "IP ban active"}
	// AtTimeout for too many calls to the API
	AtTimeout = APIError{msg: "API in timeout now"}
	// NoSymbolProvided in a call that requires one
	NoSymbolProvided = APIError{msg: "no symbol provided"}
)
