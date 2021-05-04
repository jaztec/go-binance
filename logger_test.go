package binance_test

type testLogger struct{}

func (testLogger) Log(_ ...interface{}) error {
	return nil
}
