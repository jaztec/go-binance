package binance_test

type testLogger struct{}

func (testLogger) Log(l ...interface{}) error {
	//fmt.Println(l...)
	return nil
}
