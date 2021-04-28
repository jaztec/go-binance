package model

import (
	"fmt"
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (be Error) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", be.Code, be.Msg)
}
