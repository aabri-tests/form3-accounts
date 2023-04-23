package mocks

import "time"

type MockRetrier struct{}

func (m *MockRetrier) Attempt() int {
	return 1
}

func (m *MockRetrier) NextBackOff() time.Duration {
	return -1
}

func (m *MockRetrier) Reset() {
	return
}

func (m *MockRetrier) RemainingRetries() int {
	return 0
}
