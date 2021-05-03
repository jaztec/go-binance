package binance_test

import "fmt"

type testLogger struct{}

func (testLogger) Log(keyvals ...interface{}) error {
	fmt.Println(keyvals...)
	return nil
}
