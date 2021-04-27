package model

import (
	"fmt"
)

type BinanceError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (be BinanceError) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", be.Code, be.Msg)
}
