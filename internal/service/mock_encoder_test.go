package service

import "github.com/stretchr/testify/mock"

type MockEncoder struct {
	mock.Mock
}

func (m *MockEncoder) Encode(n uint64) string {
	args := m.Called(n)
	return args.String(0)
}

func (m *MockEncoder) Decode(s string) (uint64, error) {
	args := m.Called(s)
	return args.Get(0).(uint64), args.Error(1)
}
