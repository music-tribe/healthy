package healthy

type MockPinger struct {
	Err error
}

func (mp *MockPinger) Ping() error {
	return mp.Err
}
