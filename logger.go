package binance

import "log"

// Logger is the interface used inside this library. It actually is an interface from the
// go kit logger (https://github.com/go-kit/kit/tree/master/log check it out!). It is mimicked
// here to remove the hard dependency
type Logger interface {
	Log(vals ...interface{}) error
}

type simpleLogger struct{}

func (simpleLogger) Log(vals ...interface{}) error {
	log.Println(vals)
	return nil
}
